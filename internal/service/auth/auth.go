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
)

type Auth struct {
	log          *slog.Logger
	usersStorage UsersStorage
	tokenTTL     time.Duration
	secret       string
}

type UsersStorage interface {
	Save(ctx context.Context, user *model.User) (int, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
}

func New(
	log *slog.Logger,
	usersStorage UsersStorage,
	tokenTTL time.Duration,
	secret string,
) *Auth {
	return &Auth{
		log:          log,
		usersStorage: usersStorage,
		tokenTTL:     tokenTTL,
		secret:       secret,
	}
}

func (a *Auth) Login(ctx context.Context, username string, password string) (string, error) {
	const op = "auth.Auth.Login"

	log := a.log.With(
		slog.String("username", username),
		slog.String("op", op),
	)

	log.Info("attempting login")
	user, err := a.usersStorage.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			log.Warn("user not found", username)

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

	token, err := jwt.NewToken(user, a.tokenTTL, a.secret)
	if err != nil {
		log.Error("failed to create token", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, username string, password string) (int, error) {
	const op = "auth.Auth.Register"

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
			log.Warn("user already exists", username)
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed to save user", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")
	return id, nil
}
