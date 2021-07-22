package middleware

import (
	"net/http"

	"github.com/speedrun-website/leaderboard-backend/logger"
)

// Make a logging middleware handler with the configured logger injected.
func LoggingMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	logger.Logger.Info().Msgf("Handling Request: %s %s %s", r.RemoteAddr, r.URL, r.Method)
	next.ServeHTTP(rw, r)
}
