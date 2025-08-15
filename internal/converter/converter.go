package converter

import (
	"github.com/en7ka/auth/internal/models"
	repoConverter "github.com/en7ka/auth/internal/repository/auth/converter"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserFromService конвертирует внутреннюю модель пользователя в модель для gRPC ответа.
func ToUserFromService(user *models.User) *desc.Note {
	if user == nil {
		return nil
	}

	return &desc.Note{
		Id:        user.Id,
		Info:      ToUserInfoFromService(&user.Info),
		Role:      repoConverter.RoleFromString(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToUserInfoFromService(info *models.UserInfo) *desc.NoteInfo {
	if info == nil {
		return nil
	}

	return &desc.NoteInfo{
		Username: info.Username,
		Email:    info.Email,
		Password: info.Password,
	}
}

func ToServiceModelFromDesc(userInfo *desc.NoteInfo) *models.User {
	if userInfo == nil {
		return nil
	}
	return &models.User{
		Info: models.UserInfo{
			Username: userInfo.Username,
			Email:    userInfo.Email,
			Password: userInfo.Password,
		},
	}
}
