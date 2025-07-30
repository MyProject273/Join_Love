package service

import (
	"context"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	auth_dto "github.com/MyProject273/Join_Love/internal/dto/auth"
)

type (
	AuthService interface {
		Signup(ctx context.Context, req auth_dto.SignupReq) (auth_dto.SignupRes, error)
		Login(ctx context.Context, req auth_dto.LoginReq) (auth_dto.LoginRes, error)
	}

	authService struct {
		store db.Store
	}
)

func NewAuthService(store db.Store) AuthService {
	return &authService{
		store: store,
	}
}

func (a *authService) Signup(ctx context.Context, req auth_dto.SignupReq) (res auth_dto.SignupRes, err error) {
	return auth_dto.SignupRes{}, nil
}

func (a *authService) Login(ctx context.Context, req auth_dto.LoginReq) (res auth_dto.LoginRes, err error) {
	return auth_dto.LoginRes{}, nil
}
