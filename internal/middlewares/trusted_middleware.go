package middlewares

import (
	"net"
	"net/http"
)

type TrustedMiddleware struct {
	trustedSubnet string
}

func NewTrustedMiddleware(trustedSubnet string) *TrustedMiddleware {
	return &TrustedMiddleware{trustedSubnet: trustedSubnet}
}

func (m *TrustedMiddleware) WithTrusted(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.trustedSubnet == "" {
			http.Error(w, "Trusted subnet is not configured", http.StatusForbidden)
			return
		}
		xRealIP := r.Header.Get("X-Real-IP")
		if xRealIP == "" {
			http.Error(w, "X-Real-IP header is missing", http.StatusForbidden)
			return
		}

		if !m.isIPInSubnet(xRealIP, m.trustedSubnet) {
			http.Error(w, "Forbidden: IP not in trusted subnet", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *TrustedMiddleware) isIPInSubnet(ip, subnet string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return false
	}

	return ipNet.Contains(parsedIP)
}
