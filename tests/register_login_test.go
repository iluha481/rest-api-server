package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"server/tests/suite"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	appID          = 1
	appSecret      = "secret"
	passDefaultLen = 10
)

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
func postRequest(url string, contentType string, body []byte) (*http.Response, error) {
	// Create a new HTTP client
	client := &http.Client{}

	// Create a new request with the given URL, method, and body
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	// Set the content type header
	req.Header.Set("Content-Type", contentType)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func TestRegister_Login_Happy(t *testing.T) {
	s := suite.New(t)
	url := net.JoinHostPort(s.Cfg.Host, s.Cfg.Port)
	url = "http://" + url

	email := gofakeit.Email()
	pass := randomFakePassword()
	data := map[string]interface{}{
		"email":    email,
		"password": pass,
	}
	jsonBody, err := json.Marshal(data)
	require.NoError(t, err)
	respReg, err := postRequest(url+"/api/register", "application/json", jsonBody)
	require.NoError(t, err)
	require.NotEmpty(t, respReg)

	data = map[string]interface{}{
		"email":    email,
		"password": pass,
		"app_id":   appID,
	}
	jsonBody, err = json.Marshal(data)
	require.NoError(t, err)

	respLogin, err := postRequest(url+"/api/login", "application/json", jsonBody)

	require.NoError(t, err)
	require.NotEmpty(t, respLogin)

	var structReg struct {
		Id int64 `json:"user_id"`
	}
	bodyReg, err := io.ReadAll(respReg.Body)
	require.NoError(t, err)
	err = json.Unmarshal(bodyReg, &structReg)

	require.NoError(t, err)

	var token string
	for _, cookie := range respLogin.Cookies() {
		if cookie.Name == "token" {
			token = cookie.Value
			break
		}
	}
	require.NotEmpty(t, token)

	// Отмечаем время, в которое бы выполнен логин.
	// Это понадобится для проверки TTL токена
	loginTime := time.Now()

	// Парсим и валидируем токен
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	// Если ключ окажется невалидным, мы получим соответствующую ошибку
	require.NoError(t, err)

	// Преобразуем к типу jwt.MapClaims, в котором мы сохраняли данные
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	// Проверяем содержимое токена
	assert.Equal(t, structReg.Id, int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	// Проверяем, что TTL токена примерно соответствует нашим ожиданиям.
	assert.InDelta(t, loginTime.Add(s.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}
