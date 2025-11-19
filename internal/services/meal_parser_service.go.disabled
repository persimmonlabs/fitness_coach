package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fitness-tracker/internal/adapters/external"
	"fitness-tracker/internal/core/domain"
	"fitness-tracker/internal/core/ports"

	"github.com/google/uuid"
)

// MealParserService handles parsing meals from text and photos
type MealParserService struct {
	openRouterClient *external.OpenRouterClient
	visionClient     *external.VisionClient
	foodRepository   ports.FoodRepository
}

// NewMealParserService creates a new meal parser service
func NewMealParserService(apiKey string, foodRepo ports.FoodRepository) *MealParserService {
	return &MealParserService{
		openRouterClient: external.NewOpenRouterClient(apiKey),
		visionClient:     external.NewVisionClient(apiKey),
		foodRepository:   foodRepo,
	}
}

// ExtractedFoodItem represents a food item extracted from AI
type ExtractedFoodItem struct {
	Name       string  `json:"name"`
	Quantity   float64 `json:"quantity"`
	Unit       string  `json:"unit"`
	Confidence float64 `json:"confidence"`
}

// ParseText parses meal information from text input
func (s *MealParserService) ParseText(ctx context.Context, userID uuid.UUID, text string) (*domain.ParsedMeal, error) {
	// System prompt for food extraction
	systemPrompt := `You are a nutrition expert. Extract food items, quantities, and meal type from the user's text.
Return a JSON object with:
{
  "meal_type": "breakfast|lunch|dinner|snack",
  "items": [
    {
      "name": "food name",
      "quantity": numeric_amount,
      "unit": "g|ml|cup|piece|tbsp|tsp",
      "confidence": 0.0-1.0
    }
  ]
}

If meal type cannot be determined, infer from context or time of day. Use standard units (prefer grams for solids, ml for liquids).`

	// Get AI response using OpenRouter client
	messages := []external.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: text},
	}

	resp, err := s.openRouterClient.Chat(ctx, messages, "deepseek/deepseek-chat")
	if err != nil {
		return nil, fmt.Errorf("failed to parse text with AI: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	response := resp.Choices[0].Message.Content

	// Parse AI response
	var aiResponse struct {
		MealType string              `json:"meal_type"`
		Items    []ExtractedFoodItem `json:"items"`
	}

	// Clean response (remove markdown code blocks if present)
	cleanedResponse := strings.TrimSpace(response)
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```json")
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	if err := json.Unmarshal([]byte(cleanedResponse), &aiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w (response: %s)", err, response)
	}

	// Process each food item
	parsedItems := make([]domain.ParsedFoodItem, 0, len(aiResponse.Items))
	totalConfidence := 0.0

	for _, item := range aiResponse.Items {
		parsedItem, err := s.processFoodItem(ctx, userID, item)
		if err != nil {
			// Log error but continue processing other items
			fmt.Printf("Warning: failed to process food item %s: %v\n", item.Name, err)
			continue
		}
		parsedItems = append(parsedItems, parsedItem)
		totalConfidence += item.Confidence
	}

	if len(parsedItems) == 0 {
		return nil, fmt.Errorf("no valid food items could be extracted")
	}

	avgConfidence := totalConfidence / float64(len(parsedItems))

	return &domain.ParsedMeal{
		MealType:          aiResponse.MealType,
		LoggedAt:          time.Now(),
		FoodItems:         parsedItems,
		Confidence:        avgConfidence,
		NeedsConfirmation: avgConfidence < 0.8, // Require confirmation if confidence is low
	}, nil
}

// ParsePhoto parses meal information from photo input
func (s *MealParserService) ParsePhoto(ctx context.Context, userID uuid.UUID, photoURL string) (*domain.ParsedMeal, error) {
	// Analyze image with vision AI
	result, err := s.visionClient.AnalyzeFoodPhoto(ctx, photoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	// Convert vision result to extracted items
	extractedItems := make([]ExtractedFoodItem, len(result.Items))
	for i, item := range result.Items {
		confidence := item.Confidence
		if confidence == 0 {
			confidence = 0.7 // Default confidence for vision results without explicit confidence
		}
		extractedItems[i] = ExtractedFoodItem{
			Name:       item.Name,
			Quantity:   item.Quantity,
			Unit:       item.Unit,
			Confidence: confidence,
		}
	}

	// Process each food item
	parsedItems := make([]domain.ParsedFoodItem, 0, len(extractedItems))
	totalConfidence := 0.0

	for _, item := range extractedItems {
		parsedItem, err := s.processFoodItem(ctx, userID, item)
		if err != nil {
			// Log error but continue processing other items
			fmt.Printf("Warning: failed to process food item %s: %v\n", item.Name, err)
			continue
		}
		parsedItems = append(parsedItems, parsedItem)
		totalConfidence += item.Confidence
	}

	if len(parsedItems) == 0 {
		return nil, fmt.Errorf("no valid food items could be extracted from image")
	}

	avgConfidence := totalConfidence / float64(len(parsedItems))

	// Infer meal type based on time of day
	mealType := s.inferMealType(time.Now())

	return &domain.ParsedMeal{
		MealType:          mealType,
		LoggedAt:          time.Now(),
		FoodItems:         parsedItems,
		Confidence:        avgConfidence,
		NeedsConfirmation: avgConfidence < 0.7, // Photos typically need more confirmation
	}, nil
}

// processFoodItem processes a single extracted food item
func (s *MealParserService) processFoodItem(ctx context.Context, userID uuid.UUID, item ExtractedFoodItem) (domain.ParsedFoodItem, error) {
	// Try to match food in database
	food, err := s.matchFoodInDatabase(ctx, item.Name)
	if err == nil && food != nil {
		// Found matching food in database
		return domain.ParsedFoodItem{
			FoodID:      &food.ID,
			FoodName:    food.Name,
			Quantity:    item.Quantity,
			Unit:        item.Unit,
			Confidence:  item.Confidence,
			AIGenerated: false,
		}, nil
	}

	// No match found - create AI-generated food
	aiFood, err := s.createAIFood(ctx, userID, item.Name)
	if err != nil {
		return domain.ParsedFoodItem{}, fmt.Errorf("failed to create AI food: %w", err)
	}

	return domain.ParsedFoodItem{
		FoodID:      &aiFood.ID,
		FoodName:    aiFood.Name,
		Quantity:    item.Quantity,
		Unit:        item.Unit,
		Confidence:  item.Confidence * 0.8, // Reduce confidence for AI-generated foods
		AIGenerated: true,
	}, nil
}

// matchFoodInDatabase attempts to find a matching food in the database
func (s *MealParserService) matchFoodInDatabase(ctx context.Context, name string) (*domain.Food, error) {
	// Search for food using full-text search
	foods, err := s.foodRepository.SearchFoods(ctx, name, 5)
	if err != nil {
		return nil, err
	}

	if len(foods) == 0 {
		return nil, fmt.Errorf("no matching foods found")
	}

	// Return best match (first result from search)
	// In a production system, we might want to implement fuzzy matching scoring
	return &foods[0], nil
}

// createAIFood creates a new AI-generated food with estimated nutrition
func (s *MealParserService) createAIFood(ctx context.Context, userID uuid.UUID, foodName string) (*domain.Food, error) {
	// Use AI to estimate nutrition per 100g
	systemPrompt := `You are a nutrition expert. Estimate the nutrition information per 100g for the given food.
Return a JSON object with:
{
  "calories": numeric_value,
  "protein": numeric_value_in_grams,
  "carbs": numeric_value_in_grams,
  "fat": numeric_value_in_grams,
  "fiber": numeric_value_in_grams
}

Provide realistic estimates based on typical nutrition values for this type of food.`

	messages := []external.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: fmt.Sprintf("Food: %s", foodName)},
	}

	resp, err := s.openRouterClient.Chat(ctx, messages, "deepseek/deepseek-chat")
	if err != nil {
		return nil, fmt.Errorf("failed to estimate nutrition: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	response := resp.Choices[0].Message.Content

	// Parse nutrition estimate
	var nutrition struct {
		Calories float64 `json:"calories"`
		Protein  float64 `json:"protein"`
		Carbs    float64 `json:"carbs"`
		Fat      float64 `json:"fat"`
		Fiber    float64 `json:"fiber"`
	}

	// Clean response
	cleanedResponse := strings.TrimSpace(response)
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```json")
	cleanedResponse = strings.TrimPrefix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSuffix(cleanedResponse, "```")
	cleanedResponse = strings.TrimSpace(cleanedResponse)

	if err := json.Unmarshal([]byte(cleanedResponse), &nutrition); err != nil {
		return nil, fmt.Errorf("failed to parse nutrition estimate: %w", err)
	}

	// Create food entity matching the actual Food domain model
	source := "ai_generated"
	fiber := nutrition.Fiber

	food := &domain.Food{
		ID:            uuid.New(),
		Name:          foodName,
		Brand:         nil, // AI-generated foods have no brand
		ServingSize:   100.0, // Base serving is 100g
		ServingUnit:   "g",
		Calories:      nutrition.Calories,
		Protein:       nutrition.Protein,
		Carbohydrates: nutrition.Carbs,
		Fat:           nutrition.Fat,
		Fiber:         &fiber,
		IsVerified:    false,
		Source:        &source,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.foodRepository.Create(ctx, food); err != nil {
		return nil, fmt.Errorf("failed to save AI-generated food: %w", err)
	}

	return food, nil
}

// inferMealType infers meal type based on time of day
func (s *MealParserService) inferMealType(t time.Time) string {
	hour := t.Hour()

	switch {
	case hour >= 5 && hour < 11:
		return "breakfast"
	case hour >= 11 && hour < 15:
		return "lunch"
	case hour >= 15 && hour < 18:
		return "snack"
	case hour >= 18 || hour < 5:
		return "dinner"
	default:
		return "snack"
	}
}
