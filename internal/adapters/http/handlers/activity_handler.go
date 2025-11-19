package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// ActivityHandler handles activity-related requests
type ActivityHandler struct {
	activityService ports.ActivityService
	validator       *validator.Validate
}

// NewActivityHandler creates a new activity handler
func NewActivityHandler(activityService ports.ActivityService) *ActivityHandler {
	return &ActivityHandler{
		activityService: activityService,
		validator:       validator.New(),
	}
}

// GetActivities retrieves activities for a user
// @Summary Get user activities
// @Description Retrieve activities for the authenticated user with optional date filtering
// @Tags activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param type query string false "Filter by activity type"
// @Success 200 {array} dto.ActivityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /activities [get]
func (h *ActivityHandler) GetActivities(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	activityType := c.Query("type")

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

	activities, err := h.activityService.GetActivities(c.Request.Context(), userID.(string), startDate, endDate, activityType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve activities",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, activities)
}

// GetActivity retrieves a specific activity by ID
// @Summary Get activity by ID
// @Description Retrieve detailed information about a specific activity
// @Tags activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity ID"
// @Success 200 {object} dto.ActivityResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /activities/{id} [get]
func (h *ActivityHandler) GetActivity(c *gin.Context) {
	userID, _ := c.Get("userID")
	activityID := c.Param("id")

	activity, err := h.activityService.GetActivity(c.Request.Context(), userID.(string), activityID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "RETRIEVAL_FAILED"

		if err.Error() == "activity not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to retrieve activity",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// CreateActivity creates a new activity entry
// @Summary Create a new activity
// @Description Create a new activity entry for the user
// @Tags activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateActivityRequest true "Activity data"
// @Success 201 {object} dto.ActivityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /activities [post]
func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.CreateActivityRequest

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

	activity, err := h.activityService.CreateActivity(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create activity",
			Message: err.Error(),
			Code:    "CREATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// UpdateActivity updates an existing activity
// @Summary Update activity
// @Description Update an existing activity entry
// @Tags activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity ID"
// @Param request body dto.CreateActivityRequest true "Updated activity data"
// @Success 200 {object} dto.ActivityResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /activities/{id} [put]
func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	userID, _ := c.Get("userID")
	activityID := c.Param("id")
	var req dto.CreateActivityRequest

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

	activity, err := h.activityService.UpdateActivity(c.Request.Context(), userID.(string), activityID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "UPDATE_FAILED"

		if err.Error() == "activity not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to update activity",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// DeleteActivity deletes an activity
// @Summary Delete activity
// @Description Delete an activity entry
// @Tags activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Activity ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /activities/{id} [delete]
func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	userID, _ := c.Get("userID")
	activityID := c.Param("id")

	err := h.activityService.DeleteActivity(c.Request.Context(), userID.(string), activityID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "DELETE_FAILED"

		if err.Error() == "activity not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to delete activity",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// SyncGarmin syncs activities from Garmin Connect
// @Summary Sync Garmin activities
// @Description Sync activities from Garmin Connect (Not Implemented)
// @Tags activities
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 501 {object} dto.ErrorResponse
// @Router /activities/sync/garmin [post]
func (h *ActivityHandler) SyncGarmin(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, dto.ErrorResponse{
		Error:   "Not Implemented",
		Message: "Garmin sync functionality is not yet implemented",
		Code:    "NOT_IMPLEMENTED",
	})
}
