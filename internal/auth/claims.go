package auth

import "github.com/golang-jwt/jwt/v5"

// Claims represents JWT payload used by the application.
//
// It embeds standard registered JWT claims and adds an application-specific user ID.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// NewClaims creates JWT claims with the provided user ID.
func NewClaims(userID string) *Claims {
	return &Claims{
		UserID: userID,
	}
}
