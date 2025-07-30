package service

import db "github.com/MyProject273/Join_Love/internal/db/sqlc"

type Services struct {
	AuthService AuthService
}

func NewServices(store db.Store) *Services {
	return &Services{
		AuthService: NewAuthService(store),
	}
}
