package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)

}
