package internal

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/repository/auth/model"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, id int64, info *model.UserInfo) error
	Delete(ctx context.Context, id int64) error
}

type UserCache interface {
	Set(ctx context.Context, id int64, user *models.UserInfo) error
	Delete(ctx context.Context, id int64) error
	Get(ctx context.Context, id int64) (*models.UserInfo, error)
	GetRole(ctx context.Context, username string) (*bool, error)
	CreateRole(ctx context.Context, username string, role bool) error
	CreateRoleEndpoints(ctx context.Context, isAdmin bool, endpoints []string) error
	GetRoleEndpoints(ctx context.Context, isAdmin bool) ([]string, error)
}

type AuthRepository interface {
	Login(ctx context.Context, user models.LoginRequest) (*models.UserInfoJwt, error)
	GetUserRole(ctx context.Context, username string) (bool, error)
	GetUserAccess(ctx context.Context, isAdmin bool) ([]string, error)
}
