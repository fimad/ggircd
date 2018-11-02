package irc

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	active_connections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "The number of active connections.",
		})
	command_received = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "command_received_total",
			Help: "Number of irc commands received.",
		},
		[]string{"nick", "command"},
	)
	nicks_in_channel = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "nicks_in_channel",
			Help: "The number of nicks in a channel.",
		},
		[]string{"channel"},
	)
)

func init() {
	prometheus.MustRegister(active_connections)
	prometheus.MustRegister(command_received)
	prometheus.MustRegister(nicks_in_channel)
}

// runPrometheus starts an HTTP server with a /metrics endpoint for publishing
// Prometheus metrics. This method does not return and should be run in a
// goroutine.
func runPrometheus(cfg Config) {
	http.Handle("/metrics", promhttp.Handler())
	var port string = fmt.Sprintf(":%d", cfg.Prometheus.Port)
	logf(fatal, "Http server error: %d", http.ListenAndServe(port, nil))
}
