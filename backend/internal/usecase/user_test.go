package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
)

type fakeUserRepository struct {
	createUserFn     func(ctx context.Context, user *entity.User) error
	getUserByEmailFn func(ctx context.Context, email string) (*entity.User, error)
}

func (f *fakeUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	if f.createUserFn != nil {
		return f.createUserFn(ctx, user)
	}

	return nil
}

func (f *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	if f.getUserByEmailFn != nil {
		return f.getUserByEmailFn(ctx, email)
	}

	return nil, repository.ErrUserNotFound
}

func TestRegisterUserHashesPassword(t *testing.T) {
	t.Parallel()

	var createdUser *entity.User

	userUseCase := NewUserUseCase(&fakeUserRepository{
		createUserFn: func(_ context.Context, user *entity.User) error {
			createdCopy := *user
			createdUser = &createdCopy

			user.ID = 1
			user.CreatedAt = time.Now()
			user.UpdatedAt = user.CreatedAt

			return nil
		},
	})

	user, err := userUseCase.RegisterUser(context.Background(), " Test@Example.com ", "plain-password")
	if err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}

	if createdUser == nil {
		t.Fatal("CreateUser() was not called")
	}

	if createdUser.PasswordHash == "plain-password" {
		t.Fatal("password was stored in plaintext")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(createdUser.PasswordHash), []byte("plain-password")); err != nil {
		t.Fatalf("stored password hash does not match plaintext password: %v", err)
	}

	if user.Email != "test@example.com" {
		t.Fatalf("Email = %q, want %q", user.Email, "test@example.com")
	}
}

func TestRegisterUserReturnsDuplicateEmailError(t *testing.T) {
	t.Parallel()

	createCalled := false

	userUseCase := NewUserUseCase(&fakeUserRepository{
		getUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: 1, Email: email}, nil
		},
		createUserFn: func(_ context.Context, user *entity.User) error {
			createCalled = true
			return nil
		},
	})

	_, err := userUseCase.RegisterUser(context.Background(), "test@example.com", "plain-password")
	if !errors.Is(err, repository.ErrUserAlreadyExists) {
		t.Fatalf("error = %v, want %v", err, repository.ErrUserAlreadyExists)
	}

	if createCalled {
		t.Fatal("CreateUser() was called for duplicate email")
	}
}

func TestRegisterUserReturnsCreateUserDuplicateError(t *testing.T) {
	t.Parallel()

	userUseCase := NewUserUseCase(&fakeUserRepository{
		createUserFn: func(_ context.Context, user *entity.User) error {
			return repository.ErrUserAlreadyExists
		},
	})

	_, err := userUseCase.RegisterUser(context.Background(), "test@example.com", "plain-password")
	if !errors.Is(err, repository.ErrUserAlreadyExists) {
		t.Fatalf("error = %v, want %v", err, repository.ErrUserAlreadyExists)
	}
}

func TestRegisterUserValidatesInput(t *testing.T) {
	t.Parallel()

	userUseCase := NewUserUseCase(&fakeUserRepository{})

	_, err := userUseCase.RegisterUser(context.Background(), "   ", "plain-password")
	if !errors.Is(err, ErrEmailRequired) {
		t.Fatalf("error = %v, want %v", err, ErrEmailRequired)
	}

	_, err = userUseCase.RegisterUser(context.Background(), "test@example.com", "   ")
	if !errors.Is(err, ErrPasswordRequired) {
		t.Fatalf("error = %v, want %v", err, ErrPasswordRequired)
	}
}

func TestLoginUserSuccess(t *testing.T) {
	t.Parallel()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("plain-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword() error = %v", err)
	}

	userUseCase := NewUserUseCase(&fakeUserRepository{
		getUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:           1,
				Email:        email,
				PasswordHash: string(passwordHash),
			}, nil
		},
	})

	user, err := userUseCase.LoginUser(context.Background(), " Test@Example.com ", "plain-password")
	if err != nil {
		t.Fatalf("LoginUser() error = %v", err)
	}

	if user.Email != "test@example.com" {
		t.Fatalf("Email = %q, want %q", user.Email, "test@example.com")
	}
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	t.Parallel()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("another-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword() error = %v", err)
	}

	userUseCase := NewUserUseCase(&fakeUserRepository{
		getUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{
				ID:           1,
				Email:        email,
				PasswordHash: string(passwordHash),
			}, nil
		},
	})

	_, err = userUseCase.LoginUser(context.Background(), "test@example.com", "plain-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestLoginUserUserNotFound(t *testing.T) {
	t.Parallel()

	userUseCase := NewUserUseCase(&fakeUserRepository{
		getUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return nil, repository.ErrUserNotFound
		},
	})

	_, err := userUseCase.LoginUser(context.Background(), "test@example.com", "plain-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestGetUserByEmail(t *testing.T) {
	t.Parallel()

	userUseCase := NewUserUseCase(&fakeUserRepository{
		getUserByEmailFn: func(_ context.Context, email string) (*entity.User, error) {
			return &entity.User{ID: 1, Email: email}, nil
		},
	})

	user, err := userUseCase.GetUserByEmail(context.Background(), " Test@Example.com ")
	if err != nil {
		t.Fatalf("GetUserByEmail() error = %v", err)
	}

	if user.Email != "test@example.com" {
		t.Fatalf("Email = %q, want %q", user.Email, "test@example.com")
	}
}
