package redis

import (
	"github.com/en7ka/auth/internal/client/cache/redis"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

type cache struct {
	cl redis.Client
}

func NewRedisCache(cl redis.Client) repinf.UserCache {
	return &cache{
		cl: cl,
	}
}
