# AI Agent Service Implementation

## Overview

The Agent Service provides an AI-powered fitness and nutrition coach that can interact with users through natural language and perform actions using function tools.

## Files Created

### 1. `/internal/services/agent_service.go`
Main service implementation with:
- AgentService struct with dependencies on all other services
- SendMessage method for processing user messages
- 8 tool functions for interacting with user data
- OpenRouter API integration for LLM calls

### 2. Updated `/internal/core/ports/services.go`
Added:
- AgentService interface
- AgentResponse struct

## Key Features

### 1. Conversational AI
- Maintains conversation history (last 20 messages)
- Builds user context from profile, goals, and recent activity
- Uses OpenRouter API for LLM responses

### 2. Tool Support (8 Tools)

#### Meal & Nutrition Tools
1. **log_meal** - Log meals with food items
2. **get_recent_meals** - Retrieve meal history (last N days)
3. **search_foods** - Search food database
4. **calculate_daily_macros** - Get nutrition totals for a specific date

#### Activity & Workout Tools
5. **get_recent_workouts** - Retrieve workout history
6. **get_recent_activities** - Retrieve activity logs

#### Metrics Tools
7. **log_weight** - Log weight measurements
8. **get_weight_trend** - Get weight trend over time

### 3. Context-Aware Responses

The agent builds user context including:
- User profile (name, demographics)
- Active goals (weight loss, muscle gain, etc.)
- Today's nutrition summary (calories, macros)
- Recent activity summary (last 7 days)

### 4. System Prompt

```
You are a fitness and nutrition coach assistant with access to the user's tracking data.

User Context:
{dynamic user data}

Guidelines:
- Be concise but thorough
- Reference user's actual data when relevant
- Provide evidence-based advice
- ALWAYS use tools to get accurate data before answering
- Never hallucinate meal or workout history
```

## Architecture

### Dependencies

The AgentService depends on:
- **Services**: MealService, FoodService, ActivityService, WorkoutService, MetricService, GoalService, SummaryService
- **Repositories**: ConversationRepository, UserRepository
- **External**: OpenRouterClient

### Response Flow

1. User sends message
2. Get or create conversation
3. Load last 20 messages for context
4. Build user context (profile, goals, today's data)
5. Build system prompt with context
6. Call OpenRouter API with tools
7. Execute any tool calls
8. Inject tool results back into conversation
9. Return final response
10. Save both user and assistant messages

### Tool Execution

Tools are executed through the `executeTool` method which:
1. Parses tool arguments (JSON)
2. Routes to appropriate tool function
3. Calls underlying service methods
4. Returns formatted results

## Usage Example

```go
agentService := services.NewAgentService(
    mealService,
    foodService,
    activityService,
    workoutService,
    metricService,
    goalService,
    summaryService,
    conversationRepo,
    userRepo,
    openRouterClient,
)

response, err := agentService.SendMessage(ctx, userID, "What did I eat today?")
if err != nil {
    // Handle error
}

fmt.Println(response.Message)
fmt.Println("Tools used:", response.ToolsUsed)
```

## Tool Examples

### Example 1: Search Foods
```
User: "Find nutrition info for chicken breast"
Tool: search_foods(query="chicken breast")
Result: Returns top 10 matching foods with nutrition data
```

### Example 2: Log Weight
```
User: "I weigh 75kg today"
Tool: log_weight(weight=75, date="2025-11-19")
Result: "Logged weight: 75.0 kg on 2025-11-19"
```

### Example 3: Get Weight Trend
```
User: "Show me my weight progress this month"
Tool: get_weight_trend(days=30)
Result: Weight measurements with calculated trend
```

## Configuration

- **Default Model**: deepseek/deepseek-chat (via OpenRouter)
- **Max Tool Iterations**: 5 (prevents infinite loops)
- **Context Window**: Last 20 messages
- **Confidence Score**: 0.85 (placeholder, can be enhanced)

## Future Enhancements

1. **Tool Implementations**
   - Complete log_meal implementation
   - Add date range filtering for get_recent_meals

2. **Context Enhancement**
   - Calculate actual goal values from user profile
   - Add more personalized insights

3. **Advanced Features**
   - Stream responses for real-time feedback
   - Multi-turn tool calling optimization
   - Confidence scoring based on tool usage
   - Conversation summarization

4. **Performance**
   - Cache frequently accessed user context
   - Batch tool calls where possible
   - Optimize message history loading

## Available Tools Summary

| Tool | Description | Parameters | Returns |
|------|-------------|------------|---------|
| log_meal | Log a meal with food items | food_items, meal_type, timestamp | Confirmation |
| get_recent_meals | Get recent meal history | days (default: 7) | Formatted meal list |
| search_foods | Search food database | query | Top 10 matching foods |
| calculate_daily_macros | Get daily nutrition totals | date | Macros breakdown |
| get_recent_workouts | Get workout history | days (default: 7) | Workout list |
| get_recent_activities | Get activity logs | days (default: 7) | Activity list |
| log_weight | Log weight measurement | weight, date (optional) | Confirmation |
| get_weight_trend | Get weight trend | days (default: 30) | Weight measurements + trend |

## Integration Points

### 1. HTTP Handler (TODO)
Create endpoint in `/internal/adapters/http/` for chat API

### 2. WebSocket Support (Future)
Real-time streaming responses

### 3. Mobile App Integration
RESTful API for mobile clients

## Testing Recommendations

1. **Unit Tests**
   - Test each tool function independently
   - Mock service dependencies
   - Test error handling

2. **Integration Tests**
   - Test full conversation flow
   - Test tool execution with real services
   - Test OpenRouter API integration

3. **E2E Tests**
   - Test multi-turn conversations
   - Test tool chaining
   - Test context persistence

## Security Considerations

1. **Input Validation**
   - Validate all tool parameters
   - Sanitize user messages

2. **Rate Limiting**
   - Limit messages per user per minute
   - Limit tool calls per conversation

3. **Data Access**
   - Ensure users can only access their own data
   - Validate userID in all tool calls

## Monitoring & Logging

Current logging includes:
- Message processing start
- Tool execution (name + args)
- Tool execution results/errors
- Context building warnings
- Message save failures

Recommended additions:
- Response time metrics
- Tool usage analytics
- Error rate tracking
- User engagement metrics
