# Fitness Coach API Documentation

Complete REST API reference for the Fitness Coach Backend.

## Base URL

```
http://localhost:8080/api/v1
```

## Table of Contents

- [Authentication](#authentication)
- [Error Responses](#error-responses)
- [Authentication Endpoints](#authentication-endpoints)
- [Meal Endpoints](#meal-endpoints)
- [Food Endpoints](#food-endpoints)
- [Activity Endpoints](#activity-endpoints)
- [Workout Endpoints](#workout-endpoints)
- [Exercise Endpoints](#exercise-endpoints)
- [Metric Endpoints](#metric-endpoints)
- [Goal Endpoints](#goal-endpoints)
- [Chat Endpoints](#chat-endpoints)
- [Summary Endpoints](#summary-endpoints)

## Authentication

Most endpoints require JWT authentication. Include the access token in the Authorization header:

```
Authorization: Bearer <access_token>
```

### Token Lifecycle

- **Access Token**: Expires in 24 hours
- **Refresh Token**: Expires in 7 days (168 hours)
- Use `/auth/refresh` endpoint to get new access token

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {}
  }
}
```

### Common HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate email)
- `422 Unprocessable Entity` - Validation error
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error

## Authentication Endpoints

### Register User

Create a new user account.

**Endpoint**: `POST /auth/register`

**Authentication**: None required

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response**: `201 Created`
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2025-11-19T10:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400
}
```

**Errors**:
- `400` - Invalid request format
- `409` - Email already registered
- `422` - Validation errors (weak password, invalid email)

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

---

### Login

Authenticate and receive access tokens.

**Endpoint**: `POST /auth/login`

**Authentication**: None required

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response**: `200 OK`
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2025-11-19T10:00:00Z",
    "updated_at": "2025-11-19T10:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400
}
```

**Errors**:
- `400` - Invalid request format
- `401` - Invalid credentials

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePassword123!"
  }'
```

---

### Refresh Token

Get a new access token using refresh token.

**Endpoint**: `POST /auth/refresh`

**Authentication**: None required (uses refresh token)

**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response**: `200 OK`
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400
}
```

**Errors**:
- `400` - Invalid request format
- `401` - Invalid or expired refresh token

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }'
```

---

## Meal Endpoints

### Create Meal

Log a new meal with food items.

**Endpoint**: `POST /meals`

**Authentication**: Required

**Request Body**:
```json
{
  "name": "Breakfast",
  "meal_type": "breakfast",
  "consumed_at": "2025-11-19T08:00:00Z",
  "notes": "Healthy start to the day",
  "food_items": [
    {
      "food_id": "123e4567-e89b-12d3-a456-426614174000",
      "quantity": 2.0,
      "unit": "piece"
    },
    {
      "food_id": "123e4567-e89b-12d3-a456-426614174001",
      "quantity": 1.0,
      "unit": "slice"
    }
  ]
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Breakfast",
  "meal_type": "breakfast",
  "consumed_at": "2025-11-19T08:00:00Z",
  "notes": "Healthy start to the day",
  "total_calories": 350.5,
  "total_protein": 15.2,
  "total_carbohydrates": 42.0,
  "total_fat": 12.5,
  "food_items": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174003",
      "meal_id": "123e4567-e89b-12d3-a456-426614174002",
      "food_id": "123e4567-e89b-12d3-a456-426614174000",
      "quantity": 2.0,
      "unit": "piece",
      "calories": 140.0,
      "protein": 12.0,
      "carbohydrates": 2.0,
      "fat": 10.0,
      "food": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Egg (Large)",
        "brand": null,
        "serving_size": 50.0,
        "serving_unit": "g",
        "calories": 70.0,
        "protein": 6.0,
        "carbohydrates": 1.0,
        "fat": 5.0
      }
    }
  ],
  "created_at": "2025-11-19T08:05:00Z",
  "updated_at": "2025-11-19T08:05:00Z"
}
```

**Errors**:
- `400` - Invalid request format
- `401` - Unauthorized
- `404` - Food ID not found
- `422` - Validation errors

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/meals \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Breakfast",
    "meal_type": "breakfast",
    "consumed_at": "2025-11-19T08:00:00Z",
    "food_items": [
      {
        "food_id": "123e4567-e89b-12d3-a456-426614174000",
        "quantity": 2.0,
        "unit": "piece"
      }
    ]
  }'
```

---

### List Meals

Retrieve user's meals with pagination and filtering.

**Endpoint**: `GET /meals`

**Authentication**: Required

**Query Parameters**:
- `offset` (optional, default: 0) - Number of records to skip
- `limit` (optional, default: 20, max: 100) - Number of records to return
- `start_date` (optional) - Filter by start date (ISO 8601)
- `end_date` (optional) - Filter by end date (ISO 8601)
- `meal_type` (optional) - Filter by meal type (breakfast, lunch, dinner, snack)

**Response**: `200 OK`
```json
{
  "meals": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174002",
      "user_id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Breakfast",
      "meal_type": "breakfast",
      "consumed_at": "2025-11-19T08:00:00Z",
      "total_calories": 350.5,
      "total_protein": 15.2,
      "total_carbohydrates": 42.0,
      "total_fat": 12.5,
      "created_at": "2025-11-19T08:05:00Z"
    }
  ],
  "total": 25,
  "offset": 0,
  "limit": 20
}
```

**cURL Example**:
```bash
curl -X GET "http://localhost:8080/api/v1/meals?limit=20&meal_type=breakfast" \
  -H "Authorization: Bearer <access_token>"
```

---

### Get Meal by ID

Retrieve detailed information about a specific meal.

**Endpoint**: `GET /meals/:id`

**Authentication**: Required

**Path Parameters**:
- `id` - Meal UUID

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174002",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Breakfast",
  "meal_type": "breakfast",
  "consumed_at": "2025-11-19T08:00:00Z",
  "notes": "Healthy start to the day",
  "total_calories": 350.5,
  "total_protein": 15.2,
  "total_carbohydrates": 42.0,
  "total_fat": 12.5,
  "food_items": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174003",
      "meal_id": "123e4567-e89b-12d3-a456-426614174002",
      "food_id": "123e4567-e89b-12d3-a456-426614174000",
      "quantity": 2.0,
      "unit": "piece",
      "calories": 140.0,
      "protein": 12.0,
      "carbohydrates": 2.0,
      "fat": 10.0,
      "food": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Egg (Large)",
        "serving_size": 50.0,
        "serving_unit": "g",
        "calories": 70.0,
        "protein": 6.0,
        "carbohydrates": 1.0,
        "fat": 5.0
      }
    }
  ],
  "created_at": "2025-11-19T08:05:00Z",
  "updated_at": "2025-11-19T08:05:00Z"
}
```

**Errors**:
- `401` - Unauthorized
- `404` - Meal not found or doesn't belong to user

**cURL Example**:
```bash
curl -X GET http://localhost:8080/api/v1/meals/123e4567-e89b-12d3-a456-426614174002 \
  -H "Authorization: Bearer <access_token>"
```

---

### Update Meal

Update an existing meal.

**Endpoint**: `PUT /meals/:id`

**Authentication**: Required

**Path Parameters**:
- `id` - Meal UUID

**Request Body**:
```json
{
  "name": "Updated Breakfast",
  "notes": "Added coffee"
}
```

**Response**: `200 OK` (same as Get Meal response)

**Errors**:
- `400` - Invalid request format
- `401` - Unauthorized
- `404` - Meal not found
- `422` - Validation errors

---

### Delete Meal

Delete a meal (soft delete).

**Endpoint**: `DELETE /meals/:id`

**Authentication**: Required

**Path Parameters**:
- `id` - Meal UUID

**Response**: `204 No Content`

**Errors**:
- `401` - Unauthorized
- `404` - Meal not found

**cURL Example**:
```bash
curl -X DELETE http://localhost:8080/api/v1/meals/123e4567-e89b-12d3-a456-426614174002 \
  -H "Authorization: Bearer <access_token>"
```

---

## Food Endpoints

### Create Food

Add a new food item to the database.

**Endpoint**: `POST /foods`

**Authentication**: Required

**Request Body**:
```json
{
  "name": "Chicken Breast (Grilled)",
  "description": "Skinless, boneless chicken breast",
  "brand": "Generic",
  "category": "Protein",
  "serving_size": 100.0,
  "serving_unit": "g",
  "calories": 165.0,
  "protein": 31.0,
  "carbohydrates": 0.0,
  "fat": 3.6,
  "fiber": 0.0,
  "sugar": 0.0,
  "saturated_fat": 1.0,
  "sodium": 74.0
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174010",
  "name": "Chicken Breast (Grilled)",
  "description": "Skinless, boneless chicken breast",
  "brand": "Generic",
  "category": "Protein",
  "serving_size": 100.0,
  "serving_unit": "g",
  "calories": 165.0,
  "protein": 31.0,
  "carbohydrates": 0.0,
  "fat": 3.6,
  "fiber": 0.0,
  "sugar": 0.0,
  "saturated_fat": 1.0,
  "sodium": 74.0,
  "is_verified": false,
  "source": "user",
  "created_at": "2025-11-19T10:00:00Z",
  "updated_at": "2025-11-19T10:00:00Z"
}
```

**Errors**:
- `400` - Invalid request format
- `401` - Unauthorized
- `422` - Validation errors

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/foods \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chicken Breast (Grilled)",
    "serving_size": 100.0,
    "serving_unit": "g",
    "calories": 165.0,
    "protein": 31.0,
    "carbohydrates": 0.0,
    "fat": 3.6
  }'
```

---

### List Foods

Search and list food items.

**Endpoint**: `GET /foods`

**Authentication**: Required

**Query Parameters**:
- `offset` (optional, default: 0) - Pagination offset
- `limit` (optional, default: 20, max: 100) - Results per page
- `search` (optional) - Full-text search query
- `category` (optional) - Filter by category
- `verified_only` (optional, default: false) - Only verified foods

**Response**: `200 OK`
```json
{
  "foods": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174010",
      "name": "Chicken Breast (Grilled)",
      "brand": "Generic",
      "category": "Protein",
      "serving_size": 100.0,
      "serving_unit": "g",
      "calories": 165.0,
      "protein": 31.0,
      "carbohydrates": 0.0,
      "fat": 3.6,
      "is_verified": true,
      "source": "usda"
    }
  ],
  "total": 1542,
  "offset": 0,
  "limit": 20
}
```

**cURL Example**:
```bash
curl -X GET "http://localhost:8080/api/v1/foods?search=chicken&verified_only=true" \
  -H "Authorization: Bearer <access_token>"
```

---

### Get Food by ID

Retrieve detailed food information.

**Endpoint**: `GET /foods/:id`

**Authentication**: Required

**Path Parameters**:
- `id` - Food UUID

**Response**: `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174010",
  "fdc_id": 171477,
  "name": "Chicken Breast (Grilled)",
  "description": "Skinless, boneless chicken breast",
  "brand": "Generic",
  "category": "Protein",
  "serving_size": 100.0,
  "serving_unit": "g",
  "calories": 165.0,
  "protein": 31.0,
  "carbohydrates": 0.0,
  "fat": 3.6,
  "fiber": 0.0,
  "sugar": 0.0,
  "saturated_fat": 1.0,
  "trans_fat": 0.0,
  "cholesterol": 85.0,
  "sodium": 74.0,
  "potassium": 256.0,
  "is_verified": true,
  "source": "usda",
  "created_at": "2025-11-19T10:00:00Z",
  "updated_at": "2025-11-19T10:00:00Z"
}
```

**Errors**:
- `401` - Unauthorized
- `404` - Food not found

---

### Update Food

Update food information.

**Endpoint**: `PUT /foods/:id`

**Authentication**: Required

**Path Parameters**:
- `id` - Food UUID

**Request Body**: (Partial update supported)
```json
{
  "description": "Updated description",
  "calories": 170.0
}
```

**Response**: `200 OK` (same as Get Food response)

**Errors**:
- `400` - Invalid request
- `401` - Unauthorized
- `403` - Cannot modify verified foods
- `404` - Food not found

---

### Delete Food

Delete a food item (soft delete).

**Endpoint**: `DELETE /foods/:id`

**Authentication**: Required

**Path Parameters**:
- `id` - Food UUID

**Response**: `204 No Content`

**Errors**:
- `401` - Unauthorized
- `403` - Cannot delete verified foods
- `404` - Food not found

---

## Activity Endpoints

### Create Activity

Log a cardio or general activity.

**Endpoint**: `POST /activities`

**Authentication**: Required

**Request Body**:
```json
{
  "activity_type": "running",
  "name": "Morning Run",
  "duration_minutes": 30,
  "calories_burned": 300.0,
  "distance_km": 5.0,
  "average_heart_rate": 145,
  "performed_at": "2025-11-19T06:00:00Z",
  "notes": "Felt great!"
}
```

**Response**: `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174020",
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "activity_type": "running",
  "name": "Morning Run",
  "duration_minutes": 30,
  "calories_burned": 300.0,
  "distance_km": 5.0,
  "average_heart_rate": 145,
  "performed_at": "2025-11-19T06:00:00Z",
  "notes": "Felt great!",
  "created_at": "2025-11-19T06:30:00Z",
  "updated_at": "2025-11-19T06:30:00Z"
}
```

**cURL Example**:
```bash
curl -X POST http://localhost:8080/api/v1/activities \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "activity_type": "running",
    "name": "Morning Run",
    "duration_minutes": 30,
    "calories_burned": 300.0,
    "performed_at": "2025-11-19T06:00:00Z"
  }'
```

---

### List Activities

Retrieve user activities with filtering.

**Endpoint**: `GET /activities`

**Authentication**: Required

**Query Parameters**:
- `offset`, `limit` - Pagination
- `start_date`, `end_date` - Date range filter
- `activity_type` - Filter by type

**Response**: `200 OK`
```json
{
  "activities": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174020",
      "activity_type": "running",
      "name": "Morning Run",
      "duration_minutes": 30,
      "calories_burned": 300.0,
      "distance_km": 5.0,
      "performed_at": "2025-11-19T06:00:00Z"
    }
  ],
  "total": 15,
  "offset": 0,
  "limit": 20
}
```

---

### Get/Update/Delete Activity

Similar patterns as Meals endpoints:
- `GET /activities/:id` - Get activity details
- `PUT /activities/:id` - Update activity
- `DELETE /activities/:id` - Delete activity

---

## Workout Endpoints

Workouts contain multiple exercises with sets.

### Create Workout

**Endpoint**: `POST /workouts`

**Request Body**:
```json
{
  "name": "Upper Body Strength",
  "workout_type": "strength",
  "performed_at": "2025-11-19T17:00:00Z",
  "duration_minutes": 60,
  "notes": "Great pump!",
  "exercises": [
    {
      "exercise_id": "123e4567-e89b-12d3-a456-426614174030",
      "order": 1,
      "sets": [
        {
          "set_number": 1,
          "reps": 10,
          "weight_kg": 50.0,
          "rest_seconds": 90
        },
        {
          "set_number": 2,
          "reps": 10,
          "weight_kg": 50.0,
          "rest_seconds": 90
        }
      ]
    }
  ]
}
```

**Response**: `201 Created` (includes full workout with exercises and sets)

---

## Exercise Endpoints

Exercise library management.

### List Exercises

**Endpoint**: `GET /exercises`

**Query Parameters**:
- `search` - Search exercise name
- `muscle_group` - Filter by muscle group
- `equipment` - Filter by equipment type

**Response**: `200 OK`
```json
{
  "exercises": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174030",
      "name": "Bench Press",
      "description": "Chest exercise",
      "muscle_group": "chest",
      "equipment": "barbell",
      "difficulty": "intermediate",
      "instructions": "Lie on bench, lower bar to chest, press up"
    }
  ],
  "total": 150,
  "offset": 0,
  "limit": 20
}
```

---

## Metric Endpoints

Track body metrics over time.

### Create Metric

**Endpoint**: `POST /metrics`

**Request Body**:
```json
{
  "weight_kg": 75.5,
  "body_fat_percentage": 18.5,
  "muscle_mass_kg": 58.0,
  "waist_cm": 82.0,
  "measured_at": "2025-11-19T07:00:00Z",
  "notes": "Morning measurement"
}
```

**Response**: `201 Created`

---

## Goal Endpoints

Set and track fitness goals.

### Create Goal

**Endpoint**: `POST /goals`

**Request Body**:
```json
{
  "goal_type": "weight_loss",
  "title": "Lose 5kg",
  "description": "Reach 70kg by summer",
  "target_value": 70.0,
  "current_value": 75.0,
  "target_date": "2025-06-01",
  "status": "in_progress"
}
```

**Response**: `201 Created`

---

## Chat Endpoints

AI coaching assistant.

### Send Message

**Endpoint**: `POST /chat/message`

**Request Body**:
```json
{
  "message": "What should I eat for post-workout recovery?"
}
```

**Response**: `200 OK`
```json
{
  "message_id": "123e4567-e89b-12d3-a456-426614174040",
  "user_message": "What should I eat for post-workout recovery?",
  "assistant_response": "For optimal post-workout recovery, I recommend:\n\n1. Protein (20-40g): Chicken, fish, or protein shake\n2. Carbohydrates: Rice, sweet potato, or banana\n3. Hydration: Water with electrolytes\n\nTiming: Within 30-60 minutes post-workout for best results.",
  "timestamp": "2025-11-19T18:30:00Z"
}
```

---

### Get Chat History

**Endpoint**: `GET /chat/history`

**Query Parameters**:
- `limit` (default: 50) - Number of messages
- `before` (optional) - Get messages before timestamp

**Response**: `200 OK`
```json
{
  "messages": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174040",
      "role": "user",
      "content": "What should I eat for post-workout recovery?",
      "timestamp": "2025-11-19T18:30:00Z"
    },
    {
      "id": "123e4567-e89b-12d3-a456-426614174041",
      "role": "assistant",
      "content": "For optimal post-workout recovery...",
      "timestamp": "2025-11-19T18:30:05Z"
    }
  ],
  "total": 24
}
```

---

## Summary Endpoints

Aggregated daily statistics.

### Get Daily Summary

**Endpoint**: `GET /summary/daily`

**Query Parameters**:
- `date` (optional, default: today) - Date in ISO 8601 format

**Response**: `200 OK`
```json
{
  "date": "2025-11-19",
  "nutrition": {
    "total_calories": 2150.5,
    "total_protein": 165.2,
    "total_carbohydrates": 220.0,
    "total_fat": 65.5,
    "meals_count": 4,
    "target_calories": 2200.0,
    "calories_remaining": 49.5
  },
  "activity": {
    "total_duration_minutes": 90,
    "total_calories_burned": 650.0,
    "activities_count": 2,
    "workouts_count": 1
  },
  "metrics": {
    "weight_kg": 75.5,
    "body_fat_percentage": 18.5
  },
  "goals": {
    "daily_calorie_target": 2200.0,
    "protein_target": 165.0,
    "active_goals": 3,
    "goals_on_track": 2
  }
}
```

**cURL Example**:
```bash
curl -X GET "http://localhost:8080/api/v1/summary/daily?date=2025-11-19" \
  -H "Authorization: Bearer <access_token>"
```

---

## Rate Limiting

Default rate limits (configurable):
- **100 requests per minute** per user
- **429 Too Many Requests** response when exceeded
- Rate limit headers included in response:
  - `X-RateLimit-Limit`: Maximum requests
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Unix timestamp when limit resets

## Pagination

List endpoints support pagination:

**Query Parameters**:
- `offset` - Number of records to skip (default: 0)
- `limit` - Number of records to return (default: 20, max: 100)

**Response includes**:
- `total` - Total number of records
- `offset` - Current offset
- `limit` - Current limit

## Data Types

- **UUID**: `123e4567-e89b-12d3-a456-426614174000`
- **ISO 8601 DateTime**: `2025-11-19T10:00:00Z`
- **ISO 8601 Date**: `2025-11-19`
- **Decimal**: Numeric values with precision (e.g., `75.5`)

## Best Practices

1. Always include `Authorization` header for protected endpoints
2. Handle token expiration gracefully (use refresh token)
3. Implement exponential backoff for rate limit errors
4. Use pagination for list endpoints
5. Validate data client-side before sending
6. Store refresh token securely (never in localStorage)
7. Use HTTPS in production
8. Implement request timeouts (15s recommended)

---

For more details, see:
- [Main Documentation](./README.md)
- [Architecture Guide](./ARCHITECTURE.md)
- [Meal Parsing Details](./MEAL_PARSING.md)
