package main

import (
	"context"
	"github.com/k6mil6/hackathon-game-backend/internal/app"
	"github.com/k6mil6/hackathon-game-backend/internal/lib/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Get()
	log := logger.SetupLogger(cfg.Env).With(slog.String("env", cfg.Env))
	log.Debug("logger debug mode enabled")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	storages, err := storage.New(cfg.PostgresDatabaseDSN, cfg.RedisDatabaseDSN, cfg.DBRetriesNumber, cfg.DBRetryCooldown)
	if err != nil {
		log.Error("failed to connect to database", err)

		return
	}

	defer func() {
		if err := storages.CloseAll(); err != nil {
			log.Error("failed to close storages", err)
		}
	}()

	application := app.New()

	go func() {
		application.HTTPServer.MustRun()
	}()

	<-ctx.Done()
}
