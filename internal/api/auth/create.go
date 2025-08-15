package auth

import (
	"context"
	"github.com/en7ka/auth/internal/converter"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"log"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	beach := converter.ToServiceModelFromDesc(req.GetInfo())

	id, err := i.userService.Create(ctx, &beach.Info)
	if err != nil {
		return nil, err
	}
	log.Printf("inserted user id %s", id)

	return &desc.CreateResponse{
		Id: id,
	}, nil
}
