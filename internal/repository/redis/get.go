package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/en7ka/auth/internal/models"
	redisModels "github.com/en7ka/auth/internal/repository/auth/model"
	"github.com/gomodule/redigo/redis"
)

func (c cache) Get(ctx context.Context, params models.GetUserParams) (*models.User, error) {
	var key string
	switch {
	case params.ID != nil:
		key = strconv.FormatInt(*params.ID, 10)
	case params.Username != nil:
		key = *params.Username
	}

	userCache, err := c.cl.HGetAll(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("error with get user cache: %w", err)
	}

	if userCache == nil {
		return nil, errors.New("user not found")
	}

	var userProfile redisModels.User
	err = redis.ScanStruct(userCache, &userProfile)
	if err != nil {
		return nil, fmt.Errorf("error scanning user profile: %w", err)
	}

	user, err := toServiceModels(userProfile)
	if err != nil {
		return nil, fmt.Errorf("error converting user profile: %w", err)
	}

	return user, nil
}
