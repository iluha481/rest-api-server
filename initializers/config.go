package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Host string
	Port string
}

func NewServerConfig() ServerConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	return ServerConfig{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
}
