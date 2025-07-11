package v1

import (
	"context"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	store db.Store
}

func NewAuthServer(config config.Config, store db.Store) *AuthServer {
	return &AuthServer{
		store: store,
	}
}

func (a *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (res *auth.LoginResponse, err error) {
	return nil, nil
}
