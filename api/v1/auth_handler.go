package v1

import (
	"github.com/MyProject273/Join_Love/api/helper"
	auth_dto "github.com/MyProject273/Join_Love/internal/dto/auth"
	"github.com/MyProject273/Join_Love/internal/service"
	"github.com/MyProject273/Join_Love/pkg/utils/rescode"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type (
	AuthHandler interface {
		Signup(ctx *gin.Context)
		Login(ctx *gin.Context)
	}

	authHandler struct {
		authService service.AuthService
	}
)

func NewAuthHandler(authService service.AuthService) AuthHandler {
	return &authHandler{
		authService: authService,
	}
}

func (a *authHandler) Signup(ctx *gin.Context) {
	var req auth_dto.SignupReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to bind JSON for signup")
		helper.RespondValidationError(ctx, err)
		return
	}
	res, err := a.authService.Signup(ctx, req)
	if err != nil {
		log.Err(err).
			Str("email", req.Email).
			Msg("Failed to signup")
		helper.RespondError(ctx, rescode.Internal, err)
		return
	}
	log.Info().Str("email", res.Email).Msg("Signup successful")
	helper.RespondSuccess(ctx, res, rescode.Success)
}

func (a *authHandler) Login(ctx *gin.Context) {
	var req auth_dto.LoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to bind JSON for login")
		helper.RespondValidationError(ctx, err)
		return
	}

	res, err := a.authService.Login(ctx, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("email", req.Email).
			Msg("Failed to login")

		helper.RespondError(ctx, rescode.Internal, err)
		return
	}

	log.Info().
		Str("email", req.Email).
		Str("user_id", res.UserID).
		Msg("User login successful")

	helper.RespondSuccess(ctx, res, rescode.Success)
}
