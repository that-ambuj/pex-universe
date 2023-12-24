package routes

import (
	"math"
	"pex-universe/model"
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
	user := c.Locals("user").(*user.User)

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

	var err error

	err = c.BodyParser(dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(dto)
	if err != nil {
		return err
	}

	u := c.Locals("user").(*user.User)

	_, err = s.DB.Exec(`UPDATE users SET name = ? WHERE id = ?;`, dto.Name, u.Id)
	if err != nil {
		return err
	}

	newUser := new(user.User)
	s.DB.Get(newUser, `SELECT * FROM users WHERE id = ?;`, u.Id)

	return c.Status(fiber.StatusCreated).JSON(newUser)
}

type AddressesResponse struct {
	Data        []*address.Address
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
	user := c.Locals("user").(*user.User)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	pagination := &model.PaginationDto{
		Page:  page,
		Limit: limit,
	}

	addrs, err := address.FindManyByUserId(s.DB, user.Id, pagination)
	if err != nil {
		return err
	}

	count, err := address.CountByUserId(s.DB, user.Id)

	return c.JSON(AddressesResponse{
		Data:        addrs,
		CurrentPage: page,
		TotalPages:  int(math.Ceil(float64(count) / float64(limit))),
	})
}

// addressByIdGet
//
//	@Description	Get `Address` Info By Id
//	@Tags			profile
//	@Produce		json
//	@Param			id							path		int	true	"Address ID"
//	@Success		200							{object}	address.Address
//	@Router			/v1/profile/addresses/{id} 	[get]
func (s *Controller) addressByIdGet(c *fiber.Ctx) error {
	u := c.Locals("user").(*user.User)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	addr, err := address.FindById(s.DB, uint64(id), u.Id)
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
//	@Success		201	{object}	address.Address
//	@Router			/v1/profile/addresses [post]
func (s *Controller) addressPost(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)

	dto := new(address.AddressCreateDto)

	err := c.BodyParser(dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(dto)
	if err != nil {
		return err
	}

	dto.UserId = user.Id

	lastId, err := dto.CreateNew(s.DB)
	if err != nil {
		return err
	}

	addr, err := address.FindById(s.DB, uint64(lastId), user.Id)
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
//	@Param			id	path		int	true	"Address ID"
//	@Success		200	{object}	address.Address
//	@Router			/v1/profile/addresses/{id} [put]
func (s *Controller) addressByIdPut(c *fiber.Ctx) error {

	return fiber.ErrNotImplemented
}
