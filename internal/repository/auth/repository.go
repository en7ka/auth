package auth

import (
	"context"
	"errors"

	"github.com/en7ka/auth/internal/client/db"
	"github.com/jackc/pgx/v5"

	sq "github.com/Masterminds/squirrel"
	"github.com/en7ka/auth/internal/repository/auth/converter"
	"github.com/en7ka/auth/internal/repository/auth/model"
	repoif "github.com/en7ka/auth/internal/repository/repositoryinterface"
)

const (
	tableName      = "users"
	idColumn       = "id"
	usernameColumn = "username"
	emailColumn    = "email"
	passwordColumn = "password"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repoif.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserInfo) (int64, error) {
	qb := sq.Insert(tableName).
		Columns(usernameColumn, emailColumn, passwordColumn, roleColumn).
		Values(info.Username, info.Email, info.Password, info.Role).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}

	var id int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	qb := sq.Select(idColumn, usernameColumn, emailColumn, passwordColumn, roleColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var u model.User
	if err = r.db.DB().ScanOneContext(ctx, &u, q, args...); err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&u), nil
}

func (r *repo) Update(ctx context.Context, id int64, info *model.UserInfo) error {
	qb := sq.Update(tableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	if info.Username != "" {
		qb = qb.Set(usernameColumn, info.Username)
	}
	if info.Email != "" {
		qb = qb.Set(emailColumn, info.Email)
	}
	if info.Password != "" {
		qb = qb.Set(passwordColumn, info.Password)
	}

	qb = qb.Suffix("RETURNING id")

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}
	var updatedID int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("user not found")
		}
		return err
	}

	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	qb := sq.Delete(tableName).
		Where(sq.Eq{idColumn: id}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id")

	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}

	var deleted int64
	if err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&deleted); err != nil {
		return err
	}

	return nil
}
