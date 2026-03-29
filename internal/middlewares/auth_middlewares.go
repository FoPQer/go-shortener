package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/FoPQer/go-shortener/internal/auth"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/service"
)

const secretKey string = "your_secret_key"

type AuthMiddleware struct {
	userService *service.UserService
	claimsService *auth.ClaimsService
}

func NewAuthMiddleware(userService *service.UserService, claimsService *auth.ClaimsService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService, claimsService: claimsService}
}

func (m *AuthMiddleware) WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		
		cookie, err := r.Cookie("X-Auth-Token")
		if errors.Is(err, http.ErrNoCookie) {
			newCookie, err := m.buildNewCookie(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to build new cookie: %v", err), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, newCookie)
			cookie = newCookie
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Something went wrong while getting the X-Auth-Token cookie: %v", err), http.StatusBadRequest)
			return
		}
		var errInvalidToken *auth.ErrInvalidToken
		var errMissingUserID *auth.ErrMissingUserID
		tokenString := cookie.Value
		log.Printf("Received token: %s", tokenString)
		userID, err := m.claimsService.GetUserIDFromJWTString(tokenString, []byte(secretKey))
		if errors.As(err, &errInvalidToken) {
			newCookie, err := m.buildNewCookie(r)
			if err != nil {
				http.Error(w, fmt.Sprintf("failed to build new cookie: %v", err), http.StatusInternalServerError)
				return
			}
			http.SetCookie(w, newCookie)
			cookie = newCookie
		} else if errors.As(err, &errMissingUserID) {
			http.Error(w, fmt.Sprintf("Missing user ID in token: %v", err), http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse token: %v", err), http.StatusBadRequest)
			return
		}
		log.Printf("UserID: %s", userID)
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) buildNewCookie(r *http.Request) (*http.Cookie, error) {
	userID, err := m.userService.Create(r.Context(), &model.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	claims := m.claimsService.CreateClaims(userID)
	tokenString, err := m.claimsService.BuildJWTString(claims, []byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to build JWT string: %w", err)
	}
	log.Printf("Generated new token for user %s: %s", userID, tokenString)
	return &http.Cookie{
		Name:  "X-Auth-Token",
		Value: tokenString,
		HttpOnly: true,
	}, nil
}
