package tasks

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
	"time"
)

const (
	InProgressStatusID           = 1
	WaitingForAcceptanceStatusID = 2
	CompletedStatusID            = 3
	CancelledStatusID            = 4
	AllGroupID                   = 1
	UserGroupID                  = 2
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

func (s *Storage) GetAll(ctx context.Context, userID int) ([]model.Task, error) {
	op := "tasks.GetAll"

	log := s.log.With(slog.String("op", op))

	log.Info("getting all tasks from storage")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return nil, err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	query := `SELECT id, name, amount, created_at, created_by, for_group_id FROM tasks WHERE user_id = $1 OR for_group_id = $2 ORDER BY created_at DESC`

	var tasks []dbTask
	if err := conn.SelectContext(ctx, &tasks, query, userID, AllGroupID); err != nil {
		log.Error("failed to get all tasks", slog.String("error", err.Error()))
		return nil, err
	}

	log.Info("got all tasks from storage")

	var shopTasks []model.Task
	for _, task := range tasks {
		shopTasks = append(shopTasks, model.Task(task))
	}

	return shopTasks, nil
}

func (s *Storage) Add(ctx context.Context, task model.Task) (int, error) {
	op := "tasks.Add"

	log := s.log.With(slog.String("op", op))

	log.Info("adding task to storage")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return 0, err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	var userID interface{} = task.UserID
	if task.UserID == 0 {
		userID = nil
	}

	var taskID int

	query := `INSERT INTO tasks (name, status_id, amount, created_by, for_group_id, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = conn.QueryRowxContext(ctx,
		query,
		task.Name,
		InProgressStatusID,
		task.Amount,
		task.CreatedBy,
		task.ForGroupID,
		userID,
	).Scan(&taskID)
	if err != nil {
		log.Error("failed to add task", slog.String("error", err.Error()))
		return 0, err
	}

	log.Info("added task to storage")

	return taskID, nil
}

func (s *Storage) MarkAsInProgress(ctx context.Context, taskID int) error {
	op := "tasks.MarkAsInProgress"

	log := s.log.With(slog.String("op", op))

	log.Info("marking task as in progress")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	query := `UPDATE tasks SET status_id = $1 WHERE id = $2`

	_, err = conn.ExecContext(ctx, query, InProgressStatusID, taskID)
	if err != nil {
		log.Error("failed to mark task as in progress", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as in progress")

	return nil
}

func (s *Storage) MarkAsWaitingForAcceptance(ctx context.Context, taskID int) error {
	op := "tasks.MarkAsWaitingForAcceptance"

	log := s.log.With(slog.String("op", op))

	log.Info("marking task as waiting for acceptance")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	query := `UPDATE tasks SET status_id = $1 WHERE id = $2`

	_, err = conn.ExecContext(ctx, query, WaitingForAcceptanceStatusID, taskID)

	if err != nil {
		log.Error("failed to mark task as waiting for acceptance", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as waiting for acceptance")

	return nil
}

func (s *Storage) MarkAsCompleted(ctx context.Context, taskID int) error {
	op := "tasks.MarkAsCompleted"

	log := s.log.With(slog.String("op", op))

	log.Info("marking task as completed")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	query := `UPDATE tasks SET status_id = $1 WHERE id = $2`

	_, err = conn.ExecContext(ctx, query, CompletedStatusID, taskID)
	if err != nil {
		log.Error("failed to mark task as completed", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as completed")

	return nil
}

func (s *Storage) MarkAsCancelled(ctx context.Context, taskID int) error {
	op := "tasks.MarkAsCancelled"

	log := s.log.With(slog.String("op", op))

	log.Info("marking task as cancelled")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	query := `UPDATE tasks SET status_id = $1 WHERE id = $2`

	_, err = conn.ExecContext(ctx, query, CancelledStatusID, taskID)
	if err != nil {
		log.Error("failed to mark task as cancelled", slog.String("error", err.Error()))
		return err
	}

	log.Info("marked task as cancelled")

	return nil
}

func (s *Storage) GetByID(ctx context.Context, taskID int) (model.Task, error) {
	op := "tasks.GetByID"

	log := s.log.With(slog.String("op", op))

	log.Info("getting task from storage")
	conn, err := s.db.Connx(ctx)
	if err != nil {
		log.Error("failed to get connection", slog.String("error", err.Error()))
		return model.Task{}, err
	}

	defer func(conn *sqlx.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error("failed to close connection", slog.String("error", err.Error()))
			return
		}
	}(conn)

	var task dbTask
	query := `SELECT id, name, status_id, amount, created_at, created_by, for_group_id, user_id FROM tasks WHERE id = $1`

	err = conn.GetContext(ctx, &task, query, taskID)
	if err != nil {
		log.Error("failed to get task", slog.String("error", err.Error()))
		return model.Task{}, err
	}

	log.Info("got task from storage")

	return model.Task(task), nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

type dbTask struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`
	StatusID   int       `db:"status_id"`
	Amount     float64   `db:"amount"`
	CreatedAt  time.Time `db:"created_at"`
	CreatedBy  int       `db:"created_by"`
	ForGroupID int       `db:"for_group_id"`
	UserID     int       `db:"user_id"`
}
