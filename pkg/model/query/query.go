package query

import (
	ml "jaegerperf/pkg/model"
	"time"
)

// InputConfig for tests
type InputConfig struct {
	HostURL       string    `json:"hostUrl" yaml:"host_url" `
	CurrentTimeAs time.Time `json:"currentTimeAs" yaml:"current_time_as" `
	Tags          []string  `json:"tags" yaml:"tags"`
	Query         []Config  `json:"query" yaml:"query"`
}

// Config data
type Config struct {
	Name        string                 `json:"name" yaml:"name"`
	Type        string                 `json:"type" yaml:"type"`
	RepeatCount int                    `json:"repeatCount" yaml:"repeat_count"`
	QueryParams map[string]interface{} `json:"queryParams" yaml:"query_params"`
	StatusCode  int                    `json:"statusCode" yaml:"status_code"`
}

// ExecutionReport struct
type ExecutionReport struct {
	Input  InputConfig  `json:"input"`
	Status ml.Status    `json:"status"`
	Report MetricReport `json:"report"`
}

// MetricReport struct
type MetricReport struct {
	Raw     map[string][]*MetricRaw `json:"raw"`
	Summary []MetricSummary         `json:"summary"`
}

// MetricRaw data
type MetricRaw struct {
	URL           string                 `json:"url"`
	QueryParams   map[string]interface{} `json:"queryParams"`
	StatusCode    int                    `json:"statusCode"`
	ContentLength int64                  `json:"contentLength"`
	Elapsed       int64                  `json:"elapsed"`
	ErrorMessage  string                 `json:"errorMessage"`
	Others        map[string]interface{} `json:"others"`
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

// MetricQuickReport struct
type MetricQuickReport struct {
	JobID   string          `json:"jobId"`
	Tags    []string        `json:"tags"`
	Summary []MetricSummary `json:"summary"`
	Status  ml.Status       `json:"status"`
}
