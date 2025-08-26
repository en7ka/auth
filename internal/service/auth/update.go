package auth

import (
	"context"
	"errors"
	"log"

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

	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		txErr := s.userRepository.Update(ctx, id, repoconv.ToRepoUserInfo(info))
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
