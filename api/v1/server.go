package v1

import (
	"fmt"

	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	"github.com/MyProject273/Join_Love/pkg/config"
	"github.com/MyProject273/Join_Love/pkg/utils/token"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Server struct {
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
	route      *gin.Engine
	logger     zerolog.Logger
}

func NewServer(config config.Config, store db.Store, tokenMaker token.Maker, logger zerolog.Logger) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		logger:     logger,
	}

	return server, nil
}

func (s *Server) Router() *gin.Engine {
	return s.route
}
