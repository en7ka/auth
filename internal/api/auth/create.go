package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/en7ka/auth/internal/converter"
	"github.com/en7ka/auth/internal/logger"
	desc "github.com/en7ka/auth/pkg/user_v1"
)

func (c *Controller) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	logger.Info("Starting user creation process")

	user := converter.ToServiceModelFromDesc(req.GetInfo())

	id, err := c.userService.Create(ctx, &user.Info)
	if err != nil {
		return nil, fmt.Errorf("error while creating: %w", err)
	}
	log.Printf("inserted user id %v", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
