package httpapp

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	adminLogin "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/login"
	adminRegister "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/register"
	userLogin "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/login"
	userRegister "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/register"
	"github.com/k6mil6/hackathon-game-backend/internal/http/middleware/identity"
	mwlogger "github.com/k6mil6/hackathon-game-backend/internal/http/middleware/logger"
	"log/slog"
	"net/http"
)

type App struct {
	log    *slog.Logger
	router *chi.Mux
	server *http.Server
}

func New(ctx context.Context, log *slog.Logger, port int, auth httpserver.Auth, secret string) *App {
	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/register", userRegister.New(ctx, log, auth))
	router.Post("/login", userLogin.New(ctx, log, auth))
	router.Post("/login/admin", adminLogin.New(ctx, log, auth))

	routerWithAuth := chi.NewRouter()
	routerWithAuth.Use(identity.New(secret))

	routerWithAuth.Post("/register/admin", adminRegister.New(ctx, log, auth))

	// routes with authentication

	router.Mount("/", routerWithAuth)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &App{
		log:    log,
		router: router,
		server: server,
	}
}

func (a *App) Run() error {
	a.log.Info("starting server", slog.String("address", a.server.Addr))
	return a.server.ListenAndServe()
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}
