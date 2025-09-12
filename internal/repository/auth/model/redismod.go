package model

import "time"

type (
	// UserRedis модель пользователя для redis.
	UserRedis struct {
		ID        int64     `redis:"id"`
		Name      string    `redis:"name"`
		Email     string    `redis:"email"`
		Role      string    `redis:"role"`
		Password  string    `redis:"password"`
		CreatedAt time.Time `redis:"created_at"`
		UpdatedAt time.Time `redis:"updated_at"`
	}
)
