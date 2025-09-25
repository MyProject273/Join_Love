package v1

import (
	"github.com/MyProject273/Join_Love/api/helper"
	user_dto "github.com/MyProject273/Join_Love/internal/dto/user"
	"github.com/MyProject273/Join_Love/internal/service"
	"github.com/MyProject273/Join_Love/pkg/utils/rescode"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type (
	UserHandler interface {
		GetUser(ctx *gin.Context)
		CreateUserByAdmin(ctx *gin.Context)
		UpdateUserActiveStatus(ctx *gin.Context)
		DeleteUser(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
	}

	userHandler struct {
		userService service.UserService
	}
)

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}

func (u *userHandler) GetUser(ctx *gin.Context) {
	var req user_dto.GetUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Error().Err(err).Msg("failed to should bind Uri")
		helper.RespondValidationError(ctx, err)
		return
	}
	res, err := u.userService.GetUser(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to get user")
		helper.RespondError(ctx, err)
		return
	}

	log.Info().Str("user_id", req.ID).Msg("get user successful")
	helper.RespondSuccess(ctx, res, rescode.Success)
}

func (u *userHandler) CreateUserByAdmin(ctx *gin.Context) {
	var req user_dto.CreateUserByAdminReq
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind body")
		helper.RespondValidationError(ctx, err)
		return
	}
	res, err := u.userService.CreateUserByAdmin(ctx, req)

	if err != nil {
		log.Error().Err(err).Msg("failed to create user by admin")
		helper.RespondError(ctx, err)
		return
	}

	log.Info().Str("user_id", res.ID).Msg("create user by admin successful")
	helper.RespondSuccess(ctx, res, rescode.Success)
}

func (u *userHandler) DeleteUser(ctx *gin.Context) {
	var req user_dto.DeleteUserReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Error().Err(err).Msg("failed to bind uri")
		helper.RespondValidationError(ctx, err)
		return
	}

	res, err := u.userService.DeleteUser(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("failed to delete user")
		helper.RespondError(ctx, err)
		return
	}

	log.Info().Str("user_id", req.ID).Msg("delete user successful")
	helper.RespondSuccess(ctx, res, rescode.Success)
}

func (u *userHandler) UpdateUserActiveStatus(ctx *gin.Context) {
	var uriReq user_dto.GetUserReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		log.Error().Err(err).Msg("failed to bind uri")
		helper.RespondValidationError(ctx, err)
		return
	}

	var bodyReq user_dto.UpdateUserActiveStatusReq
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		log.Error().Err(err).Msg("failed to bind body")
		helper.RespondValidationError(ctx, err)
		return
	}
	bodyReq.ID = uriReq.ID
	res, err := u.userService.UpdateUserActiveStatus(ctx, bodyReq)
	if err != nil {
		log.Error().Err(err).Str("user_id", uriReq.ID).Msg("failed to update user active status")
		helper.RespondError(ctx, err)
		return
	}

	log.Info().Str("user_id", uriReq.ID).Msg("update user active status successful")
	helper.RespondSuccess(ctx, res, rescode.Success)
}

func (u *userHandler) UpdateUser(ctx *gin.Context) {
	var uriReq user_dto.GetUserReq
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		log.Error().Err(err).Msg("failed to bind uri")
		helper.RespondValidationError(ctx, err)
		return
	}

	var bodyReq user_dto.UpdateUserReq
	if err := ctx.ShouldBindJSON(&bodyReq); err != nil {
		log.Error().Err(err).Msg("failed to bind body")
		helper.RespondValidationError(ctx, err)
		return
	}

	bodyReq.ID = uriReq.ID

	res, err := u.userService.UpdateUser(ctx, bodyReq)
	if err != nil {
		log.Error().Err(err).Str("user_id", bodyReq.ID).Msg("failed to update user")
		helper.RespondError(ctx, err)
		return
	}

	log.Info().Str("user_id", bodyReq.ID).Msg("update user successful")
	helper.RespondSuccess(ctx, res, rescode.Success)
}
