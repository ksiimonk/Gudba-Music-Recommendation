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
	userUseCase := usecase.NewUserUseCase(userRepository)
	tokenManager, err := auth.NewTokenManager(cfg.JWTSecret)
	if err != nil {
		return err
	}

	userHandler := httpcontroller.NewUserHandler(userUseCase, tokenManager)
	authMiddleware := httpcontroller.NewAuthMiddleware(tokenManager)

	router, err := httpcontroller.NewRouter(userHandler, authMiddleware)
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
