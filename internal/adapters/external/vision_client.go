package external

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const (
	visionModel = "google/gemini-2.0-flash-exp:free"
)

// VisionClient handles food photo analysis using vision models
type VisionClient struct {
	openRouter *OpenRouterClient
}

// FoodItem represents a detected food item from the image
type FoodItem struct {
	Name        string  `json:"name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
	Calories    int     `json:"calories,omitempty"`
	Protein     float64 `json:"protein,omitempty"`
	Carbs       float64 `json:"carbs,omitempty"`
	Fat         float64 `json:"fat,omitempty"`
	Confidence  float64 `json:"confidence,omitempty"`
	Description string  `json:"description,omitempty"`
}

// FoodAnalysisResult contains all detected food items
type FoodAnalysisResult struct {
	Items       []FoodItem `json:"items"`
	TotalItems  int        `json:"total_items"`
	ImageQuality string    `json:"image_quality,omitempty"`
	Notes       string     `json:"notes,omitempty"`
}

// NewVisionClient creates a new vision client
func NewVisionClient(apiKey string) *VisionClient {
	return &VisionClient{
		openRouter: NewOpenRouterClient(apiKey),
	}
}

// AnalyzeFoodPhoto analyzes a food photo and returns structured food data
func (c *VisionClient) AnalyzeFoodPhoto(ctx context.Context, imageURL string) (*FoodAnalysisResult, error) {
	log.Printf("[Vision] Analyzing food photo: %s", imageURL)

	// Create vision prompt
	prompt := `Analyze this food image and identify all food items visible. For each item, provide:
1. Name of the food item
2. Estimated quantity/portion size
3. Unit of measurement (e.g., cup, piece, gram, oz)
4. Brief description

Format your response as a JSON array of food items like this:
[
  {
    "name": "Grilled Chicken Breast",
    "quantity": 6,
    "unit": "oz",
    "description": "Grilled boneless chicken breast"
  },
  {
    "name": "Steamed Broccoli",
    "quantity": 1,
    "unit": "cup",
    "description": "Fresh steamed broccoli florets"
  }
]

Be specific about the food items and realistic about portion sizes. Only include items you can clearly identify.`

	messages := []Message{
		{
			Role: "user",
			Content: fmt.Sprintf(`[{"type": "image_url", "image_url": {"url": "%s"}}, {"type": "text", "text": "%s"}]`,
				imageURL, prompt),
		},
	}

	resp, err := c.openRouter.Chat(ctx, messages, visionModel)
	if err != nil {
		return nil, fmt.Errorf("vision API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from vision model")
	}

	content := resp.Choices[0].Message.Content
	log.Printf("[Vision] Raw response: %s", content)

	// Parse the response
	result, err := c.parseVisionResponse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse vision response: %w", err)
	}

	log.Printf("[Vision] Detected %d food items", len(result.Items))
	return result, nil
}

// parseVisionResponse parses the vision model response into structured data
func (c *VisionClient) parseVisionResponse(content string) (*FoodAnalysisResult, error) {
	// Try to extract JSON array from response
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	var items []FoodItem
	if err := json.Unmarshal([]byte(jsonStr), &items); err != nil {
		// If JSON parsing fails, try to parse manually
		items = c.parseManually(content)
		if len(items) == 0 {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
	}

	// Validate and clean items
	validItems := make([]FoodItem, 0, len(items))
	for _, item := range items {
		if item.Name != "" && item.Quantity > 0 {
			// Normalize units
			item.Unit = normalizeUnit(item.Unit)
			validItems = append(validItems, item)
		}
	}

	return &FoodAnalysisResult{
		Items:      validItems,
		TotalItems: len(validItems),
	}, nil
}

// parseManually attempts to parse the response manually if JSON parsing fails
func (c *VisionClient) parseManually(content string) []FoodItem {
	var items []FoodItem

	// Look for patterns like "- Item: quantity unit"
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Try to extract food item information
		item := parseLine(line)
		if item.Name != "" {
			items = append(items, item)
		}
	}

	return items
}

// parseLine attempts to parse a single line into a FoodItem
func parseLine(line string) FoodItem {
	// Remove common prefixes
	line = strings.TrimPrefix(line, "-")
	line = strings.TrimPrefix(line, "*")
	line = strings.TrimPrefix(line, "•")
	line = strings.TrimSpace(line)

	// Try to find quantity and unit
	re := regexp.MustCompile(`(\d+\.?\d*)\s*(cup|cups|oz|ounce|ounces|gram|grams|g|piece|pieces|tbsp|tsp|serving|servings)`)
	matches := re.FindStringSubmatch(strings.ToLower(line))

	item := FoodItem{}
	if len(matches) >= 3 {
		qty, _ := strconv.ParseFloat(matches[1], 64)
		item.Quantity = qty
		item.Unit = matches[2]
		// The name is likely before the quantity
		namePart := strings.Split(line, matches[0])[0]
		item.Name = strings.TrimSpace(namePart)
	} else {
		// Just use the whole line as name
		item.Name = line
		item.Quantity = 1
		item.Unit = "serving"
	}

	return item
}

// extractJSON finds and extracts JSON array from text
func extractJSON(text string) string {
	// Find JSON array
	start := strings.Index(text, "[")
	end := strings.LastIndex(text, "]")
	if start != -1 && end != -1 && end > start {
		return text[start : end+1]
	}
	return ""
}

// normalizeUnit normalizes unit names
func normalizeUnit(unit string) string {
	unit = strings.ToLower(strings.TrimSpace(unit))
	switch unit {
	case "cups", "cup":
		return "cup"
	case "oz", "ounce", "ounces":
		return "oz"
	case "gram", "grams", "g":
		return "g"
	case "piece", "pieces":
		return "piece"
	case "tbsp", "tablespoon", "tablespoons":
		return "tbsp"
	case "tsp", "teaspoon", "teaspoons":
		return "tsp"
	case "serving", "servings":
		return "serving"
	default:
		return unit
	}
}

// AnalyzeFoodPhotoWithNutrition analyzes food and enriches with nutrition data
func (c *VisionClient) AnalyzeFoodPhotoWithNutrition(ctx context.Context, imageURL string) (*FoodAnalysisResult, error) {
	result, err := c.AnalyzeFoodPhoto(ctx, imageURL)
	if err != nil {
		return nil, err
	}

	// Enhance with nutrition estimates
	for i := range result.Items {
		c.estimateNutrition(&result.Items[i])
	}

	return result, nil
}

// estimateNutrition adds basic nutrition estimates
func (c *VisionClient) estimateNutrition(item *FoodItem) {
	// This is a simplified estimation
	// In production, you'd query a nutrition database or use another API

	name := strings.ToLower(item.Name)

	// Very basic estimation based on common foods
	switch {
	case strings.Contains(name, "chicken"):
		item.Calories = int(float64(165) * (item.Quantity / 3.5)) // per 3.5oz
		item.Protein = float64(31) * (item.Quantity / 3.5)
		item.Carbs = 0
		item.Fat = float64(3.6) * (item.Quantity / 3.5)
	case strings.Contains(name, "broccoli"):
		item.Calories = int(float64(55) * item.Quantity) // per cup
		item.Protein = float64(3.7) * item.Quantity
		item.Carbs = float64(11) * item.Quantity
		item.Fat = float64(0.6) * item.Quantity
	case strings.Contains(name, "rice"):
		item.Calories = int(float64(205) * item.Quantity) // per cup
		item.Protein = float64(4.3) * item.Quantity
		item.Carbs = float64(45) * item.Quantity
		item.Fat = float64(0.4) * item.Quantity
	default:
		// Default moderate estimate
		item.Calories = int(float64(150) * item.Quantity)
		item.Protein = float64(10) * item.Quantity
		item.Carbs = float64(15) * item.Quantity
		item.Fat = float64(5) * item.Quantity
	}
}
