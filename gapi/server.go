package gapi

import (
	db "github.com/MyProject273/Join_Love/internal/db/sqlc"
	joinlove "github.com/MyProject273/Join_Love/pb"
	"github.com/MyProject273/Join_Love/pkg/config"
)

type Server struct {
	joinlove.UnimplementedJoinLoveServer
	config config.Config
	store  db.Store
}

func NewServer(config config.Config, store db.Store) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}
	return server, nil
}
