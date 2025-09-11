package user

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) GetRefreshToken(ctx context.Context, req models.GetRefreshTokenRequest) (*models.GetRefreshTokenResponse, error) {
	claims, err := utils.VerifyToken(req.OldToken, []byte(req.OldToken))
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "Token is invalid")
	}

	refreshToken, err := utils.GenerateToken(models.UserInfoJwt{
		Username: claims.Username,
		Role:     claims.Role == "1" || claims.Role == "true",
	},
		[]byte(s.token.RefreshToken()),
		refreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &models.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
