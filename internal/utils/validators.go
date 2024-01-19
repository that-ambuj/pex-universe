package utils

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

func ValidatePaymentMethod(f validator.FieldLevel) bool {
	switch f.Field().String() {
	case "PAYPAL_CREDIT",
		"PAYPAL_EXPRESS",
		"STRIPE_CREDIT":
		return true
	default:
		return false
	}
}

func ValidateMonth(f validator.FieldLevel) bool {
	n, _ := strconv.Atoi(f.Field().String())

	if 1 <= n && n <= 12 {
		return true
	}

	return false
}

func ValidateYear(f validator.FieldLevel) bool {
	str := f.Field().String()
	n, _ := strconv.Atoi(str)

	if 2000 <= n && n <= 2100 {
		return true
	}

	return false
}

func ValidateCardType(f validator.FieldLevel) bool {
	switch f.Field().String() {
	case
		"American Express",
		"Visa",
		"Mastercard",
		"Discover":
		return true
	default:
		return false
	}
}
