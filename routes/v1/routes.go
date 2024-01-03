package routes

import (
	"fmt"
	"strings"
	"time"

	"pex-universe/internal/database"
	"pex-universe/internal/server"
	"pex-universe/model/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/swagger"
	"github.com/iancoleman/strcase"

	_ "pex-universe/docs"
)

type Controller server.FiberServer

func (s *Controller) RegisterRoutes() {
	s.Get("/swagger/*", swagger.New(swagger.Config{
		TryItOutEnabled: true,
	}))

	s.Get("/hello", s.HelloWorldHandler)
	s.Get("/health", s.healthHandler)

	s.RegisterAuthRoutes()
	s.RegisterUtilRoutes()
	s.RegisterHomeRoutes()
	s.RegisterCategoryRoutes()
	s.RegisterProductRoutes()
	s.RegisterOrderOpenRoutes()

	s.Use("/v1/*", s.UserAuthMiddleware)

	s.RegisterProfileRoutes()
	s.RegisterOrderAuthorisedRoutes()
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

	u := user.User{}

	err = s.DB.Where("remember_token = ?", token).First(&u).Error
	if err != nil {
		log.Error(err)
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
			var (
				field = strcase.ToSnake(err.Field())
				tag   = err.Tag()
				val   = err.Value()
				param = err.Param()
			)

			if param != "" {
				tag = tag + ": " + param
			}

			errMsgs = append(errMsgs, fmt.Sprintf(
				"'%s' has failed the constraint: %s (value: '%v')",
				field,
				tag,
				val,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, ", "),
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
	db, err := s.DB.DB()
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(database.SqlxHealth(db))
}
