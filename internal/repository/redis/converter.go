package redis

import (
	"github.com/en7ka/auth/internal/models"
	redismodels "github.com/en7ka/auth/internal/repository/auth/model"
)

func toRedisModels(userInfo models.UserInfo) redismodels.User {

	return redismodels.User{
		Info: redismodels.UserInfo{
			Username: &userInfo.Username,
			Email:    &userInfo.Email,
			Password: &userInfo.Password,
		},
		Role: userInfo.Role,
	}
}

func toServiceModels(user redismodels.UserRedis) (*models.User, error) {

	var username, email, password string
	if user.Name == "" {
		username = user.Name
	}
	if user.Email == "" {
		email = user.Email
	}
	if user.Password == "" {
		password = user.Password
	}

	user1 := &models.UserInfo{
		Username: username,
		Email:    email,
		Password: password,
	}
	return &models.User{
		Id:        user.ID,
		Info:      *user1,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func toServiceModelsUserInfo(userProfile redismodels.UserRedis) (*models.UserInfo, error) {
	userInfo := &models.UserInfo{}

	userInfo.Role = userProfile.Role

	if userProfile.Name == "" {
		userInfo.Username = userProfile.Name
	}

	if userProfile.Email == "" {
		userInfo.Email = userProfile.Email
	}

	if userProfile.Password == "" {
		userInfo.Password = userProfile.Password
	}

	return userInfo, nil
}
