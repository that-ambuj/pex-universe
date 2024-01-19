package user

import (
	"time"
)

type Cart struct {
	ID        uint  `json:"id"`
	UserID    *uint `json:"-"`
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Items     []CartProduct
}

func (Cart) TableName() string {
	return "cart"
}

// TODO: add image
type CartProduct struct {
	ID            uint `json:"id"`
	CartID        uint `json:"-"`
	ProductID     int  `json:"product_id"`
	Title         string
	Manufacturer  string
	PartNumber    string
	Qty           int
	MqtyID        *int `json:"mqty_id"` // What are these three?
	Mqty          *int
	MqtyLabel     *string
	Price         float32
	StartingPrice float32 // What is this?
	Weight        float32
	WeightUnits   string
	SavedForLater bool
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

// TODO: Know and use other fields of CartProduct
type CartProductCreateDto struct {
	ProductID     int  `validate:"required" json:"product_id"`
	Qty           int  `validate:"required"`
	SavedForLater bool `validate:"omitempty" default:"false"`
}

type CartProductUpdateDto struct {
	Qty           int  `validate:"omitempty"`
	SavedForLater bool `validate:"omitempty" example:"false"`
}
