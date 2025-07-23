package v1

import (
	"github.com/MyProject273/Join_Love/gapi/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Server struct {
	GRPCServer *grpc.Server
	Config     config.Config
	Store      db.Store
	logger     zerolog.Logger
}

func NewServer(cfg config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger, taskDistributor worker.TaskDistributor) (*Server, error) {
	grpcLogger := grpc.UnaryInterceptor(helper.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	authServer := NewAuthServer(cfg, store, tokenMaker, logger, taskDistributor)

	auth.RegisterAuthServiceServer(grpcServer, authServer)

	return &Server{
		GRPCServer: grpcServer,
		Config:     cfg,
		Store:      store,
	}, nil
}
