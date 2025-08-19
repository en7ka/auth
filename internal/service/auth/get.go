package auth

import (
	"context"
	"fmt"

	"github.com/en7ka/auth/internal/models"
	repoconv "github.com/en7ka/auth/internal/repository/auth/converter"
)

func (s *serv) Get(ctx context.Context, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid id %d", id)
	}

	user, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return repoconv.ToModelUser(user), nil
}
