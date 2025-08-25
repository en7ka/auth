package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/en7ka/auth/internal/client/db"
	dbMocks "github.com/en7ka/auth/internal/client/mocks"
	"github.com/en7ka/auth/internal/models"
	repoModel "github.com/en7ka/auth/internal/repository/auth/model"
	repoMocks "github.com/en7ka/auth/internal/repository/mocks"
	repinf "github.com/en7ka/auth/internal/repository/repositoryinterface"
	"github.com/en7ka/auth/internal/service/auth"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repinf.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = int64(1)

		username = "testuser"
		email    = "test@example.com"
		password = "12345678"
		role     = "user"

		repoUser = &repoModel.User{
			Id:   id,
			Role: role,
			Info: repoModel.UserInfo{
				Username: &username,
				Email:    &email,
				Password: &password,
				Role:     role,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		serviceUserInfo = &models.UserInfo{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "12345678",
			Role:     role,
		}

		errorGet = errors.New("user not found")
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *models.UserInfo
		err                error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "successful get",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: serviceUserInfo,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(repoUser, nil)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
		{
			name: "user not found",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  errorGet,
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.GetMock.Expect(ctx, id).Return(nil, errorGet)
				return mock
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				mock := dbMocks.NewTxManagerMock(mc)
				mock.ReadCommitedMock.Set(func(ctx context.Context, f db.Handler) error {
					return f(ctx)
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userRepository := tt.userRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(userRepository, txManagerMock)

			result, err := service.Get(tt.args.ctx, tt.args.id)

			if tt.err != nil {
				require.Nil(t, result)
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, result)
			}
		})
	}
}
