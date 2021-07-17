package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/speedrun-website/leaderboard-backend/logger"
	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddlewareWritesExpectedLogs(t *testing.T) {
	builder := new(strings.Builder)

	logger.Logger = zerolog.New(builder).With().Timestamp().Logger()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req, err := http.NewRequest("method", "url", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	LoggingMiddleware(w, req, nextHandler)

	var log struct {
		Level     string `json:"level"`
		Message   string `json:"message"`
		Timestamp string `json:"time"`
	}
	json.Unmarshal([]byte(builder.String()), &log)
	assert.Equal(t, log.Level, "info")
	assert.True(t, strings.Contains(log.Message, "url"))
	assert.True(t, strings.Contains(log.Message, "method"))
	fmt.Println(builder.String())
	assert.NotEqual(t, log.Timestamp, "")
}
