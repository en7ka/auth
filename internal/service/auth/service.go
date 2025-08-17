package auth

import (
	"github.com/en7ka/auth/internal/client/db"
	repoif "github.com/en7ka/auth/internal/repository/repositoryinterface"
	userService "github.com/en7ka/auth/internal/service/servinterface"
)

type serv struct {
	userRepository repoif.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repoif.UserRepository, txManager db.TxManager) userService.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
