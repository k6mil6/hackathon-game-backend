package tasks

import (
	"context"
	"errors"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
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
}

func New(log *slog.Logger, storage Storage) *Tasks {
	return &Tasks{
		log:     log,
		storage: storage,
	}
}

func (s *Tasks) GetAll(ctx context.Context, userID int) ([]model.Task, error) {
	op := "tasks.GetAll"

	log := s.log.With(slog.String("op", op))

	log.Info("getting all tasks from storage")

	tasks, err := s.storage.GetAll(ctx, userID)
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

func (s *Tasks) Add(ctx context.Context, task model.Task) (int, error) {
	op := "tasks.Add"

	log := s.log.With(slog.String("op", op))

	log.Info("adding task to storage")

	taskID, err := s.storage.Add(ctx, task)
	if err != nil {
		log.Error("failed to add task", slog.String("error", err.Error()))
		return 0, err
	}

	log.Info("added task to storage")

	return taskID, nil
}

func (s *Tasks) MarkAsCompleted(ctx context.Context, taskID, adminID int) error {
	op := "tasks.MarkAsCompleted"

	log := s.log.With(slog.String("op", op))

	log.Info("marking task as completed")

	task, err := s.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return err
	}

	if task.CreatedBy != adminID {
		log.Error("admin does not have permission to mark task as completed")
		return ErrNotEnoughPermission
	}

	err = s.storage.MarkAsCompleted(ctx, taskID)
	if err != nil {
		log.Error("failed to mark task as completed", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as completed")

	return nil
}

func (s *Tasks) MarkAsCancelled(ctx context.Context, taskID, adminID int) error {
	op := "tasks.MarkAsCancelled"

	log := s.log.With(slog.String("op", op))

	log.Info("marking task as cancelled")

	task, err := s.storage.GetByID(ctx, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return err
	}

	if task.CreatedBy != adminID {
		log.Error("admin does not have permission to mark task as cancelled")
		return ErrNotEnoughPermission
	}

	err = s.storage.MarkAsCancelled(ctx, taskID)
	if err != nil {
		log.Error("failed to mark task as cancelled", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as cancelled")

	return nil
}
