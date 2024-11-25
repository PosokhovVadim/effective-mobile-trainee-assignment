package main

import (
	"fmt"
	"log"
	"songs_lib/config"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v\n", err)
	}

	cfg := config.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Error loading environment: %v\n", err)
	}

	fmt.Printf("%+v\n", cfg)
	// init logger
	// init app lay

	// run server

}
