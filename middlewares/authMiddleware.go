package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrorKey        = "error"
	UidKey          = "uid"
)

// TODO
// вынести ErrorKey и UidKey в отдельный пакет и сделать кастом тип
// добавить логгер

func ParseJwtToken(
	tokenStr string,
	appSecret string,
) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(appSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		if float64(time.Now().Unix()) > claims["exp"].(float64) {

			return nil, ErrInvalidToken
		}

		return claims, nil
	} else {
		return nil, ErrInvalidToken
	}
}

func NewAuthMiddleware(
	appSecret string,
) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := r.Cookie("auth")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			claims, err := ParseJwtToken(tokenStr.Value, appSecret)
			if err != nil {
				ctx := context.WithValue(r.Context(), ErrorKey, ErrInvalidToken)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			ctx := context.WithValue(r.Context(), UidKey, claims["id"])
			next.ServeHTTP(w, r.WithContext(ctx))

		})

	}

}
