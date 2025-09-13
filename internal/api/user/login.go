package user

import (
	"context"

	"github.com/en7ka/auth/internal/converter"
	authv1 "github.com/en7ka/auth/pkg/auth_v1"
)

func (c *Controller) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	resp, err := c.authService.Login(ctx, converter.ToLoginFromAuthAPI(req))
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{
		RefreshToken: resp.Token,
	}, nil
}
