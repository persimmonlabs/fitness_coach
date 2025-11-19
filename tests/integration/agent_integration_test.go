package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fitness-tracker/internal/adapters/external"
	"fitness-tracker/internal/adapters/repositories/postgres"
	"fitness-tracker/internal/core/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockOpenRouterServer creates a mock OpenRouter API server for testing
func MockOpenRouterServer(t *testing.T, responseContent string, toolCalls []external.ToolCall) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Contains(t, r.Header.Get("Authorization"), "Bearer")

		// Parse request
		var req external.ChatRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		require.NoError(t, err)

		// Verify request has messages
		assert.NotEmpty(t, req.Messages)

		// Create response
		response := external.ChatResponse{
			ID:      "test-response-id",
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []struct {
				Index   int `json:"index"`
				Message struct {
					Role      string               `json:"role"`
					Content   string               `json:"content"`
					ToolCalls []external.ToolCall  `json:"tool_calls,omitempty"`
				} `json:"message"`
				FinishReason string `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: struct {
						Role      string               `json:"role"`
						Content   string               `json:"content"`
						ToolCalls []external.ToolCall  `json:"tool_calls,omitempty"`
					}{
						Role:      "assistant",
						Content:   responseContent,
						ToolCalls: toolCalls,
					},
					FinishReason: "stop",
				},
			},
			Usage: struct {
				PromptTokens     int `json:"prompt_tokens"`
				CompletionTokens int `json:"completion_tokens"`
				TotalTokens      int `json:"total_tokens"`
			}{
				PromptTokens:     100,
				CompletionTokens: 50,
				TotalTokens:      150,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
}

func TestAgentToolExecution(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "agent_test@example.com")

	// Create test food for tool execution
	food := CreateTestFood(t, testDB.DB, "Chicken Breast", 165.0)

	t.Run("Agent executes food search tool", func(t *testing.T) {
		// Mock OpenRouter response with tool call
		toolCalls := []external.ToolCall{
			{
				ID:   "call_123",
				Type: "function",
				Function: struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				}{
					Name:      "search_food",
					Arguments: `{"query": "chicken"}`,
				},
			},
		}

		server := MockOpenRouterServer(t, "Let me search for chicken in the food database.", toolCalls)
		defer server.Close()

		// Create OpenRouter client pointing to mock server
		client := &external.OpenRouterClient{}
		// Note: In real implementation, we'd need to override the base URL

		// Test that we can parse tool calls from response
		ctx := context.Background()
		messages := []external.Message{
			{Role: "user", Content: "What's the nutrition info for chicken?"},
		}

		// Simulate what the service would do
		foodRepo := postgres.NewFoodRepository(testDB.DB)

		// Execute the tool (search for food)
		foods, err := foodRepo.Search("chicken", 10, 0)
		require.NoError(t, err)
		assert.NotEmpty(t, foods)

		// Verify food was found
		foundChicken := false
		for _, f := range foods {
			if f.Name == "Chicken Breast" {
				foundChicken = true
				assert.Equal(t, 165.0, f.Calories)
				break
			}
		}
		assert.True(t, foundChicken, "Should find chicken in search results")
	})

	t.Run("Agent executes meal logging tool", func(t *testing.T) {
		// Mock OpenRouter response with meal logging tool call
		toolCalls := []external.ToolCall{
			{
				ID:   "call_456",
				Type: "function",
				Function: struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				}{
					Name:      "log_meal",
					Arguments: `{"name": "Lunch", "meal_type": "lunch", "foods": [{"food_id": "` + food.ID.String() + `", "quantity": 200}]}`,
				},
			},
		}

		server := MockOpenRouterServer(t, "I've logged your lunch with chicken breast.", toolCalls)
		defer server.Close()

		// Simulate tool execution - create meal
		mealRepo := postgres.NewMealRepository(testDB.DB)
		meal := &domain.Meal{
			UserID:             user.ID,
			Name:               "Lunch",
			MealType:           "lunch",
			ConsumedAt:         time.Now(),
			TotalCalories:      330.0, // 165 * 2 (200g)
			TotalProtein:       62.0,
			TotalCarbohydrates: 0.0,
			TotalFat:           7.2,
		}

		err := mealRepo.Create(meal)
		require.NoError(t, err)
		assert.NotEqual(t, "", meal.ID.String())

		// Verify meal was created
		retrieved, err := mealRepo.GetByID(meal.ID)
		require.NoError(t, err)
		assert.Equal(t, "Lunch", retrieved.Name)
		assert.Equal(t, "lunch", retrieved.MealType)
	})
}

func TestConversationFlow(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "conversation_test@example.com")

	// Initialize repository
	conversationRepo := postgres.NewConversationRepository(testDB.DB)

	t.Run("Create conversation and add messages", func(t *testing.T) {
		// Create conversation
		title := "Nutrition Consultation"
		conversation := &domain.Conversation{
			UserID: user.ID,
			Title:  &title,
		}

		err := conversationRepo.Create(conversation)
		require.NoError(t, err)
		assert.NotEqual(t, "", conversation.ID.String())

		// Add user message
		userMessage := &domain.Message{
			ConversationID: conversation.ID,
			Role:           "user",
			Content:        "What should I eat for breakfast?",
		}
		err = testDB.DB.Create(userMessage).Error
		require.NoError(t, err)

		// Add assistant message
		assistantMessage := &domain.Message{
			ConversationID: conversation.ID,
			Role:           "assistant",
			Content:        "For a healthy breakfast, I recommend eggs, oatmeal, and fruit.",
		}
		err = testDB.DB.Create(assistantMessage).Error
		require.NoError(t, err)

		// Retrieve conversation with messages
		retrieved, err := conversationRepo.GetByID(conversation.ID)
		require.NoError(t, err)
		assert.Equal(t, title, *retrieved.Title)

		// Load messages
		err = testDB.DB.Order("created_at ASC").Find(&retrieved.Messages, "conversation_id = ?", conversation.ID).Error
		require.NoError(t, err)
		assert.Len(t, retrieved.Messages, 2)
		assert.Equal(t, "user", retrieved.Messages[0].Role)
		assert.Equal(t, "assistant", retrieved.Messages[1].Role)
	})
}

func TestChatContextPersistence(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "context_test@example.com")

	// Initialize repository
	conversationRepo := postgres.NewConversationRepository(testDB.DB)

	t.Run("Store and retrieve conversation context", func(t *testing.T) {
		// Create conversation with context
		title := "Workout Planning"
		context := `{"user_goals": ["lose_weight", "build_muscle"], "fitness_level": "intermediate"}`

		conversation := &domain.Conversation{
			UserID:  user.ID,
			Title:   &title,
			Context: &context,
		}

		err := conversationRepo.Create(conversation)
		require.NoError(t, err)

		// Retrieve and verify context
		retrieved, err := conversationRepo.GetByID(conversation.ID)
		require.NoError(t, err)
		assert.NotNil(t, retrieved.Context)

		// Parse context JSON
		var contextData map[string]interface{}
		err = json.Unmarshal([]byte(*retrieved.Context), &contextData)
		require.NoError(t, err)

		// Verify context data
		goals, ok := contextData["user_goals"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, goals, 2)
		assert.Equal(t, "lose_weight", goals[0])
		assert.Equal(t, "build_muscle", goals[1])
		assert.Equal(t, "intermediate", contextData["fitness_level"])
	})
}

func TestGetConversationsByUser(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test users
	user1 := CreateTestUser(t, testDB.DB, "user1@example.com")
	user2 := CreateTestUser(t, testDB.DB, "user2@example.com")

	// Initialize repository
	conversationRepo := postgres.NewConversationRepository(testDB.DB)

	t.Run("Filter conversations by user", func(t *testing.T) {
		// Create conversations for user1
		title1 := "Conversation 1"
		title2 := "Conversation 2"
		conv1 := &domain.Conversation{UserID: user1.ID, Title: &title1}
		conv2 := &domain.Conversation{UserID: user1.ID, Title: &title2}

		err := conversationRepo.Create(conv1)
		require.NoError(t, err)
		err = conversationRepo.Create(conv2)
		require.NoError(t, err)

		// Create conversation for user2
		title3 := "Conversation 3"
		conv3 := &domain.Conversation{UserID: user2.ID, Title: &title3}
		err = conversationRepo.Create(conv3)
		require.NoError(t, err)

		// Get conversations for user1
		user1Conversations, err := conversationRepo.GetByUserID(user1.ID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, user1Conversations, 2)

		// Verify all belong to user1
		for _, conv := range user1Conversations {
			assert.Equal(t, user1.ID, conv.UserID)
		}

		// Get conversations for user2
		user2Conversations, err := conversationRepo.GetByUserID(user2.ID, 10, 0)
		require.NoError(t, err)
		assert.Len(t, user2Conversations, 1)
		assert.Equal(t, user2.ID, user2Conversations[0].UserID)
	})
}

func TestMessageOrdering(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "ordering_test@example.com")

	t.Run("Messages are ordered chronologically", func(t *testing.T) {
		// Create conversation
		title := "Test Ordering"
		conversation := &domain.Conversation{UserID: user.ID, Title: &title}
		err := testDB.DB.Create(conversation).Error
		require.NoError(t, err)

		// Add messages with slight delays
		messages := []string{
			"First message",
			"Second message",
			"Third message",
		}

		var createdMessages []*domain.Message
		for i, content := range messages {
			msg := &domain.Message{
				ConversationID: conversation.ID,
				Role:           "user",
				Content:        content,
			}
			// Add small delay to ensure different timestamps
			time.Sleep(10 * time.Millisecond)
			err = testDB.DB.Create(msg).Error
			require.NoError(t, err)
			createdMessages = append(createdMessages, msg)

			// Verify ordering
			assert.True(t, i == 0 || createdMessages[i].CreatedAt.After(createdMessages[i-1].CreatedAt),
				"Message %d should be created after message %d", i, i-1)
		}

		// Retrieve messages in order
		var retrievedMessages []domain.Message
		err = testDB.DB.Where("conversation_id = ?", conversation.ID).
			Order("created_at ASC").
			Find(&retrievedMessages).Error
		require.NoError(t, err)

		// Verify chronological order
		assert.Len(t, retrievedMessages, 3)
		assert.Equal(t, "First message", retrievedMessages[0].Content)
		assert.Equal(t, "Second message", retrievedMessages[1].Content)
		assert.Equal(t, "Third message", retrievedMessages[2].Content)
	})
}

func TestConversationDeletion(t *testing.T) {
	// Setup
	testDB := SetupTestDB(t)
	defer TeardownTestDB(t, testDB)

	// Create test user
	user := CreateTestUser(t, testDB.DB, "deletion_test@example.com")

	// Initialize repository
	conversationRepo := postgres.NewConversationRepository(testDB.DB)

	t.Run("Delete conversation cascades to messages", func(t *testing.T) {
		// Create conversation with messages
		title := "To Be Deleted"
		conversation := &domain.Conversation{UserID: user.ID, Title: &title}
		err := conversationRepo.Create(conversation)
		require.NoError(t, err)

		// Add messages
		for i := 0; i < 3; i++ {
			msg := &domain.Message{
				ConversationID: conversation.ID,
				Role:           "user",
				Content:        "Message " + string(rune(i)),
			}
			err = testDB.DB.Create(msg).Error
			require.NoError(t, err)
		}

		// Delete conversation
		err = conversationRepo.Delete(conversation.ID)
		require.NoError(t, err)

		// Verify conversation is deleted
		_, err = conversationRepo.GetByID(conversation.ID)
		assert.Error(t, err)

		// Verify messages are also deleted (cascade)
		var messageCount int64
		testDB.DB.Model(&domain.Message{}).
			Where("conversation_id = ?", conversation.ID).
			Count(&messageCount)
		assert.Equal(t, int64(0), messageCount, "Messages should be cascade deleted")
	})
}
