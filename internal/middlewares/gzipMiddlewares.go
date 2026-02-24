package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

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
		if  isCompressibleContentType(r.Header.Get("Content-Type")) {
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

func supportsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func contentGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Encoding"), "gzip")
}

func isCompressibleContentType(contentType string) bool {
	return strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html")
}