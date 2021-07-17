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

	AuthMiddleware(nil, req, authedHandler)
}

func TestAuthMiddlewareNotInstalled(t *testing.T) {
	authedHandler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		user := GetAuthedUser(r.Context())
		assert.False(t, user.AuthRan)
		assert.False(t, user.IsAuthed)
	})

	req, err := http.NewRequest("method", "url", nil)
	assert.NoError(t, err)

	authedHandler.ServeHTTP(nil, req)
}
