package login

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	authService "github.com/k6mil6/hackathon-game-backend/internal/service/auth"
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
		op := "handlers.user.login.New"

		log = log.With(
			slog.String("op", op),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			log.Error("error decoding JSON request:", err)

			render.JSON(w, r, resp.Error("error decoding JSON request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if req.Username == "" {
			w.WriteHeader(http.StatusUnauthorized)

			log.Error("username is required")

			render.JSON(w, r, resp.Error("username is required"))

			return
		}

		if req.Password == "" {
			w.WriteHeader(http.StatusUnauthorized)

			log.Error("password is required")

			render.JSON(w, r, resp.Error("password is required"))

			return
		}

		token, err := auth.LoginUser(ctx, req.Username, req.Password)
		if err != nil {
			if errors.Is(err, authService.ErrInvalidCredentials) {
				w.WriteHeader(http.StatusUnauthorized)

				render.JSON(w, r, resp.Error("user not found"))

				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			log.Error("error logging in:", slog.String("error", err.Error()))

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
