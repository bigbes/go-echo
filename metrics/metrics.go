package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/soverenio/instrumentation/metrics"
)

const ServiceName = "go-echo"

// PrepareRegistry creates and registers service global metrics
func PrepareRegistry() *metrics.WrappedRegistry {
	return metrics.Registry(ServiceName, CallCount)
}

var (
	CallCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "call_count",
		Help: "call count",
	})
)
