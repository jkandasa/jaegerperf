package jaegerperf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
	"gopkg.in/yaml.v2"
)

// StartHandler to provide http api
func StartHandler() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/queryRunner", executeQueryTest)
	mux.HandleFunc("/api/spansGenerator", generateSpans)
	mux.HandleFunc("/api/jobs", listJobs)
	mux.HandleFunc("/api/jobs/delete", deleteJob)
	mux.HandleFunc("/api/status", status)

	fs := http.FileServer(http.Dir("/app/ui"))
	mux.Handle("/", fs)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	// Insert the middleware
	handler := c.Handler(mux)
	//handler := cors.Default().Handler(mux)

	fmt.Println("Listening...")
	return http.ListenAndServe(":8080", handler)
}

func getPayload(r *http.Request, out interface{}) error {
	d, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}
	ct := r.Header.Get("Content-Type")
	switch strings.ToLower(ct) {
	case "application/json":
		err = json.Unmarshal(d, out)
		if err != nil {
			return err
		}
	case "application/yaml":
		err = yaml.Unmarshal(d, out)
		if err != nil {
			return nil
		}
	default:
		return fmt.Errorf("Unknown format: %s", ct)
	}
	return nil
}

func listJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	d, err := ListJobs()
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
	err := DeleteJob(jobID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write([]byte(fmt.Sprintf("record[%s] deleted successfully", jobID)))
}

func status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s := map[string]interface{}{
		"time": time.Now(),
	}
	hn, err := os.Hostname()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	s["hostname"] = hn
	od, err := json.Marshal(&s)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(od)
}

func generateSpans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := GeneratorConfiguration{}
	err := getPayload(r, &d)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if IsGeneratorRunning() {
		http.Error(w, "A request is in progress", 503)
		return
	}
	jobID := uuid.New().String()
	go func() {
		err = ExecuteSpansGenerator(jobID, d)
		if err != nil {
			fmt.Println(err)
		}
	}()
	setJobID(jobID, w)
}

func executeQueryTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d := QueryRunnerInput{}
	err := getPayload(r, &d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if IsQueryEngineRuning() {
		http.Error(w, "A request is in progress", 503)
		return
	}

	jobID := uuid.New().String()
	go func() {
		_, err := ExecuteQueryTest(jobID, d)
		if err != nil {
			fmt.Println(err)
		}
	}()
	setJobID(jobID, w)
}

func setJobID(jobID string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	od, err := json.Marshal(map[string]string{"jobID": jobID})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(od)
}
