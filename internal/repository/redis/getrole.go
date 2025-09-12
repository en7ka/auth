package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (c *cache) GetRole(ctx context.Context, username string) (*bool, error) {
	var role bool
	err := c.pool.Execute(ctx, func(ctx context.Context, conn redis.Conn) error {
		res, err := redis.Bool(conn.Do("GET", username))
		if err != nil {
			if err == redis.ErrNil {
				return ErrNotFound
			}
			return fmt.Errorf("redis Get error: %w", err)
		}
		role = res
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &role, nil
}
