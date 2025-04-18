package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleHttp() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode("asda")
		})
}
