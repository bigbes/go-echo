package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type request struct {
	URL      string                 `json:"url"`
	Method   string                 `json:"method"`
	Headers  http.Header            `json:"headers"`
	Body     string                 `json:"body"`
	BodyJSON map[string]interface{} `json:"body_json,omitempty"`
}

func handle(rw http.ResponseWriter, r *http.Request) {
	rr := &request{}
	rr.Method = r.Method
	rr.Headers = r.Header
	rr.URL = r.URL.String()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rr.Body = string(body)

	jsonBody := make(map[string]interface{})
	err = json.Unmarshal(body, &jsonBody)
	if err == nil {
		rr.BodyJSON = jsonBody
	}

	rrb, err := json.Marshal(rr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(rrb)
	rw.Write([]byte("\n"))
}

func main() {
	http.HandleFunc("/", handle)
	http.ListenAndServe(":8000", nil)
}
