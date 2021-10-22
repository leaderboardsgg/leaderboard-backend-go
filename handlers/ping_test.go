package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/speedrun-website/leaderboard-backend/handlers"
)

func TestPingHandlerReturnsOK(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handlers.PingHandler(c)

	if w.Result().StatusCode != http.StatusOK {
		t.FailNow()
	}
}
