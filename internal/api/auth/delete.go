package auth

import (
	"context"

	"github.com/en7ka/auth/internal/logger"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	logger.Info("Starting user deletion process")

	if err := c.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
