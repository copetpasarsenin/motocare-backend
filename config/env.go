package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}
}

func GetEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func GetAppPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}

	if port := os.Getenv("APP_PORT"); port != "" {
		return port
	}

	return "8080"
}
