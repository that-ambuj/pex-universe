package product

import (
	"log"
	"net/url"
	"time"

	"gorm.io/gorm"
)

type Manufacturer struct {
	ID        uint `json:"id"`
	Name      string
	Logo      string
	Slug      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

var manufacturerLogoPrefix = "https://pexuniverse.com/uploads/manufacturers/"

func (m *Manufacturer) AfterFind(tx *gorm.DB) error {

	if m.Logo == "" {
		return nil
	}

	u, err := url.Parse(manufacturerLogoPrefix + m.Logo)
	if err != nil {
		log.Println(err)
		return err
	}

	m.Logo = u.String()

	return nil
}
