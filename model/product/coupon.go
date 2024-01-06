package product

import "time"

type Coupon struct {
	ID               uint `json:"id"`
	Name             string
	QualifyingAmount float32
	Type             int
	Amount           float32
	Code             string
	MaxUses          int
	Expire           *time.Time
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
}
