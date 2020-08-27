package query

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ml "jaegerperf/pkg/model"
	qml "jaegerperf/pkg/model/query"
	"jaegerperf/pkg/util"
	"strings"

	"go.uber.org/zap"
)

var _tags = make([]string, 0)

// ListTags returns available tags
func ListTags() []string {
	return _tags
}

// updateTags adds new tags
func updateTags(nTags ...string) {
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
	util.StoreJSON(ml.FileOthers, "tags", _tags)
}

// LoadTags from disk
func LoadTags() error {
	// load tags from disk
	util.CreateDir(ml.FileOthers)
	t := make([]string, 0)
	filename := fmt.Sprintf("%s/%s.json", ml.FileOthers, "tags")
	if !util.IsFileExists(filename) {
		zap.L().Debug("Tags file not available", zap.String("filename", filename))
		return nil
	}
	b, err := ioutil.ReadFile(filename)
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

// ListMetricQuickReport func
func ListMetricQuickReport(filterTags ...string) ([]qml.MetricQuickReport, error) {
	util.CreateDir(ml.FileMetricQuickReport)
	reports := make([]qml.MetricQuickReport, 0)
	files, err := ioutil.ReadDir(ml.FileMetricQuickReport)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		report := qml.MetricQuickReport{}
		err = util.LoadJSON(ml.FileMetricQuickReport, file.Name(), &report)
		if err != nil {
			return reports, err
		}
		if len(filterTags) == 0 {
			reports = append(reports, report)
		} else {
			for _, ft := range filterTags {
				for _, st := range report.Tags {
					if ft == st {
						reports = append(reports, report)
					}
				}
			}
		}
	}
	return reports, nil
}
