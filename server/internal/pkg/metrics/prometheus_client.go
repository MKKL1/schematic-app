package metrics

//TODO move to prometheus package

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

// PrometheusMetrics implements the MetricsClient interface using Prometheus.
type PrometheusMetrics struct {
	commandDuration *prometheus.HistogramVec
	commandCounter  *prometheus.CounterVec
}

// NewPrometheusMetrics creates and registers Prometheus metrics
// and returns an initialized PrometheusMetrics client.
func NewPrometheusMetrics() *PrometheusMetrics {
	durationVec := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "command_handler_duration_seconds",
			Help:    "Histogram of command amd query handler execution durations.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"command", "status", "type"},
	)

	counterVec := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "command_handler_requests_total",
			Help: "Total number of command and query handler requests.",
		},
		[]string{"command", "status", "type"},
	)

	return &PrometheusMetrics{
		commandDuration: durationVec,
		commandCounter:  counterVec,
	}
}

// RecordCommandDuration implements the MetricsClient interface for Prometheus.
func (p *PrometheusMetrics) RecordCommandDuration(commandName string, success bool, duration time.Duration) {
	status := getStatusLabel(success)
	p.commandDuration.WithLabelValues(commandName, status, "cmd").Observe(duration.Seconds())
	p.commandCounter.WithLabelValues(commandName, status, "cmd").Inc()
}

func (p *PrometheusMetrics) RecordQueryDuration(queryName string, success bool, duration time.Duration) {
	status := getStatusLabel(success)
	p.commandDuration.WithLabelValues(queryName, status, "q").Observe(duration.Seconds())
	p.commandCounter.WithLabelValues(queryName, status, "q").Inc()
}

// Ensure PrometheusMetrics implements MetricsClient
var _ Client = (*PrometheusMetrics)(nil)
