package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"server/initializers"
	"server/internal/storage"
	"server/internal/storage/postgresql"
	"server/routes"
	"time"

	_ "github.com/lib/pq"
)

// TODO:
// подумать над структурой
type App struct {
	httpServer *http.Server
	storage    *storage.Storage
	dbConn     *postgresql.Storage
	authClient *initializers.GrpcClient
}

var (
	ErrGracefulShutdown = errors.New("error stoping server")
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func NewMux(
	config initializers.ServerConfig,
	storage *storage.Storage,
	authClient *initializers.GrpcClient,
) http.Handler {

	mux := http.NewServeMux()
	routes.AddRoutes(mux, config, storage, authClient)

	return mux
}

func New() *App {
	// creating config and mux
	config := initializers.NewServerConfig()
	// creating server

	ctx := context.Background()

	db, err := postgresql.New(config.ConnectionString)
	if err != nil {
		log.Fatal("error connecting to DB", err)
	}
	slog := setupLogger(config.Env)

	storage := storage.New(slog, db, db)

	authClient, err := initializers.NewGrpcClient(ctx, slog, net.JoinHostPort(config.SSO_host, config.SSO_port), config.SSO_timeout, config.SSO_retriesCount)

	if err != nil {
		log.Fatal("error creating authClient")
	}
	mux := NewMux(config, storage, authClient)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: mux,
	}
	return &App{
		httpServer: httpServer,
		dbConn:     db,
		storage:    storage,
		authClient: authClient,
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	}
	return log
}

func (a *App) MustRun() {

	// running server in goroiutine
	go func() {
		log.Printf("Listening on %s\n", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stdout, "error listening and serving: %s\n", err)
		}
	}()
}

func (a *App) Stop() error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	if err := a.dbConn.Stop(); err != nil {
		log.Fatalf("DB connection closing failed: %v", err)
	}

	return nil
}
