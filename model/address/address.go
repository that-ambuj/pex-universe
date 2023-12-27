package address

import (
	"time"
)

type (
	Address struct {
		ID             uint `json:"id" gorm:"primaryKey"`
		CreatedAt      time.Time
		UpdatedAt      time.Time
		Verified       bool
		FirstName      string
		LastName       string
		Company        *string
		StreetAddress1 string `db:"street_address" gorm:"column:street_address"`
		StreetAddress2 *string
		City           string
		Zip            string
		Phone          string
		Ext            string
		Email          string

		UserID uint `json:"-"`

		StateID   uint `json:"-"`
		CountryID uint `json:"-"`

		State   State
		Country Country
	}

	State struct {
		ID              uint `gorm:"primaryKey" json:"-"`
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
		ID          uint `gorm:"primaryKey" json:"-"`
		Name        string
		DisplayName string
		PpCode      string
		Position    int64
		Locked      bool `json:"-"`
	}
)

type (
	AddressCreateDto struct {
		FirstName      string
		LastName       string
		Company        *string
		StreetAddress1 string `db:"street_address" gorm:"column:street_address"`
		StreetAddress2 *string
		City           string
		Zip            string
		Phone          string
		Ext            string
		Email          string

		UserID    uint `json:"-"`
		StateID   uint
		CountryID uint
	}

	AddressUpdateDto struct {
		FirstName      *string
		LastName       *string
		Company        *string
		StreetAddress1 *string `db:"street_address" gorm:"column:street_address"`
		StreetAddress2 *string
		City           *string
		Zip            *string
		Phone          *string
		Ext            *string
		Email          *string
		StateId        uint64
		CountryId      uint64
	}
)
