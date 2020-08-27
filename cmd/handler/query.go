package handler

import (
	"encoding/json"
	qapi "jaegerperf/pkg/api/query"
	ml "jaegerperf/pkg/model"
	qml "jaegerperf/pkg/model/query"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func registerQueryRoutes(router *mux.Router) {
	router.HandleFunc("/api/query", executeQueryTest).Methods(http.MethodPost)
	router.HandleFunc("/api/template/query", listQueryTemplates).Methods(http.MethodGet)
	router.HandleFunc("/api/template/query/{filename}", getQueryTemplates).Methods(http.MethodGet)
	router.HandleFunc("/api/template/query", saveQueryTemplates).Methods(http.MethodPost)
	router.HandleFunc("/api/query/tags", listTags).Methods(http.MethodGet)
	router.HandleFunc("/api/query/summary", listQueryMetrics).Methods(http.MethodGet)
}

func listQueryTemplates(w http.ResponseWriter, r *http.Request) {
	listFiles(w, r, ml.FileTemplatesQuery)
}

func getQueryTemplates(w http.ResponseWriter, r *http.Request) {
	readFile(w, r, ml.FileTemplatesQuery)
}

func saveQueryTemplates(w http.ResponseWriter, r *http.Request) {
	saveFile(w, r, ml.FileTemplatesQuery)
}

func executeQueryTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := qml.InputConfig{}
	err := getPayload(r, &d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if qapi.IsRunning() {
		http.Error(w, "A request is in progress", 503)
		return
	}

	jobID := uuid.New().String()
	go func() {
		err := qapi.ExecuteQueryTest(jobID, d)
		if err != nil {
			zap.L().Error("Failed to execute query runner", zap.Error(err))
		}
	}()
	setJobID(jobID, w)
}

func listTags(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_tags := qapi.ListTags()
	od, err := json.Marshal(&_tags)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(od)
}

func listQueryMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filterTags := r.URL.Query()["tags"]
	_ft := make([]string, len(filterTags))
	for _, t := range filterTags {
		_ft = append(_ft, strings.ToLower(t))
	}
	d, err := qapi.ListMetricQuickReport(_ft...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	od, err := json.Marshal(&d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(od)
}
