package query

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	qml "jaegerperf/pkg/model/query"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// client to run jaeger query
type client struct {
	BaseURL    *url.URL
	httpClient *http.Client
	Metrics    map[string][]*qml.MetricRaw
}

// newClient for jaeger query service
func newClient(rawURL string) (*client, error) {
	baseURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	var tr *http.Transport
	if strings.HasPrefix(rawURL, "https") {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	} else {
		tr = &http.Transport{}
	}
	return &client{
		BaseURL:    baseURL,
		httpClient: &http.Client{Transport: tr},
		Metrics:    map[string][]*qml.MetricRaw{},
	}, nil
}

func (c *client) timeTrack(name string, metric *qml.MetricRaw) {
	v, ok := c.Metrics[name]
	if !ok {
		c.Metrics[name] = make([]*qml.MetricRaw, 0)
	}
	c.Metrics[name] = append(v, metric)
}

func (c *client) getMetric(name string, index int) *qml.MetricRaw {
	v, ok := c.Metrics[name]
	if ok {
		return v[index]

	}
	return nil
}

func (c *client) newRequest(test, method, path string, queryParams map[string]interface{}, body interface{}, response interface{}) error {
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
	m := qml.MetricRaw{
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
func (c *client) Services(test string) (map[string]interface{}, error) {
	resp := make(map[string]interface{})
	err := c.newRequest(test, "GET", "/services", nil, nil, &resp)
	return resp, err
}

// Search traces with given filter
func (c *client) Search(test string, queryParams map[string]interface{}) (map[string]interface{}, error) {
	resp := make(map[string]interface{})
	err := c.newRequest(test, "GET", "/traces", queryParams, nil, &resp)
	return resp, err
}
