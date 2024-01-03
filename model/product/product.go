package product

import (
	"log"
	"net/url"
	"strings"
	"text/template"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID                      uint `json:"id"`
	Title                   string
	Slug                    string
	HideManufacturer        bool
	Description             *string
	Features                *string
	TechSpecs               *string
	CrossRef                *string
	Application             *string
	Documentation           *string
	Warranty                *string
	PartNumber              string
	Model                   *string
	Upc                     *string
	SellPrice               float32
	DiscountPrice           float32
	DiscountStartDate       *time.Time
	DiscountDeadline        *time.Time
	PickupPrice             float32
	PriceRangeId            int
	Weight                  *float32
	WeightUnits             string
	Length                  float32
	Width                   float32
	Height                  float32
	LwhUnit                 string
	ListPosition            int
	UsuallyShips            string
	FreeShipping            bool
	StockQuantity           int
	TemporaryUnavailable    bool
	CustomLabel             *string
	MetaTitle               string
	MetaDescription         string
	NotLeadFree             bool
	SpecOrder               bool
	Fluid                   bool
	Discontinued            bool
	Replacement             *string
	DiscontinuedReplacement *string
	Published               bool
	FreightOnly             bool
	PickupOnly              bool
	SellQty                 uint
	ShelfIdExtra            *string
	PossibleFraud           bool
	ReviewGroupLabel        *string
	FaqGroupLabel           *string
	ShowMapPrice            bool
	DeclareValue            bool
	InStockAtSupplier       bool
	CaPropWarning           bool
	ShipDedicatedBox        *uint
	ShipSeparately          bool
	ShipSelfPackaging       bool
	Note                    *string
	MadeInUsa               bool
	PickupDiscount          float32
	PricePerFoot            float32
	NonPickup               bool
	CreatedAt               *time.Time
	UpdatedAt               *time.Time
	// Foreign Key
	ManufacturerID *uint   `json:"-"`
	ShelfID        *string `json:"-"`
	ReviewGroupID  *uint   `json:"-"`
	FaqGroupId     *uint   `json:"-"`
	// Joins
	Manufacturer    *Manufacturer
	Images          []ProductImage
	Reviews         []ProductReview
	Faqs            []ProductFaq
	RelatedProducts []Product `gorm:"many2many:product_related;"`
}

type ProductImage struct {
	ID        uint `json:"id"`
	ProductID uint `json:"-"`
	Src       string
	Position  int
	CreatedAt *time.Time
	UpdatedAt *time.Time
	Delete    bool `json:"-"`
}

var productImageTemp = "https://pexuniverse.com/uploads/products/{{- .PartNumber -}}/images/458x458/{{- .Image -}}"

type imageData struct {
	PartNumber string
	Image      string
}

func (p *Product) AfterFind(tx *gorm.DB) error {

	if p.Images == nil {
		return nil
	}

	for idx, image := range p.Images {

		if image.Delete {
			// Delete the current image
			p.Images = append(p.Images[:idx], p.Images[idx+1:]...)
			continue
		}

		if image.Src == "" {
			continue
		}

		info := imageData{p.PartNumber, image.Src}
		imgStr := strings.Builder{}

		tmpl := template.Must(template.New("imageUrl").Parse(productImageTemp))

		err := tmpl.Execute(&imgStr, info)
		if err != nil {
			log.Println(err)
			return err
		}

		u, err := url.Parse(imgStr.String())
		if err != nil {
			log.Println(err)
			return err
		}

		p.Images[idx].Src = u.String()
	}

	return nil
}
