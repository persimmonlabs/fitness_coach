package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// GoalHandler handles goal-related requests
type GoalHandler struct {
	goalService ports.GoalService
	validator   *validator.Validate
}

// NewGoalHandler creates a new goal handler
func NewGoalHandler(goalService ports.GoalService) *GoalHandler {
	return &GoalHandler{
		goalService: goalService,
		validator:   validator.New(),
	}
}

// CreateGoal creates a new fitness goal
// @Summary Create a new goal
// @Description Create a new fitness goal for the user
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateGoalRequest true "Goal data"
// @Success 201 {object} dto.GoalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals [post]
func (h *GoalHandler) CreateGoal(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.CreateGoalRequest

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

	goal, err := h.goalService.CreateGoal(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create goal",
			Message: err.Error(),
			Code:    "CREATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, goal)
}

// GetGoals retrieves goals for a user
// @Summary Get user goals
// @Description Retrieve all goals for the authenticated user
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status (active, completed, failed, cancelled)"
// @Param type query string false "Filter by goal type"
// @Success 200 {array} dto.GoalResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals [get]
func (h *GoalHandler) GetGoals(c *gin.Context) {
	userID, _ := c.Get("userID")

	status := c.Query("status")
	goalType := c.Query("type")

	goals, err := h.goalService.GetGoals(c.Request.Context(), userID.(string), status, goalType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve goals",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, goals)
}

// UpdateGoal updates an existing goal
// @Summary Update goal
// @Description Update an existing fitness goal
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Goal ID"
// @Param request body dto.CreateGoalRequest true "Updated goal data"
// @Success 200 {object} dto.GoalResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals/{id} [put]
func (h *GoalHandler) UpdateGoal(c *gin.Context) {
	userID, _ := c.Get("userID")
	goalID := c.Param("id")
	var req dto.CreateGoalRequest

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

	goal, err := h.goalService.UpdateGoal(c.Request.Context(), userID.(string), goalID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "UPDATE_FAILED"

		if err.Error() == "goal not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to update goal",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, goal)
}

// DeleteGoal deletes a goal
// @Summary Delete goal
// @Description Delete a fitness goal
// @Tags goals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Goal ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /goals/{id} [delete]
func (h *GoalHandler) DeleteGoal(c *gin.Context) {
	userID, _ := c.Get("userID")
	goalID := c.Param("id")

	err := h.goalService.DeleteGoal(c.Request.Context(), userID.(string), goalID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "DELETE_FAILED"

		if err.Error() == "goal not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to delete goal",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
