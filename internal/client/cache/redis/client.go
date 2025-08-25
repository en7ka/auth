package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/en7ka/auth/internal/config"
	"github.com/gomodule/redigo/redis"
)

type handler func(ctx context.Context, conn redis.Conn) error

// Client for redis
type Client struct {
	pool   *redis.Pool
	config config.RedisConfig
}

// NewClient constructor for redis`s client
func NewClient(pool *redis.Pool, config config.RedisConfig) Client {
	return Client{
		pool:   pool,
		config: config,
	}
}

func (c *Client) execute(ctx context.Context, fn handler) error {
	conn, err := c.connect(ctx)
	if err != nil {
		return fmt.Errorf("could not connect to redis: %v", err)
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Printf("could not close redis connection: %v", err)
		}
	}()

	if err = fn(ctx, conn); err != nil {
		return fmt.Errorf("could not handle request: %w", err)
	}

	return nil
}

func (c *Client) connect(ctx context.Context) (redis.Conn, error) {
	getConnTimeoutCtx, cancel := context.WithTimeout(ctx, c.config.ConnectionTimeout())
	defer cancel()

	conn, err := c.pool.GetContext(getConnTimeoutCtx)
	if err != nil {
		log.Printf("failed to connect to redis: %v", err)

		_ = conn.Close()
		return nil, err
	}

	return conn, nil
}

// HashSet метод для сохранения структуры
func (c *Client) HashSet(ctx context.Context, key string, values interface{}) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("HSET", redis.Args{key}.AddFlat(values)...)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not hash set: %w", err)
	}

	return nil
}

// Set для записи по ключу
func (c *Client) Set(ctx context.Context, key string, values interface{}) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("SET", redis.Args{key}.Add(values)...)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not set: %w", err)
	}

	return nil
}

// Get для получения по ключу
func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	var value interface{}
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		var errEx error
		value, errEx = conn.Do("GET", key)
		if errEx != nil {
			return errEx
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not get: %w", err)
	}

	return value, nil
}

// HGetAll для получения всех значений хеш-таблицы
func (c *Client) HGetAll(ctx context.Context, key string) ([]interface{}, error) {
	var values []interface{}
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		var errEx error
		values, errEx = redis.Values(conn.Do("HGETALL", key))
		if errEx != nil {
			return errEx
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("could not get hgetall: %w", err)
	}

	return values, nil
}

// Expire установка TTL
func (c *Client) Expire(ctx context.Context, key string, duration time.Duration) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		seconds := int64(duration.Seconds())
		_, err := conn.Do("EXPIRE", key, seconds)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not expire: %w", err)
	}

	return nil
}

// Ping пингуем редис
func (c *Client) Ping(ctx context.Context) error {
	err := c.execute(ctx, func(_ context.Context, conn redis.Conn) error {
		_, err := conn.Do("PING")
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("could not ping: %w", err)
	}

	return nil
}
