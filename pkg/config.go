package pkg

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionURI    string
	JwtSecret          string
	RedisAddr          string
	LogLevel           slog.Level
	CorsAllowedOrigins string
}

func LoadConfig() (Config, error) {
	godotenv.Load(".env")

	dbConnectionURI, err := requireEnv("MYSQL_CONNECTION_URI")
	if err != nil {
		return Config{}, err
	}

	jwtSecret, err := requireEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	redisAddress, err := requireEnv("REDIS_ADDR")
	if err != nil {
		return Config{}, err
	}

	corsAllowedOrigins, err := requireEnv("CORS_ALLOWED_ORIGINS")
	if err != nil {
		corsAllowedOrigins = "*"
	}

	logLevel := slog.LevelDebug
	if ll := os.Getenv("LOG_LEVEL"); ll != "" {
		parsedLl, err := strconv.ParseInt(ll, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("invalid LOG_LEVEL %q: %w", ll, err)
		}
		logLevel = slog.Level(parsedLl)
	}

	return Config{
		DBConnectionURI:    dbConnectionURI,
		JwtSecret:          jwtSecret,
		RedisAddr:          redisAddress,
		LogLevel:           logLevel,
		CorsAllowedOrigins: corsAllowedOrigins,
	}, nil
}

func requireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s cannot be empty", key)
	}

	return value, nil
}
