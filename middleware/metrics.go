package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var metricsInfoKey = struct{}{}

var (
	RequestDurationOpts = prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Time (in seconds) spent serving this request.",
		Buckets: prometheus.DefBuckets,
	}
)

// MetricsInfo holds all of the metrics collected for a request.
type MetricsInfo struct {
	RequestDuration prometheus.Histogram
	Start           time.Time
}

// GetMetrics retuns the metrics for this current context.
func GetRequestMetrics(ctx context.Context) MetricsInfo {
	metrics, ok := ctx.Value(metricsInfoKey).(*MetricsInfo)
	if !ok || metrics == nil {
		return MetricsInfo{
			RequestDuration: nil,
			Start:           time.Time{},
		}
	}
	duration := time.Since(metrics.Start)
	metrics.RequestDuration.Observe(duration.Seconds())
	return *metrics
}

// NewMetricsMiddleware creates a middleware which produces metrics about a request, and tags the context with them.
// Metrics info can be retrieved with `GetMetrics(ctx)`.
func NewRequestMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), metricsInfoKey, &MetricsInfo{
			RequestDuration: prometheus.NewHistogram(RequestDurationOpts),
			Start:           time.Now(),
		})
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
