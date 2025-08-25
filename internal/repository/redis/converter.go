package redis

import (
	"github.com/en7ka/auth/internal/models"
	redismodels "github.com/en7ka/auth/internal/repository/auth/model"
)

func toRedisModels(user models.User) redismodels.User {

	return redismodels.User{
		Id: user.Id,
		Info: redismodels.UserInfo{
			Username: &user.Info.Username,
			Email:    &user.Info.Email,
			Password: &user.Info.Password,
		},
		Role:      user.Info.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func toServiceModels(user redismodels.User) (*models.User, error) {

	var username, email, password string
	if user.Info.Username != nil {
		username = *user.Info.Username
	}
	if user.Info.Email != nil {
		email = *user.Info.Email
	}
	if user.Info.Password != nil {
		password = *user.Info.Password
	}

	return &models.User{
		Id:        user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Info: models.UserInfo{
			Username: username,
			Email:    email,
			Password: password,
			Role:     user.Role,
		},
	}, nil
}
