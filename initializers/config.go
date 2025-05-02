package initializers

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Host string
	Port string

	SSO_host         string
	SSO_port         string
	SSO_timeout      time.Duration
	SSO_retriesCount int
}

func NewServerConfig() ServerConfig {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	timeout, err := strconv.Atoi(os.Getenv("SSO_timeout"))
	if err != nil {
		log.Fatal("can not convert SSO_timeout to time")
	}
	retries, err := strconv.Atoi(os.Getenv("SSO_retriesCount"))
	if err != nil {
		log.Fatal("can not convert SSO_retriesCount to int")
	}
	return ServerConfig{
		Host:             os.Getenv("HOST"),
		Port:             os.Getenv("PORT"),
		SSO_host:         os.Getenv("SSO_host"),
		SSO_port:         os.Getenv("SSO_port"),
		SSO_timeout:      time.Duration(timeout) * time.Second,
		SSO_retriesCount: retries,
	}
}
