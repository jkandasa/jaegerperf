package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ml "jaegerperf/pkg/model"
	jobml "jaegerperf/pkg/model/job"
	"jaegerperf/pkg/util"
	"os"
)

// ListJobs return available job details
func ListJobs() ([]jobml.Report, error) {
	util.CreateDir(ml.FileJobs)
	jobs := make([]jobml.Report, 0)
	files, err := ioutil.ReadDir(ml.FileJobs)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		d := jobml.Report{}
		b, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", ml.FileJobs, file.Name()))
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
	os.Remove(fmt.Sprintf("%s/%s.json", ml.FileMetricQuickReport, jobID))
	// remove the detailed file
	return os.Remove(fmt.Sprintf("%s/%s.json", ml.FileJobs, jobID))
}
