package routes

import (
	"pex-universe/model"
	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Controller) RegisterProfileRoutes() {
	v1 := s.App.Group("/v1")

	v1.Get("/profile", s.profileGet)
	v1.Put("/profile", s.profilePut)
}

type ProfileUpdateDto struct {
	Name *string `example:"John Doe"`
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
		log.Error(err)
		return err
	}

	u := c.Locals("user").(user.User)

	u.Name = dto.Name

	err = s.DB.Save(&u).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(u)
}

type AddressesResponse struct {
	Data []user.Address
	model.PageResponse
}
