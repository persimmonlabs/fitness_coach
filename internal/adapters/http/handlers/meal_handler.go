package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// MealHandler handles meal-related requests
type MealHandler struct {
	mealService ports.MealService
	validator   *validator.Validate
}

// NewMealHandler creates a new meal handler
func NewMealHandler(mealService ports.MealService) *MealHandler {
	return &MealHandler{
		mealService: mealService,
		validator:   validator.New(),
	}
}

// ParseMeal handles natural language meal parsing
// @Summary Parse meal from natural language
// @Description Parse meal description using AI to extract foods and nutritional information
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ParseMealRequest true "Meal description"
// @Success 200 {object} dto.ParsedMealResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals/parse [post]
func (h *MealHandler) ParseMeal(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.ParseMealRequest

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

	parsedMeal, err := h.mealService.ParseMeal(c.Request.Context(), userID.(string), req.Description, req.MealType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to parse meal",
			Message: err.Error(),
			Code:    "PARSE_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, parsedMeal)
}

// ConfirmMeal confirms and saves a parsed meal
// @Summary Confirm parsed meal
// @Description Confirm and save a previously parsed meal to the database
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ConfirmMealRequest true "Meal confirmation data"
// @Success 201 {object} dto.MealResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals/confirm [post]
func (h *MealHandler) ConfirmMeal(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.ConfirmMealRequest

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

	meal, err := h.mealService.ConfirmMeal(c.Request.Context(), userID.(string), req.ParsedMealID, req.Adjustments)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to confirm meal",
			Message: err.Error(),
			Code:    "CONFIRM_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, meal)
}

// CreateMeal handles manual meal creation
// @Summary Create a new meal
// @Description Manually create a meal entry with specified foods
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateMealRequest true "Meal data"
// @Success 201 {object} dto.MealResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals [post]
func (h *MealHandler) CreateMeal(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.CreateMealRequest

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

	meal, err := h.mealService.CreateMeal(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create meal",
			Message: err.Error(),
			Code:    "CREATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, meal)
}

// GetMeals retrieves meals for a user
// @Summary Get user meals
// @Description Retrieve meals for the authenticated user with optional date filtering
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Param meal_type query string false "Filter by meal type"
// @Success 200 {array} dto.MealResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals [get]
func (h *MealHandler) GetMeals(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Parse query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	mealType := c.Query("meal_type")

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

	meals, err := h.mealService.GetMeals(c.Request.Context(), userID.(string), startDate, endDate, mealType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve meals",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, meals)
}

// GetMeal retrieves a specific meal by ID
// @Summary Get meal by ID
// @Description Retrieve detailed information about a specific meal
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal ID"
// @Success 200 {object} dto.MealResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals/{id} [get]
func (h *MealHandler) GetMeal(c *gin.Context) {
	userID, _ := c.Get("userID")
	mealID := c.Param("id")

	meal, err := h.mealService.GetMeal(c.Request.Context(), userID.(string), mealID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "RETRIEVAL_FAILED"

		if err.Error() == "meal not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to retrieve meal",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, meal)
}

// UpdateMeal updates an existing meal
// @Summary Update meal
// @Description Update an existing meal entry
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal ID"
// @Param request body dto.CreateMealRequest true "Updated meal data"
// @Success 200 {object} dto.MealResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals/{id} [put]
func (h *MealHandler) UpdateMeal(c *gin.Context) {
	userID, _ := c.Get("userID")
	mealID := c.Param("id")
	var req dto.CreateMealRequest

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

	meal, err := h.mealService.UpdateMeal(c.Request.Context(), userID.(string), mealID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "UPDATE_FAILED"

		if err.Error() == "meal not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to update meal",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, meal)
}

// DeleteMeal deletes a meal
// @Summary Delete meal
// @Description Delete a meal entry
// @Tags meals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Meal ID"
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /meals/{id} [delete]
func (h *MealHandler) DeleteMeal(c *gin.Context) {
	userID, _ := c.Get("userID")
	mealID := c.Param("id")

	err := h.mealService.DeleteMeal(c.Request.Context(), userID.(string), mealID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "DELETE_FAILED"

		if err.Error() == "meal not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to delete meal",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.Status(http.StatusNoContent)
}
