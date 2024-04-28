package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	errs "github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/errors"
	"github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Save(ctx context.Context, user *model.User) (int, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id`

	var id int

	if err := conn.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.PasswordHash,
	).Scan(&id); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, errs.ErrUserExists
		}
		return 0, err
	}

	return id, nil
}

func (s *Storage) GetByUsername(ctx context.Context, username string) (model.User, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return model.User{}, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `SELECT id, username, password_hash, created_at FROM users WHERE username = $1`

	var user dbUser

	if err := conn.GetContext(ctx, &user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, errs.ErrUserNotFound
		}
		return model.User{}, err
	}

	return model.User(user), nil

}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbUser struct {
	ID           int
	Username     string
	PasswordHash []byte
	CreatedAt    string
}
