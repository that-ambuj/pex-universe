package server

import (
	"pex-universe/model"

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
//	@Success	200	{object}	model.User
//	@Router		/v1/profile [get]
func (s *FiberServer) profileGet(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)

	return c.JSON(user)
}

// profilePut
//
//	@Summary	Update Profile
//	@Tags		profile
//	@Produce	json
//	@Param		request	body		ProfileUpdateDto	true	"Profile Update Dto"
//	@Success	201		{object}	model.User
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

	user := c.Locals("user").(*model.User)

	_, err = s.db.Exec(`UPDATE users SET name = ? WHERE id = ?;`, dto.Name, user.Id)
	if err != nil {
		return err
	}

	newUser := new(model.User)
	s.db.Get(newUser, `SELECT * FROM users WHERE id = ?;`, user.Id)

	return c.Status(fiber.StatusCreated).JSON(newUser)
}

// addressGet
//
//	@Summary	Get List of Addresses for the current user
//	@Tags		profile
//	@Produce	json
//	@Success	200	{array}	model.Address
//	@Router		/v1/profile/addresses [get]
func (s *FiberServer) addressGet(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)

	var err error

	addresses, err := model.FindAddressesByUserId(s.db.DB, user.Id)
	if err != nil {
		return err
	}

	return c.JSON(addresses)
}

// addressPost
//
//	@Summary 	Create a new `Address` for the current `User`
//	@Tags		profile
//	@Produce	json
//	@Success	200	{array}	model.Address
//	@Router		/v1/profile/addresses [post]
func (s *FiberServer) addressPost(c *fiber.Ctx) error {
	user := c.Locals("user").(*model.User)

	dto := new(model.AddressCreateDto)

	err := c.BodyParser(dto)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	err = s.ValidateStruct(dto)
	if err != nil {
		return err
	}

	dto.UserId = user.Id

	err = dto.InsertNew(s.db.DB)
	if err != nil {
		return err
	}

	addrs, err := model.FindAddressesByUserId(s.db.DB, user.Id)
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
