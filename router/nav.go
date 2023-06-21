package router

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func AddNavGroup(app *fiber.App) {
	navGroup := app.Group("/nav")

	navGroup.Get("/management", getManagement)
	navGroup.Get("/purchase", buyStuff)
}

func getManagement(c *fiber.Ctx) error {
	return c.Render("management", fiber.Map{"url": os.Getenv("LocalURL")})
}

func buyStuff(c *fiber.Ctx) error {
	return c.Render("purchase", fiber.Map{"url": os.Getenv("LocalURL"), "product": os.Getenv("testProduct")})
}
