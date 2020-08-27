package handler

import (
	"encoding/json"
	"fmt"
	"jaegerperf/pkg/api/job"
	"net/http"

	"github.com/gorilla/mux"
)

func registerJobRoutes(router *mux.Router) {
	router.HandleFunc("/api/jobs", listJobs).Methods(http.MethodGet)
	router.HandleFunc("/api/jobs/delete", deleteJob).Methods(http.MethodDelete)
}

func listJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d, err := job.ListJobs()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}

func deleteJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/text")

	jobID := r.URL.Query().Get("jobId")
	err := job.DeleteJob(jobID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf("record[%s] deleted successfully", jobID)))
}
