package routes

import (
	"pex-universe/model/product"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func (s *Controller) RegisterHomeRoutes() {
	v1 := s.Group("/v1")

	v1.Get("/home-page", s.homepageGet)
}

type HomePageResp struct {
	Categories []product.Category
	Brands     []product.Manufacturer
}

// homepageGet
//
//	@Router		/v1/home-page [get]
//	@Tags		homepage
//	@Produce	json
//	@Success	200	{object} HomePageResp
func (s *Controller) homepageGet(c *fiber.Ctx) error {
	// These are hardcoded as taken from the design and main website
	categoryIds := []uint{372, 373, 374, 778, 382, 628, 597, 1329, 1345}
	cs := []product.Category{}

	err := s.DB.
		Where("id IN (?)", categoryIds).
		Find(&cs).Error
	if err != nil {
		log.Error(err)
		return err
	}

	// These are hardcoded as taken from the design and main website
	manufacturerIds := []uint{996, 910, 997, 724, 1030, 1022, 63, 57, 999}
	ms := []product.Manufacturer{}

	err = s.DB.
		Where("id IN (?)", manufacturerIds).
		Find(&ms).Error
	if err != nil {
		log.Error(err)
		return err
	}

	return c.JSON(HomePageResp{cs, ms})
}
