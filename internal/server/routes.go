package server

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2/middleware/session"

	_ "pex-universe/docs"
	"pex-universe/internal/database"
	"pex-universe/model"

	"github.com/gofiber/swagger"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.App.Get("/swagger/*", swagger.HandlerDefault)

	s.App.Get("/hello", s.HelloWorldHandler)
	s.App.Get("/health", s.healthHandler)

	s.RegisterAuthRoutes()

	s.App.Use("/v1/*", s.UserAuthMiddleware)

	s.RegisterProfileRoutes()
}

func (s *FiberServer) UserAuthMiddleware(c *fiber.Ctx) error {
	var err error

	sess, err := s.store.Get(c)
	if err != nil {
		return err
	}

	token := sess.ID()
	if token == "" {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Invalid User Session")
	}

	user := new(model.User)

	err = s.db.Get(user, `SELECT * FROM users WHERE remember_token = ?`, token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User Token Expired")
	}

	c.Locals("user", user)
	return c.Next()
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
	return c.JSON(database.SqlxHealth(s.db))
}
