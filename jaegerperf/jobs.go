package jaegerperf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

const fileLocation = "/tmp/jaegerperf"

// JobData for json file
type JobData struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	Data         interface{} `json:"data"`
	ModifiedTime time.Time   `json:"modifiedTime"`
}

func createRootLocation() {
	if _, err := os.Stat(fileLocation); os.IsNotExist(err) {
		os.MkdirAll(fileLocation, 0775)
	}
}

// Update stores data into disk
func (j *JobData) Update() error {
	createRootLocation()
	j.ModifiedTime = time.Now()
	file, err := json.MarshalIndent(j, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s_%s.json", fileLocation, j.ID, j.Type), file, 0644)
	return err
}

// ListJobs return available job details
func ListJobs() ([]JobData, error) {
	createRootLocation()
	jobs := make([]JobData, 0)
	files, err := ioutil.ReadDir(fileLocation)
	/*
		err := filepath.Walk(
			fmt.Sprintf("%s/", fileLocation),
			func(path string, info os.FileInfo, err error) error {
				files = append(files, path)
				return nil
			})
	*/
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		d := JobData{}
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", fileLocation, file.Name()))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &d)
		jobs = append(jobs, d)
	}
	return jobs, nil
}

// JobStatus displays last job status
type JobStatus struct {
	IsRunning      bool        `json:"isRunning"`
	JobID          string      `json:"jobId"`
	Data           interface{} `json:"data"`
	CompletionTime time.Time   `json:"completionTime"`
}

// Job actions
type Job struct {
	js      JobStatus
	mu      sync.RWMutex
	jobType string
}

// GetStatus returns status
func (j *Job) GetStatus() JobStatus {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.js
}

// IsRunning job status
func (j *Job) IsRunning() bool {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return j.js.IsRunning
}

// SetStatus sets value
func (j *Job) SetStatus(isRunning bool, jobID string, data interface{}) {
	j.mu.Lock()
	j.js.IsRunning = isRunning
	j.js.JobID = jobID
	j.js.Data = data
	if !isRunning {
		j.js.CompletionTime = time.Now()
	}
	j.mu.Unlock()
	jsc := j.GetStatus()
	jd := JobData{ID: jobID, Type: j.jobType, Data: jsc}
	err := jd.Update()
	if err != nil {
		fmt.Println(err)
	}
}

// SetCompleted update as completed
func (j *Job) SetCompleted() {
	j.mu.Lock()
	j.js.IsRunning = false
	j.js.CompletionTime = time.Now()
	j.mu.Unlock()
	jsc := j.GetStatus()
	jd := JobData{ID: jsc.JobID, Type: j.jobType, Data: jsc}
	err := jd.Update()
	if err != nil {
		fmt.Println(err)
	}
}
