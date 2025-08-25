// Файл: internal/repository/redis/delete.go

package redis

import (
	"context"
	"fmt"
)

// Delete удаляет пользователя из кэша (инвалидация).
func (c *cache) Delete(ctx context.Context, id int64) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection from redis pool: %w", err)
	}
	defer conn.Close()

	userKey := fmt.Sprintf("user:%d", id)
	_, err = conn.Do("DEL", userKey)

	return err
}
