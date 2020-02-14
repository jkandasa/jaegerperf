package jaegerperf

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const queryRunner = "QueryRunner"

var qJob = Job{js: JobStatus{}, jobType: queryRunner}

// Client to run jaeger query
type Client struct {
	BaseURL    *url.URL
	httpClient *http.Client
	Metrics    map[string][]*Metric
}

// Metric data
type Metric struct {
	URL           string                 `json:"url"`
	QueryParams   map[string]interface{} `json:"queryParams"`
	StatusCode    int                    `json:"statusCode"`
	ContentLength int64                  `json:"contentLength"`
	Elapsed       int64                  `json:"elapsed"`
	ErrorMessage  string                 `json:"errorMessage"`
	Others        map[string]interface{} `json:"others"`
}

// QueryRunnerInput for tests
type QueryRunnerInput struct {
	HostURL       string      `yaml:"hostUrl" json:"hostUrl"`
	CurrentTimeAs time.Time   `yaml:"currentTimeAs" json:"currentTimeAs"`
	Tests         []TestInput `yaml:"tests" json:"tests"`
	Tags          []string    `yaml:"tags" json:"tags"`
}

// TestInput data
type TestInput struct {
	Name        string                 `yaml:"name" json:"name"`
	Type        string                 `yaml:"type" json:"type"`
	Iteration   int                    `yaml:"iteration" json:"iteration"`
	QueryParams map[string]interface{} `yaml:"queryParams" json:"queryParams"`
	StatusCode  int                    `yaml:"statusCode" json:"statusCode"`
}

// MetricSummary data
type MetricSummary struct {
	Name            string  `json:"name"`
	Samples         int     `json:"samples"`
	Elapsed         int64   `json:"elapsed"`
	ErrorsCount     int     `json:"errorsCount"`
	ErrorPercentage float64 `json:"errorPercentage"`
	ContentLength   int64   `json:"contentLength"`
}

// JobResult stores job result data
type JobResult struct {
	Configuration QueryRunnerInput       `json:"configuration"`
	Metrics       map[string]interface{} `json:"metrics"`
}

// IsQueryEngineRuning return bool
func IsQueryEngineRuning() bool {
	return qJob.IsRunning()
}

func (c *Client) timeTrack(name string, metric *Metric) {
	v, ok := c.Metrics[name]
	if !ok {
		c.Metrics[name] = make([]*Metric, 0)
	}
	c.Metrics[name] = append(v, metric)
}

func (c *Client) getMetric(name string, index int) *Metric {
	v, ok := c.Metrics[name]
	if ok {
		return v[index]

	}
	return nil
}

// NewClient for jaeger query service
func NewClient(rawURL string) *Client {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	var tr *http.Transport
	if strings.HasPrefix(rawURL, "https") {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{}
	}
	return &Client{
		BaseURL:    baseURL,
		httpClient: &http.Client{Transport: tr},
		Metrics:    map[string][]*Metric{},
	}
}

func (c *Client) newRequest(test, method, path string, queryParams map[string]interface{}, body interface{}, response interface{}) error {
	rel := &url.URL{Path: fmt.Sprintf("/api%s", path)}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if queryParams != nil {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	}

	start := time.Now()
	resp, err := c.httpClient.Do(req)
	m := Metric{
		URL:         req.URL.String(),
		QueryParams: queryParams,
		Elapsed:     time.Since(start).Microseconds(),
	}
	defer c.timeTrack(test, &m)
	if err != nil {
		m.ErrorMessage = err.Error()
		return err
	}
	m.StatusCode = resp.StatusCode
	m.ContentLength = resp.ContentLength

	defer resp.Body.Close()
	//err = json.NewDecoder(resp.Body).Decode(v)
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	m.ContentLength = int64(len(respBytes))
	err = json.Unmarshal(respBytes, &response)
	if err != nil {
		return err
	}
	return nil
}

// Services lists available services
func (c *Client) Services(test string) (map[string]interface{}, error) {
	resp := make(map[string]interface{})
	err := c.newRequest(test, "GET", "/services", nil, nil, &resp)
	return resp, err
}

// Search traces with given filter
func (c *Client) Search(test string, queryParams map[string]interface{}) (map[string]interface{}, error) {
	resp := make(map[string]interface{})
	err := c.newRequest(test, "GET", "/traces", queryParams, nil, &resp)
	return resp, err
}

// update startTime and endTime
func updateQueryParams(currentTime time.Time, queryParams map[string]interface{}) {
	for k, v := range queryParams {
		if k == "start" || k == "end" {
			d, err := time.ParseDuration(v.(string))
			if err != nil {
				fmt.Println(err)
			} else {
				queryParams[k] = uint64(currentTime.Add(d).UnixNano() / 1000) // set it in microseconds
			}
		}
	}
}

// ExecuteQueryTest runs set of requests
func ExecuteQueryTest(jobID string, input QueryRunnerInput) (map[string]interface{}, error) {
	jResult := JobResult{Configuration: input}
	qJob.SetStatus(true, jobID, jResult)
	defer qJob.SetCompleted()
	UpdateTags(input.Tags...)
	if input.CurrentTimeAs.IsZero() {
		input.CurrentTimeAs = time.Now()
	}
	c := NewClient(input.HostURL)
	tests := make(map[string]TestInput)
	for _, t := range input.Tests {
		tests[t.Name] = t
		updateQueryParams(input.CurrentTimeAs, t.QueryParams)
		for count := 0; count < t.Iteration; count++ {
			switch t.Type {
			case "search":
				d, err := c.Search(t.Name, t.QueryParams)
				if err == nil {
					m := c.getMetric(t.Name, count)
					if m != nil {
						od := make(map[string]interface{})
						od["errors"] = d["errors"]
						if d["data"] != nil {
							od["count"] = len(d["data"].([]interface{}))
						}
						m.Others = od
					}
				}
			case "services":
				_, err := c.Services(t.Name)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	r := make(map[string]interface{})
	r["raw"] = c.Metrics
	s := make([]MetricSummary, 0)
	for k, m := range c.Metrics {
		t := tests[k]
		var el int64
		var errors int
		var cl int64
		for _, mt := range m {
			el += mt.Elapsed
			cl += mt.ContentLength
			if mt.StatusCode != t.StatusCode || mt.ErrorMessage != "" {
				errors++
			}
		}
		s = append(s,
			MetricSummary{
				Name:            k,
				Samples:         len(m),
				Elapsed:         el / int64(len(m)),
				ErrorsCount:     errors,
				ErrorPercentage: PercentOf(errors, t.Iteration),
				ContentLength:   cl / int64(len(m)),
			})
	}
	r["summary"] = s
	// update summary data
	if input.Tags != nil && len(input.Tags) > 0 {
		cd := &CustomData{
			Tags: input.Tags,
			Type: queryRunner,
			Data: s,
		}
		DumpCustom(fmt.Sprintf("summary_%s", jobID), cd)
	}
	jResult.Metrics = r
	qJob.SetStatus(true, jobID, jResult)
	return r, nil
}

// PercentOf returns percentage
func PercentOf(part int, total int) float64 {
	return (float64(part) * float64(100)) / float64(total)
}
