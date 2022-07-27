package server

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/soverenio/log"

	"github.com/soverenio/instrumentation/logger"
)

// errorLogger wrapper for error logs. Implements promhttp.Logger interface.
type errorLogger log.Logger

// Println is wrapper method for Error method.
func (e errorLogger) Println(v ...interface{}) {
	log.Logger(e).Error(v)
}

const MetricsEndpoint = "/metrics"

// metrics is a component which serve metrics data to Prometheus.
type MetricsService struct {
	ctx      context.Context
	gatherer prometheus.Gatherer
}

func NewMetricsService(ctx context.Context, gatherer prometheus.Gatherer) *MetricsService {
	return &MetricsService{
		ctx:      ctx,
		gatherer: gatherer,
	}

}

func (m MetricsService) ApplyHandlers(mux *http.ServeMux) error {
	opts := promhttp.HandlerOpts{
		ErrorLog: errorLogger(logger.FromContext(m.ctx)),
	}
	handler := promhttp.HandlerFor(m.gatherer, opts)

	mux.Handle(MetricsEndpoint, handler)

	return nil
}
