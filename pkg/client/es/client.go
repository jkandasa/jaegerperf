package es

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client struct
type Client struct {
	Config     string
	httpClient *http.Client
	BaseURL    *url.URL
}

// Init func
func Init() (*Client, error) {
	client := &Client{
		BaseURL:    &url.URL{Host: "elasticsearch:9200", Scheme: "http"},
		httpClient: http.DefaultClient,
	}
	return client, nil
}

// Info func
func (c *Client) Info() (interface{}, error) {
	out := make(map[string]interface{})
	err := c.newRequest(http.MethodGet, "", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ListIndices func
func (c *Client) ListIndices() (interface{}, error) {
	out := make([]map[string]interface{}, 0)
	err := c.newRequest(http.MethodGet, "/_cat/indices", nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecordCount func
func (c *Client) RecordCount(indexName string) (interface{}, error) {
	_, err := c.RefreshIndex(indexName)
	if err != nil {
		return nil, err
	}

	out := make(map[string]interface{})
	err = c.newRequest(http.MethodGet, fmt.Sprintf("/%s/_count", indexName), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RefreshIndex func
func (c *Client) RefreshIndex(indexName string) (interface{}, error) {
	out := make(map[string]interface{})
	err := c.newRequest(http.MethodGet, fmt.Sprintf("/%s/_refresh", indexName), nil, nil, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *Client) newRequest(method, path string, queryParams map[string]interface{}, body interface{}, out interface{}) error {
	rel := &url.URL{Path: path}
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		// failure
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(respBytes, &out)
}
