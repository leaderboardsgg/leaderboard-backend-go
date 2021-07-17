package middleware

import (
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
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
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
