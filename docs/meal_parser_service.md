# Meal Parser Service Documentation

## Overview

The Meal Parser Service provides intelligent meal parsing capabilities for the Fitness Coach application. It uses AI models (DeepSeek for text, Gemini Vision for photos) to extract food items, portions, and nutrition information from both text input and food photos.

## Location

`internal/services/meal_parser_service.go`

## Dependencies

- **OpenRouter Client**: Text-based AI parsing using DeepSeek model
- **Vision Client**: Image analysis using Google Gemini 2.0 Flash
- **Food Repository**: Database access for food matching and storage

## Core Functionality

### 1. ParseText(ctx, userID, text) -> ParsedMeal

Parses meal information from natural language text input.

**Process:**
1. Sends text to DeepSeek via OpenRouter with structured prompt
2. AI extracts: meal type, food items, quantities, units, confidence scores
3. Matches each food item against database using full-text search
4. Creates AI-generated foods for unmatched items
5. Returns ParsedMeal with all items and overall confidence

**Example Input:**
```
"I had 2 scrambled eggs, 1 slice of whole wheat toast, and a cup of coffee for breakfast"
```

**Example Output:**
```json
{
  "meal_type": "breakfast",
  "logged_at": "2025-11-19T05:00:00Z",
  "food_items": [
    {
      "food_id": "uuid-123",
      "food_name": "Scrambled Eggs",
      "quantity": 2,
      "unit": "piece",
      "confidence": 0.95,
      "ai_generated": false
    },
    {
      "food_id": "uuid-456",
      "food_name": "Whole Wheat Toast",
      "quantity": 1,
      "unit": "slice",
      "confidence": 0.90,
      "ai_generated": false
    }
  ],
  "confidence": 0.925,
  "needs_confirmation": false
}
```

### 2. ParsePhoto(ctx, userID, photoURL) -> ParsedMeal

Parses meal information from food photo URL.

**Process:**
1. Sends image to Gemini Vision via OpenRouter
2. Vision AI identifies food items with estimated portions
3. Matches each identified food against database
4. Creates AI-generated foods for unmatched items
5. Infers meal type from time of day
6. Returns ParsedMeal with visual confidence scores

**Confidence Thresholds:**
- Text parsing: Requires confirmation if < 0.8
- Photo parsing: Requires confirmation if < 0.7 (more conservative)

### 3. matchFoodInDatabase(ctx, name) -> Food

Searches for matching food in database.

**Features:**
- Full-text search on food names
- Returns top 5 matches, selects best match
- Supports fuzzy matching through PostgreSQL full-text search

**Future Enhancement:**
- Implement fuzzy matching scoring algorithm
- Consider brand matching
- Weight results by source (USDA > verified > user > AI)

### 4. createAIFood(ctx, userID, foodName) -> Food

Creates AI-generated food with estimated nutrition.

**Process:**
1. Prompts DeepSeek for nutrition estimates per 100g
2. Parses JSON response: calories, protein, carbs, fat, fiber
3. Creates Food entity with:
   - `source: "ai_generated"`
   - `is_verified: false`
   - `serving_size: 100.0g`
   - All nutrition fields populated
4. Saves to database
5. Returns created Food entity

**AI-Generated Food Properties:**
- Private visibility (not shared with other users)
- No brand association
- Unverified by default
- 100g base serving size

### 5. inferMealType(t time.Time) -> string

Infers meal type from time of day.

**Time Ranges:**
- **Breakfast**: 5:00 AM - 11:00 AM
- **Lunch**: 11:00 AM - 3:00 PM
- **Snack**: 3:00 PM - 6:00 PM
- **Dinner**: 6:00 PM - 5:00 AM

## Data Models

### ParsedMeal
```go
type ParsedMeal struct {
    MealType          string           // breakfast|lunch|dinner|snack
    LoggedAt          time.Time        // When meal was parsed
    FoodItems         []ParsedFoodItem // Extracted food items
    Confidence        float64          // Average confidence (0-1)
    NeedsConfirmation bool             // True if confidence too low
}
```

### ParsedFoodItem
```go
type ParsedFoodItem struct {
    FoodID      *uuid.UUID // nil for AI-generated foods
    FoodName    string     // Display name
    Quantity    float64    // Amount consumed
    Unit        string     // g|ml|cup|piece|tbsp|tsp
    Confidence  float64    // AI confidence (0-1)
    AIGenerated bool       // True if food was AI-created
}
```

### ExtractedFoodItem
```go
type ExtractedFoodItem struct {
    Name       string  // Food name from AI
    Quantity   float64 // Parsed quantity
    Unit       string  // Parsed unit
    Confidence float64 // AI confidence score
}
```

## AI Integration

### DeepSeek (Text Parsing)
- **Model**: `deepseek/deepseek-chat`
- **Purpose**: Text extraction, nutrition estimation
- **Response Format**: JSON
- **Retry Logic**: 3 attempts with exponential backoff

### Gemini Vision (Photo Analysis)
- **Model**: `google/gemini-2.0-flash-exp:free`
- **Purpose**: Food identification from images
- **Response Format**: JSON array of food items
- **Fallback**: Manual parsing if JSON fails

## Error Handling

### Graceful Degradation
- Individual food item failures don't fail entire parse
- Warnings logged for failed items
- Continues processing remaining items

### Validation
- Requires at least 1 valid food item
- JSON parsing with markdown cleanup
- Handles both structured and unstructured AI responses

## Usage Example

```go
// Initialize service
foodRepo := postgres.NewFoodRepository(db)
parser := services.NewMealParserService(apiKey, foodRepo)

// Parse text input
parsedMeal, err := parser.ParseText(ctx, userID, "2 eggs and toast")
if err != nil {
    return err
}

if parsedMeal.NeedsConfirmation {
    // Show confirmation UI to user
    return askUserToConfirm(parsedMeal)
}

// Save meal to database
saveMeal(parsedMeal)

// Parse photo input
photoMeal, err := parser.ParsePhoto(ctx, userID, "https://storage.url/food.jpg")
if err != nil {
    return err
}
```

## Integration Points

### Food Repository Interface
```go
type FoodRepository interface {
    Create(ctx, food) error
    GetByID(ctx, id) (*Food, error)
    SearchFoods(ctx, query, limit) ([]Food, error)
    ListFoods(ctx, offset, limit, filter) ([]Food, error)
    Update(ctx, food) error
    Delete(ctx, id) error
}
```

### External Adapters
- `external.OpenRouterClient`: AI text processing
- `external.VisionClient`: Image analysis

## Environment Variables

```env
OPENROUTER_API_KEY=your_api_key_here
```

## Future Enhancements

1. **Improved Matching**
   - Fuzzy matching algorithm with scoring
   - Brand-aware matching
   - Category-based filtering

2. **Batch Processing**
   - Parse multiple meals concurrently
   - Bulk food creation optimization

3. **Confidence Calibration**
   - Track actual vs predicted accuracy
   - Adjust confidence thresholds dynamically
   - User feedback loop for AI training

4. **Multi-Language Support**
   - International food names
   - Unit conversions (metric/imperial)
   - Regional meal type variations

5. **Advanced Vision Features**
   - Portion size estimation improvements
   - Multi-food plate analysis
   - Nutritional label OCR

## Testing

### Unit Tests
```bash
go test ./internal/services -v -run TestMealParserService
```

### Integration Tests
```bash
go test ./tests/integration -v -run TestMealParsing
```

## Performance Considerations

- **API Latency**: 1-3 seconds per parse (AI processing)
- **Database Queries**: Full-text search optimized with GIN indexes
- **Concurrent Requests**: Safe for parallel execution
- **Token Usage**: ~500-1000 tokens per text parse, ~1500-3000 per photo

## Security

- User-specific AI-generated foods (private visibility)
- Input validation on all text/URLs
- Rate limiting recommended for API endpoints
- Sanitize file uploads before sending to Vision API
