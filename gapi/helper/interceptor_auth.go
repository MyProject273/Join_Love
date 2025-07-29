package helper

import (
	"context"
	"strings"

	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contextKey string

const (
	AuthPayloadKey contextKey = "authPayload"
)

type AuthInterceptor struct {
	tokenMaker token.Maker
}

func NewAuthInterceptor(tokenMaker token.Maker) *AuthInterceptor {
	return &AuthInterceptor{tokenMaker: tokenMaker}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		authHeader := values[0]
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		accessToken := fields[1]
		payload, err := interceptor.tokenMaker.VerifyToken(accessToken, consts.TokenTypeAccessToken)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid access token")
		}

		ctx = context.WithValue(ctx, AuthPayloadKey, payload)

		return handler(ctx, req)
	}
}

func GetAuthPayload(ctx context.Context) (*token.Payload, bool) {
	payload, ok := ctx.Value(AuthPayloadKey).(*token.Payload)
	return payload, ok
}

func HasPermission(userRole string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
