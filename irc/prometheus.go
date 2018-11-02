package irc

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	command_received = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "command_received_total",
			Help: "Number of irc commands received.",
		},
		[]string{"nick", "command"},
	)
	active_connections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "The number of active connections.",
		})
)

func init() {
	prometheus.MustRegister(command_received)
	prometheus.MustRegister(active_connections)
}

// runPrometheus starts an HTTP server with a /metrics endpoint for publishing
// Prometheus metrics. This method does not return and should be run in a
// goroutine.
func runPrometheus(cfg Config) {
	http.Handle("/metrics", promhttp.Handler())
	var port string = fmt.Sprintf(":%d", cfg.Prometheus.Port)
	logf(fatal, "Http server error: %d", http.ListenAndServe(port, nil))
}
