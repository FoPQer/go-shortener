package middlewares

import (
	"context"
	"strings"

	"github.com/FoPQer/go-shortener/internal/auth"
	"github.com/FoPQer/go-shortener/internal/model"
	"github.com/FoPQer/go-shortener/internal/service"
	"github.com/FoPQer/go-shortener/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewGRPCAuthInterceptor returns a unary interceptor that mirrors HTTP auth middleware:
// valid JWT in "authorization" metadata → extract user ID;
// missing or invalid token → create new user.
func NewGRPCAuthInterceptor(userService *service.UserService, claimsService *auth.ClaimsService) grpc.UnaryServerInterceptor {
	secretKey := []byte(service.GetSecretKey())

	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		userID, err := resolveUserID(ctx, userService, claimsService, secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "auth failed: %v", err)
		}

		ctx = context.WithValue(ctx, utils.UserID("userID"), userID)
		return handler(ctx, req)
	}
}

// resolveUserID extracts or creates a user ID from incoming gRPC metadata.
func resolveUserID(ctx context.Context, userService *service.UserService, claimsService *auth.ClaimsService, secretKey []byte) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		if values := md.Get("authorization"); len(values) > 0 {
			token := strings.TrimPrefix(values[0], "Bearer ")
			userID, err := claimsService.GetUserIDFromJWTString(token, secretKey)
			if err == nil {
				return userID, nil
			}
		}
	}

	return userService.Create(ctx, &model.User{})
}
