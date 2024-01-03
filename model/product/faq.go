package product

import "time"

type ProductFaq struct {
	ID                 uint `json:"id"`
	Name               *string
	Email              *string `json:"-"` // Hide User Email for Privacy
	Question           string
	Answer             *string
	Emailed            bool
	Answered           *time.Time
	Published          bool
	ForCustomerService bool
	UpdatedAt          *time.Time
	CreatedAt          *time.Time
	// Foreign Keys
	ProductID uint  `json:"-"`
	UserID    *uint `json:"-"`
}
