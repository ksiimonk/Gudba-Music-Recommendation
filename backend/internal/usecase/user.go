package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
)

var ErrEmailRequired = errors.New("email is required")
var ErrPasswordRequired = errors.New("password is required")
var ErrInvalidCredentials = errors.New("invalid email or password")

type userRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type playlistCreator interface {
	CreateFavoritesPlaylist(ctx context.Context, userID int64) (*entity.Playlist, error)
}

type UserUseCase struct {
	userRepository  userRepository
	playlistCreator playlistCreator
}

func NewUserUseCase(userRepository userRepository, playlistCreator playlistCreator) *UserUseCase {
	return &UserUseCase{userRepository: userRepository, playlistCreator: playlistCreator}
}

func (u *UserUseCase) RegisterUser(ctx context.Context, email, password string) (*entity.User, error) {
	if u == nil || u.userRepository == nil {
		return nil, errors.New("user use case repository is nil")
	}

	email, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	existingUser, err := u.userRepository.GetUserByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, repository.ErrUserAlreadyExists
	}

	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &entity.User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	if u.playlistCreator != nil {
		if _, err := u.playlistCreator.CreateFavoritesPlaylist(ctx, user.ID); err != nil {
			if delErr := u.userRepository.DeleteUser(ctx, user.ID); delErr != nil {
				return nil, fmt.Errorf("create favorites playlist failed: %w; cleanup also failed: %v", err, delErr)
			}
			return nil, fmt.Errorf("create favorites playlist: %w; user deleted for consistency", err)
		}
	}

	return user, nil
}

func (u *UserUseCase) LoginUser(ctx context.Context, email, password string) (*entity.User, error) {
	if u == nil || u.userRepository == nil {
		return nil, errors.New("user use case repository is nil")
	}

	email, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}

	if err := validatePassword(password); err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf("compare password hash: %w", err)
	}

	return user, nil
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if u == nil || u.userRepository == nil {
		return nil, errors.New("user use case repository is nil")
	}

	email, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}

	user, err := u.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" {
		return "", ErrEmailRequired
	}

	return email, nil
}

func validatePassword(password string) error {
	if strings.TrimSpace(password) == "" {
		return ErrPasswordRequired
	}

	return nil
}
