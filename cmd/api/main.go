package main

import (
	"fmt"
	"os"
	"pex-universe/internal/server"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

var (
	env = os.Getenv("APP_ENV")
)

// @title		Pex Universe API
// @version		1.0
// @BasePath	/
func main() {
	server := server.New()

	if env != "production" {
		server.App.Use(cors.New())
	}

	server.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic("cannot start server")
	}
}
