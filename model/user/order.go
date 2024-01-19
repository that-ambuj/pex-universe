package user

import (
	"pex-universe/model/product"
	"time"
)

type Order struct {
	ID                 uint `json:"id"`
	Email              *string
	IsPickup           bool
	ShipTogether       bool
	AddressModified    bool
	DifferentBilling   bool
	Host               *string `validate:"omitempty,ip"`
	Printed            bool
	OnHold             bool
	EstimatedDelivery  string
	PrintedDate        *time.Time
	BoxType            *string
	LiftgateRequired   bool
	ShippingCost       float32
	TaxAmount          float32
	CouponCode         *string
	CouponAmount       float32
	ResidentialAddress bool
	ShipwordsFetched   bool
	TaxRate            float32
	ExtraTaxAmount     float32
	ExtraTaxRate       float32
	FraudCheckedAt     *time.Time
	VerificationSentAt *time.Time
	PhoneCall          bool
	TaxExempted        bool
	QbStatus           *string
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
	// Foreign Keys
	UserID            *uint `json:"-"`
	StatusID          *uint `json:"-"`
	ShippingAddressID *uint `json:"-"`
	BillingAddressID  *uint `json:"-"`
	ShippingMethodID  *uint `json:"-"`
	TaxID             *uint `json:"-"`
	CouponID          *uint `json:"-"`
	// Joins
	ShippingAddress *Address
	BillingAddress  *Address
	TaxState        *State
	Coupon          *product.Coupon
	ShippingMethod  *product.ShippingMethod
}
