package mux_adapter

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Middleware converts a set of negroni.Handler middlewares into a mux.MiddlewareFunc.
// This method is meant to make it easier to set up middleware in a router agnostic way.
func Middleware(middlewares ...negroni.Handler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		negroniHandler := negroni.New(middlewares...)
		negroniHandler.UseHandler(next)
		return negroniHandler
	}
}
