package router

import (
	"PBAPP/common"
	"PBAPP/models"

	"github.com/gofiber/fiber/v2"
)

func AddBIZGroup(app *fiber.App) {
	bizGroup := app.Group("/biz")
	bizGroup.Get("/:id", fetchBIZ)
	bizGroup.Post("/update/:id", updateBiz)
}

// mostly deprecated, returns a BIZ object based on the ctm id
func fetchBIZ(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "business id is required",
		})
	}

	BIZ := common.DoesBizExist(id)
	if BIZ == "Not Found" {
		return c.Status(400).JSON(fiber.Map{
			"error": "business does not exist",
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": BIZ})
}

// Pushes updated biz data to our paddle handler and then updates the DB if successful (see bizUtils)
func updateBiz(c *fiber.Ctx) error {
	id := c.Params("id")
	business := common.DoesBizExist(id)
	if business == "Not Found" {
		return c.Status(400).JSON(fiber.Map{
			"error": "business does not exist",
		})
	}
	updatedBusiness, ok := business.(models.BIZ)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	b := new(models.BIZ)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	b.PadID = updatedBusiness.PadID
	err := common.PadUpdateBiz(id, b)
	if err == "x" {
		return c.Status(400).JSON(fiber.Map{
			"data": "failed to update",
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": "Updated Successfully"})
}
