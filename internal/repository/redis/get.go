// Файл: internal/repository/redis/get.go

package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/en7ka/auth/internal/models"
	redigo "github.com/gomodule/redigo/redis"
)

var (
	ErrNotFound = errors.New("not found")
)

// Get получает пользователя из кэша.
func (c *cache) Get(ctx context.Context, id int64) (*models.UserInfo, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection from redis pool: %w", err)
	}
	defer conn.Close()

	userKey := fmt.Sprintf("user:%d", id)

	data, err := redigo.Bytes(conn.Do("GET", userKey))
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var userInfo models.UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user data from cache: %w", err)
	}

	return &userInfo, nil
}
