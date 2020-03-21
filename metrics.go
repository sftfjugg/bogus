package bogus

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
)

// Variables declared for monitoring.
var (
	HitsCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "bogus",
		Name:      "hits_count_total",
		Help:      "Counter of hits bogus.",
	}, []string{"to"})
)
