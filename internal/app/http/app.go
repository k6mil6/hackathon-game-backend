package httpapp

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpserver "github.com/k6mil6/hackathon-game-backend/internal/http"
	adminLogin "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/login"
	adminRegister "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/register"
	adminTasksAccept "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/tasks/accept"
	adminTasksALl "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/tasks/all"
	adminTasksCreate "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/tasks/create"
	adminUserAll "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/admin/user/all"
	userLogin "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/login"
	userRegister "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/register"
	userAllTasks "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/tasks/all"
	userTasksComplete "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/tasks/complete"
	userTasksDecline "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/tasks/decline"
	userTop "github.com/k6mil6/hackathon-game-backend/internal/http/handlers/user/top"
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

func New(
	ctx context.Context,
	log *slog.Logger,
	port int,
	auth httpserver.Auth,
	tasks httpserver.Tasks,
	transactions httpserver.Transactions,
	users httpserver.Users,
	secret string,
) *App {
	router := chi.NewRouter()

	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/register", userRegister.New(ctx, log, auth, users))
	router.Post("/login", userLogin.New(ctx, log, auth))
	router.Post("/admin/login", adminLogin.New(ctx, log, auth))
	router.Get("/user/top", userTop.New(ctx, log, users))

	// routes with authentication
	routerWithAuth := chi.NewRouter()
	routerWithAuth.Use(identity.New(secret))

	routerWithAuth.Post("/admin/register", adminRegister.New(ctx, log, auth))
	routerWithAuth.Post("/admin/task/create", adminTasksCreate.New(ctx, log, tasks))

	routerWithAuth.Get("/admin/user", adminUserAll.New(ctx, log, users))
	routerWithAuth.Get("/admin/task", adminTasksALl.New(ctx, log, tasks))
	routerWithAuth.Get("/admin/task/accept/{id}", adminTasksAccept.New(ctx, log, tasks, transactions))

	routerWithAuth.Get("/user/task", userAllTasks.New(ctx, log, tasks))
	routerWithAuth.Get("/user/task/decline/{id}", userTasksDecline.New(ctx, log, tasks))
	routerWithAuth.Get("/user/task/complete/{id}", userTasksComplete.New(ctx, log, tasks))

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
