package routes

import (
	"fmt"
	"pex-universe/internal/errors"
	"pex-universe/model"
	"pex-universe/model/product"
	"pex-universe/model/user"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// Authed
// TODO: checkout
// TODO: Cancel Order (DELETE)

func (s *Controller) RegisterOrderOpenRoutes() {
	v1 := s.Group("/v1")

	v1.
		Get("/cart/items", s.cartGet).
		Get("/cart/items/:id", s.cartItemByIdGet).
		Post("/cart/items", s.cartPost).
		Delete("/cart/items/:id", s.cartItemByIdDelete).
		Put("/cart/items/:id", s.cartItemByIdPut)
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

// cartItemByIdGet
//
//	@Router			/v1/cart/{id} [get]
//	@Success		200	{object}	user.CartProduct
//	@Tags			cart
//	@Description	Get a Items in the `Cart` with given `ID`
//	@Param			id	path		int	true	"Cart ID"
func (s *Controller) cartItemByIdGet(c *fiber.Ctx) error {
	cartId, _ := strconv.Atoi(c.Cookies("cart_id"))

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	item := user.CartProduct{}

	err = s.DB.
		Where(&user.CartProduct{
			CartID: uint(cartId),
			ID:     uint(id),
		}).
		Find(&item).Error

	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&item)
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

// cartItemByIdDelete
//
//	@Router			/v1/cart/{id} [delete]
//	@Success		200 {object} model.EntityDeletedResponse
//	@Tags			cart
//	@Description	Remove the Item with given `ID` from the `Cart`
//	@Param			id	path		int	true	"Cart Item ID"
func (s *Controller) cartItemByIdDelete(c *fiber.Ctx) error {
	cartId, _ := strconv.Atoi(c.Cookies("cart_id"))

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	item := user.CartProduct{
		ID:     uint(id),
		CartID: uint(cartId),
	}

	res := s.DB.Where(&item).Delete(&item)

	if res.Error != nil {
		log.Error(res.Error)
		return res.Error
	}

	if res.RowsAffected < 1 {
		return fiber.NewError(404, fmt.Sprintf("Cart Item with ID: %d does not exist.", id))
	}

	return c.JSON(model.EntityDeletedResponse{
		RowsAffected: res.RowsAffected,
		Message:      "Cart Item deleted successfully.",
	})
}

// cartItemByIdPut
//
//	@Router			/v1/cart/{id} [put]
//	@Tags			cart
//	@Description	Update the Item with given `ID` from the `Cart`
//	@Param			id	path		int	true	"Cart Item ID"
//	@Param			request	body		user.CartProductUpdateDto	true	"Cart Product Update Data"
//	@Success		200 {object} model.EntityDeletedResponse
func (s *Controller) cartItemByIdPut(c *fiber.Ctx) error {
	cartId, _ := strconv.Atoi(c.Cookies("cart_id"))

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	dto := user.CartProductUpdateDto{}

	err = c.BodyParser(&dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		log.Error(err)
		return err
	}

	item := user.CartProduct{
		ID:     uint(id),
		CartID: uint(cartId),
	}

	err = s.DB.Model(&item).Updates(&dto).Error
	if err != nil {
		return err
	}

	return c.JSON(&item)
}

func (s *Controller) RegisterOrderAuthorisedRoutes() {

}
