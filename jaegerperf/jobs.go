package jaegerperf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"
)

const flJobs = "/tmp/jaegerperf/jobs"
const flCustomData = "/tmp/jaegerperf/custom-data"
const flOthers = "/tmp/jaegerperf/others"

// JobData for json file
type JobData struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	Data         interface{} `json:"data"`
	ModifiedTime time.Time   `json:"modifiedTime"`
}

// CustomData definition
type CustomData struct {
	Tags []string    `json:"tags"`
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

func createRootLocation(rootLocation string) {
	if _, err := os.Stat(rootLocation); os.IsNotExist(err) {
		os.MkdirAll(rootLocation, 0775)
	}
}

// DumpCustom data
func dumpData(rootLocation, filename string, data interface{}) error {
	createRootLocation(rootLocation)
	file, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s.json", rootLocation, filename), file, 0644)
	return err
}

// Update stores data into disk
func (j *JobData) Update() error {
	j.ModifiedTime = time.Now()
	return dumpData(flJobs, j.ID, j)
}

// DumpCustom data
func DumpCustom(filename string, data interface{}) error {
	return dumpData(flCustomData, filename, data)
}

// ListCustomData details
func ListCustomData(filePrefix string, filterTags ...string) ([]CustomData, error) {
	createRootLocation(flCustomData)
	if filePrefix == "" {
		return nil, errors.New("file prefix not supplied")
	}
	customData := make([]CustomData, 0)
	files, err := ioutil.ReadDir(flCustomData)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), filePrefix) {
			d := CustomData{}
			b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", flCustomData, file.Name()))
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(b, &d)
			if err != nil {
				return customData, err
			}
			if len(filterTags) == 0 {
				customData = append(customData, d)
			} else {
				for _, ft := range filterTags {
					for _, st := range d.Tags {
						if ft == st {
							customData = append(customData, d)
						}
					}
				}
			}
		}
	}
	return customData, nil
}

// ListJobs return available job details
func ListJobs() ([]JobData, error) {
	createRootLocation(flJobs)
	jobs := make([]JobData, 0)
	files, err := ioutil.ReadDir(flJobs)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		d := JobData{}
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", flJobs, file.Name()))
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(b, &d)
		jobs = append(jobs, d)
	}
	return jobs, nil
}

// DeleteJob remove from disk
func DeleteJob(jobID string) error {
	// remove summary data, do not care about the error
	os.Remove(fmt.Sprintf("%s/summary_%s.json", flCustomData, jobID))
	// remove the detailed file
	return os.Remove(fmt.Sprintf("%s/%s.json", flJobs, jobID))
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

var _tags = make([]string, 0)

// ListTags returns available tags
func ListTags() []string {
	return _tags
}

// UpdateTags adds new tags
func UpdateTags(nTags ...string) {
	d := map[string]bool{}
	newTags := append(_tags, nTags...)
	for _, t := range newTags {
		d[strings.ToLower(t)] = true
	}

	_newTags := make([]string, 0)
	for k := range d {
		_newTags = append(_newTags, k)
	}
	_tags = _newTags
	// Update tags to disk
	dumpData(flOthers, "tags", _tags)
}

// InitJobData from disk
func InitJobData() error {
	// load tags from disk
	createRootLocation(flOthers)
	t := make([]string, 0)

	b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s.json", flOthers, "tags"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	_tags = t
	return nil
}
