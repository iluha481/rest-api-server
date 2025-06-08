package routes

import (
	"net/http"
	"server/initializers"
	"server/internal/storage"
	"server/middlewares"
	"server/routes/handlers"
)

func AddRoutes(
	mux *http.ServeMux,
	config initializers.ServerConfig,
	storage *storage.Storage,
	authClient *initializers.GrpcClient,
) {
	mux.Handle("/ping", middlewares.LoggingMiddleware(handlers.HandleHttp()))
	mux.Handle("/api/register", middlewares.LoggingMiddleware(handlers.HandleRegiser(authClient)))
	mux.Handle("/api/login", middlewares.LoggingMiddleware(handlers.HandleLogin(config, authClient)))
}
