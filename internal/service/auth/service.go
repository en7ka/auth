package auth

import (
	"github.com/en7ka/auth/internal/client/db"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
	"github.com/segmentio/kafka-go"
)

type serv struct {
	userRepository repinf.UserRepository
	userCache      repinf.UserCache
	txManager      db.TxManager

	producer *kafka.Writer
}

func NewService(
	userRepository repinf.UserRepository,
	userCache repinf.UserCache,
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
