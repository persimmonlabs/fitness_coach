package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// SummaryHandler handles daily summary requests
type SummaryHandler struct {
	summaryService ports.SummaryService
}

// NewSummaryHandler creates a new summary handler
func NewSummaryHandler(summaryService ports.SummaryService) *SummaryHandler {
	return &SummaryHandler{
		summaryService: summaryService,
	}
}

// GetDailySummary retrieves daily fitness summary
// @Summary Get daily summary
// @Description Retrieve a comprehensive summary of fitness data for a specific date
// @Tags summary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param date query string false "Date (YYYY-MM-DD), defaults to today"
// @Success 200 {object} dto.DailySummaryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /summary/daily [get]
func (h *SummaryHandler) GetDailySummary(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Parse date parameter
	dateStr := c.Query("date")
	var date time.Time

	if dateStr != "" {
		parsed, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid date format",
				Message: "Use YYYY-MM-DD format",
				Code:    "INVALID_DATE",
			})
			return
		}
		date = parsed
	} else {
		// Default to today
		date = time.Now()
	}

	summary, err := h.summaryService.GetDailySummary(c.Request.Context(), userID.(string), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve daily summary",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, summary)
}
