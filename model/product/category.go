package product

import "time"

type Category struct {
	ID                      uint `json:"id"`
	Title                   string
	DisplayTitle            *string
	Slug                    string
	FirstDescription        string
	SecondDescription       string
	CustomCategory          bool
	ButtonEnable            *bool
	ButtonTitle             *string
	ButtonContent           *string
	Image                   *string
	Position                *int
	DefaultSorting          *bool
	PriceAttributeEnabled   bool
	PriceAttributeCollapsed bool
	BrandAttributeEnabled   bool
	BrandAttributeCollapsed bool
	SliderID                *uint `json:"-"`
	NumOfSlides             *uint `json:"-"`
	MetaTitle               string
	MetaDescription         string
	Published               bool
	CreatedAt               *time.Time
	UpdatedAt               *time.Time
	Children                []*Category `gorm:"many2many:category_children"`
}
