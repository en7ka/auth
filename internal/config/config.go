package config

import (
	"time"
)

func Load(path string) error {
	//err := godotenv.Load(path)
	//if err != nil {
	//	return err
	//}

	return nil
}

// GRPC конфиг
type GRPCConfig interface {
	Address() string
}

// PGConfig интерфейс получения DSN для старта хранилища.
type PGConfig interface {
	DSN() string
}

// RedisConfig интерфейс для получения данных для конфига Redis.
type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
}
