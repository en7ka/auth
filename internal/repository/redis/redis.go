package redis

import (
	cl "github.com/en7ka/auth/internal/client/cache/redis"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

type cache struct {
	pool *cl.Client
}

func NewRedisCache(pool *cl.Client) repinf.UserCache {
	return &cache{
		pool: pool,
	}
}
