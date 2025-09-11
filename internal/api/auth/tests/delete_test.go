package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/en7ka/auth/internal/api/auth"
	serviceMocks "github.com/en7ka/auth/internal/service/mocks"
	"github.com/en7ka/auth/internal/service/servinterface"
	userv1 "github.com/en7ka/auth/pkg/user_v1"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type userServiceMockFunc func(mc *minimock.Controller) servinterface.UserService

	type args struct {
		ctx context.Context
		req *userv1.DeleteRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id          = gofakeit.Int64()
		errorDelete = errors.New("delete error")
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
			name: "successful deletion",
			args: args{
				ctx: ctx,
				req: &userv1.DeleteRequest{Id: id},
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(nil)
				return mock
			},
		},
		{
			name: "error during deletion",
			args: args{
				ctx: ctx,
				req: &userv1.DeleteRequest{Id: id},
			},
			want: nil,
			err:  errorDelete,
			userServiceMock: func(mc *minimock.Controller) servinterface.UserService {
				mock := serviceMocks.NewUserServiceMock(mc)
				mock.DeleteMock.Expect(ctx, id).Return(errorDelete)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			userServiceMock := tt.userServiceMock(mc)
			api := user.NewImplementation(userServiceMock)

			res, err := api.Delete(tt.args.ctx, tt.args.req)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}
