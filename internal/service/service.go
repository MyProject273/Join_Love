package service

import (
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/rs/zerolog"
)

type Services struct {
	AuthService AuthService
}

func NewServices(store db.Store, config config.Config, tokenMaker token.Maker, taskDistributor worker.TaskDistributor, logger zerolog.Logger) *Services {
	return &Services{
		AuthService: NewAuthService(store, config, tokenMaker, taskDistributor, logger),
	}
}
