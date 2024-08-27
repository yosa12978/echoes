package middleware

import (
	"net/http"
	"time"

	"github.com/yosa12978/echoes/logging"
)

func Logger(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			snapshot := time.Now()
			next.ServeHTTP(w, r)
			latencyUs := time.Since(snapshot).Microseconds()
			logger.Info(
				"incoming request",
				"method", r.Method,
				"endpoint", r.URL.String(),
				"latency_us", latencyUs,
			)
		})
	}
}
