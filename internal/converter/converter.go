package converter

import (
	"github.com/en7ka/auth/internal/models"
	repoConverter "github.com/en7ka/auth/internal/repository/auth/converter"
	authv1 "github.com/en7ka/auth/pkg/auth_v1"
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

func ToLoginFromAuthAPI(req *authv1.LoginRequest) models.LoginRequest {
	if req == nil {
		return models.LoginRequest{}
	}
	return models.LoginRequest{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	}
}

func ToGetRefreshTokenFromDesc(req *authv1.GetRefreshTokenRequest) models.GetRefreshTokenRequest {
	if req == nil {
		return models.GetRefreshTokenRequest{}
	}
	return models.GetRefreshTokenRequest{
		OldToken: req.GetOldRefreshToken(),
	}
}

func ToGetAccessTokenFromAuthAPI(req *authv1.GetAccessTokenRequest) models.GetAccessTokenRequest {
	if req == nil {
		return models.GetAccessTokenRequest{}
	}

	return models.GetAccessTokenRequest{
		RefreshToken: req.GetRefreshToken(),
	}
}

func ToCheckAccessFromAuthAPI(req *authv1.CheckRequest) models.CheckRequest {
	if req == nil {
		return models.CheckRequest{}
	}
	return models.CheckRequest{
		EndpointAddress: req.EndpointAddress,
	}
}
