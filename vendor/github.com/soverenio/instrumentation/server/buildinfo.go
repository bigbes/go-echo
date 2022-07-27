package server

import (
	"encoding/json"
	"net/http"

	"github.com/soverenio/instrumentation/buildinfo"
)

const InfoEndpoint = "/info"

type BuildInfoEndpoint struct{}

func NewBuildInfoEndpoint() *BuildInfoEndpoint {
	return &BuildInfoEndpoint{}
}

func (e BuildInfoEndpoint) ApplyHandlers(mux *http.ServeMux) error {
	mux.Handle(InfoEndpoint, http.HandlerFunc(e.infoEndpoint))

	return nil
}

func (e BuildInfoEndpoint) infoEndpoint(w http.ResponseWriter, _ *http.Request) {
	bi := buildinfo.NewInfo()
	rsp, err := json.MarshalIndent(bi, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp)
	w.Write([]byte("\n"))
}
