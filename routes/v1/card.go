package routes

import (
	"fmt"
	"pex-universe/model"
	"pex-universe/model/user"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func (s *Controller) RegisterCardRoutes() {
	v1 := s.Group("/v1/profile")

	v1.
		Get("/cards", s.cardsGet).
		Get("/cards/:id", s.cardByIdGet).
		Post("/cards", s.cardsPost).
		Delete("/cards/:id", s.cardsDelete)
}

// cardsGet
//
//	@Router			/v1/profile/cards [get]
//	@Description	Get a List of Saved `Card`s by the current `User`
//	@Tags			cards
//	@Produce		json
//	@Success		200	{array}	user.Card
func (s *Controller) cardsGet(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	cs := []user.Card{}

	err := s.DB.
		Model(&u).
		Preload("PaymentMethod").
		Association("Cards").
		Find(&cs)

	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&cs)
}

// cardByIdGet
//
//	@Router			/v1/profile/cards/{id} [get]
//	@Description	Get a List of Saved `Card`s by the current `User`
//	@Tags			cards
//	@Produce		json
//	@Param			id	path		int	true	"Card ID"
//	@Success		200	{object}	user.Card
func (s *Controller) cardByIdGet(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	if id < 1 {
		return fiber.NewError(400, "Please set the :id route parameter correctly.")
	}

	card := user.Card{}

	err = s.DB.
		Model(&u).
		Preload("PaymentMethod").
		Where(&user.Card{ID: uint(id)}).
		Association("Cards").
		Find(&card)
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&card)
}

// cardsPost
//
//	@Router			/v1/profile/cards [post]
//	@Description	Create a new Credit/Debit Card for User
//	@Tags			cards
//	@Param			request	body	user.CardCreateDto	true	"Card Create Data"
//	@Produce		json
//	@Success		201	{object}	user.Card
func (s *Controller) cardsPost(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	dto := user.CardCreateDto{}

	err := c.BodyParser(&dto)
	if err != nil {
		log.Error(err)
		return err
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		log.Error(err)
		return err
	}

	lastOrder := user.Order{}
	err = s.DB.
		Select("ID").
		Order("created_at DESC").
		Where(&user.Order{UserID: &u.ID}).
		First(&lastOrder).Error

	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusConflict,
			"Please place an order before attempting to save this card.")
	}

	if err != nil {
		log.Error(err)
		return err
	}

	cardNumber := removeWhitespace(dto.CardNumber)

	dto.CardNumberLength = len(cardNumber)
	// Take the last 4 digits of the Card Number to store in the database
	dto.CardNumber = cardNumber[len(cardNumber)-4:]

	newCard := user.Card{
		CardNumber:       dto.CardNumber,
		CardNumberLength: dto.CardNumberLength,
		ExpMonth:         dto.ExpMonth,
		ExpYear:          dto.ExpYear,
		TransactionID:    dto.TransactionID,
		OrderID:          lastOrder.ID,
		CardType:         dto.CardType.String(),
		PaymentMethodID:  dto.PaymentMethod.ID(),
		UserID:           u.ID,
	}

	err = s.DB.Create(&newCard).Error
	if err != nil {
		log.Error(err)
		return err
	}

	s.DB.Preload("PaymentMethod").Find(&newCard)

	return c.Status(201).JSON(&newCard)
}

// cardsDelete
//
//	@Router			/v1/profile/cards/{id} [delete]
//	@Param			id	path	int	true	"Card ID"
//	@Description	Delete a Card by ID
//	@Tags			cards
//	@Produce		json
//	@Success		200	{object}	model.EntityDeletedResponse
func (s *Controller) cardsDelete(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	card := user.Card{
		ID:     uint(id),
		UserID: u.ID,
	}

	res := s.DB.Where(&card).Delete(&card)

	if res.Error != nil {
		log.Error(res.Error)
		return res.Error
	}

	if res.RowsAffected < 1 {
		return fiber.NewError(404, fmt.Sprintf("Card with ID: %d does not exist.", id))
	}

	return c.JSON(model.EntityDeletedResponse{
		RowsAffected: res.RowsAffected,
		Message:      fmt.Sprintf("Card with ID: %d deleted successfully.", id),
	})
}

func removeWhitespace(s string) string {
	number := strings.Builder{}

	for _, r := range s {
		if r == ' ' {
			continue
		}

		number.WriteRune(r)
	}

	return number.String()
}
