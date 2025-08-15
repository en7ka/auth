package auth

import (
	"context"
	"github.com/en7ka/auth/internal/models"
	desc "github.com/en7ka/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	id := req.GetId()

	info := &models.UserInfo{}

	if v := req.GetInfo().GetUsername(); v != nil {
		info.Username = v.GetValue()
	}

	if v := req.GetInfo().GetEmail(); v != nil {
		info.Email = v.GetValue()
	}

	if err := i.userService.Update(ctx, id, info); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
