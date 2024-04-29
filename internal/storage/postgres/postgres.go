package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/users"
	"io"
	"reflect"
	"time"
)

type Storages struct {
	UsersStorage *users.Storage
}

func NewStorages(postgresConnectionString string, maxRetries int, retryCooldown time.Duration) (*Storages, error) {
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
		UsersStorage: users.NewStorage(db),
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
