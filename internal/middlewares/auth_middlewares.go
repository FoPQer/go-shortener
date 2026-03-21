package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/auth"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/service"
)

type AuthMiddleware struct {
	userService *service.UserService
	claimsService *auth.ClaimsService
}

const secretKey string = "your_secret_key"

func (m *AuthMiddleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		
		cookie, err := r.Cookie("X-Auth-Token")
		if errors.Is(err, http.ErrNoCookie) {
			newCookie, err := m.buildNewCookie()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to build new cookie: %v", err), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, newCookie)
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Missing X-Auth-Token cookie: %v", err), http.StatusBadRequest)
			return
		}
		var errInvalidToken *auth.ErrInvalidToken
		var errMissingUserID *auth.ErrMissingUserID
		tokenString := cookie.Value

		userID, err := m.claimsService.GetUserIDFromJWTString(tokenString, []byte(secretKey))
		if errors.As(err, &errInvalidToken) {
			newCookie, err := m.buildNewCookie()
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to build new cookie: %v", err), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, newCookie)
		} else if errors.As(err, &errMissingUserID) {
			http.Error(w, fmt.Sprintf("Missing user ID in token: %v", err), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse token: %v", err), http.StatusBadRequest)
			return
		}
		
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) buildNewCookie() (*http.Cookie, error) {
	userID, err := m.userService.Create(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	claims := m.claimsService.CreateClaims(userID)
	tokenString, err := m.claimsService.BuildJWTString(claims, []byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to build JWT string: %w", err)
	}
	return &http.Cookie{
		Name:  "X-Auth-Token",
		Value: tokenString,
	}, nil
}
