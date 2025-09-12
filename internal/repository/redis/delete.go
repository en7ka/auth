package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (c *cache) Delete(ctx context.Context, id int64) error {
	key := fmt.Sprintf("user:%d", id)
	return c.pool.Execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		_, err := conn.Do("DEL", key)
		if err != nil {
			return fmt.Errorf("failed to delete user from cache: %w", err)
		}
		return nil
	})
}
