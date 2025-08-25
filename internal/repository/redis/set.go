// Файл: internal/repository/redis/set.go

package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/en7ka/auth/internal/models"
)

const userExpirationSeconds = 900 // 15 minutes

// Set сохраняет пользователя в кэш.
func (c *cache) Set(ctx context.Context, id int64, user *models.UserInfo) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection from redis pool: %w", err)
	}
	defer conn.Close()

	userKey := fmt.Sprintf("user:%d", id)
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user for cache: %w", err)
	}

	_, err = conn.Do("SETEX", userKey, userExpirationSeconds, data)
	return err
}
