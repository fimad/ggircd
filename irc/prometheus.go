package irc

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// runPrometheus starts an HTTP server with a /metrics endpoint for publishing
// Prometheus metrics. This method does not return and should be run in a
// goroutine.
func runPrometheus(cfg Config) {
	http.Handle("/metrics", promhttp.Handler())
	var port string = fmt.Sprintf(":%d", cfg.Prometheus.Port)
	logf(fatal, "Http server error: %d", http.ListenAndServe(port, nil))
}
