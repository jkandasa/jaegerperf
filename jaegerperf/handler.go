package jaegerperf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// StartHandler to provide http api
func StartHandler() error {
	http.HandleFunc("/executeQueryTest", executeQueryTest)
	http.HandleFunc("/generateSpans", generateSpans)
	return http.ListenAndServe(":8080", nil)
}

func generateSpans(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	ExecuteSpansGenerator(string(d))
}

func executeQueryTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	m, err := ExecuteQueryTest(string(d))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	output, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(output)
}
