package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ErrInvalidToken struct {
	Err string
}

func (e *ErrInvalidToken) Error() string {
	return e.Err
}

type ClaimsService struct {
}

func NewClaimsService() *ClaimsService {
	return &ClaimsService{}
}

func (s *ClaimsService) CreateClaims(userID string) *Claims {
	return NewClaims(userID)
}

func (s *ClaimsService) GetUserIDFromClaims(claims *Claims) string {
	return claims.UserID
}

func (s *ClaimsService) GetUserIDFromJWTString(tokenString string, secretKey []byte) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return "", &ErrInvalidToken{Err: "invalid token"}
	}
	return claims.UserID, nil
}