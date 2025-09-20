package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/MyProject273/Join_Love/api/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/utils/rescode"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// AuthMiddleware xác thực JWT và lưu payload vào context
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(consts.AuthorizationHeaderKey)
		log.Debug().
			Str("path", ctx.FullPath()).
			Str("method", ctx.Request.Method).
			Str("auth_header", authorizationHeader).
			Msg("checking authorization header")

		if authorizationHeader == "" {
			err := errors.New("authorization is not provided")
			log.Warn().Err(err).Msg("missing authorization header")
			helper.AbortWithErrorResponse(ctx, rescode.UnAuthorized, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			log.Warn().Err(err).Str("auth_header", authorizationHeader).Msg("invalid authorization header format")
			helper.AbortWithErrorResponse(ctx, rescode.UnAuthorized, err)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != consts.AuthorizationType {
			err := errors.New("unsupported authorization type")
			log.Warn().Err(err).Str("auth_header", authorizationType).Msg("unsupported authorization type")
			helper.AbortWithErrorResponse(ctx, rescode.UnAuthorized, err)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken, consts.TokenTypeAccessToken)
		if err != nil {
			log.Warn().Err(err).Msg("token verification failed")
			helper.AbortWithErrorResponse(ctx, rescode.UnAuthorized, err)
			return
		}

		log.Info().
			Str("user_id", payload.UserID.String()).
			Str("token_id", payload.ID.String()).
			Msg("token verified successfully")

		// Lưu payload vào context với key chuẩn
		ctx.Set(consts.AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}

func RequirePermission(permission string, redisClient *redis.Client, store db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get(consts.AuthorizationPayloadKey)
		if !exists {
			helper.AbortWithErrorResponse(ctx, rescode.UnAuthorized, errors.New("failed to decode token"))
			return
		}

		claims := payload.(*token.Payload)
		userID := claims.UserID

		cacheKey := "user:" + userID.String() + ":permissions"

		perms, err := redisClient.SMembers(ctx, cacheKey).Result()
		log.Print(perms)
		if err != nil || len(perms) == 0 {
			dbPerms, err := store.ListPermissionsByUser(ctx, userID)
			log.Print(perms)
			if err != nil {
				log.Error().Err(err).Msg("failed to load permissions")
				helper.AbortWithErrorResponse(ctx, rescode.Internal, errors.New("failed to load permissions"))
				return
			}

			perms = make([]string, len(dbPerms))
			for i, p := range dbPerms {
				perms[i] = p.Name
			}

			if len(perms) > 0 {
				redisClient.SAdd(ctx, cacheKey, perms)
				redisClient.Expire(ctx, cacheKey, time.Minute*2)
			}
		}

		for _, p := range perms {
			if p == permission {
				ctx.Next()
				return
			}
		}

		helper.AbortWithErrorResponse(ctx, rescode.PermissionDenied, errors.New("permission denied"))
	}
}
