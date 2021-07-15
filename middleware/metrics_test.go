package middleware

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
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
	t.Log("start")

	t.Log("reggister")

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

	t.Log("assert")

	// tR collected three metrics
	assert.Equal(3, testutil.CollectAndCount(tR))
	// responseStatus collected three metrics
	assert.Equal(3, testutil.CollectAndCount(rS))
	// httpDuration collected three metrics
	assert.Equal(3, testutil.CollectAndCount(hD))
}
