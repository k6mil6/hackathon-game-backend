package decline

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"log/slog"
	"net/http"
	"strconv"
)

type Response struct {
	resp.Response
}

func New(ctx context.Context, log *slog.Logger, tasks httpserver.Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.user.tasks.decline.New"

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

		userID, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get user ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get user ID"))

			return
		}

		if err := tasks.MarkAsCancelled(ctx, taskID, userID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to decline task", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to decline task"))

			return
		}

		w.WriteHeader(http.StatusOK)

		log.Info("task declined")

	}
}
