package web

import (
	"log/slog"
	"songs_lib/internal/dto"
	songService "songs_lib/internal/service"
	web "songs_lib/internal/web/external"
	"songs_lib/pkg/logger"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SongsHandlers struct {
	songService songService.SongService
	log         *slog.Logger
	validate    *validator.Validate
}

func NewSongsHandlers(log *slog.Logger, songService songService.SongService) *SongsHandlers {
	return &SongsHandlers{
		songService: songService,
		log:         log,
		validate:    validator.New(),
	}
}

func (h *SongsHandlers) AddSong(c *fiber.Ctx) error {
	var req dto.CreateSongRequest
	if err := c.BodyParser(&req); err != nil {
		log.Debug("Failed to parse request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		log.Debug("Failed to validate request body", logger.Err(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	//TODO:
	fetchData, err := web.FetchSongs()
	if err != nil {
		log.Debug("Failed to fetch song data", logger.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch song data",
		})
	}

	if err := h.songService.AddSong(
		req.Group,
		req.Name,
		fetchData.Link,
		fetchData.ReleaseDate,
		fetchData.Text,
	); err != nil {
		log.Error("Failed to add song", logger.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add song",
		})
	}

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
