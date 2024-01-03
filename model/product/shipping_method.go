package product

import (
	"time"
)

type ShippingMethod struct {
	ID                  uint `json:"id"`
	Title               string
	Position            int
	Active              bool
	FreeShippingApplies bool
	Class               string
	Method              string
	ServiceDays         string
	ShipDays            string
	Cutoff              int
	Discount            int
	FreightDiscount     int
	MaxDiscount         int
	CreatedAt           *time.Time
	UpdatedAt           *time.Time
}
