package auth

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
)

func (s *serv) Update(ctx context.Context, id int64, info *models.UserInfo) error {
	return s.userRepository.Update(ctx, id, repoconv.ToRepoUserInfo(info))
}
