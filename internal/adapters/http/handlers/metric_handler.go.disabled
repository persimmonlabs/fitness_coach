package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// MetricHandler handles body metric-related requests
type MetricHandler struct {
	metricService ports.MetricService
	validator     *validator.Validate
}

// NewMetricHandler creates a new metric handler
func NewMetricHandler(metricService ports.MetricService) *MetricHandler {
	return &MetricHandler{
		metricService: metricService,
		validator:     validator.New(),
	}
}

// LogMetric logs a body metric
// @Summary Log body metric
// @Description Log a body metric measurement (weight, body fat, etc.)
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.LogMetricRequest true "Metric data"
// @Success 201 {object} dto.MetricResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /metrics [post]
func (h *MetricHandler) LogMetric(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.LogMetricRequest

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

	// Set current time if not provided
	if req.RecordedAt.IsZero() {
		req.RecordedAt = time.Now()
	}

	metric, err := h.metricService.LogMetric(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to log metric",
			Message: err.Error(),
			Code:    "LOG_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, metric)
}

// GetMetricTrend retrieves metric trend data
// @Summary Get metric trend
// @Description Retrieve trend data for a specific metric type over a time period
// @Tags metrics
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type path string true "Metric type (weight, body_fat, muscle_mass, bmi, waist_circumference)"
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param limit query int false "Maximum number of results" default(30)
// @Success 200 {array} dto.MetricResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /metrics/{type}/trend [get]
func (h *MetricHandler) GetMetricTrend(c *gin.Context) {
	userID, _ := c.Get("userID")
	metricType := c.Param("type")

	// Validate metric type
	validTypes := map[string]bool{
		"weight":               true,
		"body_fat":             true,
		"muscle_mass":          true,
		"bmi":                  true,
		"waist_circumference":  true,
	}

	if !validTypes[metricType] {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Invalid metric type",
			Message: "Metric type must be one of: weight, body_fat, muscle_mass, bmi, waist_circumference",
			Code:    "INVALID_METRIC_TYPE",
		})
		return
	}

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	limit := 30

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

	metrics, err := h.metricService.GetMetricTrend(c.Request.Context(), userID.(string), metricType, startDate, endDate, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve metric trend",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}
