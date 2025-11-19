package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// WorkoutHandler handles workout-related requests
type WorkoutHandler struct {
	workoutService ports.WorkoutService
	validator      *validator.Validate
}

// NewWorkoutHandler creates a new workout handler
func NewWorkoutHandler(workoutService ports.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{
		workoutService: workoutService,
		validator:      validator.New(),
	}
}

// StartWorkout starts a new workout session
// @Summary Start a new workout
// @Description Start a new workout session
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.StartWorkoutRequest true "Workout data"
// @Success 201 {object} dto.WorkoutResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts/start [post]
func (h *WorkoutHandler) StartWorkout(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.StartWorkoutRequest

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

	workout, err := h.workoutService.StartWorkout(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to start workout",
			Message: err.Error(),
			Code:    "START_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, workout)
}

// GetWorkouts retrieves workouts for a user
// @Summary Get user workouts
// @Description Retrieve workouts for the authenticated user with optional date filtering
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param status query string false "Filter by status (in_progress, completed, cancelled)"
// @Success 200 {array} dto.WorkoutResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts [get]
func (h *WorkoutHandler) GetWorkouts(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	status := c.Query("status")

	var startDate, endDate *time.Time

	if startDateStr != "" {
		parsed, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid start_date format",
				Message: "Use YYYY-MM-DD format",
				Code:    "INVALID_DATE",
			})
			return
		}
		startDate = &parsed
	}

	if endDateStr != "" {
		parsed, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid end_date format",
				Message: "Use YYYY-MM-DD format",
				Code:    "INVALID_DATE",
			})
			return
		}
		endDate = &parsed
	}

	workouts, err := h.workoutService.GetWorkouts(c.Request.Context(), userID.(string), startDate, endDate, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve workouts",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, workouts)
}

// GetWorkout retrieves a specific workout by ID
// @Summary Get workout by ID
// @Description Retrieve detailed information about a specific workout
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout ID"
// @Success 200 {object} dto.WorkoutResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts/{id} [get]
func (h *WorkoutHandler) GetWorkout(c *gin.Context) {
	userID, _ := c.Get("userID")
	workoutID := c.Param("id")

	workout, err := h.workoutService.GetWorkout(c.Request.Context(), userID.(string), workoutID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "RETRIEVAL_FAILED"

		if err.Error() == "workout not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to retrieve workout",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, workout)
}

// FinishWorkout finishes an active workout
// @Summary Finish workout
// @Description Mark an active workout as completed
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout ID"
// @Success 200 {object} dto.WorkoutResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts/{id}/finish [post]
func (h *WorkoutHandler) FinishWorkout(c *gin.Context) {
	userID, _ := c.Get("userID")
	workoutID := c.Param("id")

	workout, err := h.workoutService.FinishWorkout(c.Request.Context(), userID.(string), workoutID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "FINISH_FAILED"

		if err.Error() == "workout not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to finish workout",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, workout)
}

// AddExercise adds an exercise to a workout
// @Summary Add exercise to workout
// @Description Add an exercise to an active workout
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout ID"
// @Param exercise_id query string true "Exercise ID"
// @Success 200 {object} dto.WorkoutResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts/{id}/exercises [post]
func (h *WorkoutHandler) AddExercise(c *gin.Context) {
	userID, _ := c.Get("userID")
	workoutID := c.Param("id")
	exerciseID := c.Query("exercise_id")

	if exerciseID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Missing exercise_id",
			Message: "exercise_id query parameter is required",
			Code:    "MISSING_PARAMETER",
		})
		return
	}

	workout, err := h.workoutService.AddExercise(c.Request.Context(), userID.(string), workoutID, exerciseID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "ADD_EXERCISE_FAILED"

		if err.Error() == "workout not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to add exercise",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, workout)
}

// LogSet logs a set for an exercise in a workout
// @Summary Log exercise set
// @Description Log a set for an exercise during a workout
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogSetRequest true "Set data"
// @Success 200 {object} dto.WorkoutResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts/sets [post]
func (h *WorkoutHandler) LogSet(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.LogSetRequest

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

	workout, err := h.workoutService.LogSet(c.Request.Context(), userID.(string), &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "LOG_SET_FAILED"

		if err.Error() == "workout not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to log set",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, workout)
}

// DeleteWorkout deletes a workout
// @Summary Delete workout
// @Description Delete a workout entry
// @Tags workouts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workout ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /workouts/{id} [delete]
func (h *WorkoutHandler) DeleteWorkout(c *gin.Context) {
	userID, _ := c.Get("userID")
	workoutID := c.Param("id")

	err := h.workoutService.DeleteWorkout(c.Request.Context(), userID.(string), workoutID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "DELETE_FAILED"

		if err.Error() == "workout not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to delete workout",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
