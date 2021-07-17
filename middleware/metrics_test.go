package middleware

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/speedrun-website/leaderboard-backend/middleware/mux_adapter"
	"github.com/urfave/negroni"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/stretchr/testify/assert"
)

func clearMetrics() {
	totalRequests.Reset()
	responseStatus.Reset()
	httpDuration.Reset()
}

func TestPrometheusRegisters(t *testing.T) {
	assert := assert.New(t)
	clearMetrics()

	totalRequests.WithLabelValues("firstLabel").Inc()
	totalRequests.WithLabelValues("secondLabel").Inc()
	totalRequests.WithLabelValues("thirdLabel").Inc()
	totalRequests.WithLabelValues("thirdLabel").Inc()

	responseStatus.WithLabelValues("fourthlabel").Inc()
	responseStatus.WithLabelValues("fifthlabel").Inc()
	responseStatus.WithLabelValues("sixthlabel").Inc()

	httpDuration.WithLabelValues("seventhlabel").Observe(1)
	httpDuration.WithLabelValues("eighthlabel").Observe(2)
	httpDuration.WithLabelValues("ninelabel").Observe(2)

	// tR collected three metrics
	assert.Equal(3, testutil.CollectAndCount(totalRequests))
	// responseStatus collected three metrics
	assert.Equal(3, testutil.CollectAndCount(responseStatus))
	// httpDuration collected three metrics
	assert.Equal(3, testutil.CollectAndCount(httpDuration))
}

func TestPrometheusMiddlewareAttached(t *testing.T) {
	assert := assert.New(t)
	clearMetrics()

	t.Log("Defining router and registering prometheus targets")
	router := mux.NewRouter()
	router.Use(mux_adapter.Middleware(negroni.HandlerFunc(PrometheusMiddleware)))
	router.Path("/metrics").Handler(promhttp.Handler())

	t.Log("Launching server")
	ts := httptest.NewServer(router)
	defer ts.Close()
	metricsUrl := ts.URL + "/metrics"

	t.Log("Access the server a few times to ensure all the metrics populate")
	for i := 0; i < 10; i++ {
		_, err := http.Get(metricsUrl)
		if err != nil {
			t.Log(err)
			t.FailNow()
		}
	}
	t.Log("Retrieving metrics page")
	res, err := http.Get(metricsUrl)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	t.Log("Reading response body")
	body, err := ioutil.ReadAll(res.Body)
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
		if strings.Contains(line, "http_requests_total ") {
			t.Logf("Line contains http_requests_total: %s", line)
			totalRequestsAttached = true
		} else if strings.Contains(line, "response_status") {
			t.Logf("Line contains response_status: %s", line)
			responseStatusAttached = true
		} else if strings.Contains(line, "http_response_time_seconds") {
			t.Logf("Line contains http_response_time_seconds: %s", line)
			httpDurationAttached = true
		}
	}

	assert.Equal(true, totalRequestsAttached, "http_requests_total should appear in the metrics")
	assert.Equal(true, responseStatusAttached, "response_status should appear in the metrics")
	assert.Equal(true, httpDurationAttached, "http_response_time_seconds should appear in the metrics")
}
