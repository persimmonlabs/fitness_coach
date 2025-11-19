package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"fitness-tracker/internal/adapters/http/dto"
	"fitness-tracker/internal/core/ports"
)

// ChatHandler handles AI coach chat requests
type ChatHandler struct {
	chatService ports.ChatService
	validator   *validator.Validate
}

// NewChatHandler creates a new chat handler
func NewChatHandler(chatService ports.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		validator:   validator.New(),
	}
}

// SendMessage sends a message to the AI coach
// @Summary Send message to AI coach
// @Description Send a message to the AI fitness coach and get a response
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChatRequest true "Chat message"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /chat [post]
func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req dto.ChatRequest

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

	response, err := h.chatService.SendMessage(c.Request.Context(), userID.(string), req.Message, req.Context)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to process message",
			Message: err.Error(),
			Code:    "CHAT_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetHistory retrieves chat conversation history
// @Summary Get chat history
// @Description Retrieve conversation history with the AI coach
// @Tags chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Maximum number of messages" default(50)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} dto.ChatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /chat/history [get]
func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("userID")

	limit := 50
	offset := 0

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

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if _, err := fmt.Sscanf(offsetStr, "%d", &offset); err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "Invalid offset parameter",
				Message: "Offset must be a valid integer",
				Code:    "INVALID_OFFSET",
			})
			return
		}
	}

	history, err := h.chatService.GetHistory(c.Request.Context(), userID.(string), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Failed to retrieve chat history",
			Message: err.Error(),
			Code:    "RETRIEVAL_FAILED",
		})
		return
	}

	c.JSON(http.StatusOK, history)
}
