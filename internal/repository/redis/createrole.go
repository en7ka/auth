package redis

import (
	"context"
	"fmt"
	"time"
)

func (c *cache) CreateRole(ctx context.Context, username string, role bool) error {
	if err := c.pool.Set(ctx, username, role); err != nil {
		return fmt.Errorf("failed to set role for %s: %w", username, err)
	}
	if err := c.pool.Expire(ctx, username, time.Hour); err != nil {
		return fmt.Errorf("failed to expire role for %s: %w", username, err)
	}
	return nil
}
