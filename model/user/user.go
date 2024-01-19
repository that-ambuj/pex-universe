package user

import (
	"time"
)

type (
	User struct {
		ID              uint `json:"id" gorm:"primaryKey"`
		Name            *string
		Username        string
		Email           string
		CustomerID      string `json:"customer_id"`
		Location        *string
		Password        string `json:"-"`
		LastLoggedIn    *time.Time
		PasswordResetAt *time.Time
		LastLockoutAt   *time.Time
		RememberToken   *string `json:"-"`
		CreatedAt       *time.Time
		UpdatedAt       *time.Time
		RetAdLastSent   time.Time `json:"-" gorm:"column:RetAdLastSent"`
		// Joins
		Addresses []Address
		Cards     []Card
		Carts     []Cart
	}

	UserLoginDto struct {
		Email    string `validate:"required,email" example:"john@example.com"`
		Password string `validate:"required,min=8" example:"avEryStrongPass@93"`
	}

	UserSignUpDto struct {
		Name     *string `example:"John Doe"`
		Username string  `validate:"required,excludes=' '"`
		Email    string  `validate:"required,email" example:"john@example.com"`
		Password string  `validate:"required,min=8" example:"avEryStrongPass@93"`
	}
)

func (User) TableName() string {
	return "site_users"
}
