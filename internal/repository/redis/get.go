package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/en7ka/auth/internal/models"
	redismodels "github.com/en7ka/auth/internal/repository/auth/model"
	"github.com/gomodule/redigo/redis"
)

var (
	ErrNotFound = errors.New("not found")
)

// Get получает пользователя из кэша.
func (c *cache) Get(ctx context.Context, id int64) (*models.UserInfo, error) {
	key := strconv.FormatInt(id, 10)

	userCache, err := c.pool.HGetAll(ctx, key)
	if err != nil {
		if err == redis.ErrNil {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("redis HGetAll error: %w", err)
	}

	if len(userCache) == 0 {
		return nil, ErrNotFound
	}

	var userProfile redismodels.UserRedis
	if err := redis.ScanStruct(userCache, &userProfile); err != nil {
		return nil, fmt.Errorf("error scanning user profile: %w", err)
	}

	user, err := toServiceModelsUserInfo(userProfile)
	if err != nil {
		return nil, fmt.Errorf("error converting user profile: %w", err)
	}

	return user, nil
}
