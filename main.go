package main

import (
	"os"

	"PBAPP/common"
	"PBAPP/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func main() {
	err := run()

	if err != nil {
		panic(err)
	}
}

func run() error {
	// init env
	err := common.LoadEnv()
	if err != nil {
		return err
	}

	// init db
	err = common.InitDB()
	if err != nil {
		return err
	}

	// defer closing db
	defer common.CloseDB()

	// template engine
	engine := html.New("./views", ".tmpl")

	// create app
	app := fiber.New(
		fiber.Config{
			Views: engine,
		})

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// add routes
	router.AddCTMGroup(app)
	router.AddBIZGroup(app)
	router.AddADDGroup(app)
	router.AddNavGroup(app)
	router.AddPaddleGroup(app)

	// start server
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	app.Listen(":" + port)

	return nil
}
