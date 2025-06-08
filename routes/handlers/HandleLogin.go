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

func HandleLogin(cfg initializers.ServerConfig, authClient *initializers.GrpcClient) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var req struct {
				Email    string `json:"email"`
				Password string `json:"password"`
				App_id   int32  `json:"app_id"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			resp, err := authClient.Api.Login(r.Context(), &ssov1.LoginRequest{
				Email:    req.Email,
				Password: req.Password,
				AppId:    req.App_id,
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
			cookie := http.Cookie{
				Name:     "token",
				Value:    resp.GetToken(),
				Domain:   cfg.Host,
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, &cookie)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"login": "success",
			})
		})
}
