package all

import (
	"context"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"log/slog"
	"net/http"
	"time"
)

type Response struct {
	resp.Response
	Tasks []TaskResponse `json:"tasks"`
}

type TaskResponse struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	StatusID   int       `json:"status_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
	ForGroupID int       `json:"for_group_id"`
	UserID     int       `json:"user_id,omitempty"`
}

func New(ctx context.Context, log *slog.Logger, tasks httpserver.Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.admin.tasks.all.New"

		log = log.With(
			slog.String("op", op),
		)

		adminID, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get admin ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get admin ID"))

			return
		}

		log.Info("request received")

		tasks, err := tasks.GetAllAdminTasks(ctx, adminID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get tasks", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get tasks"))

			return
		}

		log.Info("tasks retrieved", slog.Any("tasks", tasks))

		var taskResponse []TaskResponse

		for _, task := range tasks {
			taskResponse = append(taskResponse, TaskResponse{
				ID:         task.ID,
				Name:       task.Name,
				StatusID:   task.StatusID,
				Amount:     task.Amount,
				CreatedAt:  task.CreatedAt,
				ForGroupID: task.ForGroupID,
				UserID:     task.UserID,
			})
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Tasks:    taskResponse,
		})

	}
}
