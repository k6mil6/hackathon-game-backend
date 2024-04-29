package main

import (
	"context"
	"github.com/k6mil6/hackathon-game-backend/internal/app"
	"github.com/k6mil6/hackathon-game-backend/internal/config"
	"github.com/k6mil6/hackathon-game-backend/internal/lib/logger"
	"github.com/k6mil6/hackathon-game-backend/internal/storage/postgres"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Info("connecting to db", slog.String("dsn", cfg.DB.PostgresDSN))
	storages, err := postgres.NewStorages(cfg.DB.PostgresDSN, cfg.DB.RetriesNumber, cfg.DB.RetryCooldown)
	if err != nil {
		log.Error("failed to connect to database", err)

		return
	}

	log.Info("connected to db", slog.String("dsn", cfg.DB.PostgresDSN))

	defer func() {
		if err := storages.CloseAll(); err != nil {
			log.Error("failed to close storages", err)
		}
	}()

	application := app.New(ctx, log, storages, cfg.JWT.TokenTTL, cfg.JWT.Secret, cfg.HTTPPort)

	go func() {
		application.HTTPServer.MustRun()
	}()

	<-ctx.Done()
}
