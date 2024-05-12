package tasks

import (
	"context"
	"errors"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	taskstorage "github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/tasks"
	"log/slog"
)

var (
	ErrNoTasks             = errors.New("no tasks")
	ErrNotEnoughPermission = errors.New("not enough permission")
)

type Tasks struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	GetAll(ctx context.Context, userID int) ([]model.Task, error)
	Add(ctx context.Context, task model.Task) (int, error)
	MarkAsCompleted(ctx context.Context, taskID int) error
	MarkAsCancelled(ctx context.Context, taskID int) error
	GetByID(ctx context.Context, taskID int) (model.Task, error)
	MarkAsInProgress(ctx context.Context, taskID int) error
	MarkAsWaitingForAcceptance(ctx context.Context, taskID int) error
}

func New(log *slog.Logger, storage Storage) *Tasks {
	return &Tasks{
		log:     log,
		storage: storage,
	}
}

func (t *Tasks) GetAll(ctx context.Context, userID int) ([]model.Task, error) {
	op := "tasks.GetAll"

	log := t.log.With(slog.String("op", op))

	log.Info("getting all tasks from storage")

	tasks, err := t.storage.GetAll(ctx, userID)
	if err != nil {
		log.Error("failed to get all tasks", slog.String("error", err.Error()))
		return nil, err
	}

	if len(tasks) == 0 {
		log.Info("no tasks found")
		return nil, ErrNoTasks
	}

	log.Info("got all tasks from storage")

	return tasks, nil
}

func (t *Tasks) GetByID(ctx context.Context, taskID int) (model.Task, error) {
	op := "tasks.GetByID"

	log := t.log.With(slog.String("op", op))

	log.Info("getting task from storage")

	task, err := t.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return model.Task{}, err
	}

	log.Info("got task from storage")

	return task, nil
}

func (t *Tasks) Add(ctx context.Context, task model.Task) (int, error) {
	op := "tasks.Add"

	log := t.log.With(slog.String("op", op))

	log.Info("adding task to storage")

	taskID, err := t.storage.Add(ctx, task)
	if err != nil {
		log.Error("failed to add task", slog.String("error", err.Error()))
		return 0, err
	}

	log.Info("added task to storage")

	return taskID, nil
}

func (t *Tasks) MarkAsInProgress(ctx context.Context, taskID, userID int) error {
	op := "tasks.MarkAsInProgress"

	log := t.log.With(slog.String("op", op))

	log.Info("marking task as in progress")

	task, err := t.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return err
	}

	if task.UserID != userID {
		log.Error("user does not have permission to mark this task as in progress")
		return ErrNotEnoughPermission
	}

	err = t.storage.MarkAsInProgress(ctx, taskID)
	if err != nil {
		log.Error("failed to mark task as in progress", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as in progress")

	return nil
}

func (t *Tasks) MarkAsWaitingForAcceptance(ctx context.Context, taskID, userID int) error {
	op := "tasks.MarkAsWaitingForAcceptance"

	log := t.log.With(slog.String("op", op))

	log.Info("marking task as waiting for acceptance")

	task, err := t.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return err
	}

	if task.UserID != userID {
		log.Error("user does not have permission to mark this task as waiting for acceptance")
		return ErrNotEnoughPermission
	}

	err = t.storage.MarkAsWaitingForAcceptance(ctx, taskID)
	if err != nil {
		log.Error("failed to mark task as waiting for acceptance", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as waiting for acceptance")

	return nil
}

func (t *Tasks) MarkAsCompleted(ctx context.Context, taskID, adminID int) error {
	op := "tasks.MarkAsCompleted"

	log := t.log.With(slog.String("op", op))

	log.Info("marking task as completed")

	task, err := t.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return err
	}

	if task.CreatedBy != adminID {
		log.Error("admin does not have permission to mark task as completed")
		return ErrNotEnoughPermission
	}

	if task.StatusID != taskstorage.WaitingForAcceptanceStatusID {
		log.Error("task is not waiting for acceptance")
		return ErrNotEnoughPermission
	}

	err = t.storage.MarkAsCompleted(ctx, taskID)
	if err != nil {
		log.Error("failed to mark task as completed", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as completed")

	return nil
}

func (t *Tasks) MarkAsCancelled(ctx context.Context, taskID, userID int) error {
	op := "tasks.MarkAsCancelled"

	log := t.log.With(slog.String("op", op))

	log.Info("marking task as cancelled")

	task, err := t.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return err
	}

	if task.UserID != userID {
		log.Error("user does not have permission to mark this task as cancelled")
		return ErrNotEnoughPermission
	}

	err = t.storage.MarkAsCancelled(ctx, taskID)
	if err != nil {
		log.Error("failed to mark task as cancelled", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as cancelled")

	return nil
}
