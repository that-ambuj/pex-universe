package routes

import (
	"pex-universe/model/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Controller) RegisterUtilRoutes() {
	v1 := s.Group("/v1")

	v1.Get("/states", s.statesGet)
	v1.Get("/countries", s.countriesGet)
}

// statesGet
//
//	@Description	Get List of `States`es by ID
//	@Tags			utility
//	@Produce		json
//	@Success		200	{array}	address.State
//	@Router			/v1/states [get]
func (s *Controller) statesGet(c *fiber.Ctx) error {
	states := []user.State{}

	err := s.DB.Find(&states).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&states)
}

// countriesGet
//
//	@Description	Get List of `Country`es by ID
//	@Tags			utility
//	@Produce		json
//
//	@Success		200	{array}	address.Country
//	@Router			/v1/countries [get]
func (s *Controller) countriesGet(c *fiber.Ctx) error {
	countries := []user.Country{}

	err := s.DB.Find(&countries).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(&countries)
}
