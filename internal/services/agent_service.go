package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"fitness-tracker/internal/adapters/external"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"
)

// AgentService handles AI agent interactions with tool support
type AgentService struct {
	// Service dependencies
	mealService     ports.MealService
	foodService     ports.FoodService
	activityService ports.ActivityService
	workoutService  ports.WorkoutService
	metricService   ports.MetricService
	goalService     ports.GoalService
	summaryService  ports.SummaryService

	// Repository dependencies
	conversationRepo ports.ConversationRepository
	userRepo         ports.UserRepository

	// External client
	openRouterClient *external.OpenRouterClient

	// Configuration
	defaultModel string
}

// AgentResponse represents the response from the AI agent
type AgentResponse struct {
	Message    string    `json:"message"`
	ToolsUsed  []string  `json:"tools_used"`
	Confidence float64   `json:"confidence"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewAgentService creates a new agent service
func NewAgentService(
	mealService ports.MealService,
	foodService ports.FoodService,
	activityService ports.ActivityService,
	workoutService ports.WorkoutService,
	metricService ports.MetricService,
	goalService ports.GoalService,
	summaryService ports.SummaryService,
	conversationRepo ports.ConversationRepository,
	userRepo ports.UserRepository,
	openRouterClient *external.OpenRouterClient,
) *AgentService {
	return &AgentService{
		mealService:      mealService,
		foodService:      foodService,
		activityService:  activityService,
		workoutService:   workoutService,
		metricService:    metricService,
		goalService:      goalService,
		summaryService:   summaryService,
		conversationRepo: conversationRepo,
		userRepo:         userRepo,
		openRouterClient: openRouterClient,
		defaultModel:     "deepseek/deepseek-chat",
	}
}

// SendMessage processes a user message and returns an AI response
func (s *AgentService) SendMessage(ctx context.Context, userID uuid.UUID, message string) (*AgentResponse, error) {
	log.Printf("[AgentService] Processing message for user %s", userID)

	// Get or create conversation
	conversation, err := s.getOrCreateConversation(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get conversation: %w", err)
	}

	// Load last 20 messages for context
	messages, err := s.conversationRepo.GetLatestMessages(ctx, conversation.ID, 20)
	if err != nil {
		return nil, fmt.Errorf("failed to load messages: %w", err)
	}

	// Build user context
	userContext, err := s.buildUserContext(ctx, userID)
	if err != nil {
		log.Printf("[AgentService] Warning: failed to build user context: %v", err)
		userContext = "User context unavailable"
	}

	// Build system prompt
	systemPrompt := s.buildSystemPrompt(userContext)

	// Convert messages to OpenRouter format
	chatMessages := []external.Message{
		{Role: "system", Content: systemPrompt},
	}

	for _, msg := range messages {
		chatMessages = append(chatMessages, external.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Add current user message
	chatMessages = append(chatMessages, external.Message{
		Role:    "user",
		Content: message,
	})

	// Build tool definitions
	toolDefs := s.buildToolDefinitions()

	// Execute LLM call with tools
	response, toolsUsed, err := s.executeWithTools(ctx, chatMessages, toolDefs, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute LLM: %w", err)
	}

	// Save user message
	userMsg := &domain.Message{
		ID:             uuid.New(),
		ConversationID: conversation.ID,
		Role:           "user",
		Content:        message,
		CreatedAt:      time.Now(),
	}
	if err := s.conversationRepo.AddMessage(ctx, userMsg); err != nil {
		log.Printf("[AgentService] Warning: failed to save user message: %v", err)
	}

	// Save assistant response
	assistantMsg := &domain.Message{
		ID:             uuid.New(),
		ConversationID: conversation.ID,
		Role:           "assistant",
		Content:        response,
		CreatedAt:      time.Now(),
	}
	if len(toolsUsed) > 0 {
		metadata := map[string]interface{}{
			"tools_used": toolsUsed,
		}
		metadataJSON, _ := json.Marshal(metadata)
		metadataStr := string(metadataJSON)
		assistantMsg.Metadata = &metadataStr
	}
	if err := s.conversationRepo.AddMessage(ctx, assistantMsg); err != nil {
		log.Printf("[AgentService] Warning: failed to save assistant message: %v", err)
	}

	return &AgentResponse{
		Message:    response,
		ToolsUsed:  toolsUsed,
		Confidence: 0.85,
		CreatedAt:  time.Now(),
	}, nil
}

// getOrCreateConversation gets the most recent conversation or creates a new one
func (s *AgentService) getOrCreateConversation(ctx context.Context, userID uuid.UUID) (*domain.Conversation, error) {
	conversations, err := s.conversationRepo.ListByUser(ctx, userID, 1, 0)
	if err != nil {
		return nil, err
	}

	if len(conversations) > 0 {
		return conversations[0], nil
	}

	// Create new conversation
	title := "New Conversation"
	conversation := &domain.Conversation{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     &title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.conversationRepo.Create(ctx, conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

// buildUserContext builds context information about the user
func (s *AgentService) buildUserContext(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Get active goals
	goals, err := s.goalService.GetGoals(ctx, userID.String(), stringPtr("active"))
	if err != nil {
		log.Printf("[AgentService] Warning: failed to get goals: %v", err)
		goals = []*domain.Goal{}
	}

	// Get today's summary
	today := time.Now()
	summary, err := s.summaryService.GetDailySummary(ctx, userID.String(), today)
	if err != nil {
		log.Printf("[AgentService] Warning: failed to get daily summary: %v", err)
	}

	// Get recent activities (last 7 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -7)
	activities, err := s.activityService.GetActivities(ctx, userID.String(), &startDate, &endDate)
	if err != nil {
		log.Printf("[AgentService] Warning: failed to get activities: %v", err)
		activities = []*domain.Activity{}
	}

	// Build context string
	context := fmt.Sprintf("User: %s %s\n", user.FirstName, user.LastName)

	if len(goals) > 0 {
		context += "\nActive Goals:\n"
		for _, goal := range goals {
			context += fmt.Sprintf("- %s: %s (target: %.1f %s)\n",
				goal.GoalType, goal.Description, goal.TargetValue, goal.Unit)
		}
	}

	if summary != nil {
		context += fmt.Sprintf("\nToday's Nutrition:\n")
		context += fmt.Sprintf("- Calories: %.0f / %.0f\n", summary.TotalCalories, 2000)
		context += fmt.Sprintf("- Protein: %.1fg / %.1fg\n", summary.TotalProtein, 150)
		context += fmt.Sprintf("- Carbs: %.1fg / %.1fg\n", summary.TotalCarbohydrates, 200)
		context += fmt.Sprintf("- Fat: %.1fg / %.1fg\n", summary.TotalFat, 65)
	}

	if len(activities) > 0 {
		context += fmt.Sprintf("\nRecent Activity (last 7 days): %d activities logged\n", len(activities))
	}

	return context, nil
}

// buildSystemPrompt creates the system prompt with user context
func (s *AgentService) buildSystemPrompt(userContext string) string {
	return fmt.Sprintf(`You are a fitness and nutrition coach assistant with access to the user's tracking data.

User Context:
%s

Guidelines:
- Be concise but thorough
- Reference user's actual data when relevant
- Provide evidence-based advice
- ALWAYS use tools to get accurate data before answering
- Never hallucinate meal or workout history

When user asks about progress, meals, or workouts, use the appropriate tool first.`, userContext)
}

// buildToolDefinitions creates tool definitions for function calling
func (s *AgentService) buildToolDefinitions() []external.Tool {
	return []external.Tool{
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "log_meal",
				Description: "Log a meal with food items",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"food_items": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"food_id":  map[string]string{"type": "string"},
									"quantity": map[string]string{"type": "number"},
									"unit":     map[string]string{"type": "string"},
								},
							},
						},
						"meal_type": map[string]interface{}{
							"type": "string",
							"enum": []string{"breakfast", "lunch", "dinner", "snack"},
						},
						"timestamp": map[string]string{"type": "string"},
					},
					"required": []string{"food_items", "meal_type"},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "get_recent_meals",
				Description: "Get user's recent meals",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"days": map[string]interface{}{
							"type":        "integer",
							"description": "Number of days to look back",
							"default":     7,
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "search_foods",
				Description: "Search for foods in the database",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "Search query for food name",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "calculate_daily_macros",
				Description: "Calculate daily macro totals for a specific date",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Date in YYYY-MM-DD format",
						},
					},
					"required": []string{"date"},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "get_recent_workouts",
				Description: "Get user's recent workouts",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"days": map[string]interface{}{
							"type":        "integer",
							"description": "Number of days to look back",
							"default":     7,
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "get_recent_activities",
				Description: "Get user's recent activities",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"days": map[string]interface{}{
							"type":        "integer",
							"description": "Number of days to look back",
							"default":     7,
						},
					},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "log_weight",
				Description: "Log a weight measurement",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"weight": map[string]interface{}{
							"type":        "number",
							"description": "Weight in kg",
						},
						"date": map[string]interface{}{
							"type":        "string",
							"description": "Date in YYYY-MM-DD format",
						},
					},
					"required": []string{"weight"},
				},
			},
		},
		{
			Type: "function",
			Function: external.ToolFunction{
				Name:        "get_weight_trend",
				Description: "Get weight trend data",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"days": map[string]interface{}{
							"type":        "integer",
							"description": "Number of days to look back",
							"default":     30,
						},
					},
				},
			},
		},
	}
}

// executeWithTools executes the LLM call with tool support
func (s *AgentService) executeWithTools(ctx context.Context, messages []external.Message, toolDefs []external.Tool, userID uuid.UUID) (string, []string, error) {
	toolsUsed := []string{}
	maxIterations := 5

	for i := 0; i < maxIterations; i++ {
		// Call OpenRouter with tools
		response, err := s.openRouterClient.ChatWithTools(ctx, messages, toolDefs, s.defaultModel)
		if err != nil {
			return "", toolsUsed, fmt.Errorf("OpenRouter API call failed: %w", err)
		}

		if len(response.Choices) == 0 {
			return "", toolsUsed, fmt.Errorf("no response choices returned")
		}

		choice := response.Choices[0]

		// Check if we have tool calls
		if len(choice.Message.ToolCalls) == 0 {
			// No more tool calls, return final response
			return choice.Message.Content, toolsUsed, nil
		}

		// Execute tool calls
		for _, toolCall := range choice.Message.ToolCalls {
			log.Printf("[AgentService] Executing tool: %s with args: %s", toolCall.Function.Name, toolCall.Function.Arguments)

			result, err := s.executeTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments, userID)
			if err != nil {
				log.Printf("[AgentService] Tool execution failed: %v", err)
				result = fmt.Sprintf("Error: %v", err)
			}

			toolsUsed = append(toolsUsed, toolCall.Function.Name)

			// Add tool result to messages
			messages = append(messages, external.Message{
				Role:    "assistant",
				Content: "",
			})
			messages = append(messages, external.Message{
				Role:    "tool",
				Content: result,
			})
		}
	}

	return "Maximum tool iterations reached", toolsUsed, nil
}

// executeTool executes a specific tool function
func (s *AgentService) executeTool(ctx context.Context, toolName, arguments string, userID uuid.UUID) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	switch toolName {
	case "log_meal":
		return s.toolLogMeal(ctx, args, userID)
	case "get_recent_meals":
		return s.toolGetRecentMeals(ctx, args, userID)
	case "search_foods":
		return s.toolSearchFoods(ctx, args, userID)
	case "calculate_daily_macros":
		return s.toolCalculateDailyMacros(ctx, args, userID)
	case "get_recent_workouts":
		return s.toolGetRecentWorkouts(ctx, args, userID)
	case "get_recent_activities":
		return s.toolGetRecentActivities(ctx, args, userID)
	case "log_weight":
		return s.toolLogWeight(ctx, args, userID)
	case "get_weight_trend":
		return s.toolGetWeightTrend(ctx, args, userID)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

// Tool implementations

func (s *AgentService) toolLogMeal(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	// Implementation would create a meal using MealService
	return "Meal logging not yet implemented", nil
}

func (s *AgentService) toolGetRecentMeals(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	days := 7
	if d, ok := args["days"].(float64); ok {
		days = int(d)
	}

	endDate := time.Now()
	_ = endDate.AddDate(0, 0, -days)

	// Note: The actual implementation would need a GetMeals method that accepts date range
	// For now, we'll return a placeholder
	return fmt.Sprintf("Retrieved meals from last %d days", days), nil
}

func (s *AgentService) toolSearchFoods(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	query, ok := args["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter required")
	}

	foods, err := s.foodService.SearchFoods(ctx, query, nil, 10)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Found %d foods matching '%s':\n", len(foods), query)
	for i, food := range foods {
		if i >= 10 {
			break
		}
		result += fmt.Sprintf("- %s (%.0f cal, %.1fg protein, %.1fg carbs, %.1fg fat per %s)\n",
			food.Name, food.Calories, food.Protein, food.Carbohydrates, food.Fat, food.ServingUnit)
	}

	return result, nil
}

func (s *AgentService) toolCalculateDailyMacros(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	dateStr, ok := args["date"].(string)
	if !ok {
		return "", fmt.Errorf("date parameter required")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %w", err)
	}

	summary, err := s.summaryService.GetDailySummary(ctx, userID.String(), date)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Daily macros for %s:\n", dateStr)
	result += fmt.Sprintf("- Calories: %.0f / %.0f\n", summary.TotalCalories, 2000)
	result += fmt.Sprintf("- Protein: %.1fg / %.1fg\n", summary.TotalProtein, 150)
	result += fmt.Sprintf("- Carbs: %.1fg / %.1fg\n", summary.TotalCarbohydrates, 200)
	result += fmt.Sprintf("- Fat: %.1fg / %.1fg\n", summary.TotalFat, 65)

	return result, nil
}

func (s *AgentService) toolGetRecentWorkouts(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	days := 7
	if d, ok := args["days"].(float64); ok {
		days = int(d)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	workouts, err := s.workoutService.GetWorkouts(ctx, userID.String(), &startDate, &endDate)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Found %d workouts in the last %d days:\n", len(workouts), days)
	for _, workout := range workouts {
		duration := ""
		if workout.DurationMinutes != nil {
			duration = fmt.Sprintf(" (%d min)", *workout.DurationMinutes)
		}
		result += fmt.Sprintf("- %s on %s%s\n", workout.Name, workout.StartTime.Format("2006-01-02"), duration)
	}

	return result, nil
}

func (s *AgentService) toolGetRecentActivities(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	days := 7
	if d, ok := args["days"].(float64); ok {
		days = int(d)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	activities, err := s.activityService.GetActivities(ctx, userID.String(), &startDate, &endDate)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Found %d activities in the last %d days:\n", len(activities), days)
	for _, activity := range activities {
		duration := ""
		if activity.DurationMinutes != nil {
			duration = fmt.Sprintf(" (%d min)", *activity.DurationMinutes)
		}
		calories := ""
		if activity.CaloriesBurned != nil {
			calories = fmt.Sprintf(", %.0f cal", *activity.CaloriesBurned)
		}
		result += fmt.Sprintf("- %s on %s%s%s\n", activity.ActivityType, activity.StartTime.Format("2006-01-02"), duration, calories)
	}

	return result, nil
}

func (s *AgentService) toolLogWeight(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	weight, ok := args["weight"].(float64)
	if !ok {
		return "", fmt.Errorf("weight parameter required")
	}

	date := time.Now()
	if dateStr, ok := args["date"].(string); ok {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = parsedDate
		}
	}

	_, err := s.metricService.LogMetric(ctx, userID.String(), "weight", weight, "kg", date)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Logged weight: %.1f kg on %s", weight, date.Format("2006-01-02")), nil
}

func (s *AgentService) toolGetWeightTrend(ctx context.Context, args map[string]interface{}, userID uuid.UUID) (string, error) {
	days := 30
	if d, ok := args["days"].(float64); ok {
		days = int(d)
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	metrics, err := s.metricService.GetMetricTrend(ctx, userID.String(), "weight", &startDate, &endDate)
	if err != nil {
		return "", err
	}

	if len(metrics) == 0 {
		return fmt.Sprintf("No weight data found for the last %d days", days), nil
	}

	result := fmt.Sprintf("Weight trend (last %d days, %d measurements):\n", days, len(metrics))
	for _, metric := range metrics {
		result += fmt.Sprintf("- %s: %.1f kg\n", metric.MeasuredAt.Format("2006-01-02"), metric.Value)
	}

	// Calculate trend
	if len(metrics) >= 2 {
		first := metrics[0].Value
		last := metrics[len(metrics)-1].Value
		change := last - first
		result += fmt.Sprintf("\nChange: %.1f kg", change)
	}

	return result, nil
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}
