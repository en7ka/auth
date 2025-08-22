package auth

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/repository/auth/converter"
	repoModel "github.com/en7ka/auth/internal/repository/auth/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*models.UserInfo, error) {
	var repoUser *repoModel.User
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var txErr error
		repoUser, txErr = s.userRepository.Get(ctx, id)
		if txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	serviceUser := converter.ToServiceUserInfo(repoUser)
	return serviceUser, nil
}
