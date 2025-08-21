package converter

import (
	"github.com/en7ka/auth/internal/models"
	repoConverter "github.com/en7ka/auth/internal/repository/auth/converter"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToUserFromService конвертирует внутреннюю модель пользователя в модель для gRPC ответа.
func ToUserFromService(user *models.User) *desc.User {
	if user == nil {
		return nil
	}

	return &desc.User{
		Id:        user.Id,
		Info:      ToUserInfoFromService(&user.Info),
		Role:      repoConverter.RoleFromString(user.Info.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToUserInfoFromService(info *models.UserInfo) *desc.UserInfo {
	if info == nil {
		return nil
	}

	return &desc.UserInfo{
		Username: info.Username,
		Email:    info.Email,
		Password: info.Password,
	}
}

func ToServiceModelFromDesc(userInfo *desc.UserInfo) *models.User {
	if userInfo == nil {
		return nil
	}
	return &models.User{
		Info: models.UserInfo{
			Username: userInfo.Username,
			Email:    userInfo.Email,
			Password: userInfo.Password,
			Role:     userInfo.Role.String(),
		},
	}
}

func UpdateRequestToUserInfo(req *desc.UpdateRequest) *models.UserInfo {
	info := &models.UserInfo{}

	if v := req.GetInfo().GetUsername(); v != nil {
		info.Username = v.GetValue()
	}

	if v := req.GetInfo().GetEmail(); v != nil {
		info.Email = v.GetValue()
	}

	return info
}
