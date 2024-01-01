package middleware

import (
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// RateLimiterMiddleware creates a new rate limiter middleware.
func RateLimiterMiddleware(r rate.Limit, b int) func(http.Handler) http.Handler {
	var mu sync.Mutex
	limiters := make(map[string]*rate.Limiter)

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		limiter, exists := limiters[ip]
		if !exists {
			limiter = rate.NewLimiter(r, b)
			limiters[ip] = limiter
		}

		return limiter
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter := getLimiter(r.RemoteAddr)

			if !limiter.Allow() {
				log.Error("rate limited...")
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
