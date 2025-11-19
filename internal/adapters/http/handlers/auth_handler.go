package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
    authService ports.AuthService
    validator   *validator.Validate
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Message: "Invalid input data",
			Code:    "VALIDATION_ERROR",
			Details: validationErrors,
		})
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "REGISTRATION_FAILED"

		// Check for specific errors
		if err.Error() == "user already exists" {
			statusCode = http.StatusConflict
			errorCode = "USER_EXISTS"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Registration failed",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusCreated, dto.AuthResponse{
		User: dto.UserData{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	})
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = err.Tag()
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Message: "Invalid input data",
			Code:    "VALIDATION_ERROR",
			Details: validationErrors,
		})
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "LOGIN_FAILED"

		// Check for authentication errors
		if err.Error() == "invalid credentials" {
			statusCode = http.StatusUnauthorized
			errorCode = "INVALID_CREDENTIALS"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Login failed",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		User: dto.UserData{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	})
}

// RefreshToken handles token refresh
// @Summary Refresh authentication token
// @Description Refresh an expired or soon-to-expire authentication token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AuthResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "Unauthorized",
			Message: "User ID not found in context",
			Code:    "UNAUTHORIZED",
		})
		return
	}

	user, token, err := h.authService.RefreshToken(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Token refresh failed",
			Message: err.Error(),
			Code:    "REFRESH_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		User: dto.UserData{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
		},
		Token: token,
	})
}

// CompleteOnboarding handles completion of onboarding profile
// Note: Service implementation to persist fields may be pending in current codebase.
func (h *AuthHandler) CompleteOnboarding(c *gin.Context) {
    var req dto.CompleteOnboardingRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, dto.ErrorResponse{
            Error:   "Invalid request body",
            Message: err.Error(),
            Code:    "INVALID_REQUEST",
        })
        return
    }

    // Try to extract userID from context (JWT middleware should set it)
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
            Error:   "Unauthorized",
            Message: "User ID not found in context",
            Code:    "UNAUTHORIZED",
        })
        return
    }

    // Until service method is available, respond with 501 to indicate pending implementation
    _ = userID // placeholder to avoid unused warning
    c.JSON(http.StatusNotImplemented, dto.ErrorResponse{
        Error:   "Not Implemented",
        Message: "Onboarding completion not yet implemented",
        Code:    "NOT_IMPLEMENTED",
    })
}
