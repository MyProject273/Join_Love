package helper

import (
	"context"
	"net/http"
	"strings"

	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
)

func HttpAuthMiddleware(tokenMaker token.Maker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			fields := strings.Fields(authHeader)
			if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
				http.Error(w, "invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenStr := fields[1]
			payload, err := tokenMaker.VerifyToken(tokenStr, consts.TokenTypeAccessToken)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), AuthPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
