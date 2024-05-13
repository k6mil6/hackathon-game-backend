package all

import (
	"context"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Users []ResponseUser `json:"users"`
}

type ResponseUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func New(ctx context.Context, log *slog.Logger, users httpserver.Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.admin.user.all.New"

		log = log.With(
			slog.String("op", op),
		)

		log.Info("request received")

		_, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get admin ID", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get admin ID"))

			return
		}

		usersRes := make([]ResponseUser, 0)

		users, err := users.GetAll(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get users", slog.String("error", err.Error()))

			render.JSON(w, r, resp.Error("failed to get users"))

			return
		}

		for _, user := range users {
			usersRes = append(usersRes, ResponseUser{
				ID:       user.ID,
				Username: user.Username,
			})
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Users:    usersRes,
		})
	}
}
