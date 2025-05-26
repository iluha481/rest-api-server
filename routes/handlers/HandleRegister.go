package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"server/initializers"

	ssov1 "github.com/iluha481/protos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func HandleRegiser(authClient *initializers.GrpcClient) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}
			resp, err := authClient.Api.Register(r.Context(), &ssov1.RegisterRequest{
				Email:    req.Email,
				Password: req.Password,
			})
			if err != nil {
				st, ok := status.FromError(err)
				if !ok {
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

				log.Print("gRPC status", "code", st.Code(), "message", st.Message())
				switch st.Code() {
				case codes.InvalidArgument:
					http.Error(w, "invalid email or password", http.StatusBadRequest)
					return
				case codes.AlreadyExists:
					http.Error(w, "user already exists", http.StatusConflict)
					return
				default:
					http.Error(w, "internal server error", http.StatusInternalServerError)
					return
				}

			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"user_id": resp.GetUserId(),
			})
		})

}
