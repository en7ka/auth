package auth

import (
	"context"
	"errors"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var txErr error
		txErr = s.userRepository.Delete(ctx, id)
		if txErr != nil {
			return txErr
		}

		return nil
	})

	return err
}
