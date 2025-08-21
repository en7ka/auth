package tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/en7ka/auth/internal/api/auth"
	"github.com/en7ka/auth/internal/models"
	serviceMocks "github.com/en7ka/auth/internal/service/mocks"
	"github.com/en7ka/auth/internal/service/servinterface"
	userv1 "github.com/en7ka/auth/pkg/user_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	type userServiceMockFunc func(mc *minimock.Controller) servinterface.UserService

	type args struct {
		ctx context.Context
		req *userv1.CreateRequest
	}

	var (
		id       = gofakeit.Int64()
		ctx      = context.Background()
		mc       = minimock.NewController(t)
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, false, false, 1)

		errorCreate = errors.New("error with service")
		res         = &userv1.CreateResponse{
			Id: id,
		}
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *userv1.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case: check convert role user",
			args: args{
				ctx: ctx,
				req: &userv1.CreateRequest{
					Info: &userv1.UserInfo{
						Username: name,
						Email:    email,
						Password: password,
						Role:     userv1.Role_user,
					},
				},
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, &models.UserInfo{
					Username: name,
					Email:    email,
					Password: password,
					Role:     "user",
				}).Return(id, nil)
				return mock
			},
		},
		{
			name: "success case: check convert role admin",
			args: args{
				ctx: ctx,
				req: &userv1.CreateRequest{
					Info: &userv1.UserInfo{
						Username: name,
						Email:    email,
						Password: password,
						Role:     userv1.Role_admin,
					},
				},
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, &models.UserInfo{
					Username: name,
					Email:    email,
					Password: password,
					Role:     "admin",
				}).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case:",
			args: args{
				ctx: ctx,
				req: &userv1.CreateRequest{
					Info: &userv1.UserInfo{
						Username: name,
						Email:    email,
						Password: password,
						Role:     userv1.Role_user,
					},
				},
			},
			want: nil,
			err:  fmt.Errorf("error while creating: %w", errorCreate),
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, &models.UserInfo{
					Username: name,
					Email:    email,
					Password: password,
					Role:     "user",
				}).Return(0, errorCreate)
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

			newID, err := api.Create(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, newID)
		})
	}
}
