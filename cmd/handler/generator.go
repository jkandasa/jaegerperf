package handler

import (
	gapi "jaegerperf/pkg/api/generator"
	ml "jaegerperf/pkg/model"
	gml "jaegerperf/pkg/model/generator"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func registerGeneratorRoutes(router *mux.Router) {
	router.HandleFunc("/api/generator", generateSpans).Methods(http.MethodPost)
	router.HandleFunc("/api/template/generator", listGeneratorTemplates).Methods(http.MethodGet)
	router.HandleFunc("/api/template/generator/{filename}", getGeneratorTemplates).Methods(http.MethodGet)
	router.HandleFunc("/api/template/generator", saveGeneratorTemplates).Methods(http.MethodPost)
}

func generateSpans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := gml.InputConfig{}
	err := getPayload(r, &d)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if gapi.IsRunning() {
		http.Error(w, "A request is in progress", 503)
		return
	}
	jobID := uuid.New().String()
	go func() {
		err = gapi.ExecuteSpansGenerator(jobID, d)
		if err != nil {
			zap.L().Error("Failed to execute spans generator", zap.Error(err))
		}
	}()
	setJobID(jobID, w)
}

func listGeneratorTemplates(w http.ResponseWriter, r *http.Request) {
	listFiles(w, r, ml.FileTemplatesGenerator)
}

func getGeneratorTemplates(w http.ResponseWriter, r *http.Request) {
	readFile(w, r, ml.FileTemplatesGenerator)
}

func saveGeneratorTemplates(w http.ResponseWriter, r *http.Request) {
	saveFile(w, r, ml.FileTemplatesGenerator)
}
