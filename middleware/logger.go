package middleware

import (
	"fmt"
	"net/http"

	"github.com/speedrun-website/leaderboard-backend/logger"
)

// Make a logging middleware handler with the configured logger injected.
func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	requestLogMessage := fmt.Sprintf("Handling Request: %s %s %s", r.RemoteAddr, r.URL, r.Method)
	logger.Logger.Info().Msg(requestLogMessage)
	next.ServeHTTP(rw, r)
}
