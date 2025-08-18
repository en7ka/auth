package auth

import (
	usserv "github.com/en7ka/auth/internal/service/servinterface"
	desc "github.com/en7ka/auth/pkg/user_v1"
)

type Controller struct {
	desc.UnimplementedUserAPIServer
	userService usserv.UserService
}

func NewImplementation(userService usserv.UserService) *Controller {
	return &Controller{userService: userService}
}
