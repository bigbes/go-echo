package server

import (
	"net/http"

	"github.com/insolar/vanilla/atomickit"
)

// kubernetes endpoints
const LivenessEndpoint = "/alive"
const ReadinessEndpoint = "/ready"

type ProbeService struct {
	ready atomickit.Bool
}

func NewProbeService() *ProbeService {
	return &ProbeService{}
}

func (s *ProbeService) ApplyHandlers(mux *http.ServeMux) error {
	mux.Handle(LivenessEndpoint, http.HandlerFunc(s.livenessProbe))
	mux.Handle(ReadinessEndpoint, http.HandlerFunc(s.readinessProbe))

	return nil
}

func (s ProbeService) livenessProbe(w http.ResponseWriter, _ *http.Request) {
	// server is live after start
	w.WriteHeader(http.StatusOK)
}

func (s *ProbeService) readinessProbe(w http.ResponseWriter, _ *http.Request) {
	if s.ready.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (s *ProbeService) SetReadyFlag(status bool) {
	s.ready.Store(status)
}
