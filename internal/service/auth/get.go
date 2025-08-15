package auth

import (
	"context"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
)

func (s *serv) Get(ctx context.Context, id int64) (*models.User, error) {
	u, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return repoconv.ToModelUser(u), nil
}
