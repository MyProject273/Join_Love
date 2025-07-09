package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/MyProject273/Join_Love/doc/statik"
	"github.com/MyProject273/Join_Love/gapi"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	pb "github.com/MyProject273/Join_Love/pb"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	runGateWayServer(ctx, config, store, waitGroup)
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

func runGateWayServer(
	ctx context.Context,
	config config.Config,
	store db.Store,
	waitGroup *errgroup.Group,
) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc gateway server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	err = pb.RegisterJoinLoveHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: config.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodOptions,
			http.MethodGet,
			http.MethodPost,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
	})
	handler := c.Handler(gapi.HttpLogger(mux))

	httpServer := &http.Server{
		Handler: handler,
		Addr:    config.HttpServerAddress,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP gateway server at %s", httpServer.Addr)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("http gateway server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP gateway server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP gateway server")
		}

		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})
}
