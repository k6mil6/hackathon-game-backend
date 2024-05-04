package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	errs "github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/errors"
	"github.com/lib/pq"
	"log/slog"
	"time"
)

type Storage struct {
	db  *sqlx.DB
	log *slog.Logger
}

func NewStorage(db *sqlx.DB, log *slog.Logger) *Storage {
	return &Storage{
		db:  db,
		log: log,
	}
}

func (s *Storage) Save(ctx context.Context, user *model.User) (int, error) {
	op := "users.Save"

	log := s.log.With("op", op)

	log.Info("saving user", user)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
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
			log.Error("user already exists", slog.String("error", err.Error()))
			return 0, errs.ErrUserExists
		}

		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, err
	}

	log.Debug("saved user", id)
	return id, nil
}

func (s *Storage) GetByUsername(ctx context.Context, username string) (model.User, error) {
	op := "users.GetByUsername"

	log := s.log.With("op", op)

	log.Info("getting user by username", username)
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

	query := `SELECT id, username, password_hash, registered_at, hired_at FROM users WHERE username = $1`

	var user dbUser

	if err := conn.GetContext(ctx, &user, query, username); err != nil {
		log.Error("failed to get user", slog.String("error", err.Error()))
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("user not found", slog.String("error", err.Error()))
			return model.User{}, errs.ErrUserNotFound
		}
		return model.User{}, err
	}

	return model.User(user), nil

}

func (s *Storage) GetByID(ctx context.Context, id int) (model.User, error) {
	op := "users.GetByID"

	log := s.log.With("op", op)

	log.Info("getting user by id", id)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return model.User{}, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `SELECT id, username, email, password_hash, registered_at, hired_at FROM users WHERE id = $1`

	var user dbUser

	if err := conn.GetContext(ctx, &user, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("user not found", slog.String("error", err.Error()))
			return model.User{}, errs.ErrUserNotFound
		}
		log.Error("failed to get user", slog.String("error", err.Error()))
		return model.User{}, err
	}

	log.Info("got user", user)
	return model.User(user), nil
}

func (s *Storage) UpdateHiredAt(ctx context.Context, id int, hiredAt time.Time) error {
	op := "users.UpdateHiredAt"

	log := s.log.With("op", op)

	log.Info("updating hired_at", id, hiredAt)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `UPDATE users SET hired_at = $1 WHERE id = $2`

	if _, err := conn.ExecContext(ctx, query, hiredAt, id); err != nil {
		log.Error("failed to update hired_at", slog.String("error", err.Error()))
		return err
	}

	log.Info("updated hired_at", id, hiredAt)
	return nil
}

func (s *Storage) UpdateEmail(ctx context.Context, id int, email string) error {
	op := "users.UpdateEmail"

	log := s.log.With("op", op)

	log.Info("updating email", id, email)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `UPDATE users SET email = $1 WHERE id = $2`

	if _, err := conn.ExecContext(ctx, query, email, id); err != nil {
		log.Error("failed to update email", slog.String("error", err.Error()))
		return err
	}

	log.Info("updated email", id, email)
	return nil
}

func (s *Storage) UpdatePassword(ctx context.Context, id int, passwordHash []byte) error {
	op := "users.UpdatePassword"

	log := s.log.With("op", op)

	log.Info("updating password", id, passwordHash)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `UPDATE users SET password_hash = $1 WHERE id = $2`

	if _, err := conn.ExecContext(ctx, query, passwordHash, id); err != nil {
		log.Error("failed to update password", slog.String("error", err.Error()))
		return err
	}

	log.Info("updated password", id, passwordHash)
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbUser struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	RegisteredAt time.Time `db:"registered_at"`
	HiredAt      time.Time `db:"hired_at"`
}
