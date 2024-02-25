package server

import (
	"pex-universe/internal/database"
	"pex-universe/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/session"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/storage/sqlite3/v2"
	json "github.com/mixcode/golib-json-snake"
)

type FiberServer struct {
	*fiber.App
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
	storage := sqlite3.New()
	sessionConfig.Storage = storage

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: true,
		ErrorHandler:      ErrorHandler,
		JSONEncoder:       json.MarshalSnakeCase,
		JSONDecoder:       json.UnmarshalSnakeCase,
	})

	app.Use(fiberLogger.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(limiter.New(limiter.Config{
		SkipFailedRequests: true,
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Storage:    storage,
		Max:        20,
		Expiration: 30 * time.Second,
	}))

	app.Get("/metrics", monitor.New())

	db := database.New()

	v := validator.New(validator.WithRequiredStructEnabled())
	//nolint:errcheck
	v.RegisterValidation("payment-method", utils.ValidatePaymentMethod)
	//nolint:errcheck
	v.RegisterValidation("month", utils.ValidateMonth)
	//nolint:errcheck
	v.RegisterValidation("year", utils.ValidateYear)
	//nolint:errcheck
	v.RegisterValidation("card-type", utils.ValidateCardType)

	server := &FiberServer{
		App:   app,
		DB:    db,
		V:     v,
		Store: session.New(sessionConfig),
	}

	return server
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		"message": "Hello World",
	})
}
