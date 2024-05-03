package admins

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

const (
	AdminRoleID = 1
)

func (s *Storage) Save(ctx context.Context, admin *model.Admin) (int, error) {
	op := "admins.Save"

	log := s.log.With("op", op)

	log.Info("saving admin", admin)

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

	query := `INSERT INTO admins (username, password_hash, registered_by, role_id) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int

	if err := conn.QueryRowContext(
		ctx,
		query,
		admin.Username,
		admin.PasswordHash,
		admin.RegisteredBy,
		admin.RoleID,
	).Scan(&id); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			log.Error("admin already exists", slog.String("error", err.Error()))
			return 0, errs.ErrAdminExists
		}
		log.Error("failed to save admin", slog.String("error", err.Error()))
		return 0, err
	}

	log.Info("saved admin", admin)

	return id, nil
}

func (s *Storage) GetByUsername(ctx context.Context, username string) (model.Admin, error) {
	op := "admins.GetByUsername"

	log := s.log.With("op", op)

	log.Info("getting admin by username", username)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return model.Admin{}, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	var admin dbAdmin

	query := `SELECT id, username, password_hash, role_id FROM admins WHERE username = $1`

	err = conn.GetContext(ctx, &admin, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("admin not found", slog.String("error", err.Error()))
			return model.Admin{}, errs.ErrAdminNotFound
		}
		log.Error("failed to get admin", slog.String("error", err.Error()))
		return model.Admin{}, err
	}

	log.Info("admin found", admin)
	return model.Admin(admin), nil
}

func (s *Storage) GetByID(ctx context.Context, id int) (model.Admin, error) {
	op := "admins.GetByID"

	log := s.log.With("op", op)

	log.Info("getting admin by id", id)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return model.Admin{}, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	var admin dbAdmin

	query := `SELECT id, username, email, password_hash, role_id FROM admins WHERE id = $1`

	err = conn.GetContext(ctx, &admin, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("admin not found", slog.String("error", err.Error()))
			return model.Admin{}, errs.ErrAdminNotFound
		}
		log.Error("failed to get admin", slog.String("error", err.Error()))
		return model.Admin{}, err
	}

	log.Info("admin found", admin)
	return model.Admin(admin), nil
}

func (s *Storage) UpdateEmail(ctx context.Context, id int, email string) error {
	op := "admins.UpdateEmail"

	log := s.log.With("op", op)

	log.Info("updating admin email", email)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `UPDATE admins SET email = $1 WHERE id = $2`

	_, err = conn.ExecContext(ctx, query, email, id)
	if err != nil {
		log.Error("failed to update admin email", slog.String("error", err.Error()))
		return err
	}

	log.Info("updated admin email", email)
	return nil
}

func (s *Storage) UpdatePassword(ctx context.Context, id int, passwordHash string) error {
	op := "admins.UpdatePassword"

	log := s.log.With("op", op)

	log.Info("updating admin password", passwordHash)
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	query := `UPDATE admins SET password_hash = $1 WHERE id = $2`

	_, err = conn.ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		log.Error("failed to update admin password", slog.String("error", err.Error()))
		return err
	}

	log.Info("updated admin password", passwordHash)
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbAdmin struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash []byte    `db:"password_hash"`
	RegisteredAt time.Time `db:"registered_at"`
	RegisteredBy int       `db:"registered_by"`
	RoleID       int       `db:"role_id"`
}
