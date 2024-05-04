package http

import (
	"context"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
)

type Auth interface {
	LoginUser(ctx context.Context, username string, password string) (string, error)
	RegisterUser(ctx context.Context, username string, password string) (int, error)
	LoginAdmin(ctx context.Context, username string, password string) (string, error)
	RegisterAdmin(ctx context.Context, username, password string, registrantID, roleID int) (int, error)
}

type Tasks interface {
	GetAll(ctx context.Context, userID int) ([]model.Task, error)
	Add(ctx context.Context, task model.Task) (int, error)
	MarkAsCompleted(ctx context.Context, taskID, adminID int) error
	MarkAsCancelled(ctx context.Context, taskID, adminID int) error
}
