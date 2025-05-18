package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	SearchCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "search_requests_total",
			Help: "Total number of search requests received",
		},
	)
)

func Init() {
	prometheus.MustRegister(SearchCounter)
}
