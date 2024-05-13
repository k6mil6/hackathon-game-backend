package top

import (
	"context"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Users []ResponseUser `json:"users"`
}

type ResponseUser struct {
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

func New(ctx context.Context, log *slog.Logger, users httpserver.Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.user.top.New"

		log = log.With(
			slog.String("op", op),
		)

		log.Info("request received")

		users, err := users.GetTopByBalance(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("failed to get users", slog.String("error", err.Error()))

			return
		}

		usersRes := make([]ResponseUser, 0)

		for _, user := range users {
			usersRes = append(usersRes, ResponseUser{
				Username: user.Username,
				Balance:  user.Balance,
			})
		}

		log.Info("users returned")

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Users:    usersRes,
		})
	}
}
