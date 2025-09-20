package service

import (
	"errors"
	"time"

	"github.com/MyProject273/Join_Love/api/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	auth_dto "github.com/MyProject273/Join_Love/internal/dto/auth"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pkg/config"
	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/MyProject273/Join_Love/pkg/utils/rescode"
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
		VerifyEmail(ctx *gin.Context, req auth_dto.VerifyEmailReq) (auth_dto.VerifyEmailRes, error)
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
		return res, rescode.Internal
	}
	if user.Email != "" {
		a.logger.Warn().Str("email", req.Email).Msg("User already exists")
		return res, rescode.UserAlreadyExists
	}
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		a.logger.Error().Err(err).Msg("Failed to hash password")
		return res, rescode.Internal
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
				Lang:  helper.GetLang(ctx),
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
		return res, rescode.UserCreateFailed
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
			return res, rescode.UserNotFound
		}
		a.logger.Error().Err(err).Str("email", req.Email).Msg("Failed to get user")
		return res, rescode.Internal
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		a.logger.Warn().Str("email", req.Email).Msg("Incorrect password")
		return res, rescode.InvalidCredentials
	}

	isVerified, err := a.store.CheckUserVerified(ctx, user.ID)
	if err != nil {
		a.logger.Warn().Str("email", req.Email).Msg("Failed to check user verified")
		return res, rescode.Internal
	}
	if isVerified == false {
		return res, rescode.UserNotVerified
	}

	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
		user.UserName,
		a.config.AccessTokenDuration,
		consts.TokenTypeAccessToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user", user.UserName).Msg("Failed to create access token")
		return res, rescode.Internal
	}

	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(
		user.UserName,
		a.config.RefreshTokenDuration,
		consts.TokenTypeRefreshToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user", user.UserName).Msg("Failed to create refresh token")
		return res, rescode.Internal
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
		return res, rescode.Internal
	}

	err = a.store.UpdateUserLastLogin(ctx, db.UpdateUserLastLoginParams{
		LastLogin: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ID:        user.ID,
	})
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to update last login")
		return res, rescode.Internal
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

func (a *authService) VerifyEmail(ctx *gin.Context, req auth_dto.VerifyEmailReq) (res auth_dto.VerifyEmailRes, err error) {
	txResult, err := a.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		VerifyEmailId: req.VerifyEmailID,
		SecretCode:    req.SecretCode,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			a.logger.Warn().Str("verify_email_id", req.VerifyEmailID).Msg("VerifyEmail record not found")
			return res, rescode.Internal
		}
	}

	if txResult.User.IsVerified {
		res.IsVerified = "true"
	} else {
		res.IsVerified = "false"
	}

	return res, nil
}
