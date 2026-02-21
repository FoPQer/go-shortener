package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	writer      io.Writer
	gz          *gzip.Writer
	status      int
	wroteHeader bool
	enableGzip  bool
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.writer.Write(b)
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.status = statusCode

	if w.enableGzip && !isRedirectStatus(statusCode) && isCompressibleContentType(w.Header().Get("Content-Type")) {
		gz, err := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestCompression)
		if err == nil {
			w.gz = gz
			w.writer = gz
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Del("Content-Length")
		}
	}

	w.ResponseWriter.WriteHeader(statusCode)
}

func WithGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !supportsGzip(r) {
			next.ServeHTTP(w, r)
			return
		}

		gzipWriter := &gzipResponseWriter{
			ResponseWriter: w,
			writer:         w,
			enableGzip:     true,
		}
		next.ServeHTTP(gzipWriter, r)
		if gzipWriter.gz != nil {
			gzipWriter.gz.Close()
		}
	})
}

func supportsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func isCompressibleContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html")
}

func isRedirectStatus(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}