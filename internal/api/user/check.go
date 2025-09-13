package user

import (
	"context"

	"github.com/en7ka/auth/internal/converter"
	authv1 "github.com/en7ka/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Check(ctx context.Context, req *authv1.CheckRequest) (*emptypb.Empty, error) {
	if err := c.authService.Check(ctx, converter.ToCheckAccessFromAuthAPI(req)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
