package auth

import (
	"context"
	"github.com/en7ka/auth/internal/converter"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"log"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	noteObj, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Printf("noteObj: %v", noteObj)

	return &desc.GetResponse{
		Note: converter.ToUserFromService(noteObj),
	}, nil
}
