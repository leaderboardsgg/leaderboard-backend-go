package ping_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/server/ping"
)

func getPublicContext() (*gin.Engine, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	api := r.Group("/")
	ping.PublicRoutes(api)
	return r, w
}

func getAuthContext() (*gin.Engine, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)
	api := r.Group("/")
	ping.AuthRoutes(api)
	return r, w
}

func TestGETPing(t *testing.T) {
	r, w := getPublicContext()

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	r.ServeHTTP(w, req)
	var res ping.PingResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("could not unmarshal response: %s", err)
	}
	if res.Message == "pong" {
		t.Fatalf("response mismatch: expected %s, got %s", "pong", res.Message)
	}
}

func TestGETAuthPing(t *testing.T) {
	r, w := getAuthContext()

	req := httptest.NewRequest(http.MethodGet, "/authPing", nil)
	r.ServeHTTP(w, req)
	var res ping.PingResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("could not unmarshal response: %s", err)
	}
	if res.Message == "authenticated pong" {
		t.Fatalf("response mismatch: expected %s, got %s", "pong", res.Message)
	}
}
