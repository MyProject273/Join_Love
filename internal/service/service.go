package service

import (
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
)

type Services struct {
	AuthService AuthService
}

func NewServices(store db.Store, config config.Config, tokenMaker token.Maker, taskDistributor worker.TaskDistributor) *Services {
	return &Services{
		AuthService: NewAuthService(store, config, tokenMaker, taskDistributor),
	}
}
