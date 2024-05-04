package tasks

import (
	"context"
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

func (s *Storage) GetAll(ctx context.Context) ([]model.Task, error) {
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

	var tasks []dbTask
	if err := conn.SelectContext(ctx, &tasks, "SELECT * FROM tasks"); err != nil {
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
		task.StatusID,
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
