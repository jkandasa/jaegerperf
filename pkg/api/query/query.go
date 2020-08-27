package query

import (
	"errors"
	ml "jaegerperf/pkg/model"
	jml "jaegerperf/pkg/model/job"
	qml "jaegerperf/pkg/model/query"
	"jaegerperf/pkg/util"
	"time"

	"go.uber.org/zap"
)

var qJob *jml.Job

// IsRunning return bool
func IsRunning() bool {
	if qJob != nil {
		return qJob.IsRunning()
	}
	return false
}

// update startTime and endTime
func updateQueryParams(currentTime time.Time, queryParams map[string]interface{}) {
	for k, v := range queryParams {
		if k == "start" || k == "end" {
			d, err := time.ParseDuration(v.(string))
			if err != nil {
				zap.L().Error("Failed to parse duration", zap.Error(err))
			} else {
				queryParams[k] = uint64(currentTime.Add(d).UnixNano() / 1000) // set it in microseconds
			}
		}
	}
}

// ExecuteQueryTest runs set of requests
func ExecuteQueryTest(jobID string, input qml.InputConfig) error {
	if IsRunning() {
		return errors.New("A query job is in running state")
	}
	qJob = jml.New(jobID, jml.JobTypeQueryRunner)
	report := qml.ExecutionReport{Input: input}
	// update job status on exit
	defer qJob.Update(false, &report)

	execStart := time.Now()

	updateStatus := func(err error) {
		report.Status.StartTime = execStart
		report.Status.EndTime = time.Now()
		report.Status.TimeTaken = time.Since(execStart).String()
		if err != nil {
			report.Status.IsSuccess = false
			report.Status.Message = err.Error()
		} else {
			report.Status.IsSuccess = true
		}
	}

	updateTags(input.Tags...)
	if input.CurrentTimeAs.IsZero() {
		input.CurrentTimeAs = time.Now()
	}
	c, err := newClient(input.HostURL)
	if err != nil {
		updateStatus(err)
		return err
	}
	queries := make(map[string]qml.Config)
	for _, qry := range input.Query {
		queries[qry.Name] = qry
		updateQueryParams(input.CurrentTimeAs, qry.QueryParams)
		for count := 0; count < qry.RepeatCount; count++ {
			switch qry.Type {
			case "search":
				d, err := c.Search(qry.Name, qry.QueryParams)
				if err == nil {
					m := c.getMetric(qry.Name, count)
					if m != nil {
						od := make(map[string]interface{})
						od["errors"] = d["errors"]
						if d["data"] != nil {
							od["count"] = len(d["data"].([]interface{}))
						}
						m.Others = od
					}
				}
			case "services":
				_, err := c.Services(qry.Name)
				if err != nil {
					updateStatus(err)
					return err
				}
			}
		}
	}
	report.Report.Raw = c.Metrics
	s := make([]qml.MetricSummary, 0)
	for k, m := range c.Metrics {
		t := queries[k]
		var el int64
		var errors int
		var cl int64
		for _, mt := range m {
			el += mt.Elapsed
			cl += mt.ContentLength
			if mt.StatusCode != t.StatusCode || mt.ErrorMessage != "" {
				errors++
			}
		}
		s = append(s,
			qml.MetricSummary{
				Name:            k,
				Samples:         len(m),
				Elapsed:         el / int64(len(m)),
				ErrorsCount:     errors,
				ErrorPercentage: util.PercentOf(errors, t.RepeatCount),
				ContentLength:   cl / int64(len(m)),
			})
	}
	report.Report.Summary = s

	// update status
	updateStatus(nil)

	// keep metricQuickSummary copy
	if input.Tags != nil && len(input.Tags) > 0 {
		qr := &qml.MetricQuickReport{
			JobID:   jobID,
			Tags:    input.Tags,
			Summary: s,
			Status:  report.Status,
		}
		err := util.StoreJSON(ml.FileMetricQuickReport, jobID, qr)
		if err != nil {
			zap.L().Error("Failed to dump metric quick summary", zap.Error(err))
			return nil
		}
	}
	return nil
}
