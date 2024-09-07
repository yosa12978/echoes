package middleware

import (
	"net/http"
	"time"

	"github.com/yosa12978/echoes/logging"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logger(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			writer := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}
			snapshot := time.Now()
			next.ServeHTTP(writer, r)
			latencyUs := time.Since(snapshot).Microseconds()
			logger.Info(
				"incoming request",
				"method", r.Method,
				"endpoint", r.URL.String(),
				"status_code", writer.statusCode,
				"latency_us", latencyUs,
			)
		})
	}
}
