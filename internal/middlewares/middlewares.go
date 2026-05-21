package middlewares

import (
	"net/http"
	"time"

	"github.com/FoPQer/go-shortener/internal/logger"
)

type (
	// responseData stores HTTP response metadata collected by logging middleware.
	responseData struct {
		status int
		size   int
	}

	// loggingResponseWriter captures response status and size for request logging.
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write writes response body and accumulates written byte count.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader writes HTTP status code and stores it for logging.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// WithLogging logs incoming request metadata and outgoing response metrics.
func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		logger.GetSugar().Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"header", r.Header,
		)

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.GetSugar().Infoln(
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logFn)
}
