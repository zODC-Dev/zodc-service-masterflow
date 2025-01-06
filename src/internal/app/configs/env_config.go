package configs

import (
	"log/slog"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type envConfig struct {
	SERVER_ADDRESS       string `env:"SERVER_ADDRESS,required" envDefault:"localhost:8080"`
	DATABASE_POSTGRE_DSN string `env:"DATABASE_POSTGRE_DSN,required"`
	DATABASE_REDIS_URL   string `env:"DATABASE_REDIS_URL,required"`
}

// Export global
var Env envConfig

func init() {
	//Load .env file
	if err := godotenv.Load(); err != nil {
		slog.Error("Warning: Error loading .env file", slog.Any("error", err))
	}

	//Check .env file
	if err := env.Parse(&Env); err != nil {
		slog.Error("Failed to parse environment variables", slog.Any("error", err))
		os.Exit(1)
	}
}
