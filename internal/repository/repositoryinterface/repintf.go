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
}
