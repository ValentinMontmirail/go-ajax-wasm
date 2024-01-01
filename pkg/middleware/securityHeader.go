package middleware

import "net/http"

// SecurityHeadersMiddleware adds security-related headers to responses
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// X-Content-Type-Options header tells the browser not to sniff the content type
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options header tells the browser whether you want to allow your site to be framed or not
		w.Header().Set("X-Frame-Options", "DENY") // DENY or SAMEORIGIN

		// X-XSS-Protection header is a feature of Internet Explorer, Chrome and Safari that stops pages from loading when they detect reflected cross-site scripting (XSS) attacks
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Feature-Policy header allows a site to control which features and APIs can be used in the browser
		w.Header().Set("Feature-Policy", "camera 'none'; microphone 'none'; geolocation 'none'") // Example policy

		// Add other headers here if needed

		next.ServeHTTP(w, r)
	})
}
