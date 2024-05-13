package register

import (
	"context"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

func New(ctx context.Context, log *slog.Logger, auth httpserver.Auth, users httpserver.Users) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.register.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Username == "" {
			log.Error("username is required")

			render.JSON(w, r, resp.Error("username is required"))

			return
		}

		if req.Password == "" {
			log.Error("password is required")

			render.JSON(w, r, resp.Error("password is required"))

			return
		}

		log.Info("registering user")

		id, err := auth.RegisterUser(ctx, req.Username, req.Password)
		if err != nil {
			log.Error("error registering user:", err)

			render.JSON(w, r, resp.Error("error registering user"))

			return
		}

		log.Info("user registered with id", slog.Int("id", id))

		err = users.CreateBalance(ctx, id)
		if err != nil {
			log.Error("error creating balance:", err)
			render.JSON(w, r, resp.Error("error creating balance"))
			return
		}

		render.JSON(w, r, Response{
			resp.OK(),
		})
	}
}
