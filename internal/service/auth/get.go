package auth

import (
	"context"
	"errors"
	"log"

	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/repository/auth/converter"
	repoModel "github.com/en7ka/auth/internal/repository/auth/model"
)

var (
	ErrNotFound = errors.New("not found")
)

func (s *serv) Get(ctx context.Context, id int64) (*models.UserInfo, error) {
	cachedUser, err := s.userCache.Get(ctx, id)

	if err != nil && !errors.Is(err, ErrNotFound) {
		log.Printf("cache Get error (non-blocking): %v", err)
	}

	if err == nil && cachedUser != nil {
		return cachedUser, nil
	}

	var repoUser *repoModel.User

	err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		var txErr error
		repoUser, txErr = s.userRepository.Get(ctx, id)
		return txErr
	})

	if err != nil {
		return nil, err
	}

	serviceUser := converter.ToServiceUserInfo(repoUser)

	go func() {
		cacheCtx := context.Background()
		if cacheErr := s.userCache.Set(cacheCtx, id, serviceUser); cacheErr != nil {
			log.Printf("cache Set error (non-blocking): %v", cacheErr)
		}
	}()

	return serviceUser, nil
}
