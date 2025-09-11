package user

import (
	"time"

	"github.com/en7ka/auth/internal/client/db"
	"github.com/en7ka/auth/internal/config"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

const (
	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 5 * time.Minute
)

type serv struct {
	authRepository repinf.AuthRepository
	txManager      db.TxManager
	authCache      repinf.UserCache
	token          config.JWTConfig
	access         config.AccessConfig
}

func NewService(
	authRepository repinf.AuthRepository,
	authCache repinf.UserCache,
	txManager db.TxManager,
	token config.JWTConfig,
	access config.AccessConfig,
) *serv {
	return &serv{
		authRepository: authRepository,
		authCache:      authCache,
		txManager:      txManager,
		token:          token,
		access:         access,
	}
}
