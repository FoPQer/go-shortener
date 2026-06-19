package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrustedMiddleware_WithTrusted_AllowsIPInSubnet(t *testing.T) {
	middleware := NewTrustedMiddleware("192.168.1.0/24")
	handler := middleware.WithTrusted(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.1.42")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestTrustedMiddleware_WithTrusted_RejectsIPOutsideSubnet(t *testing.T) {
	middleware := NewTrustedMiddleware("192.168.1.0/24")
	handler := middleware.WithTrusted(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "192.168.2.42")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestTrustedMiddleware_WithTrusted_RejectsInvalidIP(t *testing.T) {
	middleware := NewTrustedMiddleware("192.168.1.0/24")
	handler := middleware.WithTrusted(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Real-IP", "not-an-ip")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}
