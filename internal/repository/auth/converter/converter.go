package converter

import (
	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/repository/auth/model"
	repom "github.com/en7ka/auth/internal/repository/auth/model"
	userv1 "github.com/en7ka/auth/pkg/user_v1"
)

func ToUserFromRepo(user *model.User) *model.User {
	return &model.User{
		Id:        user.Id,
		Info:      ToUserInfoFromRepo(user.Info),
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserInfoFromRepo(user model.UserInfo) model.UserInfo {
	return model.UserInfo{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
}

func RoleFromString(s string) userv1.Role {
	switch s {
	case "admin":
		return userv1.Role_admin
	default:
		return userv1.Role_user
	}
}

func RoleToString(r userv1.Role) string {
	switch r {
	case userv1.Role_admin:
		return "admin"
	default:
		return "user"
	}
}

func ToRepoUserInfo(info *models.UserInfo) *repom.UserInfo {
	repoInfo := &repom.UserInfo{
		Role: info.Role,
	}

	if info.Username != "" {
		repoInfo.Username = &info.Username
	}
	if info.Email != "" {
		repoInfo.Email = &info.Email
	}
	if info.Password != "" {
		repoInfo.Password = &info.Password
	}

	return repoInfo
}

func ToModelUser(r *repom.User) *models.User {
	if r == nil {
		return nil
	}

	var username, email, password string
	if r.Info.Username != nil {
		username = *r.Info.Username
	}
	if r.Info.Email != nil {
		email = *r.Info.Email
	}
	if r.Info.Password != nil {
		password = *r.Info.Password
	}

	return &models.User{
		Id: r.Id,
		Info: models.UserInfo{
			Username: username,
			Email:    email,
			Password: password,
			Role:     r.Role,
		},
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

func ToServiceUserInfo(user *repom.User) *models.UserInfo {
	if user == nil {
		return nil
	}

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

	return &models.UserInfo{
		Username: username,
		Email:    email,
		Password: password,
		Role:     user.Role,
	}
}
