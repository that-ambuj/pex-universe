package server

import (
	"pex-universe/model/address"
	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterProfileRoutes() {
	v1 := s.App.Group("/v1")

	v1.Get("/profile", s.profileGet)
	v1.Put("/profile", s.profilePut)

	v1.Get("/profile/addresses", s.addressGet)
	v1.Get("/profile/addresses/:id", s.addressByIdGet)
	v1.Post("/profile/addresses", s.addressPost)
	v1.Put("/profile/addresses/:id", s.addressByIdPut)
}

type ProfileUpdateDto struct {
	Name string `validate:"required" json:"name" example:"John Doe"`
}

// profileGet
//
//	@Summary	Get Profile Info
//	@Tags		profile
//	@Produce	json
//	@Success	200	{object}	user.User
//	@Router		/v1/profile [get]
func (s *FiberServer) profileGet(c *fiber.Ctx) error {
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
func (s *FiberServer) profilePut(c *fiber.Ctx) error {
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

	_, err = s.db.Exec(`UPDATE users SET name = ? WHERE id = ?;`, dto.Name, u.Id)
	if err != nil {
		return err
	}

	newUser := new(user.User)
	s.db.Get(newUser, `SELECT * FROM users WHERE id = ?;`, u.Id)

	return c.Status(fiber.StatusCreated).JSON(newUser)
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
func (s *FiberServer) addressGet(c *fiber.Ctx) error {
	user := c.Locals("user").(*user.User)

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	pagination := model.PaginationDto{
		Page:  page,
		Limit: limit,
	}

	addrs, err := address.FindManyByUserId(s.db, user.Id, pagination)
	if err != nil {
		return err
	}

	count, err := address.CountByUserId(s.db, user.Id)

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
func (s *FiberServer) addressByIdGet(c *fiber.Ctx) error {
	u := c.Locals("user").(*user.User)

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	addr, err := address.FindById(s.db, uint64(id), u.Id)
	if err != nil {
		return err
	}

	return c.JSON(addr)
}

// addressPost
//
//	@Summary 	Create a new `Address` for the current `User`
//	@Tags		profile
//	@Produce	json
//	@Success	200	{array}	address.Address
//	@Router		/v1/profile/addresses [post]
func (s *FiberServer) addressPost(c *fiber.Ctx) error {
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

	err = dto.CreateNew(s.db.DB)
	if err != nil {
		return err
	}

	addrs, err := address.FindManyByUserId(s.db.DB, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(addrs)
}

func (s *FiberServer) addressByIdGet(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}

func (s *FiberServer) addressByIdPut(c *fiber.Ctx) error {
	return fiber.ErrNotImplemented
}
