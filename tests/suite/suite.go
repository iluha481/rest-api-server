package suite

import (
	"net/http"
	"server/initializers"
	"testing"
	"time"
)

type Suite struct {
	*testing.T
	Cfg        initializers.ServerConfig
	HttpClient *http.Client
}

func New(t *testing.T) *Suite {
	t.Helper()   // Функция будет восприниматься как вспомогательная для тестов
	t.Parallel() // Разрешаем параллельный запуск тестов

	// Читаем конфиг из файла
	config := initializers.NewServerConfig()

	return &Suite{
		T:          t,
		Cfg:        config,
		HttpClient: &http.Client{Timeout: 10 * time.Second},
	}
}
