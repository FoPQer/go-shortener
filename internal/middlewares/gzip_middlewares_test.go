package middlewares_test

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FoPQer/go-shortener/internal/middlewares"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithGzip(t *testing.T) {
	tests := []struct {
		name                string
		acceptEncoding      string
		contentType         string
		contentEncoding     string
		requestBody         string
		method              string
		expectedEncoding    string
		expectCompression   bool
		expectDecompression bool
	}{
		{
			name:              "no gzip support",
			acceptEncoding:    "",
			contentType:       "application/json",
			method:            http.MethodGet,
			requestBody:       `{"test":"data"}`,
			expectedEncoding:  "",
			expectCompression: false,
		},
		{
			name:              "gzip support with json content",
			acceptEncoding:    "gzip",
			contentType:       "application/json",
			method:            http.MethodGet,
			requestBody:       `{"test":"data"}`,
			expectedEncoding:  "gzip",
			expectCompression: true,
		},
		{
			name:              "gzip support with html content",
			acceptEncoding:    "gzip",
			contentType:       "text/html",
			method:            http.MethodGet,
			requestBody:       `<html><body>test</body></html>`,
			expectedEncoding:  "gzip",
			expectCompression: true,
		},
		{
			name:                "post with gzip compressed body",
			acceptEncoding:      "gzip",
			contentType:         "application/json",
			contentEncoding:     "gzip",
			method:              http.MethodPost,
			requestBody:         `{"test":"data"}`,
			expectedEncoding:    "gzip",
			expectCompression:   true,
			expectDecompression: true,
		},
		{
			name:              "multiple accept encodings including gzip",
			acceptEncoding:    "deflate, gzip, br",
			contentType:       "application/json",
			method:            http.MethodGet,
			requestBody:       `{"test":"data"}`,
			expectedEncoding:  "gzip",
			expectCompression: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectDecompression {
					body, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					assert.Equal(t, tt.requestBody, string(body))
				}

				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.requestBody))
			})

			wrappedHandler := middlewares.WithGzip(handler)

			var reqBody io.Reader
			if tt.contentEncoding == "gzip" {
				var buf bytes.Buffer
				gzw := gzip.NewWriter(&buf)
				_, err := gzw.Write([]byte(tt.requestBody))
				require.NoError(t, err)
				gzw.Close()
				reqBody = &buf
			} else {
				reqBody = bytes.NewBufferString(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, "/test", reqBody)
			if tt.acceptEncoding != "" {
				req.Header.Set("Accept-Encoding", tt.acceptEncoding)
			}
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			if tt.contentEncoding != "" {
				req.Header.Set("Content-Encoding", tt.contentEncoding)
			}

			rec := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			assert.Equal(t, tt.expectedEncoding, rec.Header().Get("Content-Encoding"))

			if tt.expectCompression {
				gzr, err := gzip.NewReader(rec.Body)
				require.NoError(t, err)
				defer gzr.Close()

				decompressed, err := io.ReadAll(gzr)
				require.NoError(t, err)
				assert.Equal(t, tt.requestBody, string(decompressed))
			} else {
				assert.Equal(t, tt.requestBody, rec.Body.String())
			}
		})
	}
}

func TestWithGzip_ErrorCases(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    []byte
		contentGzip    bool
		expectedStatus int
	}{
		{
			name:           "invalid gzip body",
			requestBody:    []byte("not a gzip data"),
			contentGzip:    true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrappedHandler := middlewares.WithGzip(handler)

			req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(tt.requestBody))
			req.Header.Set("Accept-Encoding", "gzip")
			if tt.contentGzip {
				req.Header.Set("Content-Encoding", "gzip")
			}

			rec := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
