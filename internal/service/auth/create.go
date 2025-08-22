package auth

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
)

func (s *serv) Create(ctx context.Context, info *models.UserInfo) (int64, error) {
	if info == nil {
		return 0, nil
	}

	var userID int64
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var txErr error

		userID, txErr = s.userRepository.Create(ctx, repoconv.ToRepoUserInfo(info))
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return userID, nil
}
