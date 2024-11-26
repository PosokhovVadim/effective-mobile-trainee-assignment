package app

import (
	"fmt"
	"log/slog"
	"songs_lib/config"
	"songs_lib/internal/storage/postgresql"
	"songs_lib/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type App struct {
	log   *slog.Logger
	port  int
	fiber *fiber.App
	DB    *postgresql.PostgresStorage
}

func NewApp(log *slog.Logger, http config.HTTP, storage config.Storage) (*App, error) {
	psStorage, err := postgresql.NewPostgresStorage(log, storage.Path)
	if err != nil {
		log.Error("error creating storage: %v", logger.Err(err))
	}
	log.Debug("Storage setup successfully by path ", slog.String("path", storage.Path))

	// init service layer
	_ = psStorage

	fiber := SetupFiber(http)
	// init handler layer

	return &App{
		log:   log,
		port:  http.Port,
		fiber: fiber,
		DB:    psStorage,
	}, nil
}

func SetupFiber(http config.HTTP) *fiber.App {
	app := fiber.New(
		fiber.Config{
			ReadTimeout:  http.Timeout,
			IdleTimeout:  http.Timeout,
			WriteTimeout: http.Timeout,
		},
	)
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	return app
}

func (a *App) Run() error {
	a.log.Info("Starting http server", slog.Int("port", a.port))

	if err := a.fiber.Listen(fmt.Sprintf(":%d", a.port)); err != nil {
		a.log.Error("Failed to run app:", logger.Err(err))
		return err
	}
	return nil
}
