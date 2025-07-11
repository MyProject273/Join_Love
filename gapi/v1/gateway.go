package v1

import (
	"context"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func RegisterAllHandlers(ctx context.Context, mux *runtime.ServeMux, cfg config.Config, store db.Store) error {
	authServer := NewAuthServer(cfg, store)

	if err := auth.RegisterAuthServiceHandlerServer(ctx, mux, authServer); err != nil {
		return err
	}
	return nil
}
