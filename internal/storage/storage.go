package storage

import (
	"log/slog"
)

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type Storage struct {
	log 			*slog.Logger
	userProvider 	UserProvider
	appProvider 	AppProvider

}

func New(
	log 			*slog.Logger, 
	userProvider 	UserProvider, 
	appProvider 	AppProvider) *Storage {
	return &Storage {
		log: 			log,
		userProvider: 	userProvider,
		appProvider: 	appProvider
	} 

}