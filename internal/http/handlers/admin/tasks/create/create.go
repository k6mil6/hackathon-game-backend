package create

import (
	"context"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	"log/slog"
	"net/http"
)

type Request struct {
	Name       string  `json:"name"`
	Amount     float64 `json:"amount"`
	ForGroupID int     `json:"for_group_id"`
	UserID     int     `json:"user_id,omitempty"`
}

type Response struct {
	resp.Response
	ID int `json:"id"`
}

func New(ctx context.Context, log *slog.Logger, tasks httpserver.Tasks) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.admin.tasks.create.New"

		log = log.With(
			slog.String("op", op),
		)

		log.Info("request received")

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			log.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Name == "" {
			w.WriteHeader(http.StatusUnauthorized)

			log.Error("name is required")

			render.JSON(w, r, resp.Error("name is required"))

			return
		}

		if req.Amount == 0 {
			w.WriteHeader(http.StatusUnauthorized)

			log.Error("amount is required")

			render.JSON(w, r, resp.Error("amount is required"))

			return
		}

		if req.ForGroupID == 0 {
			w.WriteHeader(http.StatusUnauthorized)

			log.Error("for_group_id is required")

			render.JSON(w, r, resp.Error("for_group_id is required"))

			return
		}

		adminID, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get admin ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get admin ID"))

			return
		}

		id, err := tasks.Add(ctx, model.Task{
			Name:       req.Name,
			Amount:     req.Amount,
			CreatedBy:  adminID,
			ForGroupID: req.ForGroupID,
			UserID:     req.UserID,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to create task", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to create task"))

			return
		}

		log.Info("response sent")

		responseOK(w, r, id)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, id int) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.Response{Status: "ok"},
		ID:       id,
	})
}
