package main

import (
	"fmt"
	"log"
	"os"
	"pex-universe/internal/server"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/joho/godotenv/autoload"
)

var (
	env = os.Getenv("APP_ENV")

	certfile = os.Getenv("CERT_FILE")
	certkey  = os.Getenv("CERT_KEY")
)

//	@title		Pex Universe API
//	@version	1.0
//	@BasePath	/
func main() {
	server := server.New()
	server.App.Use(cors.New())

	server.RegisterFiberRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	var err error
	if env == "production" {
		err = server.ListenTLS(":443", certfile, certkey)
	} else {
		err = server.Listen(fmt.Sprintf(":%d", port))
	}

	if err != nil {
		log.Panic(err)
	}
}
