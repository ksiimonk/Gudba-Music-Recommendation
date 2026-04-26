package http

import (
	"context"
	"errors"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"music-recommender-backend/internal/entity"
	"music-recommender-backend/internal/repository"
	"music-recommender-backend/internal/usecase"
)

type userRegistrationUseCase interface {
	RegisterUser(ctx context.Context, email, password string) (*entity.User, error)
	LoginUser(ctx context.Context, email, password string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type tokenManager interface {
	GenerateToken(user *entity.User) (string, error)
}

type UserHandler struct {
	userUseCase  userRegistrationUseCase
	tokenManager tokenManager
}

type registerUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerUserResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type registerUserSuccessResponse struct {
	Message string               `json:"message"`
	User    registerUserResponse `json:"user"`
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUserHandler(userUseCase userRegistrationUseCase, tokenManager tokenManager) *UserHandler {
	return &UserHandler{
		userUseCase:  userUseCase,
		tokenManager: tokenManager,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var request registerUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	validationErrors := validateRegisterUserRequest(request)
	if len(validationErrors) > 0 {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	user, err := h.userUseCase.RegisterUser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmailRequired), errors.Is(err, usecase.ErrPasswordRequired):
			c.JSON(nethttp.StatusBadRequest, gin.H{
				"message": "validation failed",
				"errors":  mapValidationError(err),
			})
		case errors.Is(err, repository.ErrUserAlreadyExists):
			c.JSON(nethttp.StatusConflict, gin.H{
				"message": "user with this email already exists",
			})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{
				"message": "failed to register user",
			})
		}
		return
	}

	c.JSON(nethttp.StatusCreated, registerUserSuccessResponse{
		Message: "user registered successfully",
		User: registerUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var request loginUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	validationErrors := validateAuthRequest(request.Email, request.Password)
	if len(validationErrors) > 0 {
		c.JSON(nethttp.StatusBadRequest, gin.H{
			"message": "validation failed",
			"errors":  validationErrors,
		})
		return
	}

	user, err := h.userUseCase.LoginUser(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrEmailRequired), errors.Is(err, usecase.ErrPasswordRequired):
			c.JSON(nethttp.StatusBadRequest, gin.H{
				"message": "validation failed",
				"errors":  mapValidationError(err),
			})
		case errors.Is(err, usecase.ErrInvalidCredentials):
			c.JSON(nethttp.StatusUnauthorized, gin.H{
				"message": "invalid email or password",
			})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{
				"message": "failed to login user",
			})
		}
		return
	}

	token, err := h.tokenManager.GenerateToken(user)
	if err != nil {
		c.JSON(nethttp.StatusInternalServerError, gin.H{
			"message": "failed to generate token",
		})
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func (h *UserHandler) Me(c *gin.Context) {
	claims, ok := getAuthClaims(c)
	if !ok {
		c.JSON(nethttp.StatusUnauthorized, gin.H{
			"message": "authorization required",
		})
		return
	}

	user, err := h.userUseCase.GetUserByEmail(c.Request.Context(), claims.Email)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			c.JSON(nethttp.StatusNotFound, gin.H{
				"message": "user not found",
			})
		default:
			c.JSON(nethttp.StatusInternalServerError, gin.H{
				"message": "failed to load current user",
			})
		}
		return
	}

	c.JSON(nethttp.StatusOK, gin.H{
		"user": registerUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func validateRegisterUserRequest(request registerUserRequest) map[string]string {
	return validateAuthRequest(request.Email, request.Password)
}

func validateAuthRequest(email, password string) map[string]string {
	validationErrors := make(map[string]string)

	if strings.TrimSpace(email) == "" {
		validationErrors["email"] = "email is required"
	}

	if strings.TrimSpace(password) == "" {
		validationErrors["password"] = "password is required"
	}

	return validationErrors
}

func mapValidationError(err error) map[string]string {
	switch {
	case errors.Is(err, usecase.ErrEmailRequired):
		return map[string]string{"email": err.Error()}
	case errors.Is(err, usecase.ErrPasswordRequired):
		return map[string]string{"password": err.Error()}
	default:
		return map[string]string{}
	}
}
