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
}

// QueryRunnerInput for tests
type QueryRunnerInput struct {
	HostURL       string      `yaml:"hostUrl" json:"hostUrl"`
	CurrentTimeAs time.Time   `yaml:"currentTimeAs" json:"currentTimeAs"`
	Tests         []TestInput `yaml:"tests" json:"tests"`
}

// TestInput data
type TestInput struct {
	Name        string                 `yaml:"name" json:"name"`
	Type        string                 `yaml:"type" json:"type"`
	Iteration   int                    `yaml:"iteration" json:"iteration"`
	QueryParams map[string]interface{} `yaml:"queryParams" json:"queryParams"`
}

// MetricSummary data
type MetricSummary struct {
	Name    string `json:"name"`
	Samples int    `json:"samples"`
	Elapsed int64  `json:"elapsed"`
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
	if err != nil {
		return err
	}
	m := Metric{URL: req.URL.String(), QueryParams: queryParams, StatusCode: resp.StatusCode, ContentLength: resp.ContentLength, Elapsed: time.Since(start).Microseconds()}
	c.timeTrack(test, &m)

	if err != nil {
		return err
	}
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
	if input.CurrentTimeAs.IsZero() {
		input.CurrentTimeAs = time.Now()
	}
	c := NewClient(input.HostURL)
	for _, t := range input.Tests {
		updateQueryParams(input.CurrentTimeAs, t.QueryParams)
		for count := 0; count < t.Iteration; count++ {
			switch t.Type {
			case "search":
				_, err := c.Search(t.Name, t.QueryParams)
				if err != nil {
					return nil, err
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
		var el int64
		for _, mt := range m {
			el += mt.Elapsed
		}
		s = append(s, MetricSummary{Name: k, Samples: len(m), Elapsed: el / int64(len(m))})
	}
	r["summary"] = s
	jResult.Metrics = r
	qJob.SetStatus(true, jobID, jResult)
	return r, nil
}
