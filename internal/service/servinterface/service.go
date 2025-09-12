package servinterface

import (
	"context"

	"github.com/en7ka/auth/internal/models"
)

type UserService interface {
	Create(ctx context.Context, info *models.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*models.UserInfo, error)
	Update(ctx context.Context, id int64, info *models.UserInfo) error
	Delete(ctx context.Context, id int64) error
}

type AuthService interface {
	Login(ctx context.Context, user models.LoginRequest) (*models.LoginResponse, error)
	Check(ctx context.Context, request models.CheckRequest) error
	GetRefreshToken(ctx context.Context, request models.GetRefreshTokenRequest) (*models.GetRefreshTokenResponse, error)
	GetAccessToken(ctx context.Context, req models.GetAccessTokenRequest) (*models.GetAccessTokenResponse, error)
}
