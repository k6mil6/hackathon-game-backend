package register

import (
	"context"
	"github.com/go-chi/render"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}

type Response struct {
	resp.Response
}

func New(ctx context.Context, log *slog.Logger, auth httpserver.Auth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.register.New"

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

		roleID := 1

		if req.RoleID != 0 {
			roleID = int(req.RoleID)
		}

		log.Info("registering user")

		registrantID, err := identity.GetID(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("error getting registrant id:", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		id, err := auth.RegisterAdmin(ctx, req.Username, req.Password, registrantID, roleID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.Error("error registering admin:", err)

			render.JSON(w, r, resp.Error("error registering admin"))

			return
		}

		log.Info("admin registered with id", slog.Int("id", id))

		render.JSON(w, r, Response{
			resp.OK(),
		})
	}
}
