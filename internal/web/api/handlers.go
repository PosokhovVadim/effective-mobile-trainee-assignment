package web

import (
	"log/slog"
	"songs_lib/internal/dto"
	songService "songs_lib/internal/service"
	web "songs_lib/internal/web/external"
	"songs_lib/pkg/logger"
	"strconv"
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

	id, err := h.songService.AddSong(
		req.Group,
		req.Name,
		fetchData.Link,
		releaseDate,
		fetchData.Text,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add song",
		})
	}

	return c.Status(fiber.StatusOK).JSON(dto.CreateSongResponse{
		ID:          id,
		Group:       req.Group,
		Name:        req.Name,
		ReleaseDate: releaseDate,
		Link:        fetchData.Link,
		Text:        fetchData.Text,
	})
}

func (h *SongsHandlers) DeleteSong(c *fiber.Ctx) error {
	param := c.Params("song_id")

	songID, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid song ID",
		})
	}
	if err := h.songService.DeleteSong((uint)(songID)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete song",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Song deleted successfully",
	})
}

func (h *SongsHandlers) GetLyrics(c *fiber.Ctx) error {
	param := c.Params("song_id")

	songID, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid song ID",
		})
	}

	queryParams := c.Queries()
	lyrics, err := h.songService.GetLyrics(
		(uint)(songID),
		queryParams["limit"],
		queryParams["offset"],
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get lyrics",
		})
	}

	return c.Status(fiber.StatusOK).JSON(lyrics)
}

func (h *SongsHandlers) UpdateSong(c *fiber.Ctx) error {
	return nil
}

func (h *SongsHandlers) GetLibrary(c *fiber.Ctx) error {
	return nil
}
