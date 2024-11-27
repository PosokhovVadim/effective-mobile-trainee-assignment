package web

import (
	"log/slog"
	"songs_lib/internal/dto"
	"songs_lib/internal/model"
	songService "songs_lib/internal/service"
	web "songs_lib/internal/web/external"
	"songs_lib/pkg/logger"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	_ "songs_lib/docs"

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

// @Summary Добавление песни
// @Description Добавление песни с указаым названием и группой
// @ID add-song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body dto.CreateSongRequest true "Song"
// @Success 201 {object} dto.CreateSongResponse
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/song [post]
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

// @Summary Удаление песни
// @Description Удалене песни по id
// @ID delete-song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song_id path int true "Song id"
// @Success 204 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/songs/{song_id} [delete]
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

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Получение текста песни
// @Description Получение текста песни с пагинацией по куплетам
// @ID get-lyrics
// @Tags Lyrics
// @Param song_id path int true "Song ID"
// @Param limit query int false "Количество куплетов"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {array} dto.LyricsInfo
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/lyrics/{id} [get]
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

// @Summary Обновление песни
// @Description Обновление полей песни и текста куплетов
// @ID update-song
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param updates body model.SongUpdate true "Updated Fields"
// @Success 204 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/song/{id} [put]
func (h *SongsHandlers) UpdateSong(c *fiber.Ctx) error {
	param := c.Params("id")

	var updates model.SongUpdate
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	songID, err := strconv.Atoi(param)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid song ID",
		})
	}

	_, err = h.songService.UpdateSong(uint(songID), updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update song",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)

}

// @Summary Получение библиотеки песен
// @Description Получение списка песен с фильтрацией и пагинацией
// @ID get-library
// @Tags Songs
// @Param group query string false "Название группы"
// @Param name query string false "Название песни"
// @Param release_date query string false "Дата релиза"
// @Param limit query int false "Количество записей на странице"
// @Param offset query int false "Смещение для пагинации"
// @Success 200 {object} dto.LibraryResponse
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /api/v1/library [get]
func (h *SongsHandlers) GetLibrary(c *fiber.Ctx) error {
	queryParams := c.Queries()
	library, err := h.songService.GetLibrary(
		map[string]string{
			"group":        queryParams["group"],
			"name":         queryParams["name"],
			"release_date": queryParams["release_date"],
		},
		queryParams["limit"],
		queryParams["offset"],
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get library",
		})
	}
	return c.Status(fiber.StatusOK).JSON(library)
}
