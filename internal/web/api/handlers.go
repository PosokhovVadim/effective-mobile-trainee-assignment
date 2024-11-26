package web

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

type SongsHandlers struct {
	//service
	log *slog.Logger
}

func NewSongsHandlers(log *slog.Logger) *SongsHandlers {
	return &SongsHandlers{
		log: log,
	}
}

func (h *SongsHandlers) AddSong(c *fiber.Ctx) error {
	

	return nil
}

func (h *SongsHandlers) DeleteSong(c *fiber.Ctx) error {
	return nil
}

func (h *SongsHandlers) GetLyrics(c *fiber.Ctx) error {
	return nil
}

func (h *SongsHandlers) UpdateSong(c *fiber.Ctx) error {
	return nil
}

func (h *SongsHandlers) GetLibrary(c *fiber.Ctx) error {
	return nil
}
