package identity

import (
	"context"
	"errors"
	"github.com/go-chi/render"
	resp "github.com/k6mil6/hackathon-game-backend/internal/http/response"
	"github.com/k6mil6/hackathon-game-backend/internal/lib/jwt"
	"net/http"
	"strings"
)

func New(secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				render.JSON(w, r, resp.Error("Authorization header is empty"))

				return
			}
			headerParts := strings.Split(header, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				render.JSON(w, r, resp.Error("Authorization header is invalid"))
				return
			}

			if len(headerParts[1]) == 0 {
				render.JSON(w, r, resp.Error("Authorization token is empty"))
				return
			}

			id, err := jwt.GetID(headerParts[1], secret)
			if err != nil {
				render.JSON(w, r, resp.Error(err.Error()))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "id", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func GetID(ctx context.Context) (int, error) {
	if ctx.Value("id") == nil {
		return 0, errors.New("id not found in context")
	}
	return ctx.Value("id").(int), nil
}
