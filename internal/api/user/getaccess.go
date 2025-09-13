package user

import (
	"context"

	"github.com/en7ka/auth/internal/converter"
	authv1 "github.com/en7ka/auth/pkg/auth_v1"
)

func (c *Controller) GetAccessToken(ctx context.Context, req *authv1.GetAccessTokenRequest) (*authv1.GetAccessTokenResponse, error) {
	resp, err := c.authService.GetAccessToken(ctx, converter.ToGetAccessTokenFromAuthAPI(req))
	if err != nil {
		return nil, err
	}

	return &authv1.GetAccessTokenResponse{
		AccessToken: resp.AccessToken,
	}, nil
}
