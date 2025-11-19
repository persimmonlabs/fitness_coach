package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// FoodHandler handles food-related requests
type FoodHandler struct {
	foodService ports.FoodService
	validator   *validator.Validate
}

// NewFoodHandler creates a new food handler
func NewFoodHandler(foodService ports.FoodService) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
		validator:   validator.New(),
	}
}

// SearchFoods searches for foods by name or barcode
// @Summary Search foods
// @Description Search for foods by name or barcode in the database
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param query query string false "Search query (name or barcode)"
// @Param limit query int false "Results limit" default(20)
// @Success 200 {array} dto.FoodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /foods/search [get]
func (h *FoodHandler) SearchFoods(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Missing query parameter",
			Message: "Query parameter is required for search",
			Code:    "MISSING_QUERY",
		})
		return
	}

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

	foods, err := h.foodService.SearchFoods(c.Request.Context(), query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Search failed",
			Message: err.Error(),
			Code:    "SEARCH_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, foods)
}

// GetFood retrieves a specific food by ID
// @Summary Get food by ID
// @Description Retrieve detailed information about a specific food item
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID"
// @Success 200 {object} dto.FoodResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /foods/{id} [get]
func (h *FoodHandler) GetFood(c *gin.Context) {
	foodID := c.Param("id")

	food, err := h.foodService.GetFood(c.Request.Context(), foodID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "RETRIEVAL_FAILED"

		if err.Error() == "food not found" {
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to retrieve food",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, food)
}

// CreateFood creates a custom food entry
// @Summary Create custom food
// @Description Create a new custom food entry for the user
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateFoodRequest true "Food data"
// @Success 201 {object} dto.FoodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /foods [post]
func (h *FoodHandler) CreateFood(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.CreateFoodRequest

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

	food, err := h.foodService.CreateFood(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to create food",
			Message: err.Error(),
			Code:    "CREATE_FAILED",
		})
		return
	}

	c.JSON(http.StatusCreated, food)
}

// UpdateFood updates a custom food entry
// @Summary Update custom food
// @Description Update an existing custom food entry
// @Tags foods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Food ID"
// @Param request body dto.CreateFoodRequest true "Updated food data"
// @Success 200 {object} dto.FoodResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /foods/{id} [put]
func (h *FoodHandler) UpdateFood(c *gin.Context) {
	userID, _ := c.Get("userID")
	foodID := c.Param("id")
	var req dto.CreateFoodRequest

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

	food, err := h.foodService.UpdateFood(c.Request.Context(), userID.(string), foodID, &req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorCode := "UPDATE_FAILED"

		switch err.Error() {
		case "food not found":
			statusCode = http.StatusNotFound
			errorCode = "NOT_FOUND"
		case "cannot update non-custom food":
			statusCode = http.StatusForbidden
			errorCode = "FORBIDDEN"
		}

		c.JSON(statusCode, dto.ErrorResponse{
			Error:   "Failed to update food",
			Message: err.Error(),
			Code:    errorCode,
		})
		return
	}

	c.JSON(http.StatusOK, food)
}
