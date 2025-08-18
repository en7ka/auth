package auth

import (
	"context"
	"log"

	"github.com/en7ka/auth/internal/converter"
	desc "github.com/en7ka/auth/pkg/user_v1"
)

func (c *Controller) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := c.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("noteObj: %v", user)

	return &desc.GetResponse{
		Note: converter.ToUserFromService(user),
	}, nil
}
