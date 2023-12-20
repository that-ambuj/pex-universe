package server

import (
	"fmt"
	"strings"

	"pex-universe/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/iancoleman/strcase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/jmoiron/sqlx"
	json "github.com/mixcode/golib-json-snake"
)

type FiberServer struct {
	*fiber.App
	db    *sqlx.DB
	v     *validator.Validate
	store *session.Store
}

type ValidationErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type ErrorResp struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	log.Error(err)

	switch t := err.(type) {
	case *fiber.Error:
		return c.Status(t.Code).JSON(ErrorResp{
			Success: false,
			Status:  t.Code,
			Message: t.Error(),
		})
	default:
		return c.Status(500).JSON(ErrorResp{
			Success: false,
			Status:  500,
			Message: "Something Unexpected Happened",
		})
	}

}

func New() *FiberServer {
	sessionConfig := session.ConfigDefault

	sessionConfig.CookieHTTPOnly = true
	sessionConfig.CookieSecure = true
	sessionConfig.CookieSameSite = "Strict"

	// TODO: Set CookieDomain to deployed URL for security
	//sessionConfig.CookieDomain = ""

	// TODO: Use Redis for storage
	sessionConfig.Storage = sqlite3.New()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
		JSONEncoder:  json.MarshalSnakeCase,
		JSONDecoder:  json.UnmarshalSnakeCase,
	})

	app.Hooks().OnRoute(func(r fiber.Route) error {
		if r.Method != "HEAD" && r.Path != "/" {
			fmt.Printf("Mapped Route: [%s] %s\n", r.Method, r.Path)
		}

		return nil
	})

	db := database.New()
	db.MapperFunc(strcase.ToSnake)

	server := &FiberServer{
		App:   app,
		db:    db,
		v:     validator.New(),
		store: session.New(sessionConfig),
	}

	return server
}

func (s *FiberServer) ValidateStruct(data interface{}) error {
	errs := s.v.Struct(data)

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
