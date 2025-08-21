package model

import (
	"time"
)

type User struct {
	Id        int64     `db:"id"`
	Info      UserInfo  `db:""`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserInfo struct {
	Username *string `db:"username"`
	Email    *string `db:"email"`
	Password *string `db:"password"`
	Role     string  `db:"role"`
}
