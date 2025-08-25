package tests

import (
	"context"
	"errors"
	"testing"

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

func TestUpdate(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repinf.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx  context.Context
		id   int64
		info *models.UserInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)
		id  = int64(1)

		username = "new_testuser"
		email    = "new_test@example.com"

		serviceUpdateInfo = &models.UserInfo{
			Username: username,
			Email:    email,
		}

		repoUpdateInfo = &repoModel.UserInfo{
			Username: &username,
			Email:    &email,
		}

		repoErr = errors.New("repository error")
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "successful update",
			args: args{
				ctx:  ctx,
				id:   id,
				info: serviceUpdateInfo,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, repoUpdateInfo).Return(nil)
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
			name: "repository update error",
			args: args{
				ctx:  ctx,
				id:   id,
				info: serviceUpdateInfo,
			},
			err: repoErr,
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.UpdateMock.Expect(ctx, id, repoUpdateInfo).Return(repoErr)
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
			name: "invalid id",
			args: args{
				ctx:  ctx,
				id:   0,
				info: serviceUpdateInfo,
			},
			err: errors.New("user ID must be positive"),
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				return repoMocks.NewUserRepositoryMock(mc)
			},
			txManagerMock: func(mc *minimock.Controller) db.TxManager {
				return dbMocks.NewTxManagerMock(mc)
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

			err := service.Update(tt.args.ctx, tt.args.id, tt.args.info)

			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
