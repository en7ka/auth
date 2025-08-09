package cmd

import (
	"time"

	auth "github.com/en7ka/auth/pkg/user_v1"
)

type (
	User struct {
		ID        int64
		Username  string
		Email     string
		Password  string // Добавлено
		Role      auth.Role
		CreatedAt time.Time
		UpdatedAt time.Time // Добавлено
	}

	UpdateUser struct {
		ID       int64
		Username string // Исправлено с Name на Username
		Email    string
	}

	DeleteID int64

	GetUserPar struct {
		ID       *int64
		Username *string
	}
)
