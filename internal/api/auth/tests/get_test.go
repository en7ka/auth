package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/en7ka/auth/internal/api/auth"
	"github.com/en7ka/auth/internal/models"
	serviceMocks "github.com/en7ka/auth/internal/service/mocks"
	"github.com/en7ka/auth/internal/service/servinterface"
	userv1 "github.com/en7ka/auth/pkg/user_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGet(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) servinterface.UserService

	type args struct {
		ctx context.Context
		req *userv1.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id        = gofakeit.Int64()
		name      = gofakeit.Name()
		email     = gofakeit.Email()
		password  = gofakeit.Password(true, true, true, false, false, 1)
		createdAt = gofakeit.Date()
		updatedAt = gofakeit.Date()

		errorGet = errors.New("service get error")

		req = &userv1.GetRequest{
			Id: id,
		}

		res = &userv1.GetResponse{
			Note: &userv1.User{
				Id: id,
				Info: &userv1.UserInfo{
					Username: name,
					Email:    email,
					Password: password,
					Role:     userv1.Role_user,
				},
				CreatedAt: timestamppb.New(createdAt),
				UpdatedAt: timestamppb.New(updatedAt),
			},
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *userv1.GetResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(&models.User{
					Id:        id,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
					Info: models.UserInfo{
						Username: name,
						Email:    email,
						Password: password,
						Role:     "user",
					},
				}, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  errorGet,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, errorGet)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			api := auth.NewImplementation(userServiceMock)

			response, err := api.Get(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, response)
		})
	}
}
