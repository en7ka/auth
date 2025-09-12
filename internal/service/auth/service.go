package auth

import (
	"github.com/en7ka/auth/internal/client/db"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

type serv struct {
	userRepository repinf.UserRepository
	userCache      repinf.UserCache
	txManager      db.TxManager
}

func NewService(
	userRepository repinf.UserRepository,
	userCache repinf.UserCache,
	txManager db.TxManager,

) *serv {
	return &serv{
		userRepository: userRepository,
		userCache:      userCache,
		txManager:      txManager,
	}
}
