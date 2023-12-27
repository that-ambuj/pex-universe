package server

import (
	"fmt"
	"os"

	"pex-universe/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/iancoleman/strcase"

	"gorm.io/gorm"
	"gorm.io/gorm/mysql"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/storage/sqlite3/v2"
	"github.com/jmoiron/sqlx"
	json "github.com/mixcode/golib-json-snake"
)

type FiberServer struct {
	*fiber.App
	OldDB *sqlx.DB
	DB    *gorm.DB
	V     *validator.Validate
	Store *session.Store
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
	// sessionConfig.CookieDomain = ""

	// TODO: Use Redis for storage
	sessionConfig.Storage = sqlite3.New()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
		JSONEncoder:  json.MarshalSnakeCase,
		JSONDecoder:  json.UnmarshalSnakeCase,
	})

	env := os.Getenv("APP_ENV")

	if env != "test" {
		app.Hooks().OnRoute(func(r fiber.Route) error {
			fmt.Printf("Mapped Route: [%s] %s\n", r.Method, r.Path)

			return nil
		})
	}

	db := database.New()
	db.MapperFunc(strcase.ToSnake)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db.DB,
	}))
	if err != nil {
		log.Fatal(err)
	}

	server := &FiberServer{
		App:   app,
		OldDB: db,
		DB:    gormDB,
		V:     validator.New(),
		Store: session.New(sessionConfig),
	}

	return server
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"message": "Hello World",
	})
}
