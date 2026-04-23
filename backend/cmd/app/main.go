package main

import (
	"log"

	"github.com/joho/godotenv"

	"music-recommender-backend/config"
	"music-recommender-backend/internal/app"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("startup error: %v", err)
	}
}
