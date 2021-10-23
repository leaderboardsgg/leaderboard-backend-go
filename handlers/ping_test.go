package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/handlers"
)

func TestPingHandlerReturnsOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.Ping(c)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf(
			"Expected Ping to return status code %d, got %d",
			http.StatusOK,
			w.Result().StatusCode,
		)
	}

	var response struct {
		Message string `json:"message"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf(
			"Ping response should have expectedPingResponse format, unmarshal failed with %s",
			err,
		)
	}
	if response.Message != "pong" {
		t.Fatalf(
			"Ping should pong, instead responded with %s",
			response.Message,
		)
	}
}
