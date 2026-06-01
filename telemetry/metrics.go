package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	HttpRequestsTotal    *prometheus.CounterVec
	HttpRequestDuration  *prometheus.HistogramVec
	DatabaseConnections  prometheus.Gauge
	WorkerTasksProcessed *prometheus.CounterVec
}

var GlobalMetrics *Metrics

// InitMetrics sets up the Prometheus metrics definitions
func InitMetrics(namespace string) {
	GlobalMetrics = &Metrics{
		HttpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_requests_total",
				Help:      "Total number of incoming HTTP requests.",
			},
			[]string{"method", "path", "status"},
		),
		HttpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "Duration histogram of incoming HTTP requests.",
				Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
			},
			[]string{"method", "path"},
		),
		DatabaseConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "db_open_connections",
				Help:      "Current open database pool connections.",
			},
		),
		WorkerTasksProcessed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "worker_tasks_processed_total",
				Help:      "Total number of background worker tasks executed.",
			},
			[]string{"task_type", "status"},
		),
	}
}
