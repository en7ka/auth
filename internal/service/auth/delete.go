package auth

import (
	"context"
	"errors"
	"log"
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

	if err != nil {
		return err
	}

	go func() {
		cacheCtx := context.Background()
		if cacheErr := s.userCache.Delete(cacheCtx, id); cacheErr != nil {
			log.Printf("cache Delete error (non-blocking): %v", cacheErr)
		}
	}()

	return nil
}
