package routes

import (
	"math"
	"pex-universe/model/address"
	"pex-universe/model/user"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Controller) RegisterProfileRoutes() {
	v1 := s.App.Group("/v1")

	v1.Get("/profile", s.profileGet)
	v1.Put("/profile", s.profilePut)

	v1.Get("/profile/addresses", s.addressGet)
	v1.Get("/profile/addresses/:id", s.addressByIdGet)
	v1.Post("/profile/addresses", s.addressPost)
	v1.Put("/profile/addresses/:id", s.addressByIdPut)
}

type ProfileUpdateDto struct {
	Name string `validate:"required" example:"John Doe"`
}

// profileGet
//
//	@Summary	Get Profile Info
//	@Tags		profile
//	@Produce	json
//	@Success	200	{object}	user.User
//	@Router		/v1/profile [get]
func (s *Controller) profileGet(c *fiber.Ctx) error {
	user := c.Locals("user")

	return c.JSON(user)
}

// profilePut
//
//	@Summary	Update Profile
//	@Tags		profile
//	@Produce	json
//	@Param		request	body		ProfileUpdateDto	true	"Profile Update Dto"
//	@Success	201		{object}	user.User
//	@Router		/v1/profile [put]
func (s *Controller) profilePut(c *fiber.Ctx) error {
	dto := new(ProfileUpdateDto)

	err := c.BodyParser(dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(dto)
	if err != nil {
		return err
	}

	u := c.Locals("user").(user.User)

	u.Name = dto.Name

	err = s.DB.Save(&u).Error
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

type AddressesResponse struct {
	Data        []address.Address
	CurrentPage int
	TotalPages  int
}

// addressGet
//
//	@Description	Get List of `Address`es for the current `User`
//	@Tags			profile
//	@Produce		json
//	@Param			page	query	int	false	"page number"		default(1)
//	@Param			limit	query	int	false	"limit of results"	default(10)
//	@Success		200		{array}	AddressesResponse
//	@Router			/v1/profile/addresses [get]
func (s *Controller) addressGet(c *fiber.Ctx) error {
	user := c.Locals("user").(user.User)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	addrs := []address.Address{}
	count := int64(0)

	err := s.DB.
		Joins("State").
		Joins("Country").
		Where(&address.Address{UserID: user.ID}).
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&addrs).
		Count(&count).Error

	if err != nil {
		return err
	}

	totalPages := int(math.Ceil(float64(count) / float64(limit)))

	return c.JSON(AddressesResponse{
		Data:        addrs,
		CurrentPage: page,
		TotalPages:  totalPages,
	})
}

// addressByIdGet
//
//	@Description	Get `Address` Info By ID
//	@Tags			profile
//	@Produce		json
//	@Param			id							path		int	true	"Address ID"
//	@Success		200							{object}	address.Address
//	@Router			/v1/profile/addresses/{id} 	[get]
func (s *Controller) addressByIdGet(c *fiber.Ctx) error {
	u := c.Locals("user").(user.User)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	addr := address.Address{}

	err = s.DB.
		Joins("State").
		Joins("Country").
		Where(&address.Address{
			ID:     uint(id),
			UserID: u.ID,
		}).
		First(&addr).Error

	if err != nil {
		return err
	}

	return c.JSON(addr)
}

// addressPost
//
//	@Description	Create a new `Address` for the current `User`
//	@Tags			profile
//	@Produce		json
//	@Param			request	body	address.AddressCreateDto	true	"Request Body"
//	@Success		201	{object}	address.Address
//	@Router			/v1/profile/addresses [post]
func (s *Controller) addressPost(c *fiber.Ctx) error {
	user := c.Locals("user").(user.User)

	dto := address.AddressCreateDto{}

	err := c.BodyParser(&dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		return err
	}

	dto.UserID = user.ID

	addr := address.Address{
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
		UserID:         user.ID,
	}

	err = s.DB.Create(&addr).Error
	if err != nil {
		return err
	}

	return c.Status(201).JSON(addr)
}

// addressByIdPut
//
//	@Description	Update an `Address` by it's `ID`.
//	@Tags			profile
//	@Produce		json
//	@Param			id		path		int							true	"Address ID"
//	@Param			request	body		address.AddressUpdateDto	true	"Request Body"
//	@Success		200		{object}	address.Address
//	@Router			/v1/profile/addresses/{id} [put]
func (s *Controller) addressByIdPut(c *fiber.Ctx) error {
	user := c.Locals("user").(user.User)

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	dto := address.AddressUpdateDto{}

	err = c.BodyParser(&dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(&dto)
	if err != nil {
		return err
	}

	addr := address.Address{ID: uint(id), UserID: user.ID}

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
