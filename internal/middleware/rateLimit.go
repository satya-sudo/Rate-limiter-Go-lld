package middleware

import (
	"github.com/satya-sudo/Rate-limiter-Go-lld/internal/limiter"
	"net/http"
)

func RateLimiterMiddleware(lim *limiter.TokenBucketRateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract the user-id
		key := r.Header.Get("X-User-Id")
		if key == "" {
			// fall back
			key = r.RemoteAddr
		}
		if !lim.Allow(key) {
			w.WriteHeader(http.StatusTooManyRequests)
			_, err := w.Write([]byte("too many requests"))
			if err != nil {
				return
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}
