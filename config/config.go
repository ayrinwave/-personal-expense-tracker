package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func LoadConfig() (*Config, error) {
	// Загружаем .env файл. Если его нет, ничего страшного, просто будут использованы переменные окружения.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		Port:      getEnv("PORT", "8080"), // Изменено на порт без двоеточия
		DBUrl:     getEnv("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "secret"),
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
