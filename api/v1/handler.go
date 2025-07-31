package v1

import "github.com/MyProject273/Join_Love/internal/service"

type Handlers struct {
	AuthHandler AuthHandler
}

func NewHandlers(s service.Services) *Handlers {
	return &Handlers{
		AuthHandler: NewAuthHandler(s.AuthService),
	}
}
