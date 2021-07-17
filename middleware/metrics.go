package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/prometheus/client_golang/prometheus"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func RegisterPrometheus() {
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

// PrometheusMiddleware is a middleware which produces metrics about requests.
func PrometheusMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	route := mux.CurrentRoute(r)
	path, err := route.GetPathTemplate()
	if err != nil {
		// TODO: Log here once we have some standard way to log.
		path = "UNKNOWN"
	}

	timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
	defer timer.ObserveDuration()

	nrw := negroni.NewResponseWriter(rw)
	next(nrw, r)
	statusCode := nrw.Status()

	responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
	totalRequests.WithLabelValues(path).Inc()
}
