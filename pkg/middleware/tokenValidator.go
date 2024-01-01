package middleware

import (
	"net/http"
	"wasm/pkg/token"

	log "github.com/sirupsen/logrus"
)

func tokenValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for the custom header
		if r.Header.Get("X-Requested-By") != token.REQUESTED_BY {
			// If the header is not what we expect, return an error response
			http.Error(w, "Invalid request source", http.StatusForbidden)
			log.Error("Invalid request source")
			return
		}

		// Extract the token from the cookie
		cookie, err := r.Cookie("AuthToken")
		if err != nil {
			// Handle missing or invalid cookie
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			log.Error("Authentication required")
			return
		}

		receivedToken := cookie.Value
		if !token.ValidateToken(receivedToken) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			log.Error("Invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}
