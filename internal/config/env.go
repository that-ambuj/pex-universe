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

	if env != "test" {
		err = godotenv.Load(".env.local")
		if err != nil {
			panic(err)
		}
	}

	godotenv.Load(".env." + env)
	if err != nil {
		panic(err)
	}

	godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
}
