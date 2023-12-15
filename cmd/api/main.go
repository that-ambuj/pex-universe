package main

import (
	"fmt"
	"os"
	"pex-universe/internal/server"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// @title		Pex Universe API
// @version	1.0
// @host		localhost:8080
// @BasePath	/
func main() {
	server := server.New()

	server.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic("cannot start server")
	}
}
