package cmd

import (
	"context"
	"fmt"
	userv1 "github.com/en7ka/auth/pkg/user_v1"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

const (
	dbDSN = "DEFAULT_DSN"
)

type Storage struct {
	con *pgx.Conn
	ctx context.Context
}

func InitStorage() (*Storage, error) {
	ctx := context.Background()
	dbDSN := getDSN()
	conn, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	return &Storage{con: conn, ctx: ctx}, nil
}

func getDSN() string {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = dbDSN
	}

	return dsn
}

func (s *Storage) CloseCon() error {
	err := s.con.Close(s.ctx)
	if err != nil {
		log.Printf("failed to close connection: %v", err)
		return err
	}

	return nil
}

// PostgresInterface
type PostgresInterface interface {
	Save(user User) (int64, error)
	Update(update UpdateUser) error
	Delete(id DeleteID) error
	GetUser(params GetUserPar) (*User, error)
}

// Save
func (s *Storage) Save(user User) (int64, error) {
	var id int64
	query := "INSERT INTO users (username, email, password, role) VALUES (1,2, 3,4) RETURNING id"
	// Возвращаем ID после вставки
	if err := s.con.QueryRow(s.ctx, query, user.Username, user.Email, user.Password, user.Role).Scan(&id); err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return 0, err
	}
	log.Printf("Inserted user into database with id: %d", id)

	return id, nil
}

// Update
func (s *Storage) Update(update UpdateUser) error {
	query := "UPDATE users SET username = $1, email = $2 WHERE id = $3"
	res, err := s.con.Exec(s.ctx, query, update.Username, update.Email, update.ID)
	if err != nil {
		log.Printf("Error updating user into database: %v", err)
		return err
	}
	log.Printf("Updated user into database: %d", res.RowsAffected())

	return nil
}

// Delete
func (s *Storage) Delete(id DeleteID) error {
	res, err := s.con.Exec(s.ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting user into database: %v", err)
		return err
	}
	log.Printf("Deleted user into database: %d", res.RowsAffected())

	return nil
}

// GetUser
func (s *Storage) GetUser(params GetUserPar) (*User, error) {
	var user User
	var err error

	query := sq.Select("id", "username", "email", "role", "created_at", "updated_at").From("users")
	switch {
	case params.ID != nil:
		query = query.Where(sq.Eq{"id": *params.ID})
	case params.Username != nil:
		query = query.Where(sq.Eq{"username": *params.Username})
	default:
		return nil, fmt.Errorf("no username provided")
	}
	dbQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := s.con.QueryRow(s.ctx, dbQuery, args...)
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) GetByUsernames(ctx context.Context, req *userv1.GetByUsernameRequest) (*userv1.GetByUsernamesResponse, error) {
	usernames := req.GetUsername()
	if len(usernames) == 0 {
		return &userv1.GetByUsernamesResponse{}, nil
	}
	query := "SELECT id, username FROM users WHERE username = ANY($1)"

	rows, err := s.con.Query(ctx, query, usernames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*userv1.UserInfo
	for rows.Next() {
		var user userv1.UserInfo
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return &userv1.GetByUsernamesResponse{Users: users}, nil
}
