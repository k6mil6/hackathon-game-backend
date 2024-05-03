package login

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
	JWTToken string `json:"jwt_token"`
	resp.Response
}

func New(ctx context.Context, log *slog.Logger, auth httpserver.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.admin.login.New"

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

		token, err := auth.LoginAdmin(ctx, req.Username, req.Password)
		if err != nil {
			log.Error("error logging in:", err)

			render.JSON(w, r, resp.Error("internal server error"))

			return
		}

		responseOK(w, r, token)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, token string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		JWTToken: token,
	})
}
