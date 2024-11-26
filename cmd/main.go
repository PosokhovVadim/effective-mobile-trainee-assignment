package main

import (
	"fmt"
	"os"
	"songs_lib/config"
	"songs_lib/internal/app"
	"songs_lib/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func run() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading environment: %v", err)
	}

	cfg := config.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("error loading environment: %v", err)
	}

	fmt.Printf("%+v\n", cfg)

	log, err := logger.SetupLogger(&cfg)
	if err != nil {
		return fmt.Errorf("error setup logger: %v", err)
	}

	log.Info("Config read success")

	app, err := app.NewApp(log, cfg.HTTP, cfg.Storage)
	if err != nil {
		log.Error("error creating app: %v", err)
		return err
	}

	if err := app.Run(); err != nil {
		log.Error("error running app: %v", err)
		return err
	}

	return nil

}

func main() {
	if err := run(); err != nil {
		fmt.Errorf("%v\n", err)
		os.Exit(1)
	}
}
