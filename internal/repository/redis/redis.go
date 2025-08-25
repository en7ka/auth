package redis

import (
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
	redigo "github.com/gomodule/redigo/redis"
)

type cache struct {
	pool *redigo.Pool
}

func NewRedisCache(pool *redigo.Pool) repinf.UserCache {
	return &cache{
		pool: pool,
	}
}
