package auth

import (
	"context"

	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (c *Controller) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := c.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
