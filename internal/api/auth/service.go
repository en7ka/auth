package auth

import (
	usserv "github.com/en7ka/auth/internal/service/servinterface"
	desc "github.com/en7ka/auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserAPIServer
	userService usserv.UserService
}

func NewImplementation(userService usserv.UserService) *Implementation {

	return &Implementation{userService: userService}
}
