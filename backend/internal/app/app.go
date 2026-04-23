package app

import (
	"errors"
	"log"
	"net/http"
	"time"

	"music-recommender-backend/config"
	httpcontroller "music-recommender-backend/internal/controller/http"
)

func Run(cfg *config.Config) error {
	db, err := newPostgresDB(cfg.PGURL)
	if err != nil {
		return err
	}
	defer db.Close()

	router, err := httpcontroller.NewRouter()
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:              cfg.Address(),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("starting app=%s port=%s log_level=%s", cfg.AppName, cfg.HTTPPort, cfg.LogLevel)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
