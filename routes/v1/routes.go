package routes

import (
	"fmt"
	"os/user"
	"pex-universe/internal/database"
	"pex-universe/internal/server"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type Controller server.FiberServer

func (c *Controller) RegisterRoutes() {
	c.Get("/swagger/*", swagger.HandlerDefault)

	c.Get("/hello", c.HelloWorldHandler)
	c.Get("/health", c.healthHandler)

	c.RegisterAuthRoutes()

	c.Use("/v1/*", c.UserAuthMiddleware)

	c.RegisterProfileRoutes()

}

func (s *Controller) UserAuthMiddleware(c *fiber.Ctx) error {
	var err error

	sess, err := s.Store.Get(c)
	if err != nil {
		return err
	}

	token := sess.ID()
	if token == "" {
		return fiber.NewError(fiber.ErrUnauthorized.Code, "Invalid User Session")
	}

	u := new(user.User)

	err = s.DB.Get(u, `SELECT * FROM users WHERE remember_token = ?`, token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User Token Expired")
	}

	sess.SetExpiry(24 * time.Hour)
	sess.Save()

	c.Locals("user", u)
	return c.Next()
}

func (s *Controller) ValidateStruct(data interface{}) error {
	errs := s.V.Struct(data)

	if errs != nil {
		errMsgs := make([]string, 0)

		for _, err := range errs.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s' (param: '%s')",
				err.Field(),
				err.Value(),
				err.Tag(),
				err.Param(),
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}

	return nil
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
func (s *Controller) HelloWorldHandler(c *fiber.Ctx) error {
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
func (s *Controller) healthHandler(c *fiber.Ctx) error {
	return c.JSON(database.SqlxHealth(s.DB))
}
