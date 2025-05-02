package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"server/initializers"
	"time"
)

type App struct {
	httpServer *http.Server
	//Storage    *sqlite.Storage
	authClient *initializers.GrpcClient
}

var (
	ErrGracefulShutdown = errors.New("error stoping server")
)

func New() *App {
	// creating config and mux
	config := initializers.NewServerConfig()
	mux := initializers.NewServer(config)
	// creating server
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: mux,
	}
	ctx := context.Background()

	// nil -> logger
	authClient, err := initializers.NewGrpcClient(ctx, nil, net.JoinHostPort(config.SSO_host, config.SSO_port), config.SSO_timeout, config.SSO_retriesCount)
	if err != nil {
		log.Fatal("error creating authClient")
	}

	return &App{
		httpServer: httpServer,
		authClient: authClient,
	}
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
	// closing DB connection
	//
	//
	return nil
}
