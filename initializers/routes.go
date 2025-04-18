package initializers

import (
	"net/http"
	"server/handlers"
	"server/middlewares"
)

func AddRoutes(
	mux *http.ServeMux,
	config ServerConfig,

) {
	mux.Handle("/ping", middlewares.LoggingMiddleware(handlers.HandleHttp()))
}
