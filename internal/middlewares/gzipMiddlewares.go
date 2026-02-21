package middlewares

import (
	"compress/gzip"
	"io"
	"log"
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
		if !supportsGzip(r) || (!isJSONContent(r) && !isHTMLContent(r)) {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		
		gzipResponseWriter := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gz,
		}
		next.ServeHTTP(gzipResponseWriter, r)
	})
}

func supportsGzip(r *http.Request) bool {
	log.Printf("%s", r.Header.Get("Accept-Encoding"))
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func isJSONContent(r *http.Request) bool {
	log.Printf("%s", r.Header.Get("Content-Type"))
	return r.Header.Get("Content-Type") == "application/json"
}

func isHTMLContent(r *http.Request) bool {
	log.Printf("%s", r.Header.Get("Content-Type"))
	return r.Header.Get("Content-Type") == "text/html"
}