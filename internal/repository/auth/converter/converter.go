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

func ToRepoUserInfo(m *models.UserInfo) *repom.UserInfo {
	if m == nil {
		return nil
	}

	return &repom.UserInfo{
		Username: &m.Username,
		Email:    &m.Email,
		Password: &m.Password,
	}
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
