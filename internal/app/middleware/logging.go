package middleware

import (
	"log"
	"net/http"
	"time"
)

// Create a response writer wrapper to capture the status code
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w}

		// Call the next handler
		next.ServeHTTP(recorder, r)

		// Log the request details with status code
		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			recorder.status,
			time.Since(start),
		)
	})
}
