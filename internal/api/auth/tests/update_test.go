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
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) servinterface.UserService

	type args struct {
		ctx context.Context
		req *userv1.UpdateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		username = gofakeit.Name()
		email    = gofakeit.Email()

		errorUpdate = errors.New("error with service")
	)
	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name            string
		args            args
		want            *emptypb.Empty
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: &userv1.UpdateRequest{
					Id: id,
					Info: &userv1.UpdateUserInfo{
						Username: &wrapperspb.StringValue{Value: username},
						Email:    &wrapperspb.StringValue{Value: email},
					},
				},
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, &models.UserInfo{
					Username: username,
					Email:    email,
				}).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: &userv1.UpdateRequest{
					Id: id,
					Info: &userv1.UpdateUserInfo{
						Username: &wrapperspb.StringValue{Value: username},
						Email:    &wrapperspb.StringValue{Value: email},
					},
				},
			},
			want: nil,
			err:  errorUpdate,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.UpdateMock.Expect(ctx, id, &models.UserInfo{
					Username: username,
					Email:    email,
				}).Return(errorUpdate)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			api := auth.NewController(userServiceMock)

			res, err := api.Update(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}
