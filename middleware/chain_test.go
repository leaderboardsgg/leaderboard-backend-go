package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainMiddlewareEmpty(t *testing.T) {
	var outputStrs []string

	finalHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		outputStrs = append(outputStrs, "final")
	})

	handler := NewChainMiddlewareHandler(nil, finalHandler)
	handler.ServeHTTP(nil, nil)
	assert.Equal(t, []string{"final"}, outputStrs)
}

func TestChainMiddlewareOrder(t *testing.T) {
	var outputStrs []string

	finalHandler := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		outputStrs = append(outputStrs, "final")
	})
	firstHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			outputStrs = append(outputStrs, "first - in")
			next.ServeHTTP(rw, r)
			outputStrs = append(outputStrs, "first - out")
		})
	}
	secondHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			outputStrs = append(outputStrs, "second - in")
			next.ServeHTTP(rw, r)
			outputStrs = append(outputStrs, "second - out")
		})
	}

	handler := NewChainMiddlewareHandler([]ChainableMiddleware{firstHandler, secondHandler}, finalHandler)
	handler.ServeHTTP(nil, nil)
	assert.Equal(t, []string{"first - in", "second - in", "final", "second - out", "first - out"}, outputStrs)
}
