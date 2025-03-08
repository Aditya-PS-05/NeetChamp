// jwt_middleware.go (Middleware to Validate JWT Tokens)
package middleware

import (
	"context"
	"errors"
	"strings"

	"github.com/Aditya-PS-05/NeetChamp/auth-service/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Define contextKey type to avoid collisions
type contextKey string

const (
	userEmailKey contextKey = "userEmail"
	userRoleKey  contextKey = "userRole"
)

// JWTAuthInterceptor - gRPC interceptor for token verification
func JWTAuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Extract metadata from request
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("missing metadata")
		}

		// Extract token from metadata
		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, errors.New("authorization token not provided")
		}

		token := strings.TrimPrefix(tokens[0], "Bearer ")

		// Verify token validity
		claims, err := utils.VerifyToken(token)
		if err != nil {
			return nil, errors.New("invalid or expired token")
		}

		// Add user data to context for further processing
		ctx = context.WithValue(ctx, userEmailKey, claims.Email)
		ctx = context.WithValue(ctx, userRoleKey, claims.Role)

		return handler(ctx, req)
	}
}
