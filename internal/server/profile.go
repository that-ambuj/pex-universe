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

	user.Addresses = []*model.Address{}

	query := `
SELECT a.id,
	a.created_at,
	a.updated_at,
	a.verified,
	a.first_name,
	a.last_name,
	a.company,
	a.street_address,
	a.street_address2,
	a.city,
	a.zip,
	a.phone,
	a.ext,
	a.email,

	s.name,
	s.full_name,
	s.tax,
	s.info,

	c.name,
	c.display_name,
	c.pp_code,
	c.locked
FROM addresses a
         JOIN states s ON a.state_id = s.id
         JOIN countries c ON a.country_id = c.id
WHERE a.user_id = ?;
`

	rows, err := s.db.Queryx(query, user.Id)
	if err != nil {
		return err
	}

	for rows.Next() {
		addr := &model.Address{}

		addr.State = &model.State{}
		addr.Country = &model.Country{}

		err = rows.Scan(
			&addr.Id,
			&addr.CreatedAt,
			&addr.UpdatedAt,
			&addr.Verified,
			&addr.FirstName,
			&addr.LastName,
			&addr.Company,
			&addr.StreetAddress1,
			&addr.StreetAddress2,
			&addr.City,
			&addr.Zip,
			&addr.Phone,
			&addr.Ext,
			&addr.Email,

			&addr.State.Name,
			&addr.State.FullName,
			&addr.State.Tax,
			&addr.State.Info,

			&addr.Country.Id,
			&addr.Country.Name,
			&addr.Country.DisplayName,
			&addr.Country.PpCode,
		)

		if err != nil {
			return err
		}

		user.Addresses = append(user.Addresses, addr)
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
