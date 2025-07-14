package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/MyProject273/Join_Love/doc/statik"
	"github.com/MyProject273/Join_Love/gapi/helper"
	v1 "github.com/MyProject273/Join_Love/gapi/v1"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/i18n"
	logg "github.com/MyProject273/Join_Love/pkg/logger"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
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

	config, store, tokenMaker, logger := initializeApp(ctx)
	waitGroup, ctx := errgroup.WithContext(ctx)

	runGrpcServer(ctx, config, store, waitGroup, tokenMaker, logger)
	runGateWayServer(ctx, config, store, waitGroup, tokenMaker, logger)
	if err := waitGroup.Wait(); err != nil {
		log.Error().Err(err).Msg("application exited with error")
	} else {
		log.Info().Msg("application exited gracefully")
	}

}

func initializeApp(ctx context.Context) (cfg config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger) {
	var err error
	cfg, err = config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger = logg.NewLogger(cfg.Environment)

	connPool, err := pgxpool.New(ctx, cfg.DB_URL)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect db")
	}

	if err := connPool.Ping(ctx); err != nil {
		log.Fatal().Err(err).Msg("database not reachable")
	}

	tokenMaker, err = token.NewJWTMaker(cfg.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create token maker")
	}

	store = db.NewStore(connPool)

	if err := i18n.LoadI18nMessages("pkg/i18n"); err != nil {
		log.Fatal().Err(err).Msg("failed to load i18n")
	}

	return
}

func runGrpcServer(
	ctx context.Context,
	config config.Config,
	store db.Store,
	waitGroup *errgroup.Group,
	tokenMaker token.Maker,
	logger zerolog.Logger,
) {
	server, err := v1.NewServer(config, store, tokenMaker, logger)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create grpc server")
	}

	log.Info().Msgf("registered grpc services: %+v", server.GRPCServer.GetServiceInfo())
	reflection.Register(server.GRPCServer)

	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen on grpc address")
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("gRPC server running at %s", config.GrpcServerAddress)
		err := server.GRPCServer.Serve(listener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Error().Err(err).Msg("gRPC server failed")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Shutting down gRPC server")
		server.GRPCServer.GracefulStop()
		log.Info().Msg("gRPC server stopped")
		return nil
	})
}

func runGateWayServer(
	ctx context.Context,
	config config.Config,
	store db.Store,
	waitGroup *errgroup.Group,
	tokenMaker token.Maker,
	logger zerolog.Logger,
) {
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Custom header matcher allow accept-language
	headerMatcher := runtime.WithIncomingHeaderMatcher(helper.CustomMatcher)

	grpcMux := runtime.NewServeMux(jsonOption, headerMatcher)

	err := v1.RegisterAllHandlers(ctx, grpcMux, config, store, tokenMaker, logger)
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
	handler := c.Handler(helper.HttpLogger(mux))

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
