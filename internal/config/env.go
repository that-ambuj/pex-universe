package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	//nolint:errcheck
	godotenv.Load(".env." + env)
	//nolint:errcheck
	godotenv.Load(".env")
}
