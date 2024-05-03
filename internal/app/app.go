package app

import (
	"context"
	httpapp "github.com/k6mil6/hackathon-game-backend/internal/app/http"
	authservice "github.com/k6mil6/hackathon-game-backend/internal/service/auth"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres"
	"log/slog"
	"time"
)

type App struct {
	HTTPServer *httpapp.App
}

func New(
	ctx context.Context,
	log *slog.Logger,
	storages *postgres.Storages,
	tokenTTL time.Duration,
	secret string,
	port int,
) *App {
	auth := authservice.New(log, storages.UsersStorage, storages.AdminsStorage, tokenTTL, secret)

	httpApp := httpapp.New(ctx, log, port, auth, secret)

	return &App{
		HTTPServer: httpApp,
	}
}
