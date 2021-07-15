package middleware

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequestMiddlewareInstalled(t *testing.T) {
	metricsHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		metrics := GetRequestMetrics(r.Context())
		// Slightly worried about this test becoming flakey if the test runs slow.
		assert.Equal(t, metrics.Start.Format("2006.01.02 15:04:05"), time.Now().Format("2006.01.02 15:04:05"))
	})

	req, err := http.NewRequest("method", "url", nil)
	assert.NoError(t, err)

	handler := NewChainMiddlewareHandler([]ChainableMiddleware{NewRequestMetricsMiddleware}, metricsHandler)
	handler.ServeHTTP(nil, req)
}

func TestRequestMiddlewareNotInstalled(t *testing.T) {
	metricsHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		metrics := GetRequestMetrics(r.Context())
		assert.Nil(t, metrics.RequestDuration)
		assert.Equal(t, metrics.Start, time.Time{})
	})

	req, err := http.NewRequest("method", "url", nil)
	assert.NoError(t, err)

	handler := NewChainMiddlewareHandler(nil, metricsHandler)
	handler.ServeHTTP(nil, req)
}
