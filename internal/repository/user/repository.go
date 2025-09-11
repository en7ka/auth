package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/en7ka/auth/internal/client/db"
	"github.com/en7ka/auth/internal/models"
	repoif "github.com/en7ka/auth/internal/repository/repositoryinterface"
	"github.com/en7ka/auth/internal/utils"
)

const (
	tableName       = "users"
	tableAccessName = "access"

	idColumn       = "id"
	usernameColumn = "username"
	emailColumn    = "email"
	passwordColumn = "password"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repoif.AuthRepository {
	return &repo{db: db}
}

func (r *repo) Login(ctx context.Context, user models.LoginRequest) (*models.UserInfoJwt, error) {
	if user.Username == "" || user.Password == "" {
		return nil, fmt.Errorf("invalid username or password")
	}

	builder := sq.Select(passwordColumn, roleColumn).
		From(tableName).
		Where(sq.Eq{usernameColumn: user.Username}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.Login",
		QueryRaw: query,
	}

	var userInfo models.UserInfoJwt
	var password string
	var role bool

	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&password, &role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	if !utils.VerifyPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}
	userInfo.Username = user.Username
	userInfo.Role = role

	return &userInfo, nil
}

func (r *repo) GetUserRole(ctx context.Context, username string) (bool, error) {
	builder := sq.Select(roleColumn).
		From(tableName).
		Where(sq.Eq{roleColumn: username}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.GetUserRole",
		QueryRaw: query,
	}

	var role bool
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("user not found")
		}
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return role, nil
}

func (r *repo) GetUserAccess(ctx context.Context, isAdmin bool) ([]string, error) {
	builder := sq.Select(roleColumn).
		From(tableName).
		Where(sq.Eq{roleColumn: isAdmin}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "auth_repository.GetUserAccess",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	endpoints := make([]string, 0)
	for rows.Next() {
		var endpoint string
		if err = rows.Scan(&endpoint); err != nil {
			return nil, fmt.Errorf("failed to scan endpoint: %w", err)
		}
		endpoints = append(endpoints, endpoint)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	if len(endpoints) == 0 {
		return endpoints, nil
	}

	return endpoints, nil
}
