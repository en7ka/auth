package servinterface

import (
	"context"
	"github.com/en7ka/auth/internal/models"
)

type UserService interface {
	Create(ctx context.Context, info *models.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*models.User, error)
	Update(ctx context.Context, id int64, info *models.UserInfo) error
	Delete(ctx context.Context, id int64) error
}
