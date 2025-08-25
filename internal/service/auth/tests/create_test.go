// Файл: internal/service/auth/tests/create_test.go

package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
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

func TestCreate(t *testing.T) {
	t.Parallel()

	type userRepositoryMockFunc func(mc *minimock.Controller) repinf.UserRepository
	type txManagerMockFunc func(mc *minimock.Controller) db.TxManager

	type args struct {
		ctx      context.Context
		userInfo *models.UserInfo
	}

	var (
		ctx         = context.Background()
		mc          = minimock.NewController(t)
		id          = gofakeit.Int64()
		name        = gofakeit.Name()
		email       = gofakeit.Email()
		password    = gofakeit.Password(true, true, true, false, false, 10)
		role        = "user"
		errorCreate = errors.New("repo create error")

		serviceUserInfo = &models.UserInfo{
			Username: name,
			Email:    email,
			Password: password,
			Role:     role,
		}

		repoUserInfo = &repoModel.UserInfo{
			Username: &name,
			Email:    &email,
			Password: &password,
			Role:     role,
		}
	)

	defer t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
		txManagerMock      txManagerMockFunc
	}{
		{
			name: "successful create",
			args: args{
				ctx:      ctx,
				userInfo: serviceUserInfo,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, repoUserInfo).Return(id, nil)
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
			name: "error during repository creation",
			args: args{
				ctx:      ctx,
				userInfo: serviceUserInfo,
			},
			want: 0,
			err:  errorCreate,
			userRepositoryMock: func(mc *minimock.Controller) repinf.UserRepository {
				mock := repoMocks.NewUserRepositoryMock(mc)
				mock.CreateMock.Expect(ctx, repoUserInfo).Return(int64(0), errorCreate)
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

			userRepositoryMock := tt.userRepositoryMock(mc)
			txManagerMock := tt.txManagerMock(mc)

			service := auth.NewService(userRepositoryMock, txManagerMock)

			newID, err := service.Create(tt.args.ctx, tt.args.userInfo)

			if tt.err != nil {
				require.ErrorContains(t, err, tt.err.Error())
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
