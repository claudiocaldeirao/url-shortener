package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()

	// @todo: set some default values here

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
