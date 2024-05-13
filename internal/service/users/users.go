package users

import (
	"context"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
)

type Users struct {
	log *slog.Logger

	storage Storage

	balanceStorage BalanceStorage
}

type Storage interface {
	GetAll(ctx context.Context) ([]model.User, error)
	GetTopByBalance(ctx context.Context) ([]model.User, error)
}

type BalanceStorage interface {
	CreateBalance(ctx context.Context, userID int) error
}

func New(log *slog.Logger, storage Storage, balanceStorage BalanceStorage) *Users {
	return &Users{
		log:            log,
		storage:        storage,
		balanceStorage: balanceStorage,
	}
}

func (u *Users) GetAll(ctx context.Context) ([]model.User, error) {
	op := "users.GetAll"

	log := u.log.With(
		slog.String("op", op),
	)

	log.Info("request received")

	users, err := u.storage.GetAll(ctx)
	if err != nil {
		log.Error("failed to get users", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("users returned")

	return users, nil
}

func (u *Users) CreateBalance(ctx context.Context, userID int) error {
	op := "users.CreateBalance"

	log := u.log.With(
		slog.String("op", op),
		slog.Int("userID", userID),
	)

	log.Info("request received")

	err := u.balanceStorage.CreateBalance(ctx, userID)
	if err != nil {
		log.Error("failed to create balance", slog.String("error", err.Error()))
		return err
	}

	log.Info("balance created")

	return nil
}

func (u *Users) GetTopByBalance(ctx context.Context) ([]model.User, error) {
	op := "users.GetTopByBalance"

	log := u.log.With(
		slog.String("op", op),
	)

	log.Info("request received")

	users, err := u.storage.GetTopByBalance(ctx)
	if err != nil {
		log.Error("failed to get top users by balance", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("users returned")

	return users, nil
}
