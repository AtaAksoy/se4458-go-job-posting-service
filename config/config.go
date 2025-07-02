package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDSN     string
	Port      string
	RedisAddr string
	RedisDB   int
	RedisPass string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}
	dsn := getEnv("DB_DSN", "")
	if dsn == "" {
		log.Fatal("DB_DSN must be set in environment or .env file")
	}
	return &Config{
		DBDSN:     dsn,
		Port:      getEnv("PORT", "8080"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:   getEnvAsInt("REDIS_DB", 0),
		RedisPass: getEnv("REDIS_PASSWORD", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
