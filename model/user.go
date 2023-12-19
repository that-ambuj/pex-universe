package model

import (
	"database/sql"
	"pex-universe/types"
	"time"
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

	Address struct {
		Id             uint64 `json:"-"`
		CreatedAt      *time.Time
		UpdatedAt      *time.Time
		Verified       bool
		FirstName      string
		LastName       string
		Company        types.NullString `swaggertype:"string"`
		StreetAddress1 string           `db:"street_address"`
		StreetAddress2 types.NullString `swaggertype:"string"`
		City           string
		Zip            string
		Phone          string
		Ext            string
		Email          string
		// Foreign Keys
		CountryId sql.NullInt64 `db:"country_id" json:"-"`
		StateId   sql.NullInt64 `db:"state_id" json:"-"`
		UserId    sql.NullInt64 `db:"user_id" json:"-"`
		State     *State
		Country   *Country
	}

	State struct {
		Name            string
		FullName        string
		Tax             float32
		Info            string
		Locked          bool `json:"-"`
		TaxEnabled      bool `json:"-"`
		FixedRate       bool `json:"-"`
		ShippingTaxable bool `json:"-"`
	}

	Country struct {
		Id          string `json:"-"`
		Name        string
		DisplayName string
		PpCode      string
		Position    int64
		Locked      bool `db:"locked" json:"-"`
	}
)
