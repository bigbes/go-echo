package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type DefaultServer struct {
	Server

	probe      *ProbeService
	migrations *MigrationsInfoEndpoint
}

func (s *DefaultServer) SetReadyFlag(flag bool) {
	s.probe.SetReadyFlag(flag)
}

func (s *DefaultServer) SetVersion(version func() (int64, int64, error)) {
	s.migrations.SetVersion(version)
}

func (s *DefaultServer) AddMetrics(ctx context.Context, gatherer prometheus.Gatherer) {
	s.AddEndpoints(NewMetricsService(ctx, gatherer))
}

func NewDefaultServer(cfg Configuration) DefaultServer {
	srv := DefaultServer{
		Server:     NewServer(cfg),
		probe:      NewProbeService(),
		migrations: NewMigrationsInfoEndpoint(),
	}
	srv.AddEndpoints(srv.probe, NewPProfEndpoints(), NewBuildInfoEndpoint(), srv.migrations)
	return srv
}
