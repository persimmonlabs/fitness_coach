package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// ExerciseHandler handles exercise-related requests
type ExerciseHandler struct {
	exerciseService ports.ExerciseService
	validator       *validator.Validate
}

// NewExerciseHandler creates a new exercise handler
func NewExerciseHandler(exerciseService ports.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
		validator:       validator.New(),
	}
}

// SearchExercises searches for exercises
// @Summary Search exercises
// @Description Search for exercises by name, category, or muscle group
// @Tags exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search query"
// @Param category query string false "Filter by category"
// @Param muscle_group query string false "Filter by muscle group"
// @Param limit query int false "Results limit" default(20)
// @Success 200 {array} dto.ExerciseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /exercises/search [get]
func (h *ExerciseHandler) SearchExercises(c *gin.Context) {
	query := c.Query("query")
	category := c.Query("category")
	muscleGroup := c.Query("muscle_group")
	limit := 20

	if limitStr := c.Query("limit"); limitStr != "" {
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid limit parameter",
				Message: "Limit must be a valid integer",
				Code:    "INVALID_LIMIT",
			})
			return
		}
	}

	exercises, err := h.exerciseService.SearchExercises(c.Request.Context(), query, category, muscleGroup, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Search failed",
			Message: err.Error(),
			Code:    "SEARCH_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

// CreateExercise creates a custom exercise
// @Summary Create custom exercise
// @Description Create a new custom exercise for the user
// @Tags exercises
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateExerciseRequest true "Exercise data"
// @Success 201 {object} dto.ExerciseResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /exercises [post]
func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req CreateExerciseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    "INVALID_REQUEST",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    "VALIDATION_ERROR",
		})
		return
	}

	exercise, err := h.exerciseService.CreateExercise(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create exercise",
			Message: err.Error(),
			Code:    "CREATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

// CreateExerciseRequest represents a new exercise entry
type CreateExerciseRequest struct {
	Name         string   `json:"name" validate:"required"`
	Category     string   `json:"category" validate:"required"`
	MuscleGroup  string   `json:"muscle_group" validate:"required"`
	Equipment    string   `json:"equipment,omitempty"`
	Description  string   `json:"description,omitempty"`
	Instructions []string `json:"instructions,omitempty"`
}
