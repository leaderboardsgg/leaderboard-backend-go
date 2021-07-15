package middleware

import "net/http"

// ChainableMiddleware is a wrapper type to make it easy to nest many middleware.
type ChainableMiddleware func(next http.Handler) http.Handler

// NewChainMiddlewareHandler creates a new http.Handler which calls each middleware in order before calling finalHandler.
// This operates from left-to-right on the middlewares.
func NewChainMiddlewareHandler(middlewares []ChainableMiddleware, finalHandler http.Handler) http.Handler {
	handler := finalHandler
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
