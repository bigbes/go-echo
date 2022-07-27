package metrics

import (
	"github.com/insolar/vanilla/throw"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	prometheusclient "github.com/prometheus/client_model/go"
	"github.com/prometheus/procfs"
	"github.com/soverenio/log/global"
)

var _ prometheus.Gatherer = (*WrappedRegistry)(nil)

type WrappedRegistry struct {
	prometheus.Registerer

	gatherer prometheus.Gatherer
}

func (w *WrappedRegistry) Gather() ([]*prometheusclient.MetricFamily, error) {
	return w.gatherer.Gather()
}

// Registry creates and registers service global metrics
// Custom collectors could be added via extraCollectors
// Registry("myService, APIMetrics, NetworkMetrics, SomeMoreMetrics)
func Registry(serviceName string, extraCollectors ...prometheus.Collector) *WrappedRegistry {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWith(prometheus.Labels{"service": serviceName}, registry)

	// Default system collectors
	if _, err := procfs.NewDefaultFS(); err == nil {
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})
	} else {
		global.Warn(throw.W(err, "failed to initialize process metrics collector"))
	}
	registerer.MustRegister(collectors.NewGoCollector())

	for _, collector := range extraCollectors {
		registerer.MustRegister(collector)
	}

	return &WrappedRegistry{Registerer: registerer, gatherer: registry}
}
