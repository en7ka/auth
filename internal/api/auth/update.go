package auth

import (
	"context"

	"github.com/en7ka/auth/internal/converter"
	"github.com/en7ka/auth/internal/logger"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info("Starting user update process")

	id := req.GetId()

	info := converter.UpdateRequestToUserInfo(req)

	if err := c.userService.Update(ctx, id, info); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
