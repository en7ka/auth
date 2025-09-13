package auth

import (
	"github.com/en7ka/auth/internal/client/db"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
	"github.com/segmentio/kafka-go"
)

type serv struct {
	userRepository repoinf.UserRepository
	userCache      repoinf.UserCache
	txManager      db.TxManager
	producer       *kafka.Writer
}

func NewService(
	userRepository repoinf.UserRepository,
	userCache repoinf.UserCache,
	txManager db.TxManager,
	producer *kafka.Writer,
) *serv {
	return &serv{
		userRepository: userRepository,
		userCache:      userCache,
		txManager:      txManager,
		producer:       producer,
	}
}