package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/gorilla/mux"
)

func registerStatusRoutes(router *mux.Router) {
	router.HandleFunc("/api/status", status).Methods(http.MethodGet)
	router.HandleFunc("/api/tags", listTags).Methods(http.MethodGet)
}

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type status struct {
		Time      time.Time `json:"time"`
		Hostname  string    `json:"hostname"`
		GoVersion string    `json:"goVersion"`
	}

	hn, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	s := status{
		Time:      time.Now(),
		Hostname:  hn,
		GoVersion: runtime.Version(),
	}

	od, err := json.Marshal(&s)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(od)
}
