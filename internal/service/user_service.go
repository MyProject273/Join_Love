package service

import (
	"context"
	"errors"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	mapper "github.com/MyProject273/Join_Love/internal/dto"
	user_dto "github.com/MyProject273/Join_Love/internal/dto/user"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/MyProject273/Join_Love/pkg/utils/rescode"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

type (
	UserService interface {
		GetUser(ctx context.Context, req user_dto.GetUserReq) (user_dto.GetUserRes, error)
		CreateUserByAdmin(ctx context.Context, req user_dto.CreateUserByAdminReq) (user_dto.CreateUserByAdminRes, error)
		UpdateUserActiveStatus(ctx context.Context, req user_dto.UpdateUserActiveStatusReq) (user_dto.UpdateUserActiveStatusRes, error)
		DeleteUser(ctx context.Context, req user_dto.DeleteUserReq) (user_dto.DeleteUserRes, error)
		UpdateUser(ctx context.Context, req user_dto.UpdateUserReq) (user_dto.UpdateUserRes, error)
	}

	userService struct {
		store  db.Store
		logger zerolog.Logger
	}
)

func NewUserService(store db.Store, logger zerolog.Logger) UserService {
	return &userService{store: store, logger: logger}
}

func (u *userService) mapError(err error, userID string, action string) error {
	if errors.Is(err, pgx.ErrNoRows) {
		u.logger.Warn().Str("user_id", userID).Msg("user not found")
		return rescode.UserNotFound
	}
	u.logger.Error().Err(err).Str("user_id", userID).Msgf("failed to %s user", action)
	return rescode.Internal
}

func (u *userService) CreateUserByAdmin(ctx context.Context, req user_dto.CreateUserByAdminReq) (res user_dto.CreateUserByAdminRes, err error) {
	user, err := u.store.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		u.logger.Error().Err(err).Str("email", req.Email).Msg("failed to get user by email")
		return res, rescode.Internal
	}
	if user.Email != "" {
		u.logger.Warn().Str("email", req.Email).Msg("user already exists")
		return res, rescode.UserAlreadyExists
	}

	if len(req.Roles) == 0 {
		return res, rescode.Invalid
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return res, rescode.Internal
	}

	birthDate, err := utils.StringToPgDate(*req.Birthdate)
	if err != nil {
		return res, rescode.Invalid
	}

	result, err := u.store.CreateUserByAdminTx(ctx, db.CreateUserByAdminTxParams{
		CreateUserByAdminParams: db.CreateUserByAdminParams{
			Email:        req.Email,
			UserName:     req.UserName,
			Phone:        req.Phone,
			PasswordHash: hashedPassword,
			FullName:     req.FullName,
			Gender:       req.Gender,
			Birthdate:    birthDate,
			AvatarUrl:    req.AvatarUrl,
			Bio:          req.Bio,
			IsActive:     req.IsActive,
			IsVerified:   req.IsVerified,
		},
		Roles: req.Roles,
	})
	if err != nil {
		return res, rescode.Internal
	}

	res.UserRes = mapper.ConvertUser(result.User)
	return res, nil
}

func (u *userService) GetUser(ctx context.Context, req user_dto.GetUserReq) (res user_dto.GetUserRes, err error) {
	userID, err := utils.StringToPgUUID(req.ID)
	if err != nil {
		return res, rescode.Invalid
	}

	user, err := u.store.GetUser(ctx, userID)
	if err != nil {
		return res, u.mapError(err, userID.String(), "get")
	}

	res.UserRes = mapper.ConvertUser(user)
	return res, nil
}

func (u *userService) UpdateUserActiveStatus(ctx context.Context, req user_dto.UpdateUserActiveStatusReq) (res user_dto.UpdateUserActiveStatusRes, err error) {
	user_id, err := utils.StringToPgUUID(req.ID)
	if err != nil {
		return res, u.mapError(err, req.ID, " update active status")
	}
	user, err := u.store.UpdateUserActiveStatus(ctx, db.UpdateUserActiveStatusParams{
		IsActive: req.IsActive,
		ID:       user_id,
	})
	if err != nil {
		return res, u.mapError(err, req.ID, "update active status")
	}

	res.UserRes = mapper.ConvertUser(user)
	return res, nil
}

func (u *userService) DeleteUser(ctx context.Context, req user_dto.DeleteUserReq) (res user_dto.DeleteUserRes, err error) {
	userID, err := utils.StringToPgUUID(req.ID)
	if err != nil {
		return res, rescode.Invalid
	}

	err = u.store.DeleteUser(ctx, userID)
	if err != nil {
		return res, u.mapError(err, req.ID, "delete")
	}

	res.Success = true
	res.Message = "delete user successful"
	return res, nil
}

func (u *userService) UpdateUser(ctx context.Context, req user_dto.UpdateUserReq) (res user_dto.UpdateUserRes, err error) {
	userID, err := utils.StringToPgUUID(req.ID)
	if err != nil {
		return res, rescode.Invalid
	}

	if _, err := u.store.GetUser(ctx, userID); err != nil {
		return res, u.mapError(err, req.ID, "update")
	}

	birthDate, err := utils.StringToPgDate(*req.Birthdate)
	if err != nil {
		return res, rescode.Invalid
	}

	result, err := u.store.UpdateUserTx(ctx, db.UpdateUserTxParams{
		UpdateUserParams: db.UpdateUserParams{
			UserName:   req.UserName,
			Phone:      req.Phone,
			FullName:   req.FullName,
			Gender:     req.Gender,
			Birthdate:  birthDate,
			AvatarUrl:  req.AvatarUrl,
			Bio:        req.Bio,
			IsActive:   req.IsActive,
			IsVerified: req.IsVerified,
		},
		Roles: *req.Roles,
	})
	if err != nil {
		return res, rescode.Internal
	}

	res.UserRes = mapper.ConvertUser(result.User)
	return res, nil
}
