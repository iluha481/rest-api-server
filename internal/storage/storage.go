package storage

import (
	"context"
	"errors"
	"log/slog"
	"server/internal/domain/models"
)

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type Storage struct {
	log          *slog.Logger
	userProvider UserProvider
	appProvider  AppProvider
}

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not foud")
)

func New(
	log *slog.Logger,
	userProvider UserProvider,
	appProvider AppProvider) *Storage {
	return &Storage{
		log:          log,
		userProvider: userProvider,
		appProvider:  appProvider,
	}

}
