package app

import (
	"log"

	"music-recommender-backend/config"
	httpcontroller "music-recommender-backend/internal/controller/http"
)

func Run(cfg *config.Config) {
	router := httpcontroller.NewRouter()

	log.Printf("starting app=%s port=%s log_level=%s", cfg.AppName, cfg.HTTPPort, cfg.LogLevel)

	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
