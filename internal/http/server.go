package http

import (
	"context"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
)

type Auth interface {
	LoginUser(ctx context.Context, username string, password string) (string, error)
	RegisterUser(ctx context.Context, username string, password string, classID int) (int, error)
	LoginAdmin(ctx context.Context, username string, password string) (string, error)
	RegisterAdmin(ctx context.Context, username, password string, registrantID, roleID int) (int, error)
}

type Tasks interface {
	GetAllUserTasks(ctx context.Context, userID int) ([]model.Task, error)
	GetAllAdminTasks(ctx context.Context, adminID int) ([]model.Task, error)
	GetByID(ctx context.Context, taskID int) (model.Task, error)
	Add(ctx context.Context, task model.Task) (int, error)
	MarkAsCompleted(ctx context.Context, taskID, adminID int) (model.Task, error)
	MarkAsCancelled(ctx context.Context, taskID, userID int) error
	MarkAsWaitingForAcceptance(ctx context.Context, taskID, userID int) error
}

type Transactions interface {
	AddUserTransaction(ctx context.Context, transaction *model.Transaction) error
	AddAdminTransaction(ctx context.Context, transaction *model.Transaction) error
	GetUserTransactions(ctx context.Context, userID int) ([]model.Transaction, error)
}

type Users interface {
	GetAll(ctx context.Context) ([]model.User, error)
	CreateBalance(ctx context.Context, userID int) error
	GetTopByBalance(ctx context.Context) ([]model.User, error)
}
