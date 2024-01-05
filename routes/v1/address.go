package routes

import (
	"fmt"
	"math"
	"strconv"

	"pex-universe/model"
	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func (s *Controller) RegisterAddressRoutes() {
	v1 := s.Group("/v1/profile")

	v1.
		Get("/addresses", s.addressGet).
		Get("/addresses/:id", s.addressByIdGet).
		Post("/addresses", s.addressPost).
		Put("/addresses/:id", s.addressByIdPut).
		Delete("/addresses/:id", s.addressByIdDelete)
}

// addressGet
//
//	@Description	Get List of `Address`es for the current `User`
//	@Tags			addresses
//	@Produce		json
//	@Param			page	query		int	false	"page number"		default(1)
//	@Param			limit	query		int	false	"limit of results"	default(10)
//	@Success		200		{object}	AddressesResponse
//	@Router			/v1/profile/addresses [get]
func (s *Controller) addressGet(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	addrs := []user.Address{}
	count := int64(0)

	err := s.DB.
		Joins("State").
		Joins("Country").
		Where(&user.Address{UserID: u.ID}).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&addrs).
		Count(&count).Error

	if err != nil {
		log.Error(err)
		return err
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return c.JSON(AddressesResponse{
		Data: addrs,
		PageResponse: model.PageResponse{
			CurrentPage: page,
			TotalPages:  totalPages,
			Count:       int(count),
		},
	})
}

// addressByIdGet
//
//	@Description	Get `Address` Info By ID
//	@Tags			addresses
//	@Produce		json
//	@Param			id							path		int	true	"Address ID"
//	@Success		200							{object}	user.Address
//	@Router			/v1/profile/addresses/{id} 	[get]
func (s *Controller) addressByIdGet(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	addr := user.Address{}

	err = s.DB.
		Joins("State").
		Joins("Country").
		Where(&user.Address{
			ID:     uint(id),
			UserID: u.ID,
		}).
		First(&addr).Error

	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(404,
			fmt.Sprintf("Address with ID: %d does not exist.", id),
		)
	}

	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&addr)
}

// addressPost
//
//	@Description	Create a new `Address` for the current `User`
//	@Tags			addresses
//	@Produce		json
//	@Param			request	body		user.AddressCreateDto	true	"Request Body"
//	@Success		201		{object}	user.Address
//	@Failure		400		{object}	model.ErrorResponse
//	@Router			/v1/profile/addresses [post]
func (s *Controller) addressPost(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	dto := user.AddressCreateDto{}

	err := c.BodyParser(&dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		log.Error(err)
		return err
	}

	dto.UserID = u.ID

	addr := user.Address{
		FirstName:      dto.FirstName,
		LastName:       dto.LastName,
		Company:        dto.Company,
		StreetAddress1: dto.StreetAddress1,
		StreetAddress2: dto.StreetAddress2,
		City:           dto.City,
		Zip:            dto.Zip,
		Phone:          dto.Phone,
		Ext:            dto.Ext,
		Email:          dto.Email,
		StateID:        dto.StateID,
		CountryID:      dto.CountryID,
		UserID:         u.ID,
	}

	err = s.DB.Create(&addr).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return c.Status(201).JSON(addr)
}

// addressByIdPut
//
//	@Description	Update an `Address` by it's `ID`.
//	@Tags			addresses
//	@Produce		json
//	@Param			id		path		int						true	"Address ID"
//	@Param			request	body		user.AddressUpdateDto	true	"Request Body"
//	@Success		200		{object}	user.Address
//	@Router			/v1/profile/addresses/{id} [put]
func (s *Controller) addressByIdPut(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	dto := user.AddressUpdateDto{}

	err = c.BodyParser(&dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		log.Error(err)
		return err
	}

	addr := user.Address{ID: uint(id), UserID: u.ID}

	err = s.DB.Model(&addr).Updates(dto).Error
	if err != nil {
		return err
	}

	s.DB.
		Joins("State").
		Joins("Country").
		First(&addr)

	return c.JSON(addr)
}

// addressByIdDelete
//
//	@Description	Get `Address` Info By ID
//	@Tags			addresses
//	@Produce		json
//	@Param			id							path		int	true	"Address ID"
//	@Success		200							{object}	model.EntityDeletedResponse
//	@Router			/v1/profile/addresses/{id} 	[delete]
func (s *Controller) addressByIdDelete(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	addr := user.Address{
		ID:     uint(id),
		UserID: u.ID,
	}

	res := s.DB.Where(&addr).Delete(&addr)

	if res.Error != nil {
		log.Error(res.Error)
		return res.Error
	}

	if res.RowsAffected < 1 {
		return fiber.NewError(404, fmt.Sprintf("Address with ID: %d does not exist.", id))
	}

	return c.JSON(model.EntityDeletedResponse{
		RowsAffected: res.RowsAffected,
		Message:      fmt.Sprintf("Address with ID: %d deleted successfully.", id),
	})

}
