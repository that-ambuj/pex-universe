package user

import (
	"time"
)

type (
	User struct {
		ID            uint `json:"id" gorm:"primaryKey"`
		Name          string
		Email         string
		Password      string  `json:"-"`
		RememberToken *string `json:"-"`
		CreatedAt     time.Time
		UpdatedAt     time.Time
		// Joins
		Addresses []Address
		Cards     []UserSavedCard
	}

	UserLoginDto struct {
		Email    string `validate:"required,email" example:"john@example.com"`
		Password string `validate:"required,min=8" example:"Very Strong Password"`
	}

	UserSignUpDto struct {
		Name     string `validate:"required" example:"John Doe"`
		Email    string `validate:"required,email" example:"john@example.com"`
		Password string `validate:"required,min=8" example:"avEryStrongPass@93"`
	}
)
