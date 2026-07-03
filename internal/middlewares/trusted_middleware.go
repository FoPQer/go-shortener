package middlewares

import (
	"net"
	"net/http"
)

type TrustedMiddleware struct {
	ipNet *net.IPNet
}

func NewTrustedMiddleware(trustedSubnet string) *TrustedMiddleware {
	if trustedSubnet == "" {
		return &TrustedMiddleware{}
	}
	_, ipNet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return &TrustedMiddleware{}
	}
	return &TrustedMiddleware{ipNet: ipNet}
}

func (m *TrustedMiddleware) WithTrusted(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.ipNet == nil {
			http.Error(w, "Trusted subnet is not configured", http.StatusForbidden)
			return
		}
		xRealIP := r.Header.Get("X-Real-IP")
		if xRealIP == "" {
			http.Error(w, "X-Real-IP header is missing", http.StatusForbidden)
			return
		}
		ip := net.ParseIP(xRealIP)
		if ip == nil || !m.ipNet.Contains(ip) {
			http.Error(w, "Forbidden: IP not in trusted subnet", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
