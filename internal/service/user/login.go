package user

import (
	"context"
	"errors"

	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/utils"
)

func (s *serv) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	var resp models.LoginResponse
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		user, txErr := s.authRepository.Login(ctx, req)
		if txErr != nil {
			return txErr
		}

		refreshToken, err := utils.GenerateToken(*user,
			[]byte(s.token.RefreshToken()),
			refreshTokenExpiration,
		)

		if err != nil {
			return errors.New("failed to generate refresh token")
		}

		resp.Token = refreshToken

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
