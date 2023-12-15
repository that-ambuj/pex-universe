package model

import (
	"database/sql"
	"time"
)

type (
	User struct {
		Id            int            `db:"id" json:"id"`
		Name          string         `db:"name" json:"name"`
		Email         string         `db:"email" json:"email"`
		Password      string         `db:"password" json:"-"`
		RememberToken sql.NullString `db:"remember_token" json:"-"`
		CreatedAt     *time.Time     `db:"created_at" json:"created_at"`
		UpdatedAt     *time.Time     `db:"updated_at" json:"updated_at"`
	}

	UserLoginDto struct {
		Email    string `validate:"required,email" example:"john@example.com"`
		Password string `validate:"required,min=8" example:"avEryStrongPass@93"`
	}

	UserSignUpDto struct {
		Name     string `validate:"required" example:"John Doe"`
		Email    string `validate:"required,email" example:"john@example.com"`
		Password string `validate:"required,min=8" example:"avEryStrongPass@93"`
	}
)
