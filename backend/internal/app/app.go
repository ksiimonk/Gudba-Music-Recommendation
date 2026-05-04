package app

import (
	"errors"
	"log"
	"net/http"
	"time"

	"music-recommender-backend/config"
	"music-recommender-backend/internal/auth"
	httpcontroller "music-recommender-backend/internal/controller/http"
	"music-recommender-backend/internal/repository"
	"music-recommender-backend/internal/usecase"
)

func Run(cfg *config.Config) error {
	db, err := newPostgresDB(cfg.PGURL)
	if err != nil {
		return err
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	trackRepository := repository.NewTrackRepository(db)
	playlistRepository := repository.NewPlaylistRepository(db)
	genreRepository := repository.NewGenreRepository(db)
	artistRepository := repository.NewArtistRepository(db)
	onboardingRepository := repository.NewOnboardingRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository)
	trackUseCase := usecase.NewTrackUseCase(trackRepository)
	playlistUseCase := usecase.NewPlaylistUseCase(playlistRepository)
	genreUseCase := usecase.NewGenreUseCase(genreRepository)
	artistUseCase := usecase.NewArtistUseCase(artistRepository)
	onboardingUseCase := usecase.NewOnboardingUseCase(onboardingRepository)
	tokenManager, err := auth.NewTokenManager(cfg.JWTSecret)
	if err != nil {
		return err
	}

	userHandler := httpcontroller.NewUserHandler(userUseCase, tokenManager)
	trackHandler := httpcontroller.NewTrackHandler(trackUseCase)
	playlistHandler := httpcontroller.NewPlaylistHandler(playlistUseCase)
	genreHandler := httpcontroller.NewGenreHandler(genreUseCase)
	artistHandler := httpcontroller.NewArtistHandler(artistUseCase)
	onboardingHandler := httpcontroller.NewOnboardingHandler(onboardingUseCase)
	authMiddleware := httpcontroller.NewAuthMiddleware(tokenManager)

	router, err := httpcontroller.NewRouter(userHandler, authMiddleware, trackHandler, playlistHandler, genreHandler, artistHandler, onboardingHandler)
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
