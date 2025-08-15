package auth

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
)

func (s *serv) Create(ctx context.Context, info *models.UserInfo) (int64, error) {
	return s.userRepository.Create(ctx, repoconv.ToRepoUserInfo(info))
}
