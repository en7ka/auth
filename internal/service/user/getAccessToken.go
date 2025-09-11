package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/en7ka/auth/internal/models"
	"github.com/en7ka/auth/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

func (s *serv) GetAccessToken(ctx context.Context, req models.GetAccessTokenRequest) (*models.GetAccessTokenResponse, error) {
	claims, err := utils.VerifyToken(req.RefreshToken, []byte(s.token.RefreshToken()))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	var userRole bool
	rolePtr, errCache := s.authCache.GetRole(ctx, claims.Username)

	if errCache != nil {
		if errors.Is(errCache, ErrUserNotFound) {
			err = s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
				var txErr error

				userRole, txErr = s.authRepository.GetUserRole(ctx, claims.Username)
				if txErr != nil {
					return txErr
				}

				if txErr = s.authCache.CreateRole(ctx, claims.Username, userRole); txErr != nil {
					return fmt.Errorf("create role failed: %v", txErr)
				}

				return nil
			})

			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("get role failed: %v", err)
		}
	} else if rolePtr == nil {
		userRole = *rolePtr
	} else {
		return nil, fmt.Errorf("get role failed: %v", rolePtr)
	}

	accessToken, err := utils.GenerateToken(models.UserInfoJwt{
		Username: claims.Username,
		Role:     userRole,
	},
		[]byte(s.token.AccessToken()),
		accessTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return &models.GetAccessTokenResponse{
		AccessToken: accessToken,
	}, nil
}
