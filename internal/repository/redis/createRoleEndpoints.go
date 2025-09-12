package redis

import (
	"context"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// CreateRoleEndpoints сохраняет список эндпоинтов для указанной роли
func (c *cache) CreateRoleEndpoints(ctx context.Context, isAdmin bool, endpoints []string) error {
	if len(endpoints) == 0 {
		return nil
	}

	// Формируем ключ в формате "role:admin" или "role:user"
	roleKey := "role:user"
	if isAdmin {
		roleKey = "role:admin"
	}

	// Преобразуем слайс строк в массив интерфейсов для Redis
	args := make([]interface{}, 0, len(endpoints)+1)
	args = append(args, roleKey)
	for _, ep := range endpoints {
		args = append(args, ep)
	}

	if err := c.pool.Execute(ctx, func(_ context.Context, conn redis.Conn) error {
		if _, err := conn.Do("DEL", roleKey); err != nil {
			return fmt.Errorf("del old endpoints: %w", err)
		}
		if len(endpoints) > 0 {
			if _, err := conn.Do("SADD", args...); err != nil {
				return fmt.Errorf("sadd endpoints: %w", err)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("save endpoints (admin=%t): %w", isAdmin, err)
	}

	return nil
}

// GetRoleEndpoints получает список эндпоинтов для указанной роли
func (c *cache) GetRoleEndpoints(ctx context.Context, isAdmin bool) ([]string, error) {
	roleKey := "role:user"
	if isAdmin {
		roleKey = "role:admin"
	}

	var endpoints []string

	err := c.pool.Execute(ctx, func(_ context.Context, conn redis.Conn) error {
		values, err := redis.Strings(conn.Do("SMEMBERS", roleKey))
		if err != nil {
			if err == redis.ErrNil {
				endpoints = []string{} // Возвращаем пустой срез, если ключ не найден
				return nil
			}
			return err
		}

		endpoints = values
		return nil
	})

	return endpoints, err
}
