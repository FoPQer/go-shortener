package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testSecret = []byte("test-secret-key")

func newService() *ClaimsService {
	return NewClaimsService()
}

func TestCreateClaims(t *testing.T) {
	s := newService()
	claims := s.CreateClaims("user-123")
	assert.Equal(t, "user-123", claims.UserID)
}

func TestGetUserIDFromClaims(t *testing.T) {
	s := newService()
	claims := s.CreateClaims("user-456")
	assert.Equal(t, "user-456", s.GetUserIDFromClaims(claims))
}

func TestBuildJWTString(t *testing.T) {
	s := newService()

	t.Run("valid token produced", func(t *testing.T) {
		claims := s.CreateClaims("user-789")
		tokenStr, err := s.BuildJWTString(claims, testSecret)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenStr)
	})
}

func TestGetUserIDFromJWTString(t *testing.T) {
	s := newService()

	tests := []struct {
		name        string
		setup       func() string
		secretKey   []byte
		wantUserID  string
		wantErrType interface{}
	}{
		{
			name: "valid token",
			setup: func() string {
				claims := s.CreateClaims("user-abc")
				token, err := s.BuildJWTString(claims, testSecret)
				require.NoError(t, err)
				return token
			},
			secretKey:  testSecret,
			wantUserID: "user-abc",
		},
		{
			name: "wrong secret",
			setup: func() string {
				claims := s.CreateClaims("user-abc")
				token, err := s.BuildJWTString(claims, testSecret)
				require.NoError(t, err)
				return token
			},
			secretKey:   []byte("wrong-secret"),
			wantErrType: new(error),
		},
		{
			name: "malformed token",
			setup: func() string {
				return "not.a.valid.jwt"
			},
			secretKey:   testSecret,
			wantErrType: new(error),
		},
		{
			name: "token without userID",
			setup: func() string {
				claims := NewClaims("")
				token, err := s.BuildJWTString(claims, testSecret)
				require.NoError(t, err)
				return token
			},
			secretKey:   testSecret,
			wantErrType: &ErrMissingUserID{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr := tt.setup()
			userID, err := s.GetUserIDFromJWTString(tokenStr, tt.secretKey)

			if tt.wantErrType == nil {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUserID, userID)
			} else {
				require.Error(t, err)
				assert.Empty(t, userID)
			}
		})
	}
}

func TestErrMissingUserID_Error(t *testing.T) {
	claims := NewClaims("")
	err := &ErrMissingUserID{Claims: claims}
	assert.Contains(t, err.Error(), "missing user ID")
}

func TestErrInvalidToken_Error(t *testing.T) {
	err := &ErrInvalidToken{Err: "invalid token"}
	assert.Equal(t, "invalid token", err.Error())
}
