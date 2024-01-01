package middleware

import "net/http"

func ApplyMiddleware(next http.Handler) http.Handler {

	rateLimiter := RateLimiterMiddleware(5, 10)

	return SecurityHeadersMiddleware((rateLimiter(setCSP(next))))
}

func ApplyTokenValidator(next http.Handler) http.Handler {
	return tokenValidator(ApplyMiddleware(next))
}
