package v1

import (
	"github.com/MyProject273/Join_Love/gapi/helper"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pb/auth"
	"github.com/MyProject273/Join_Love/pb/user"
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
	authInterceptor := helper.NewAuthInterceptor(tokenMaker)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			helper.GrpcLogger,
			authInterceptor.Unary(),
		),
	)

	authServer := NewAuthServer(cfg, store, tokenMaker, logger, taskDistributor)
	userServer := NewUserServer(cfg, store, logger)

	auth.RegisterAuthServiceServer(grpcServer, authServer)
	user.RegisterUserServiceServer(grpcServer, userServer)

	return &Server{
		GRPCServer: grpcServer,
		Config:     cfg,
		Store:      store,
	}, nil
}
