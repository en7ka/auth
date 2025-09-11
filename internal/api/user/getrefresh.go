package user

import (
	"context"

	"github.com/en7ka/auth/internal/converter"
	authv1 "github.com/en7ka/auth/pkg/auth_v1"
)

func (c *Controller) GetRefreshToken(ctx context.Context, req *authv1.GetRefreshTokenRequest) (*authv1.GetRefreshTokenResponse, error) {
	resp, err := c.authService.GetRefreshToken(ctx, *converter.ToGetRefreshTokenFromDesc(req))
	if err != nil {
		return nil, err
	}

	return &authv1.GetRefreshTokenResponse{
		RefreshToken: resp.RefreshToken,
	}, nil
}
