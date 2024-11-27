package web

import (
	_ "songs_lib/docs"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(app *fiber.App, handlers *SongsHandlers) {
	app.Get("/swagger/*", swagger.WrapHandler)
	app.Post("/api/v1/song", handlers.AddSong)
	app.Delete("/api/v1/song/:id", handlers.DeleteSong)
	app.Get("/api/v1/lyrics/:id", handlers.GetLyrics)
	app.Put("/api/v1/song/:id", handlers.UpdateSong)
	app.Get("/api/v1/library", handlers.GetLibrary)
}
