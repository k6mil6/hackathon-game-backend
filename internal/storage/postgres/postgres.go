package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/admins"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/balances"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/purchases"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/shop/items"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/tasks"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/transactions"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/users"
	"io"
	"log/slog"
	"reflect"
	"time"
)

type Storages struct {
	UsersStorage        *users.Storage
	BalancesStorage     *balances.Storage
	TransactionsStorage *transactions.Storage
	ShopItemsStorage    *items.Storage
	PurchasesStorage    *purchases.Storage
	AdminsStorage       *admins.Storage
	TasksStorage        *tasks.Storage
}

func NewStorages(
	postgresConnectionString string,
	maxRetries int,
	retryCooldown time.Duration,
	log *slog.Logger,
) (*Storages, error) {
	var db *sqlx.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = sqlx.Connect("postgres", postgresConnectionString)
		if err == nil {
			break
		}
		time.Sleep(retryCooldown)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres after %d retries: %w", maxRetries, err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis after %d retries: %w", maxRetries, err)
	}

	return &Storages{
		UsersStorage:        users.NewStorage(db, log),
		BalancesStorage:     balances.NewStorage(db, log),
		TransactionsStorage: transactions.NewStorage(db, log),
		ShopItemsStorage:    items.NewStorage(db, log),
		PurchasesStorage:    purchases.NewStorage(db, log),
		AdminsStorage:       admins.NewStorage(db, log),
		TasksStorage:        tasks.NewStorage(db, log),
	}, nil
}

func (s *Storages) CloseAll() error {
	val := reflect.ValueOf(s).Elem()
	var errList []error

	for i := 0; i < val.NumField(); i++ {
		storage := val.Field(i).Interface()
		if closer, ok := storage.(io.Closer); ok {
			err := closer.Close()
			if err != nil {
				errList = append(errList, err)
			}
		}
	}

	if len(errList) > 0 {
		return fmt.Errorf("failed to close all storages: %v", errList)
	}

	return nil
}
