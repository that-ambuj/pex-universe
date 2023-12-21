package server

import (
	"pex-universe/model"

	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) RegisterProfileRoutes() {
	v1 := s.App.Group("/v1")

	v1.Get("/profile", s.profileGet)
	v1.Put("/profile", s.profilePut)
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

	if addrs, err := model.FindAddressesByUserId(s.db.DB, user.Id); err == nil {
		user.Addresses = addrs
	} else {
		return err
	}

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
