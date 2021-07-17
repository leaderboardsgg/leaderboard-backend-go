package middleware

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/prometheus/client_golang/prometheus"
)

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
