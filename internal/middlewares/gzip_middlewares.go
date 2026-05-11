package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// gzipResponseWriter writes HTTP response data through the configured writer.
type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write writes response bytes to the wrapped writer.
func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// WithGzip enables gzip request decompression and response compression when supported.
func WithGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !supportsGzip(r) {
			next.ServeHTTP(w, r)
			return
		}

		var writer io.Writer = w

		if r.Method == http.MethodPost && contentGzip(r) {
			gzr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer gzr.Close()

			r.Body = gzr
		}
		if isCompressibleContentType(r.Header.Get("Content-Type")) {
			gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
			if err != nil {
				io.WriteString(w, err.Error())
				return
			}
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")

			writer = gz
		}

		gzipResponseWriter := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         writer,
		}
		next.ServeHTTP(gzipResponseWriter, r)
	})
}

// supportsGzip reports whether client accepts gzip-encoded responses.
func supportsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

// contentGzip reports whether request body is gzip-encoded.
func contentGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}

// isCompressibleContentType reports whether content type should be gzip-compressed.
func isCompressibleContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html")
}
