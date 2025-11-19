# AI Agent Documentation

Comprehensive documentation for AI-powered features in the Fitness Coach Backend.

## Table of Contents

- [Overview](#overview)
- [AI Capabilities](#ai-capabilities)
- [Meal Parsing AI](#meal-parsing-ai)
- [Chat Assistant](#chat-assistant)
- [AI-Generated Foods](#ai-generated-foods)
- [Tool Execution Flow](#tool-execution-flow)
- [AI Models Used](#ai-models-used)
- [Prompt Engineering](#prompt-engineering)
- [Error Handling](#error-handling)
- [Performance Considerations](#performance-considerations)
- [Adding New Tools](#adding-new-tools)
- [Future Enhancements](#future-enhancements)

## Overview

The Fitness Coach Backend leverages multiple AI models to provide intelligent features:

1. **Text-Based Meal Parsing**: Convert natural language food descriptions to structured meal data
2. **Photo-Based Meal Parsing**: Identify foods and estimate portions from photos
3. **AI-Generated Foods**: Automatic nutrition estimation for unknown foods
4. **Chat Assistant**: Provide fitness and nutrition guidance

## AI Capabilities

### 1. Natural Language Understanding

**Capability**: Parse complex food descriptions

**Examples:**
```
Input: "2 scrambled eggs with cheese, 1 slice of whole wheat toast with butter, and black coffee"
Output: Structured meal with individual food items, quantities, and units
```

**Supported Formats:**
- Simple lists: "apple, banana, orange"
- Quantities: "2 eggs, 100g chicken, 1 cup rice"
- Complex descriptions: "grilled chicken breast with roasted vegetables"
- Meal context: "for breakfast I had..."

### 2. Vision Recognition

**Capability**: Identify foods from photos

**What it can recognize:**
- Common foods and dishes
- Multiple items on a plate
- Estimated portion sizes
- Food categories

**Limitations:**
- Accuracy varies with photo quality
- May struggle with complex/mixed dishes
- Portion estimation is approximate

### 3. Nutrition Estimation

**Capability**: Generate nutritional information for unknown foods

**Estimated Values:**
- Calories (per 100g)
- Protein (g)
- Carbohydrates (g)
- Fat (g)
- Fiber (g)

**Accuracy:**
- Based on similar foods
- Uses USDA database patterns
- Marked as "AI-generated" (unverified)

### 4. Conversational AI

**Capability**: Answer fitness and nutrition questions

**Supported Topics:**
- Meal planning and nutrition
- Workout advice
- Goal setting
- Progress tracking
- Food substitutions
- Macro calculations

## Meal Parsing AI

See [MEAL_PARSING.md](./MEAL_PARSING.md) for detailed meal parsing documentation.

### Text Parsing

**Model**: DeepSeek via OpenRouter API

**Process:**
```
User Text Input
    ↓
Structured Prompt → DeepSeek
    ↓
JSON Response (Food Items)
    ↓
Database Matching
    ↓
AI Food Generation (if needed)
    ↓
ParsedMeal Result
```

### Photo Parsing

**Model**: Google Gemini 2.0 Flash Exp (via OpenRouter)

**Process:**
```
Photo Upload → Storage
    ↓
Photo URL → Vision AI
    ↓
Food Identification
    ↓
Portion Estimation
    ↓
Database Matching
    ↓
ParsedMeal Result
```

## Chat Assistant

### Capabilities

The AI chat assistant provides:

**Nutrition Guidance:**
- Macro calculations
- Meal planning advice
- Food substitutions
- Dietary recommendations

**Workout Advice:**
- Exercise selection
- Form tips
- Program design
- Recovery strategies

**Goal Support:**
- Progress tracking insights
- Motivation and tips
- Realistic goal setting
- Plateau breaking strategies

### Available Tools

The chat assistant has access to:

1. **User Profile Data**
   - Current weight, height
   - Activity level
   - Dietary preferences
   - Goals

2. **Historical Data**
   - Recent meals
   - Workout history
   - Progress metrics
   - Compliance patterns

3. **Knowledge Base**
   - Nutrition database
   - Exercise library
   - Scientific literature (embedded knowledge)

### Example Conversations

**Macro Calculation:**
```
User: "How many calories should I eat to lose weight?"