package main

import (
	"fmt"
	"os"
	"songs_lib/config"
	"songs_lib/internal/app"
	"songs_lib/pkg/logger"

	_ "songs_lib/docs"

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

	app, err := app.NewApp(log, cfg.HTTP, cfg.Storage, cfg.ExternalAPIURL)
	if err != nil {
		log.Error("error creating app: %v", err)
		return err
	}
	defer app.DB.Close()

	if err := app.Run(); err != nil {
		log.Error("error running app: %v", err)
		return err
	}

	return nil

}

//	@title			Song Library Api
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	if err := run(); err != nil {
		fmt.Printf("error running service: %v\n", err)
		os.Exit(1)
	}
}
