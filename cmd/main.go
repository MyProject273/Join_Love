package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MyProject273/Join_Love/gapi"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	pb "github.com/MyProject273/Join_Love/pb"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	config, store := initializeApp(ctx)
	waitGroup, ctx := errgroup.WithContext(ctx)

	runGrpcServer(ctx, config, store, waitGroup)
	if err := waitGroup.Wait(); err != nil {
		log.Error().Err(err).Msg("application exited with error")
	} else {
		log.Info().Msg("application exited gracefully")
	}

}

func initializeApp(ctx context.Context) (cfg config.Config, store db.Store) {
	var err error
	cfg, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if cfg.Environment == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	connPool, err := pgxpool.New(ctx, cfg.DB_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect db")
	}

	if err := connPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("database not reachable")
	}

	store = db.NewStore(connPool)

	return
}

func runGrpcServer(
	ctx context.Context,
	config config.Config,
	store db.Store,
	waitGroup *errgroup.Group,
) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)

	pb.RegisterJoinLoveServer(grpcServer, server)
	log.Info().Msgf("registered grpc services: %+v", grpcServer.GetServiceInfo())
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on grpc address")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("gRPC server running at %s", config.GrpcServerAddress)
		err := grpcServer.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error().Err(err).Msg("gRPC server failed")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Shutting down gRPC server")
		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server stopped")
		return nil
	})

}
