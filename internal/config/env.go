package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := *new(error)
	env := os.Getenv("APP_ENV")

	if env == "" {
		env = "development"
	}

	godotenv.Load(".env." + env)
	if err != nil {
	}

	godotenv.Load(".env")
	if err != nil {
	}
}
