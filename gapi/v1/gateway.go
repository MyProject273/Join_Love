package v1

import (
	"context"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
)

func RegisterAllHandlers(ctx context.Context, mux *runtime.ServeMux, cfg config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger, taskDistributor worker.TaskDistributor) error {
	authServer := NewAuthServer(cfg, store, tokenMaker, logger, taskDistributor)

	if err := auth.RegisterAuthServiceHandlerServer(ctx, mux, authServer); err != nil {
		return err
	}
	return nil
}
