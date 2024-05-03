package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/k6mil6/hackathon-game-backend/internal/lib/jwt"
	"github.com/k6mil6/hackathon-game-backend/internal/model"
	errs "github.com/k6mil6/hackathon-game-backend/internal/storage/postgres/errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrAdminNotFound      = errors.New("admin not found")
	ErrAdminExists        = errors.New("admin already exists")
)

type Auth struct {
	log           *slog.Logger
	usersStorage  UsersStorage
	adminsStorage AdminsStorage
	tokenTTL      time.Duration
	secret        string
}

type UsersStorage interface {
	Save(ctx context.Context, user *model.User) (int, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

type AdminsStorage interface {
	Save(ctx context.Context, admin *model.Admin) (int, error)
	GetByUsername(ctx context.Context, username string) (model.Admin, error)
	GetByID(ctx context.Context, id int) (model.Admin, error)
}

func New(
	log *slog.Logger,
	usersStorage UsersStorage,
	adminsStorage AdminsStorage,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		log:           log,
		usersStorage:  usersStorage,
		adminsStorage: adminsStorage,
		tokenTTL:      tokenTTL,
		secret:        secret,
	}
}

func (a *Auth) LoginUser(ctx context.Context, username string, password string) (string, error) {
	const op = "auth.Auth.LoginUser"

	log := a.log.With(
		slog.String("username", username),
		slog.String("op", op),
	)

	log.Info("attempting login")
	user, err := a.usersStorage.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Error("user not found", username)

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user by login", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in")

	token, err := jwt.NewToken(user.ID, user.Username, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed to create token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterUser(ctx context.Context, username string, password string) (int, error) {
	const op = "auth.Auth.RegisterUser"

	log := a.log.With(
		slog.String("username", username),
		slog.String("op", op),
	)

	log.Info("attempting registration")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	user := model.User{
		Username:     username,
		PasswordHash: passwordHash,
	}

	id, err := a.usersStorage.Save(ctx, &user)
	if err != nil {
		if errors.Is(err, errs.ErrUserExists) {
			log.Error("user already exists", username)
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")
	return id, nil
}

func (a *Auth) LoginAdmin(ctx context.Context, username string, password string) (string, error) {
	const op = "auth.Auth.LoginAdmin"

	log := a.log.With(
		slog.String("username", username),
		slog.String("op", op),
	)

	log.Info("attempting login")
	admin, err := a.adminsStorage.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrAdminNotFound) {
			log.Error("admin not found", username)
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(admin.PasswordHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("admin logged in")

	token, err := jwt.NewToken(admin.ID, admin.Username, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed to create token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterAdmin(ctx context.Context, username, password string, registrantID, roleID int) (int, error) {
	const op = "auth.Auth.RegisterAdmin"

	log := a.log.With(
		slog.String("username", username),
		slog.String("op", op),
	)

	registrant, err := a.adminsStorage.GetByID(ctx, registrantID)
	if err != nil {
		if errors.Is(err, errs.ErrAdminNotFound) {
			log.Error("admin not found", username)
			return 0, fmt.Errorf("%s: %w", op, ErrAdminNotFound)
		}
		log.Error("failed to get registrant", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("attempting registration")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	admin := model.Admin{
		Username:     username,
		PasswordHash: passwordHash,
		RoleID:       roleID,
		RegisteredBy: registrant.ID,
	}

	id, err := a.adminsStorage.Save(ctx, &admin)
	if err != nil {
		if errors.Is(err, errs.ErrAdminExists) {
			log.Error("admin already exists", username)
			return 0, fmt.Errorf("%s: %w", op, ErrAdminExists)
		}

		log.Error("failed to save admin", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("admin registered")
	return id, nil
}
