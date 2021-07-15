package middleware

import (
	"context"
	"net/http"
)

var authedUserInfoKey = struct{}{}

// AuthedUserInfo holds all of the metadata about an authentication context.
type AuthedUserInfo struct {
	AuthRan  bool
	IsAuthed bool
	UserName string // UserName is an example, we may have very different data in here long-term.
}

// GetAuthedUser returns the authentication info for this current context.
func GetAuthedUser(ctx context.Context) AuthedUserInfo {
	user, ok := ctx.Value(authedUserInfoKey).(*AuthedUserInfo)
	if !ok || user == nil {
		return AuthedUserInfo{
			AuthRan:  false,
			IsAuthed: false,
		}
	}

	return *user
}

// NewAuthMiddleware creates a middleware which authenticates a request, and tags the context with info about the user.
// Authentication info can be retrieved with `GetAuthedUser(ctx)`.
func NewAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// We have not authed anyone, so record that.
		ctx := context.WithValue(r.Context(), authedUserInfoKey, &AuthedUserInfo{
			AuthRan:  true,
			IsAuthed: false,
		})
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
