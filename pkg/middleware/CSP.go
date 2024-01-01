package middleware

import "net/http"

// Middleware to set Content Security Policy
func setCSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := "default-src 'self'; script-src 'self' 'wasm-eval' 'unsafe-eval'; style-src 'self' 'unsafe-inline';"
		w.Header().Set("Content-Security-Policy", csp)
		next.ServeHTTP(w, r)
	})
}
