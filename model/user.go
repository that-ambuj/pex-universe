package model

import (
	"database/sql"
)

type (
	User struct {
		Name          string         `db:"name" json:"name"`
		Email         string         `db:"email" json:"email"`
		Password      string         `db:"password" json:"-"`
		RememberToken sql.NullString `db:"remember_token" json:"-"`
		Addresses     []*Address     `json:"addresses"`
		AutoInc
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
