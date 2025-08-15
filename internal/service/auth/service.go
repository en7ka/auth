package auth

import (
	repoif "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

type serv struct {
	userRepository repoif.UserRepository
}

func NewService(userRepository repoif.UserRepository) *serv {
	return &serv{userRepository: userRepository}
}
