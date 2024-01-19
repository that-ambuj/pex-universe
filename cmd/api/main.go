package main

import (
	"fmt"
	"log"
	"os"
	"pex-universe/internal/server"
	routes "pex-universe/routes/v1"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

var (
	env = os.Getenv("APP_ENV")

	certfile = os.Getenv("CERT_FILE")
	certkey  = os.Getenv("CERT_KEY")
)

// main
//
//	@title		Pex Universe API
//	@version	1.0
//	@BasePath	/
func main() {
	s := routes.Controller(*server.New())

	s.RegisterRoutes()

	port, _ := strconv.Atoi(os.Getenv("PORT"))

	var err error
	if env == "production" {
		err = s.ListenTLS(":443", certfile, certkey)
	} else {
		err = s.Listen(fmt.Sprintf(":%d", port))
	}

	if err != nil {
		log.Panic(err)
	}
}
