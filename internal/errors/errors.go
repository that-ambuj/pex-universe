package errors

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func NotFoundEntity(entityName string,
	field string, value string) *fiber.Error {
	return &fiber.Error{
		Code: fiber.StatusNotFound,
		Message: fmt.Sprintf("%s with %s: '%s' does not exist",
			entityName, field, value),
	}
}

func BadRequestErr(err error) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusBadRequest,
		Message: err.Error(),
	}
}

func BadRequestMsg(msg string) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusBadRequest,
		Message: msg,
	}
}

func UnauthorizedMsg(msg string) *fiber.Error {
	return &fiber.Error{
		Code:    fiber.StatusUnauthorized,
		Message: msg,
	}
}
