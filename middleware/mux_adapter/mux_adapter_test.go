package mux_adapter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
)

func TestMuxAdapterMiddleware(t *testing.T) {
	middlewareRunCount := 0
	testMiddleware := func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		middlewareRunCount++
		next(rw, r)
	}

	handlerRunCount := 0
	testHandler := func(rw http.ResponseWriter, r *http.Request) {
		handlerRunCount++
		rw.WriteHeader(http.StatusOK)
	}

	router := mux.NewRouter()
	router.Use(Middleware(negroni.HandlerFunc(testMiddleware)))
	router.Path("/test").HandlerFunc(testHandler)

	t.Log("Launching server")
	ts := httptest.NewServer(router)
	defer ts.Close()
	testUrl := ts.URL + "/test"

	_, err := http.Get(testUrl)
	assert.NoError(t, err)
	assert.Equal(t, 1, middlewareRunCount)
	assert.Equal(t, 1, handlerRunCount)

	_, err = http.Get(testUrl)
	assert.NoError(t, err)
	assert.Equal(t, 2, middlewareRunCount)
	assert.Equal(t, 2, handlerRunCount)
}
