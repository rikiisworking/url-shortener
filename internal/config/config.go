package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	RedisAddr       string
	ServerPort      string
	ShortCodeLength int
}

func Load() Config {
	shortLen, _ := strconv.Atoi(getEnv("SHORT_CODE_LENGTH", "7"))

	return Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "password"),
		DBName:          getEnv("DB_NAME", "urlshortener"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort:      getEnv("PORT", "8080"),
		ShortCodeLength: shortLen,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
