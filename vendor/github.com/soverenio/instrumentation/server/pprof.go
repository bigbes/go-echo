package server

import (
	"net/http"
	"net/http/pprof"
)

// init disables default handlers registered by importing net/http/pprof.
func init() {
	http.DefaultServeMux = http.NewServeMux()
}

type PProfEndpointService struct{}

func NewPProfEndpoints() *PProfEndpointService {
	return &PProfEndpointService{}
}

// ApplyHandlers adds standard pprof handlers to mux.
func (p *PProfEndpointService) ApplyHandlers(mux *http.ServeMux) error {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	return nil
}
