package accept

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	taskservice "github.com/k6mil6/hackathon-game-backend/internal/service/tasks"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	resp.Response
}

func New(ctx context.Context, log *slog.Logger, tasks httpserver.Tasks, transactions httpserver.Transactions) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.admin.tasks.accept.New"

		log = log.With(
			slog.String("op", op),
		)

		log.Info("request received")

		urlParam := chi.URLParam(r, "id")
		if urlParam == "" {
			w.WriteHeader(http.StatusBadRequest)

			log.Error("no id")

			render.JSON(w, r, resp.Error("no id"))

			return
		}

		taskID, err := strconv.Atoi(urlParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to parse id", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to parse id"))

			return
		}

		adminID, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get admin ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get admin ID"))

			return
		}

		task, err := tasks.MarkAsCompleted(ctx, taskID, adminID)
		if err != nil {
			if errors.Is(err, taskservice.ErrNotEnoughPermission) {
				w.WriteHeader(http.StatusBadRequest)

				log.Error("not enough permission", slog.String("error", err.Error()))

				render.JSON(w, r, resp.Error("not enough permission"))

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to mark task as accepted", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to mark task as accepted"))

			return
		}

		err = transactions.AddAdminTransaction(ctx, &model.Transaction{
			Amount:     task.Amount,
			SenderID:   adminID,
			ReceiverID: task.UserID,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to add admin transaction", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to add admin transaction"))

			return
		}

		render.JSON(w, r, resp.OK())
	}
}
