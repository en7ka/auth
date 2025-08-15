package auth

import (
	repoif "github.com/en7ka/auth/internal/repository/repositoryinterface"
	"github.com/en7ka/auth/internal/service/servinterface"
)

type serv struct {
	userRepository repoif.UserRepository
}

func NewService(userRepository repoif.UserRepository) servinterface.UserService {
	return &serv{userRepository: userRepository}
}
