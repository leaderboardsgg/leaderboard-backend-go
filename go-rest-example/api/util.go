package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/speedrun-website/speedrun-rest/auth"
)

type requestFunction func(w http.ResponseWriter, r *http.Request)

func writeJsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func authenticatedHandler(handler requestFunction) requestFunction {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitAuthHeader := strings.Split(authHeader, " ")
		if len(splitAuthHeader) < 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		apiKey := splitAuthHeader[1]
		if !auth.VerifyAPIKey(apiKey) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}
