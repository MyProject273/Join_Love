package api

import (
	"fmt"

	v1 "github.com/MyProject273/Join_Love/api/v1"
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	routes "github.com/MyProject273/Join_Love/internal/route"
	"github.com/MyProject273/Join_Love/internal/service"
	"github.com/MyProject273/Join_Love/internal/worker"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	config          config.Config
	store           db.Store
	tokenMaker      token.Maker
	route           *gin.Engine
	logger          zerolog.Logger
	services        service.Services
	handlers        v1.Handlers
	taskDistributor worker.TaskDistributor
}

func NewServer(config config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	services := service.NewServices(store, config, tokenMaker, taskDistributor, logger)
	handlers := v1.NewHandlers(*services)
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		logger:          logger,
		services:        *services,
		handlers:        *handlers,
		taskDistributor: taskDistributor,
	}
	server.setupRoute()
	return server, nil
}

func (s *Server) setupRoute() {
	router := gin.Default()
	api := router.Group("/api/v1")
	// Public routes
	s.registerPublicRoutes(api)

	// Protected routes
	s.route = router
}

func (s *Server) registerPublicRoutes(r *gin.RouterGroup) {
	routes.RegisterAuthRoutes(r, s.handlers.AuthHandler)
}

func (s *Server) Router() *gin.Engine {
	return s.route
}
