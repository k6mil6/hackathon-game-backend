package transactions

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
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
	TransferTypeID    = 1
	PurchaseTypeID    = 2
	DepositTypeID     = 3
	RefundTypeID      = 4
	PendingStatusID   = 1
	CompletedStatusID = 2
	CancelledStatusID = 3
)

func (s *Storage) AddTransaction(ctx context.Context, transaction *model.Transaction) error {
	op := "transactions.AddTransaction"

	log := s.log.With("op", op)

	log.Info("adding transaction", transaction)

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

	tx, err := conn.BeginTxx(ctx, nil)
	if err != nil {
		log.Error("failed to begin transaction:", slog.String("error", err.Error()))
		return err
	}

	var balance float64
	err = tx.GetContext(ctx, &balance, "SELECT balance FROM balances WHERE user_id = $1 FOR UPDATE", transaction.SenderID)
	if err != nil {
		log.Error("failed to get sender balance", slog.String("error", err.Error()))
		tx.Rollback()
		return err
	}

	if balance < transaction.Amount {
		log.Error("insufficient funds for the transaction")
		tx.Rollback()
		return sql.ErrNoRows
	}

	var transactionID int64
	err = tx.QueryRowContext(ctx, "INSERT INTO transactions (sender_id, receiver_id, amount, type_id, status_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		transaction.SenderID, transaction.ReceiverID, transaction.Amount, transaction.TypeID, PendingStatusID).Scan(&transactionID)
	if err != nil {
		log.Error("failed to insert transaction record", slog.String("error", err.Error()))
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE balances SET balance = balance - $1 WHERE user_id = $2", transaction.Amount, transaction.SenderID)
	if err != nil {
		log.Error("failed to update sender balance", slog.String("error", err.Error()))
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE balances SET balance = balance + $1 WHERE user_id = $2", transaction.Amount, transaction.ReceiverID)
	if err != nil {
		log.Error("failed to update receiver balance", slog.String("error", err.Error()))
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE transactions SET status_id = $1 WHERE id = $2", CompletedStatusID, transactionID)
	if err != nil {
		log.Error("failed to update transaction status", slog.String("error", err.Error()))
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Error("failed to commit transaction", slog.String("error", err.Error()))
		tx.Rollback()
		return err
	}

	return nil
}

func (s *Storage) GetUserTransactions(ctx context.Context, userID int) ([]model.Transaction, error) {
	op := "transactions.GetTransactions"

	log := s.log.With("op", op, "userID", userID)

	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return nil, err
	}
	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	var transactions []dbTransaction
	err = conn.SelectContext(ctx, &transactions, "SELECT * FROM transactions WHERE sender_id = $1 OR receiver_id = $1", userID)
	if err != nil {
		log.Error("failed to get transactions", slog.String("error", err.Error()))
		return nil, err
	}

	var result []model.Transaction
	for _, transaction := range transactions {
		result = append(result, model.Transaction{
			ID:         transaction.ID,
			SenderID:   transaction.SenderID,
			ReceiverID: transaction.ReceiverID,
			Amount:     transaction.Amount,
			TypeID:     transaction.TypeID,
			StatusID:   transaction.StatusID,
			CreatedAt:  transaction.CreatedAt,
		})
	}

	return result, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbTransaction struct {
	ID         int       `db:"id"`
	SenderID   int       `db:"sender_id"`
	ReceiverID int       `db:"receiver_id"`
	Amount     float64   `db:"amount"`
	TypeID     int       `db:"type_id"`
	StatusID   int       `db:"status_id"`
	CreatedAt  time.Time `db:"created_at"`
}
