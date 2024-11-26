package app

import (
	"fmt"
	"log/slog"
	"songs_lib/config"
	"songs_lib/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type App struct {
	log   *slog.Logger
	host  string
	port  int
	fiber *fiber.App
}

func NewApp(log *slog.Logger, http config.HTTP, storage config.Storage) (*App, error) {
	// init storage layer

	// init service layer

	fiber := SetupFiber(http)
	// init handler layer

	return &App{
		log:   log,
		port:  http.Port,
		fiber: fiber,
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