package user

import (
	"time"
)

type Card struct {
	ID               uint `json:"id"`
	CardNumber       string
	CardNumberLength int
	CardType         string
	ExpMonth         string
	ExpYear          string
	TransactionID    *string `json:"transaction_id"`
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	// Foreign Keys
	UserID          uint `json:"-"`
	OrderID         uint `json:"-"`
	PaymentMethodID uint `json:"-"`
	// Joins
	PaymentMethod PaymentMethod
}

func (Card) TableName() string {
	return "user_saved_cards"
}

type PaymentMethod struct {
	ID          uint `json:"-"`
	Name        string
	DisplayName string
	Position    int
	Provider    string
	Method      string
	Locked      bool
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
}

type PaymentMethodName string

const (
	PaypalCreditCard PaymentMethodName = "PAYPAL_CREDIT"
	PaypalExpress    PaymentMethodName = "PAYPAL_EXPRESS"
	StripeCreditCard PaymentMethodName = "STRIPE_CREDIT"
)

func (p PaymentMethodName) ID() uint {
	switch p {
	case PaypalCreditCard:
		return 1
	case PaypalExpress:
		return 2
	case StripeCreditCard:
		return 3
	default:
		return 0
	}
}

type CardType string

const (
	Amex       CardType = "American Express"
	Visa       CardType = "Visa"
	MasterCard CardType = "Mastercard"
	Discover   CardType = "Discover"
)

func (c CardType) String() string {
	switch c {
	case Amex:
		return "American Express"
	case Visa:
		return "Visa"
	case MasterCard:
		return "Mastercard"
	case Discover:
		return "Discover"
	}

	return "unknown"
}

type CardCreateDto struct {
	CardNumber       string            `validate:"credit_card" example:"4242424242424242"`
	CardType         CardType          `validate:"card-type"`
	ExpMonth         string            `validate:"numeric,month" example:"06"`
	ExpYear          string            `validate:"numeric,year" example:"2026"`
	PaymentMethod    PaymentMethodName `validate:"payment-method"`
	CardNumberLength int               `json:"-"`
	TransactionID    *string           `json:"transaction_id"`
}
