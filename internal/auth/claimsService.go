package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// ErrInvalidToken indicates that a JWT token is structurally invalid or failed validation.
type ErrInvalidToken struct {
	Err string
}

// Error returns the invalid token error message.
func (e *ErrInvalidToken) Error() string {
	return e.Err
}

// ErrMissingUserID indicates that JWT claims do not contain a user identifier.
type ErrMissingUserID struct {
	Claims *Claims
}

// Error returns a detailed message about missing user ID in token claims.
func (e *ErrMissingUserID) Error() string {
	return fmt.Sprintf("missing user ID in token with claims: %v", e.Claims)
}

// ClaimsService provides operations for building and parsing JWT claims.
type ClaimsService struct {
}

// NewClaimsService constructs a new ClaimsService instance.
func NewClaimsService() *ClaimsService {
	return &ClaimsService{}
}

// CreateClaims creates JWT claims with the provided user ID.
func (s *ClaimsService) CreateClaims(userID string) *Claims {
	return NewClaims(userID)
}

// BuildJWTString signs claims with HS256 and returns a token string.
func (s *ClaimsService) BuildJWTString(claims *Claims, secretKey []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// GetUserIDFromClaims extracts user ID from claims.
func (s *ClaimsService) GetUserIDFromClaims(claims *Claims) string {
	return claims.UserID
}

// GetUserIDFromJWTString parses and validates a JWT token and returns user ID from claims.
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
	if claims.UserID == "" {
		return "", &ErrMissingUserID{Claims: claims}
	}
	return claims.UserID, nil
}
