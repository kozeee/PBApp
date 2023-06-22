package router

import (
	"PBAPP/common"
	"PBAPP/models"

	"github.com/gofiber/fiber/v2"
)

func AddADDGroup(app *fiber.App) {
	addGroup := app.Group("/add")
	addGroup.Get("/:id", fetchADD)
	addGroup.Post("/update/:id", updateAdd)
}

// mostly deprecated way to fetch address data based on ctmid
func fetchADD(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "address id is required",
		})
	}

	ADD := common.DoesAddExist(id)
	if ADD == "Not Found" {
		return c.Status(400).JSON(fiber.Map{
			"error": "address does not exist",
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": ADD})
}

// pushes Address data to the paddle handler and updates the db if successful (see addUtils)
func updateAdd(c *fiber.Ctx) error {
	id := c.Params("id")
	address := common.DoesAddExist(id)
	if address == "Not Found" {
		return c.Status(400).JSON(fiber.Map{
			"error": "address does not exist",
		})
	}
	updatedAddress, ok := address.(models.ADD)
	if !ok {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	b := new(models.ADD)
	if err := c.BodyParser(b); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid body",
		})
	}

	b.PadID = updatedAddress.PadID
	err := common.PadUpdateAddress(id, b)
	if err == "x" {
		return c.Status(400).JSON(fiber.Map{
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{"data": "Updated Successfully"})
}
