package server

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/swagger"
	_ "pex-universe/docs"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/swagger/*", swagger.HandlerDefault)

	s.App.Get("/hello", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)

	s.RegisterAuthRoutes()
}

type Hello struct {
	Message string `json:"message"`
}

// HelloWorldHandler godoc
//
//	@Summary	Hello World
//	@Tags		default
//	@Produce	json
//	@Success	200	{object}	Hello
//	@Router		/hello [get]
func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := map[string]string{
		"message": "Hello World",
	}
	return c.JSON(resp)
}

// healthHandler godoc
//
//	@Summary	Database Health Indicator
//	@Tags		default
//	@Produce	json
//	@Success	200	{object}	Hello
//	@Router		/health [get]
func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}
