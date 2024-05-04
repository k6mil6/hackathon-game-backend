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
	Tasks []ResponseTask `json:"tasks"`
}

type ResponseTask struct {
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func New(ctx context.Context, log *slog.Logger, tasks httpserver.Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.admin.tasks.all.New"

		log = log.With(
			slog.String("op", op),
		)

		log.Info("request received")

		userID, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			log.Error("failed to get user ID", slog.String("error", err.Error()))

			return
		}

		tasksRes := make([]ResponseTask, 0)

		tasks, err := tasks.GetAll(ctx, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get tasks", slog.String("error", err.Error()))

			return
		}

		for _, task := range tasks {
			tasksRes = append(tasksRes, ResponseTask{
				Name:      task.Name,
				Amount:    task.Amount,
				CreatedAt: task.CreatedAt,
			})
		}

		log.Info("response sent")

		responseOK(w, r, tasksRes)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, tasks []ResponseTask) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Tasks:    tasks,
	})
}
