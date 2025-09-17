package v1

import (
	"context"
	"errors"
	"time"

	"github.com/MyProject273/Join_Love/gapi/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/mail"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/i18n"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	store           db.Store
	config          config.Config
	tokenMaker      token.Maker
	logger          zerolog.Logger
	taskDistributor worker.TaskDistributor
}

func NewAuthServer(config config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger, taskDistributor worker.TaskDistributor) *AuthServer {
	return &AuthServer{
		store:           store,
		config:          config,
		tokenMaker:      tokenMaker,
		logger:          logger,
		taskDistributor: taskDistributor,
	}
}

func (a *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	lang := helper.ExtractMetadata(ctx).Lang

	if err := helper.ValidateAll(req); err != nil {
		a.logger.Warn().Err(err).Str("email", req.GetEmail()).Msg("invalid login request")
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	user, err := a.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			a.logger.Info().Str("email", req.GetEmail()).Msg("user not found")
			return nil, status.Errorf(codes.NotFound, "%s", i18n.GetI18nMessage("user_not_found", lang))
		}

		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to get user by email")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		a.logger.Info().Str("email", req.GetEmail()).Msg("incorrect password")
		return nil, status.Errorf(codes.NotFound, "%s", i18n.GetI18nMessage("incorrect_password", lang))
	}

	userID, err := utils.PgUUIDToString(user.ID)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to convert UUID to string")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
		userID, user.Role, a.config.AccessTokenDuration, consts.TokenTypeAccessToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create access token")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(
		userID, user.Role, a.config.RefreshTokenDuration, consts.TokenTypeRefreshToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create refresh token")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	err = a.store.UpdateUserLastLogin(ctx, db.UpdateUserLastLoginParams{
		LastLogin: pgtype.Timestamp{Time: time.Now(), Valid: true},
		ID:        user.ID,
	})

	if err != nil {
		a.logger.Error().Err(err).Str("user_id", userID).Msg("failed to update last login")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	a.logger.Info().Str("user_id", userID).Msg("user logged in successfully")

	return &auth.LoginResponse{
		User:                  helper.ConvertUser(&user),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}, nil
}

func (a *AuthServer) Signup(ctx context.Context, req *auth.SignupRequest) (*auth.SignupResponse, error) {
	lang := helper.ExtractMetadata(ctx).Lang

	if err := helper.ValidateAll(req); err != nil {
		a.logger.Warn().Err(err).Str("email", req.GetEmail()).Msg("invalid signup request")
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}
	verifyEmail, err := mail.ValidateEmailAddress(req.Email, a.config.AbstractApiKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "email verify failed: %v", err)
	}
	if verifyEmail.Deliverability != "DELIVERABLE" ||
		!verifyEmail.IsValidFormat.Value ||
		!verifyEmail.IsMxFound.Value ||
		!verifyEmail.IsSmtpValid.Value ||
		verifyEmail.IsDisposableEmail.Value {
		return nil, status.Errorf(codes.InvalidArgument, "%s", i18n.GetI18nMessage("invalid_email", lang))
	}

	if _, err := a.store.GetUserByEmail(ctx, req.GetEmail()); err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "%s", i18n.GetI18nMessage("user_already_exists", lang))
	} else {
		var pgErr *pgx.PgError
		if errors.As(err, &pgErr) && pgErr.Message == pgx.ErrNoRows.Error() {
			a.logger.Error().Err(err).Str("email", req.GetEmail()).Msgf("failed to check existing user (type: %T)", err)
			return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
		}
	}
	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to hash password")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Email:        req.GetEmail(),
			PasswordHash: hashedPassword,
			UserName:     req.GetUserName(),
		},
		AfterCreate: func(user db.User) error {
			taskPayload := worker.PayloadSendVerifyEmail{
				Email: user.Email,
				Lang:  lang,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(10 * time.Second),
				asynq.Queue(worker.QueueCritical),
			}
			return a.taskDistributor.DistributeTaskSendVerifyEmail(ctx, &taskPayload, opts...)

		},
	}

	result, err := a.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == consts.UniqueViolation {
			switch db.ErrorConstraint(err) {
			case "users_user_name_key":
				return nil, status.Errorf(codes.AlreadyExists, "%s", i18n.GetI18nMessage("users_user_name_key", lang))
			case "users_email_key":
				return nil, status.Errorf(codes.AlreadyExists, "%s", i18n.GetI18nMessage("users_email_key", lang))
			default:
				return nil, status.Errorf(codes.AlreadyExists, "%s", i18n.GetI18nMessage("user_already_exists", lang))
			}
		}

		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to create user")
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("internal_error", lang))
	}

	a.logger.Info().Str("user_id", result.User.ID.String()).Msg("user signed up successfully")

	return &auth.SignupResponse{
		User: helper.ConvertUser(&result.User),
	}, nil
}

func (a *AuthServer) VerifyEmail(ctx context.Context, req *auth.VerifyEmailRequest) (res *auth.VerifyEmailResponse, err error) {
	lang := helper.ExtractMetadata(ctx).Lang

	if err = helper.ValidateAll(req); err != nil {
		a.logger.Warn().Err(err).Str("email_id", req.EmailId).Msg("invalid verify email request")
		return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
	}

	txResult, err := a.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		VerifyEmailId:    req.EmailId,
		SecretCode: req.SecretCode,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%s", i18n.GetI18nMessage("verify_email", lang))
	}

	rsp := &auth.VerifyEmailResponse{
		IsVerified: txResult.User.IsVerified,
	}
	return rsp, nil
}
