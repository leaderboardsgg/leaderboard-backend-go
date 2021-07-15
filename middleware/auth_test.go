package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthMiddlewareInstalled(t *testing.T) {
	authedHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		user := GetAuthedUser(r.Context())
		assert.True(t, user.AuthRan)
		assert.False(t, user.IsAuthed)
	})

	req, err := http.NewRequest("method", "url", nil)
	assert.NoError(t, err)

	handler := NewChainMiddlewareHandler([]ChainableMiddleware{NewAuthMiddleware}, authedHandler)
	handler.ServeHTTP(nil, req)
}

func TestAuthMiddlewareNotInstalled(t *testing.T) {
	authedHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		user := GetAuthedUser(r.Context())
		assert.False(t, user.AuthRan)
		assert.False(t, user.IsAuthed)
	})

	req, err := http.NewRequest("method", "url", nil)
	assert.NoError(t, err)

	handler := NewChainMiddlewareHandler(nil, authedHandler)
	handler.ServeHTTP(nil, req)
}
