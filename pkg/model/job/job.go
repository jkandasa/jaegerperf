package job

import (
	"jaegerperf/pkg/model"
	"jaegerperf/pkg/util"
	"sync"
	"time"

	"go.uber.org/zap"
)

// constants
const (
	// Job types
	JobTypeSpansGenerator = "spans_generator"
	JobTypeQueryRunner    = "query_runner"
)

// Job struct
type Job struct {
	ID           string
	Type         string
	isRunning    bool
	modifiedTime time.Time
	data         interface{}
	mutex        sync.RWMutex
}

// New creates new job
func New(ID, Type string) *Job {
	job := &Job{
		ID:           ID,
		Type:         Type,
		isRunning:    true,
		modifiedTime: time.Now(),
	}
	job.store()
	return job
}

// IsRunning status
func (j *Job) IsRunning() bool {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	return j.isRunning
}

// Update job details
func (j *Job) Update(isRunning bool, data interface{}) {
	j.mutex.Lock()
	defer j.mutex.Unlock()
	j.isRunning = isRunning
	if data != nil {
		j.data = data
	}
	j.modifiedTime = time.Now()
	j.store()
}

// store it into disk
func (j *Job) store() {
	// update in to disk
	report := Report{
		ID:           j.ID,
		Type:         j.Type,
		IsRunning:    j.isRunning,
		ModifiedTime: j.modifiedTime,
		Data:         j.data,
	}
	err := util.StoreJSON(model.FileJobs, j.ID, &report)
	if err != nil {
		zap.L().Error("Failed to save job details on disk", zap.Error(err))
	}
}

// Report used to keep job data on disk
type Report struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	IsRunning    bool        `json:"isRunning"`
	ModifiedTime time.Time   `json:"modifiedTime"`
	Data         interface{} `json:"data"`
}
