package routes

import (
	"fmt"
	"math"
	"pex-universe/model"
	"pex-universe/model/product"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Controller) RegisterProductRoutes() {
	v1 := s.Group("/v1")

	v1.Get("/products", s.productsGet)
	v1.Get("/products/:id", s.productByIdGet)
}

type ProductsResp struct {
	Data []product.Product
	model.PageResponse
}

// productsGet
//
//	@Router		/v1/products [get]
//	@Tags		products
//	@Produce	json
//	@Success	200			{object}	ProductsResp
//	@Param		category_id	query		int	false	"Category ID"
//	@Param		page		query		int	false	"page number"		default(1)
//	@Param		limit		query		int	false	"limit of results"	default(10)
func (s *Controller) productsGet(c *fiber.Ctx) error {
	categoryId := c.QueryInt("category_id")
	if categoryId < 1 {
		return fiber.NewError(400, "Please set category_id correctly")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	products := []product.Product{}
	count := int64(0)

	err := s.DB.
		Where("id IN (?)", s.DB.
			Table("products_categories").
			Where("category_id = ?", categoryId).
			Select("product_id"),
		).
		Find(&products).
		Count(&count).
		Error

	if err != nil {
		log.Error(err)
		return err
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return c.JSON(&ProductsResp{
		Data: products,
		PageResponse: model.PageResponse{
			CurrentPage: page,
			TotalPages:  totalPages,
			Count:       int(count),
		},
	})
}

// productByIdGet
//
//	@Router		/v1/products/{id} [get]
//	@Tags		products
//	@Produce	json
//	@Success	200	{object}	product.Product
//	@Param		id	query		int	true	"Product ID"
func (s *Controller) productByIdGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	p := product.Product{}

	err = s.DB.
		Preload(clause.Associations).
		Preload("Reviews.Contents").
		Where(&product.Product{ID: uint(id)}).
		First(&p).Error

	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(404,
			fmt.Sprintf("Product with ID: %d does not exists", id))
	}

	if !p.Published {
		return fiber.NewError(403,
			fmt.Sprintf("Product with ID: %d is not allowed to be accessed by the current user.", p.ID))
	}

	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&p)
}
