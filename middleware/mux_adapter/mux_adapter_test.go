package mux_adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestMuxAdapterMiddleware(t *testing.T) {
	didRun := false
	testMiddleware := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		didRun = true
		next(rw, r)
	}

	router := mux.NewRouter()
	router.Use(Middleware(negroni.HandlerFunc(testMiddleware)))
	router.Path("/test").Handler(promhttp.Handler())

	t.Log("Launching server")
	ts := httptest.NewServer(router)
	defer ts.Close()

	_, err := http.Get(ts.URL + "/test")
	assert.NoError(t, err)
	assert.True(t, didRun)
}
