package generator

import (
	ml "jaegerperf/pkg/model"
	"time"
)

// constants
const (
	EndpointTypeAgent     = "agent"
	EndpointTypeCollector = "collector"
)

// InputConfig to send spans
type InputConfig struct {
	Endpoint    EndpointConfig `json:"endpoint" yaml:"endpoint"`
	StartTime   time.Time      `json:"startTime" yaml:"start_time"`
	Mode        string         `json:"mode" yaml:"mode"`
	Realtime    RealtimeConfig `json:"realtime" yaml:"realtime"`
	History     HistoryConfig  `json:"history" yaml:"history"`
	SpansConfig SpansConfig    `json:"spansConfig" yaml:"spans_config"`
}

// EndpointConfig struct
type EndpointConfig struct {
	Type string `json:"type" yaml:"type"`
	URL  string `json:"url" yaml:"url"`
}

// RealtimeConfig struct
type RealtimeConfig struct {
	Duration string `json:"duration" yaml:"duration"`
}

// HistoryConfig struct
type HistoryConfig struct {
	Days        int `json:"days" yaml:"days"`
	SpansPerDay int `json:"spansPerDay" yaml:"spans_per_day"`
}

// SpansConfig struct
type SpansConfig struct {
	ServiceName    string                 `json:"serviceName" yaml:"service_name"`
	SpansPerSecond int                    `json:"spansPerSecond" yaml:"spans_per_second"`
	ChildDepth     int                    `json:"childDepth" yaml:"child_depth"`
	Tags           map[string]interface{} `json:"tags" yaml:"tags"`
}

// Report struct
type Report struct {
	SpansSent int `json:"spansSent"`
}

// ExecutionReport struct
type ExecutionReport struct {
	Input  InputConfig `json:"input"`
	Status ml.Status   `json:"status"`
	Report Report      `json:"report"`
}
