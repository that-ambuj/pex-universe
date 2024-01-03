package product

import "time"

type ProductReview struct {
	ID                uint `json:"id"`
	Verified          bool
	Rating            int
	Published         bool
	FromEmail         bool
	ReasonNotApproved *string
	CreatedAt         *time.Time
	UpdatedAt         *time.Time
	// Foreign Keys
	ProductID uint  `json:"-"`
	UserID    *uint `json:"-"`
	OrderID   *uint `json:"-"`
	// Joins
	Contents ProductReviewContent `gorm:"foreignKey:ReviewID"`
}

type ProductReviewContent struct {
	ID           uint    `json:"id"`
	ReviewID     uint    `json:"-"`
	Email        *string `json:"-"` // Hide Email For Privacy
	Name         *string
	Location     *string
	Title        string
	Comment      string
	OriginalText *string
	Reply        *string
	CreatedAt    *time.Time
	UpdatedAt    *time.Time
}

func (ProductReviewContent) TableName() string {
	return "product_reviews_contents"
}
