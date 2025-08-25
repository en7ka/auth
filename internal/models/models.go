package models

import (
	"time"
)

type User struct {
	Id        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserInfo struct {
	Username string
	Email    string
	Password string
	Role     string
}

type GetUserParams struct {
	ID       *int64
	Username *string
}
