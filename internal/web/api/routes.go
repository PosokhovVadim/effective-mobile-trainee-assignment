package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App, handlers *SongsHandlers) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("/api/v1/song", handlers.AddSong)
	app.Delete("/api/v1/song/:id", handlers.DeleteSong)
	app.Get("/api/v1/lyrics/:id", handlers.GetLyrics)
	app.Put("/api/v1/song/:id", handlers.UpdateSong)
	app.Get("/api/v1/library", handlers.GetLibrary)
}
