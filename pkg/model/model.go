package model

import "time"

// Files used to strore job data
const (
	FileResources          = "/app/resources"
	FileBase               = "/app/data"
	FileJobs               = "/app/data/jobs"
	FileMetricQuickReport  = "/app/data/metric-quick-report"
	FileOthers             = "/app/data/others"
	FileTemplatesGenerator = "/app/data/templates/generator"
	FileTemplatesQuery     = "/app/data/templates/query"
)

// Status struct
type Status struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	TimeTaken string    `json:"timeTaken"`
	IsSuccess bool      `json:"isSuccess"`
	Message   string    `json:"message"`
}

// File struct
type File struct {
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	ModifiedTime time.Time `json:"modifiedTime"`
	Data         string    `json:"data"`
}
