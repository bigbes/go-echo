package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/insolar/vanilla/throw"

	"github.com/soverenio/instrumentation/logger"
)

type ServiceEndpoints interface {
	ApplyHandlers(mux *http.ServeMux) error
}

type Server struct {
	mux       *http.ServeMux
	server    *http.Server
	listener  net.Listener
	config    Configuration
	closeChan <-chan struct{}

	endpoints []ServiceEndpoints
}

func NewServer(cfg Configuration) Server {
	return Server{
		config: cfg,
		mux:    http.NewServeMux(),
	}
}

func (s *Server) AddEndpoints(endpoints ...ServiceEndpoints) {
	s.endpoints = append(s.endpoints, endpoints...)
}

// Run starts instrumentation server.
func (s *Server) Run(ctx context.Context) error {
	for _, ep := range s.endpoints {
		if err := ep.ApplyHandlers(s.mux); err != nil {
			return throw.W(err, "failed to apply handlers to mux")
		}
	}

	log := logger.FromContext(ctx)
	s.server = &http.Server{
		Addr:    s.config.ListenAddress,
		Handler: s.mux,
	}

	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return throw.W(err, "Failed to listen at address")
	}
	s.listener = listener

	closeChan := make(chan struct{})
	s.closeChan = closeChan

	go func() {
		log.Info("Starting instrumentation server: ", s.listener.Addr().String())
		err := s.server.Serve(listener)
		if err != http.ErrServerClosed {
			log.Error("Failed to start instrumentation server ", err)
		}
		close(closeChan)
	}()

	return nil
}

const StopTimeout = 10 * time.Second

// Stop is implementation of insolar.Component interface.
func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	logger.FromContext(ctx).Info("Shutting down instrumentation server")
	ctx, cancel := context.WithTimeout(ctx, StopTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return throw.W(err, "Can't gracefully stop instrumentation server")
	}

	select {
	case <-s.closeChan:
		return nil
	case <-ctx.Done():
		return throw.W(err, "Can't gracefully stop instrumentation server")
	}
}
