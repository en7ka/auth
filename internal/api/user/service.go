package user

import (
	usserv "github.com/en7ka/auth/internal/service/servinterface"
	desc "github.com/en7ka/auth/pkg/auth_v1"
)

type Controller struct {
	authService usserv.AuthService
	desc.UnimplementedAuthApiServer
}

func NewImplementation(authService usserv.AuthService) *Controller {
	return &Controller{authService: authService}
}
