package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
)

func (s *serv) Update(ctx context.Context, id int64, info *models.UserInfo) error {
	if id <= 0 {
		return errors.New("user ID must be positive")
	}
	if info == nil {
		return errors.New("user info cannot be nil")
	}

	if err := s.userRepository.Update(ctx, id, repoconv.ToRepoUserInfo(info)); err != nil {
		return fmt.Errorf("failed to update user in repository: %w", err)
	}

	return nil
}
