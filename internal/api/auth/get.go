package auth

import (
	"context"
	"log"

	"github.com/en7ka/auth/internal/converter"
	"github.com/en7ka/auth/internal/logger"
	repoConverter "github.com/en7ka/auth/internal/repository/auth/converter"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"go.uber.org/zap"
)

func (c *Controller) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	logger.Info("Getting note...", zap.Int64("id", req.GetId()))

	userInfo, err := c.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("noteObj: %v", userInfo)

	return &desc.GetResponse{
		Note: &desc.User{
			Id: req.GetId(),

			Info: converter.ToUserInfoFromService(userInfo),

			Role: repoConverter.RoleFromString(userInfo.Role),
		},
	}, nil
}
