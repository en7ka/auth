package auth

import (
	"context"
	"errors"
	"fmt"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	if err := s.userRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	return nil
}
