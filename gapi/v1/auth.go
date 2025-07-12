package v1

import (
	"context"
	"errors"

	"github.com/MyProject273/Join_Love/gapi/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
	consts "github.com/MyProject273/Join_Love/pkg/const"
	"github.com/MyProject273/Join_Love/pkg/utils"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	store      db.Store
	config     config.Config
	tokenMaker token.Maker
}

func NewAuthServer(config config.Config, store db.Store) *AuthServer {
	return &AuthServer{
		store: store,
	}
}

func (a *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (res *auth.LoginResponse, err error) {
	if err := helper.ValidateAll(req); err != nil {
		return nil, err
	}
	user, err := a.store.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.NotFound, "failed to find user")
	}

	err = utils.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	user_id, err := utils.PgUUIDToString(user.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert UUID to string")
	}
	accessToken, accessPayload, err := a.tokenMaker.CreateToken(
		user_id,
		user.Role,
		a.config.AccessTokenDuration,
		consts.TokenTypeAccessToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := a.tokenMaker.CreateToken(
		user_id,
		user.Role,
		a.config.RefreshTokenDuration,
		consts.TokenTypeRefreshToken,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	rsp := &auth.LoginResponse{
		User:                  helper.ConvertUser(user),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}
	return rsp, nil
}
