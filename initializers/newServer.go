package initializers

import (
	"net/http"
)

func NewServer(
	config ServerConfig,
) http.Handler {

	mux := http.NewServeMux()
	AddRoutes(mux, config)

	return mux
}
