package transactions

import (
	"context"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
)

type Transactions struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	AddUserTransaction(ctx context.Context, transaction *model.Transaction) error
	AddAdminTransaction(ctx context.Context, transaction *model.Transaction) error
	GetUserTransactions(ctx context.Context, userID int) ([]model.Transaction, error)
}

func New(log *slog.Logger, storage Storage) *Transactions {
	return &Transactions{
		log:     log,
		storage: storage,
	}
}

func (t *Transactions) AddUserTransaction(ctx context.Context, transaction *model.Transaction) error {
	op := "transactions.AddUserTransaction"

	log := t.log.With("op", op)

	log.Info("adding user transaction", transaction)

	err := t.storage.AddUserTransaction(ctx, transaction)
	if err != nil {
		log.Error("failed to add user transaction", slog.String("error", err.Error()))
		return err
	}

	log.Info("added user transaction")

	return nil
}

func (t *Transactions) AddAdminTransaction(ctx context.Context, transaction *model.Transaction) error {
	op := "transactions.AddAdminTransaction"

	log := t.log.With("op", op)

	log.Info("adding admin transaction", transaction)

	err := t.storage.AddAdminTransaction(ctx, transaction)
	if err != nil {
		log.Error("failed to add admin transaction", slog.String("error", err.Error()))
		return err
	}

	log.Info("added admin transaction")

	return nil
}

func (t *Transactions) GetUserTransactions(ctx context.Context, userID int) ([]model.Transaction, error) {
	op := "transactions.GetUserTransactions"

	log := t.log.With("op", op, "userID", userID)

	log.Info("getting user transactions")

	transactions, err := t.storage.GetUserTransactions(ctx, userID)
	if err != nil {
		log.Error("failed to get user transactions", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("got user transactions")

	return transactions, nil
}
