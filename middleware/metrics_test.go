package middleware

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/stretchr/testify/assert"
)

var tR = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var rS = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var hD = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func TestPrometheusRegisters(t *testing.T) {
	assert := assert.New(t)

	t.Log("Registering Prometheus Metrics")

	prometheus.Register(tR)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)

	tR.WithLabelValues("firstLabel").Inc()
	tR.WithLabelValues("secondLabel").Inc()
	tR.WithLabelValues("thirdLabel").Inc()
	tR.WithLabelValues("thirdLabel").Inc()

	rS.WithLabelValues("fourthlabel").Inc()
	rS.WithLabelValues("fifthlabel").Inc()
	rS.WithLabelValues("sixthlabel").Inc()

	hD.WithLabelValues("seventhlabel").Observe(1)
	hD.WithLabelValues("eighthlabel").Observe(2)
	hD.WithLabelValues("ninelabel").Observe(2)

	// tR collected three metrics
	assert.Equal(3, testutil.CollectAndCount(tR))
	// responseStatus collected three metrics
	assert.Equal(3, testutil.CollectAndCount(rS))
	// httpDuration collected three metrics
	assert.Equal(3, testutil.CollectAndCount(hD))
}

func TestPrometheusMiddlewareAttached(t *testing.T) {
	assert := assert.New(t)

	t.Log("Launching server in go routine")
	go func() {
		t.Log("Defining router and registering prometheus targets")
		router := mux.NewRouter()
		router.Use(NewPrometheusMiddleware)
		RegisterPrometheus()

		router.Path("/metrics").Handler(promhttp.Handler())

		t.Log("Launching server")
		if err := http.ListenAndServe(":3030", router); err != nil {
			t.Log("Panic! ListenAndServe: " + err.Error())
			panic("ListenAndServe: " + err.Error())
		}
	}()

	t.Log("Waiting for server to launch before proceeding")
	time.Sleep(5 * time.Second)

	t.Log("Access the server a few times to ensure all the metrics populate")
	for i := 0; i < 10; i++ {
		_, err := http.Get("http://localhost:3030/metrics")
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
	t.Log("Retrieving metrics page")
	res, err := http.Get("http://localhost:3030/metrics")

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log("Reading response body")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.FailNow()
	}

	metrics := string(body)
	metricsArray := strings.Split(metrics, "\n")

	totalRequestsAttached := false
	responseStatusAttached := false
	httpDurationAttached := false

	t.Log("Scanning for attached metrics")
	for _, line := range metricsArray {
		t.Log(line)
		if strings.Contains(line, "http_requests_total ") {
			t.Log("Line contains http_requests_total")
			totalRequestsAttached = true
		} else if strings.Contains(line, "response_status") {
			t.Log("Line contains response_status")
			responseStatusAttached = true
		} else if strings.Contains(line, "http_response_time_seconds") {
			t.Log("Line contains http_response_time_seconds")
			httpDurationAttached = true
		}
	}

	assert.Equal(true, totalRequestsAttached, "http_requests_total should appear in the metrics")
	assert.Equal(true, responseStatusAttached, "response_status should appear in the metrics")
	assert.Equal(true, httpDurationAttached, "http_response_time_seconds should appear in the metrics")
}
