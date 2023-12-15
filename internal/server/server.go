package server

import (
	"fmt"
	"pex-universe/internal/database"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
	db database.Service
	v  *validator.Validate
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
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		}),
		db: database.New(),
		v:  validator.New(),
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
				err.StructField(),
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
