package web

import (
	"log/slog"
	"songs_lib/internal/dto"
	songService "songs_lib/internal/service"
	web "songs_lib/internal/web/external"
	"songs_lib/pkg/logger"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type SongsHandlers struct {
	songService songService.ISong
	log         *slog.Logger
	validate    *validator.Validate
	externalAPI string
}

func NewSongsHandlers(log *slog.Logger,
	songService songService.ISong,
	externalAPI string,
) *SongsHandlers {
	return &SongsHandlers{
		songService: songService,
		log:         log,
		validate:    validator.New(),
		externalAPI: externalAPI,
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

	fetchData, err := web.FetchSong(h.externalAPI, req.Group, req.Name)
	if err != nil {
		log.Debug("Failed to fetch song data", logger.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch song data",
		})
	}

	releaseDate, err := time.Parse("2006-01-02", fetchData.ReleaseDate)
	if err != nil {
		log.Debug("Failed to parse release date", logger.Err(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse release date",
		})
	}

	if err := h.songService.AddSong(
		req.Group,
		req.Name,
		fetchData.Link,
		releaseDate,
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
