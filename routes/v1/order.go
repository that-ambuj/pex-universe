package routes

import (
	"pex-universe/internal/errors"
	"pex-universe/model/product"
	"pex-universe/model/user"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func (s *Controller) RegisterOrderOpenRoutes() {
	v1 := s.Group("/v1")

	v1.
		Get("/cart", s.cartGet).
		Post("/cart", s.cartPost)
}

// cartGet
//
//	@Router			/v1/cart [get]
//	@Success		200	{array}	user.CartProduct
//	@Tags			cart
//	@Description	Get a List of Items in the `Cart`
func (s *Controller) cartGet(c *fiber.Ctx) error {
	// TODO: Pagination
	cartId, _ := strconv.Atoi(c.Cookies("cart_id"))

	items := []user.CartProduct{}

	err := s.DB.
		Model(&user.Cart{ID: uint(cartId)}).
		Association("Items").
		Find(&items)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&items)
}

// cartPost
//
//	@Router			/v1/cart [post]
//	@Param			request	body		user.CartProductCreateDto	true	"Cart Product Data"
//	@Success		201		{object}	user.CartProduct
//	@Produce		json
//	@Tags			cart
//	@Description	Add a new Item to the cart
func (s *Controller) cartPost(c *fiber.Ctx) error {
	cartId, _ := strconv.Atoi(c.Cookies("cart_id"))

	dto := user.CartProductCreateDto{}

	err := c.BodyParser(&dto)
	if err != nil {
		log.Error(err)
		return errors.BadRequestErr(err)
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		return err
	}

	p := product.Product{ID: uint(dto.ProductID)}

	err = s.DB.
		Preload("Manufacturer").
		First(&p).Error
	if err == gorm.ErrRecordNotFound {
		return errors.NotFoundEntity("Product", "ID",
			strconv.Itoa(dto.ProductID))
	}

	if err != nil {
		log.Error(err)
		return err
	}

	cp := user.CartProduct{
		CartID:        uint(cartId),
		ProductID:     dto.ProductID,
		Qty:           dto.Qty,
		Title:         p.Title,
		Manufacturer:  p.Manufacturer.Name,
		PartNumber:    p.PartNumber,
		Price:         p.DiscountPrice,
		SavedForLater: dto.SavedForLater,
	}

	err = s.DB.Create(&cp).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return c.Status(201).JSON(&cp)
}

func (s *Controller) RegisterOrderAuthorisedRoutes() {

}
