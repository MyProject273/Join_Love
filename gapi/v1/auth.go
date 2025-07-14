package v1

import (
	"context"
	"errors"

	"github.com/MyProject273/Join_Love/gapi/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/i18n"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	store      db.Store
	config     config.Config
	tokenMaker token.Maker
	logger     zerolog.Logger
}

func NewAuthServer(config config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger) *AuthServer {
	return &AuthServer{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
		logger:     logger,
	}
}

func (a *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (res *auth.LoginResponse, err error) {

	if err := helper.ValidateAll(req); err != nil {
		a.logger.Warn().Err(err).Str("email", req.GetEmail()).Msg("invalid login request")
		return nil, err
	}

	user, err := a.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			a.logger.Info().Str("email", req.GetEmail()).Msg("user not found")
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to get user by email")
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		a.logger.Info().Str("email", req.GetEmail()).Msg("incorrect password")
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	userID, err := utils.PgUUIDToString(user.ID)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to convert UUID to string")
		return nil, status.Errorf(codes.Internal, "failed to convert UUID to string")
	}

	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
		userID, user.Role, a.config.AccessTokenDuration, consts.TokenTypeAccessToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create access token")
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(
		userID, user.Role, a.config.RefreshTokenDuration, consts.TokenTypeRefreshToken,
	)
	if err != nil {
		a.logger.Error().Err(err).Str("user_id", userID).Msg("failed to create refresh token")
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
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
		return nil, err
	}

	_, err := a.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil && err != pgx.ErrNoRows {
		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to check existing user")
		return nil, status.Errorf(codes.Internal, "failed to get user by email")
	}
	if err == nil {
		a.logger.Info().Str("email", req.GetEmail()).Msg("user already exists")
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	hashedPassword, err := utils.HashPassword(req.GetPassword())
	if err != nil {
		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to hash password")
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Email:        req.GetEmail(),
			PasswordHash: hashedPassword,
			UserName:     req.GetUserName(),
		},
		AfterCreate: func(user db.User) error {
			return nil
		},
	}

	result, err := a.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == consts.UniqueViolation {
			switch db.ErrorConstraint(err) {
			case "users_user_name_key":
				return nil, status.Errorf(codes.AlreadyExists, "%s", i18n.GetI18nMessage("users_user_name_key", lang))

			case "users_email_key":
				return nil, status.Errorf(codes.AlreadyExists, "email already exists")

			default:
				return nil, status.Errorf(codes.AlreadyExists, "user already exists")
			}

		}
		a.logger.Error().Err(err).Str("email", req.GetEmail()).Msg("failed to create user")
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	a.logger.Info().Str("user_id", result.User.ID.String()).Msg("user signed up successfully")

	return &auth.SignupResponse{
		User: helper.ConvertUser(&result.User),
	}, nil
}
