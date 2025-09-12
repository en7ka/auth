package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/en7ka/auth/internal/models"
)

// Set сохраняет пользователя в кэш.
func (c *cache) Set(ctx context.Context, id int64, user *models.UserInfo) error {
	idFormatted := strconv.FormatInt(id, 10)
	rm := toRedisModels(*user)
	if err := c.pool.HashSet(ctx, idFormatted, rm); err != nil {
		return fmt.Errorf("failed to hash user: %w", err)
	}
	if err := c.pool.Expire(ctx, idFormatted, 5*time.Minute); err != nil {
		return fmt.Errorf("failed to set expiration for user: %w", err)
	}
	return nil
}
