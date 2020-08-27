package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	ml "jaegerperf/pkg/model"
	"jaegerperf/pkg/util"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func listFiles(w http.ResponseWriter, r *http.Request, dir string) {
	w.Header().Set("Content-Type", "application/json")

	d, err := util.ListFiles(dir)
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

func readFile(w http.ResponseWriter, r *http.Request, dir string) {
	w.Header().Set("Content-Type", "application/json")
	f := mux.Vars(r)
	filename, ok := f["filename"]
	if !ok {
		http.Error(w, "Filename should be supplied on the url", 500)
		return
	}

	data, err := util.ReadFile(dir, filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	file := ml.File{Name: filename, Data: string(data)}

	out, err := json.Marshal(file)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.Write(out)
}

func saveFile(w http.ResponseWriter, r *http.Request, dir string) {
	w.Header().Set("Content-Type", "application/json")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out := &ml.File{}
	err = json.Unmarshal(d, out)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = util.WriteFile(dir, out.Name, []byte(out.Data))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
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

func setJobID(jobID string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	od, err := json.Marshal(map[string]string{"jobID": jobID})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(od)
}
