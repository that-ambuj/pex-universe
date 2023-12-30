package routes

import (
	"fmt"
	"math"
	"pex-universe/model"
	"pex-universe/model/product"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func (s *Controller) RegisterCategoryRoutes() {
	v1 := s.Group("/v1")

	v1.Get("/categories", s.categoryGet)
	v1.Get("/categories/:id", s.categoryByIdGet)
	v1.Get("/categories/:id/children", s.categoriesChildren)
}

type CategoriesResp struct {
	Data []product.Category
	model.PageResponse
}

// categoryGet
//
//	@Description	Get List of `Categories`
//	@Tags			products
//	@Produce		json
//	@Param			page	query	int	false	"page number"		default(1)
//	@Param			limit	query	int	false	"limit of results"	default(10)
//	@Success		200		{array}	CategoriesResp
//	@Router			/v1/categories [get]
func (s *Controller) categoryGet(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	categories := []product.Category{}
	count := int64(0)

	err := s.DB.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&categories).
		Count(&count).Error
	if err != nil {
		log.Error(err)
		return err
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return c.JSON(CategoriesResp{
		Data: categories,
		PageResponse: model.PageResponse{
			CurrentPage: page,
			TotalPages:  totalPages,
			Count:       int(count),
		},
	})
}

// categoryByIdGet
//
//	@Description	Get `Category` Info By ID
//	@Tags			products
//	@Produce		json
//	@Param			id						path		int	true	"Category ID"
//	@Success		200						{object}	product.Category
//	@Router			/v1/categories/{id} 	[get]
func (s *Controller) categoryByIdGet(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	category := product.Category{}

	err = s.DB.
		Preload("Children").
		Where(&product.Category{ID: uint(id)}).
		First(&category).Error

	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(404, fmt.Sprintf("Category with ID: %d does not exists", id))
	}

	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&category)
}

// categoriesChildren
//
//	@Description	Get List of `Categories` as Children of a `Category` with ID
//	@Tags			products
//	@Produce		json
//	@Param			id								path		int	true	"Category ID"
//	@Success		200								{object}	product.Category
//	@Router			/v1/categories/{id}/children 	[get]
func (s *Controller) categoriesChildren(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	categories := []product.Category{}

	err = s.DB.
		Where("id IN (?)", s.DB.
			Table("category_children").
			Where("category_id = ?", id).
			Select("child_id")).
		Find(&categories).Error

	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&categories)
}
