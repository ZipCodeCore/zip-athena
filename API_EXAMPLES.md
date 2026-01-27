# API Usage Examples

This document provides examples of how to interact with the Athena API.

## Starting the Server

```bash
go run main.go
# Server starts on http://localhost:8080
```

## API Examples

### 1. Create a User

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "is_admin": false,
    "total_points": 0
  }'
```

Response:
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "is_admin": false,
  "total_points": 0,
  "created_at": "2024-01-27T15:30:00Z"
}
```

### 2. Get a User

```bash
curl http://localhost:8080/users?id=1
```

Response:
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "is_admin": false,
  "total_points": 0,
  "created_at": "2024-01-27T15:30:00Z"
}
```

### 3. Create a Topic

```bash
curl -X POST http://localhost:8080/topics \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Go Programming",
    "description": "Learn the fundamentals of Go programming language"
  }'
```

Response:
```json
{
  "id": 1,
  "name": "Go Programming",
  "description": "Learn the fundamentals of Go programming language",
  "created_at": "2024-01-27T15:30:00Z"
}
```

### 4. List All Topics

```bash
curl http://localhost:8080/topics
```

Response:
```json
[
  {
    "id": 1,
    "name": "Go Programming",
    "description": "Learn the fundamentals of Go programming language",
    "created_at": "2024-01-27T15:30:00Z"
  },
  {
    "id": 2,
    "name": "JavaScript",
    "description": "Master JavaScript basics",
    "created_at": "2024-01-27T15:31:00Z"
  }
]
```

### 5. Create a Question

```bash
curl -X POST http://localhost:8080/questions \
  -H "Content-Type: application/json" \
  -d '{
    "topic_id": 1,
    "question_text": "What keyword is used to declare a variable in Go?",
    "difficulty": 1,
    "tags": ["variables", "syntax", "basics"],
    "topic_question_id": "1.1"
  }'
```

Response:
```json
{
  "id": 1,
  "topic_id": 1,
  "question_text": "What keyword is used to declare a variable in Go?",
  "difficulty": 1,
  "tags": ["variables", "syntax", "basics"],
  "topic_question_id": "1.1",
  "created_at": "2024-01-27T15:32:00Z"
}
```

### 6. Get Questions by Topic

```bash
curl http://localhost:8080/questions?topic_id=1
```

Response:
```json
[
  {
    "id": 1,
    "topic_id": 1,
    "question_text": "What keyword is used to declare a variable in Go?",
    "difficulty": 1,
    "tags": ["variables", "syntax", "basics"],
    "topic_question_id": "1.1",
    "created_at": "2024-01-27T15:32:00Z"
  },
  {
    "id": 2,
    "topic_id": 1,
    "question_text": "What is a goroutine?",
    "difficulty": 3,
    "tags": ["concurrency", "advanced"],
    "topic_question_id": "1.2",
    "created_at": "2024-01-27T15:33:00Z"
  }
]
```

### 7. Create a Quiz Attempt

```bash
curl -X POST http://localhost:8080/quiz-attempts \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": 1,
    "total_questions": 5,
    "correct_answers": 3,
    "question_ids": ["1.1", "1.2", "2.1", "2.2", "3.1"],
    "question_results": {
      "1.1": true,
      "1.2": false,
      "2.1": true,
      "2.2": true,
      "3.1": false
    }
  }'
```

Response:
```json
{
  "id": 1,
  "user_id": 1,
  "start_time": "2024-01-27T15:34:00Z",
  "total_questions": 5,
  "correct_answers": 3,
  "question_ids": ["1.1", "1.2", "2.1", "2.2", "3.1"],
  "question_results": {
    "1.1": true,
    "1.2": false,
    "2.1": true,
    "2.2": true,
    "3.1": false
  }
}
```

### 8. Get Quiz Attempts by User

```bash
curl http://localhost:8080/quiz-attempts?user_id=1
```

Response:
```json
[
  {
    "id": 1,
    "user_id": 1,
    "start_time": "2024-01-27T15:34:00Z",
    "total_questions": 5,
    "correct_answers": 3,
    "question_ids": ["1.1", "1.2", "2.1", "2.2", "3.1"],
    "question_results": {
      "1.1": true,
      "1.2": false,
      "2.1": true,
      "2.2": true,
      "3.1": false
    }
  }
]
```

## Testing the API

Use the included test suite to verify all endpoints:

```bash
go test -v ./api
```

## Error Responses

### 404 Not Found
```json
User not found
```

### 400 Bad Request
```json
Invalid user ID
```

### 405 Method Not Allowed
```json
Method not allowed
```

### 500 Internal Server Error
```json
Internal server error message
```

## Next Steps

1. Add authentication (JWT tokens)
2. Implement rate limiting
3. Add request validation
4. Implement CORS for frontend
5. Add database persistence (PostgreSQL/MySQL)
6. Implement spaced repetition algorithm
7. Add more endpoints for answers, question attempts, etc.
