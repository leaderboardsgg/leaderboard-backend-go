package middleware

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	if err := prometheus.Register(totalRequests); err != nil {
		// TODO: Log properly here once we have some standard way to log.
		fmt.Printf("error registering totalRequests counter: %s", err.Error())
	}
	if err := prometheus.Register(responseStatus); err != nil {
		// TODO: Log properly here once we have some standard way to log.
		fmt.Printf("error registering responseStatus counter: %s", err.Error())
	}
	if err := prometheus.Register(httpDuration); err != nil {
		// TODO: Log properly here once we have some standard way to log.
		fmt.Printf("error registering httpDuration histogram: %s", err.Error())
	}
}

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)
	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP response",
		},
		[]string{"status"},
	)
	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
)
