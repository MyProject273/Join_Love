package service

import (
	"errors"
	"fmt"
	"time"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	auth_dto "github.com/MyProject273/Join_Love/internal/dto/auth"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pkg/config"
	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/utils"
	apperr "github.com/MyProject273/Join_Love/pkg/utils/error"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
)

type (
	AuthService interface {
		Signup(ctx *gin.Context, req auth_dto.SignupReq) (auth_dto.SignupRes, error)
		Login(ctx *gin.Context, req auth_dto.LoginReq) (auth_dto.LoginRes, error)
	}

	authService struct {
		store           db.Store
		config          config.Config
		tokenMaker      token.Maker
		taskDistributor worker.TaskDistributor
		logger          zerolog.Logger
	}
)

func NewAuthService(store db.Store, config config.Config, tokenMaker token.Maker, taskDistributor worker.TaskDistributor, logger zerolog.Logger) AuthService {
	return &authService{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
		logger:          logger,
	}
}

func (a *authService) Signup(ctx *gin.Context, req auth_dto.SignupReq) (res auth_dto.SignupRes, err error) {
	user, err := a.store.GetUserByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		a.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to get user by email during signup")
		return res, fmt.Errorf("get user error: %w", err)
	}
	if user.Email != "" {
		a.logger.Warn().Str("email", req.Email).Msg("User already exists")
		return res, fmt.Errorf("%w: user with email %s already exists", apperr.ErrUserAlreadyExists, req.Email)
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to hash password")
		return res, fmt.Errorf("%w: %v", apperr.ErrHashPassword, err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			UserName:     req.UserName,
			PasswordHash: hashedPassword,
			Email:        req.Email,
		},
		AfterCreate: func(user db.User) error {
			taskPayload := worker.PayloadSendVerifyEmail{
				Email: user.Email,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}

			err := a.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &taskPayload, opts...)
			if err != nil {
				a.logger.Error().Err(err).Str("email", user.Email).Msg("Failed to enqueue verification email")
			} else {
				a.logger.Info().Str("email", user.Email).Msg("Verification email enqueued")
			}
			return err
		},
	}

	txResult, err := a.store.CreateUserTx(ctx, arg)
	if err != nil {
		a.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to create user transaction")
		return res, fmt.Errorf("%w: %s", apperr.ErrCreateUser, err)
	}

	res = auth_dto.SignupRes{
		UserID:   txResult.User.ID.String(),
		UserName: txResult.User.UserName,
		Email:    txResult.User.Email,
	}
	return res, nil
}

func (a *authService) Login(ctx *gin.Context, req auth_dto.LoginReq) (res auth_dto.LoginRes, err error) {
	user, err := a.store.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Warn().Str("email", req.Email).Msg("User not found")
			return res, fmt.Errorf("%w: %v", apperr.ErrUserNotFound, err)
		}
		a.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to get user")
		return res, fmt.Errorf("get user error: %w", err)
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		a.logger.Warn().Str("email", req.Email).Msg("Incorrect password")
		return res, apperr.ErrMismatchPassword
	}

	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
		user.UserName,
		user.Role,
		a.config.AccessTokenDuration,
		consts.TokenTypeAccessToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user", user.UserName).Msg("Failed to create access token")
		return res, fmt.Errorf("%w: %v", apperr.ErrCreateAccessToken, err)
	}

	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(
		user.UserName,
		user.Role,
		a.config.RefreshTokenDuration,
		consts.TokenTypeRefreshToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user", user.UserName).Msg("Failed to create refresh token")
		return res, fmt.Errorf("%w: %v", apperr.ErrCreateRefreshToken, err)
	}

	expiresAt := pgtype.Timestamptz{
		Time:  refreshPayload.ExpiredAt,
		Valid: true,
	}

	session, err := a.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           utils.UUIDToPgUUIDSafe(refreshPayload.ID),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    expiresAt,
	})
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to create session")
		return res, fmt.Errorf("%w: %v", apperr.ErrCreateSession, err)
	}

	err = a.store.UpdateUserLastLogin(ctx, db.UpdateUserLastLoginParams{
		LastLogin: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ID:        user.ID,
	})
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to update last login")
		return res, fmt.Errorf("%w: %v", apperr.ErrInternalServer, err)
	}

	res = auth_dto.LoginRes{
		Email:     user.Email,
		UserID:    user.ID.String(),
		UserName:  user.UserName,
		SessionID: session.ID.String(),
		Token: auth_dto.Token{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresIn:  accessPayload.ExpiredAt,
			RefreshTokenExpiresIn: refreshPayload.ExpiredAt,
		},
	}
	return res, nil
}
