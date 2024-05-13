package balances

import (
	"context"
	"github.com/jmoiron/sqlx"
	"log/slog"
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

func (s *Storage) CreateBalance(ctx context.Context, userID int) error {
	op := "balances.CreateBalance"

	log := s.log.With("op", op, "userID", userID)

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

	_, err = conn.ExecContext(ctx, "INSERT INTO balances (user_id, balance) VALUES ($1, 0)", userID)
	if err != nil {
		log.Error("failed to create balance", slog.String("error", err.Error()))
		return err
	}

	log.Info("created balance")
	return nil
}

func (s *Storage) GetBalance(ctx context.Context, userID int) (float64, error) {
	op := "balances.GetBalance"

	log := s.log.With("op", op, "userID", userID)

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

	var balance float64
	err = conn.GetContext(ctx, &balance, "SELECT balance FROM balances WHERE user_id = $1 FOR UPDATE", userID)
	if err != nil {
		log.Error("failed to get balance", slog.String("error", err.Error()))
		return 0, err
	}

	log.Debug("balance", balance)
	return balance, nil
}

func (s *Storage) AddBalance(ctx context.Context, userID int, amount float64) error {
	op := "balances.AddBalance"

	log := s.log.With("op", op, "userID", userID)

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

	_, err = conn.ExecContext(ctx, "UPDATE balances SET balance = balance + $1 WHERE user_id = $2", amount, userID)
	if err != nil {
		log.Error("failed to add balance", slog.String("error", err.Error()))
		return err
	}

	log.Debug("balance", amount)
	return nil
}

func (s *Storage) SubtractBalance(ctx context.Context, userID int, amount float64) error {
	op := "balances.SubtractBalance"

	log := s.log.With("op", op, "userID", userID)

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

	_, err = conn.ExecContext(ctx, "UPDATE balances SET balance = balance - $1 WHERE user_id = $2", amount, userID)
	if err != nil {
		log.Error("failed to subtract balance", slog.String("error", err.Error()))
		return err
	}

	log.Debug("balance", amount)
	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
