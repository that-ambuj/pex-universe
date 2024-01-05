package user

import (
	"time"
)

type (
	Address struct {
		ID             uint `json:"id" gorm:"primaryKey"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		Verified       bool
		FirstName      string `example:"John"`
		LastName       string `example:"Doe"`
		Company        *string
		StreetAddress1 string `gorm:"column:street_address"`
		StreetAddress2 *string
		City           string `example:"Tokyo"`
		Zip            string
		Phone          string `example:"+11349503120"`
		Ext            string
		Email          string `example:"john@example.com"`

		UserID uint `json:"-"`

		StateID   uint `json:"-"`
		CountryID uint `json:"-"`

		State   State
		Country Country
	}

	State struct {
		ID              uint `gorm:"primaryKey" json:"id"`
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
		ID          uint `gorm:"primaryKey" json:"id"`
		Name        string
		DisplayName string
		PpCode      string
		Position    int64
		Locked      bool `json:"-"`
	}
)

type (
	AddressCreateDto struct {
		FirstName      string `example:"John"`
		LastName       string `example:"Doe"`
		Company        *string
		StreetAddress1 string `gorm:"column:street_address"`
		StreetAddress2 *string
		City           string `example:"Los Angeles"`
		Zip            string `validate:"numeric"`
		Phone          string `validate:"e164" example:"+12380941034"`
		Ext            string
		Email          string `validate:"email" example:"john@example.com"`

		UserID    uint `json:"-"`
		StateID   uint `validate:"min=1"`
		CountryID uint `validate:"min=1"`
	}

	AddressUpdateDto struct {
		FirstName      *string `example:"Jane"`
		LastName       *string `example:"Doe"`
		Company        *string
		StreetAddress1 *string `gorm:"column:street_address"`
		StreetAddress2 *string
		City           *string `example:"New York"`
		Zip            *string `validate:"omitempty,numeric"`
		Phone          *string `validate:"omitempty,e164" example:"+12380941034"`
		Ext            *string
		Email          *string `validate:"omitempty,email" example:"john@example.com"`
		StateId        uint64  `validate:"omitempty,min=1"`
		CountryId      uint64  `validate:"omitempty,min=1"`
	}
)
