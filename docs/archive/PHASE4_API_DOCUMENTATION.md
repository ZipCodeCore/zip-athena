# Phase 4 API Documentation

This document describes all new API endpoints added in Phase 4: Advanced Student Features & Testing.

## Table of Contents

1. [Quiz Attempt History](#quiz-attempt-history)
2. [Timed Test Taking](#timed-test-taking)
3. [Tag-Based Review](#tag-based-review)
4. [Test Management (Instructor)](#test-management-instructor)
5. [Question Performance Analytics](#question-performance-analytics)

---

## Quiz Attempt History

### GET /students/quiz-history

Retrieve a student's quiz attempt history with pagination and optional topic filtering.

**Query Parameters:**
- `user_id` (required): Student's user ID
- `topic_id` (optional): Filter by topic ID
- `page` (optional): Page number (default: 1)
- `per_page` (optional): Results per page (default: 10, max: 100)

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 123,
    "quiz_id": 5,
    "start_time": "2026-01-29T10:00:00Z",
    "completed_time": "2026-01-29T10:15:00Z",
    "total_questions": 10,
    "correct_answers": 7,
    "question_ids": ["1.1", "1.2", "1.3"],
    "question_results": {
      "1.1": true,
      "1.2": false,
      "1.3": true
    },
    "attempt_number": 1,
    "is_retake": false
  }
]
```

**Example:**
```bash
curl "http://localhost:9206/students/quiz-history?user_id=1&page=1&per_page=10"
curl "http://localhost:9206/students/quiz-history?user_id=1&topic_id=2&page=1&per_page=10"
```

### GET /students/quiz-history/{id}

Get detailed information about a specific quiz attempt.

**Path Parameters:**
- `id`: Quiz attempt ID

**Response:**
```json
{
  "id": 1,
  "user_id": 123,
  "quiz_id": 5,
  "start_time": "2026-01-29T10:00:00Z",
  "completed_time": "2026-01-29T10:15:00Z",
  "total_questions": 10,
  "correct_answers": 7,
  "question_ids": ["1.1", "1.2", "1.3", "1.4", "1.5"],
  "question_results": {
    "1.1": true,
    "1.2": false,
    "1.3": true,
    "1.4": true,
    "1.5": false
  }
}
```

**Example:**
```bash
curl "http://localhost:9206/students/quiz-history/1"
```

---

## Timed Test Taking

### GET /tests/available

Retrieve all available tests for a student (published and within availability window).

**Query Parameters:**
- `user_id` (required): Student's user ID

**Response:**
```json
[
  {
    "id": 1,
    "name": "Java Midterm",
    "topic_id": 2,
    "description": "Comprehensive Java midterm examination",
    "time_limit": 60,
    "start_date": "2026-01-29T00:00:00Z",
    "end_date": "2026-02-05T23:59:59Z",
    "is_published": true,
    "created_at": "2026-01-20T10:00:00Z",
    "updated_at": "2026-01-25T15:00:00Z"
  }
]
```

**Example:**
```bash
curl "http://localhost:9206/tests/available?user_id=1"
```

### GET /tests/{id}

Get details about a specific test, including question count.

**Path Parameters:**
- `id`: Test ID

**Response:**
```json
{
  "test": {
    "id": 1,
    "name": "Java Midterm",
    "topic_id": 2,
    "description": "Comprehensive Java midterm examination",
    "time_limit": 60,
    "is_published": true
  },
  "question_count": 25
}
```

**Example:**
```bash
curl "http://localhost:9206/tests/1"
```

### POST /tests/{id}/start

Start a new test attempt for a student.

**Path Parameters:**
- `id`: Test ID

**Request Body:**
```json
{
  "user_id": 1
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "test_id": 1,
  "start_time": "2026-01-29T14:00:00Z",
  "answers": {},
  "question_results": {}
}
```

**Error Responses:**
- `409 Conflict`: Student already has an incomplete attempt for this test

**Example:**
```bash
curl -X POST "http://localhost:9206/tests/1/start" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'
```

### PUT /test-attempts/{id}/answer

Submit or update an answer during a test attempt.

**Path Parameters:**
- `id`: Test attempt ID

**Request Body:**
```json
{
  "question_id": "1",
  "answer_id": 42
}
```

**Response:**
```json
{
  "status": "success"
}
```

**Example:**
```bash
curl -X PUT "http://localhost:9206/test-attempts/1/answer" \
  -H "Content-Type: application/json" \
  -d '{"question_id": "1", "answer_id": 42}'
```

### POST /test-attempts/{id}/submit

Submit a test attempt (manual or auto-submit when timer expires).

**Path Parameters:**
- `id`: Test attempt ID

**Request Body:**
```json
{
  "is_auto_submit": false
}
```

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "test_id": 1,
  "start_time": "2026-01-29T14:00:00Z",
  "submit_time": "2026-01-29T15:00:00Z",
  "time_remaining": 0,
  "is_auto_submit": false,
  "total_questions": 25,
  "correct_answers": 18,
  "score": 72.0,
  "answers": {...},
  "question_results": {...}
}
```

**Example:**
```bash
curl -X POST "http://localhost:9206/test-attempts/1/submit" \
  -H "Content-Type: application/json" \
  -d '{"is_auto_submit": false}'
```

### GET /test-attempts/{id}/results

Get the results of a completed test attempt.

**Path Parameters:**
- `id`: Test attempt ID

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "test_id": 1,
  "start_time": "2026-01-29T14:00:00Z",
  "submit_time": "2026-01-29T15:00:00Z",
  "time_remaining": 0,
  "is_auto_submit": false,
  "total_questions": 25,
  "correct_answers": 18,
  "score": 72.0,
  "answers": {...},
  "question_results": {
    "1": true,
    "2": false,
    "3": true
  }
}
```

**Example:**
```bash
curl "http://localhost:9206/test-attempts/1/results"
```

---

## Tag-Based Review

### GET /tags

Get all unique tags with question counts.

**Response:**
```json
{
  "loops": 15,
  "arrays": 12,
  "functions": 8,
  "recursion": 5,
  "oop": 20
}
```

**Example:**
```bash
curl "http://localhost:9206/tags"
```

### GET /tags/stats

Get accuracy statistics for each tag for a specific user.

**Query Parameters:**
- `user_id` (required): User ID

**Response:**
```json
[
  {
    "tag": "loops",
    "question_count": 15,
    "total_attempts": 45,
    "correct_attempts": 35,
    "accuracy": 77.78
  },
  {
    "tag": "arrays",
    "question_count": 12,
    "total_attempts": 36,
    "correct_attempts": 20,
    "accuracy": 55.56
  }
]
```

**Example:**
```bash
curl "http://localhost:9206/tags/stats?user_id=1"
```

### POST /quiz/generate/by-tags

Generate a quiz filtered by specific tags.

**Request Body:**
```json
{
  "tags": ["loops", "arrays"],
  "use_and": false,
  "user_id": 1,
  "count": 10
}
```

**Parameters:**
- `tags`: Array of tags to filter by
- `use_and`: If `true`, questions must have ALL tags (AND logic). If `false`, questions must have AT LEAST ONE tag (OR logic)
- `user_id`: User ID for the quiz attempt
- `count`: Number of questions to include (default: 10)

**Response:**
```json
{
  "attempt_id": 42,
  "questions": [
    {
      "id": 1,
      "topic_id": 2,
      "question_text": "What is a for loop?",
      "difficulty": 3,
      "tags": ["loops"],
      "topic_question_id": "1.5"
    },
    {
      "id": 2,
      "topic_id": 2,
      "question_text": "How do you declare an array?",
      "difficulty": 2,
      "tags": ["arrays"],
      "topic_question_id": "1.8"
    }
  ]
}
```

**Example:**
```bash
# OR logic: questions with loops OR arrays
curl -X POST "http://localhost:9206/quiz/generate/by-tags" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["loops", "arrays"], "use_and": false, "user_id": 1, "count": 10}'

# AND logic: questions with BOTH loops AND arrays
curl -X POST "http://localhost:9206/quiz/generate/by-tags" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["loops", "arrays"], "use_and": true, "user_id": 1, "count": 5}'
```

---

## Test Management (Instructor)

**Note:** All instructor endpoints require admin authentication.

### POST /instructor/tests/{id}/publish

Publish a test to make it available to students.

**Path Parameters:**
- `id`: Test ID

**Response:**
```json
{
  "status": "success",
  "message": "Test published successfully"
}
```

**Example:**
```bash
curl -X POST "http://localhost:9206/instructor/tests/1/publish" \
  -H "Authorization: Bearer <admin_token>"
```

### GET /instructor/tests/{id}/results

Get all student results for a specific test.

**Path Parameters:**
- `id`: Test ID

**Response:**
```json
{
  "test_id": 1,
  "test_name": "Java Midterm",
  "student_results": [
    {
      "user_id": 1,
      "username": "student1",
      "score": 85.0,
      "correct_answers": 17,
      "total_questions": 20,
      "submit_time": "2026-01-29T15:30:00Z",
      "time_taken": 3300
    },
    {
      "user_id": 2,
      "username": "student2",
      "score": 72.0,
      "correct_answers": 14,
      "total_questions": 20,
      "submit_time": "2026-01-29T15:45:00Z",
      "time_taken": 3600
    }
  ]
}
```

**Example:**
```bash
curl "http://localhost:9206/instructor/tests/1/results" \
  -H "Authorization: Bearer <admin_token>"
```

### GET /instructor/tests/{id}/export

Export test results as a CSV file.

**Path Parameters:**
- `id`: Test ID

**Response:**
CSV file download with headers: `Username,Score,Correct Answers,Total Questions,Submit Time,Time Taken (seconds)`

**Example:**
```bash
curl "http://localhost:9206/instructor/tests/1/export" \
  -H "Authorization: Bearer <admin_token>" \
  -o test_1_results.csv
```

---

## Question Performance Analytics

**Note:** All analytics endpoints require admin authentication.

### GET /instructor/questions/analytics

Get analytics for all questions, identifying which are too easy, too hard, or appropriate.

**Response:**
```json
[
  {
    "question_id": 1,
    "topic_question_id": "1.1",
    "question_text": "What is a variable?",
    "topic_id": 2,
    "topic_name": "Java",
    "tags": ["basics", "variables"],
    "current_difficulty": 2,
    "total_attempts": 100,
    "correct_attempts": 95,
    "accuracy": 95.0,
    "recommended_difficulty": 3,
    "flag_for_review": true,
    "answer_distribution": {}
  },
  {
    "question_id": 2,
    "topic_question_id": "1.2",
    "question_text": "Explain polymorphism",
    "topic_id": 2,
    "topic_name": "Java",
    "tags": ["oop", "advanced"],
    "current_difficulty": 6,
    "total_attempts": 80,
    "correct_attempts": 20,
    "accuracy": 25.0,
    "recommended_difficulty": 5,
    "flag_for_review": true,
    "answer_distribution": {}
  }
]
```

**Flagging Logic:**
- `flag_for_review` is `true` if accuracy > 90% (too easy) or < 30% (too hard)
- Recommended difficulty adjusts based on accuracy

**Example:**
```bash
curl "http://localhost:9206/instructor/questions/analytics" \
  -H "Authorization: Bearer <admin_token>"
```

### GET /instructor/questions/{id}/analytics

Get detailed analytics for a single question.

**Path Parameters:**
- `id`: Question ID

**Response:**
```json
{
  "question_id": 1,
  "topic_question_id": "1.1",
  "question_text": "What is a variable?",
  "topic_id": 2,
  "topic_name": "Java",
  "tags": ["basics", "variables"],
  "current_difficulty": 2,
  "total_attempts": 100,
  "correct_attempts": 95,
  "accuracy": 95.0,
  "recommended_difficulty": 3,
  "flag_for_review": true,
  "answer_distribution": {
    "1": 95,
    "2": 3,
    "3": 1,
    "4": 1
  }
}
```

**Example:**
```bash
curl "http://localhost:9206/instructor/questions/1/analytics" \
  -H "Authorization: Bearer <admin_token>"
```

### GET /instructor/questions/{id}/answer-distribution

Get distribution of which answers were chosen for a question.

**Path Parameters:**
- `id`: Question ID

**Response:**
```json
{
  "1": 95,
  "2": 3,
  "3": 1,
  "4": 1
}
```

**Note:** Currently returns empty map as answer tracking in quiz attempts is not fully implemented. This would require tracking individual answer selections in quiz/test attempts.

**Example:**
```bash
curl "http://localhost:9206/instructor/questions/1/answer-distribution" \
  -H "Authorization: Bearer <admin_token>"
```

### PUT /instructor/questions/{id}/difficulty

Update a question's difficulty rating based on analytics.

**Path Parameters:**
- `id`: Question ID

**Request Body:**
```json
{
  "difficulty": 5
}
```

**Validation:**
- Difficulty must be between 1 and 7

**Response:**
```json
{
  "status": "success",
  "message": "Difficulty updated successfully"
}
```

**Example:**
```bash
curl -X PUT "http://localhost:9206/instructor/questions/1/difficulty" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"difficulty": 5}'
```

---

## Error Responses

All endpoints follow standard HTTP status codes:

- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., test already started)
- `500 Internal Server Error`: Server error

Error response format:
```json
{
  "error": "Error message describing what went wrong"
}
```

---

## Authentication

Instructor endpoints (`/instructor/*`) require admin authentication via the `AdminMiddleware`. This is implemented in the API layer and checks that the authenticated user has `is_admin` set to `true`.

For development/testing purposes, admin status can be set directly in the database.

---

## Database Schema Changes

### New Tables

#### test_attempts
```sql
CREATE TABLE test_attempts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    test_id INTEGER NOT NULL,
    start_time DATETIME NOT NULL,
    submit_time DATETIME,
    time_remaining INTEGER DEFAULT 0,
    is_auto_submit INTEGER DEFAULT 0,
    total_questions INTEGER DEFAULT 0,
    correct_answers INTEGER DEFAULT 0,
    score REAL DEFAULT 0.0,
    answers TEXT DEFAULT '{}',
    question_results TEXT DEFAULT '{}',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
);
```

### Modified Tables

#### tests
Added fields:
- `start_date DATETIME`: When test becomes available
- `end_date DATETIME`: When test is no longer available
- `is_published INTEGER DEFAULT 0`: Whether test is published to students

---

## Known Limitations

1. **Answer Distribution**: The `GetAnswerDistribution` endpoint currently returns empty data as individual answer selections are not tracked in quiz/test attempts. To implement this fully, the `QuizAttempt` and `TestAttempt` models would need to track which specific answer IDs were selected for each question.

2. **Percentile Calculation**: Test result percentiles are not yet implemented. The submission returns the score, but percentile ranking compared to other students is not calculated.

3. **Cohort Support**: Test assignment to specific cohorts is not yet implemented. All published tests are available to all students.

4. **Client-Side Timer**: The timer functionality is implemented server-side (time tracking), but requires client-side JavaScript implementation for countdown display and auto-submit trigger.

5. **Local Storage Backup**: The requirement for local storage backup of answers (to prevent data loss on page refresh) requires client-side implementation.

---

## Integration Examples

### Complete Test Taking Flow

```bash
# 1. Student views available tests
curl "http://localhost:9206/tests/available?user_id=1"

# 2. Student gets test details
curl "http://localhost:9206/tests/1"

# 3. Student starts test
curl -X POST "http://localhost:9206/tests/1/start" \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1}'

# 4. Student submits answers (one at a time)
curl -X PUT "http://localhost:9206/test-attempts/1/answer" \
  -H "Content-Type: application/json" \
  -d '{"question_id": "1", "answer_id": 5}'

# 5. Student submits test
curl -X POST "http://localhost:9206/test-attempts/1/submit" \
  -H "Content-Type: application/json" \
  -d '{"is_auto_submit": false}'

# 6. Student views results
curl "http://localhost:9206/test-attempts/1/results"
```

### Tag-Based Practice Flow

```bash
# 1. Student views all available tags
curl "http://localhost:9206/tags"

# 2. Student checks their accuracy for each tag
curl "http://localhost:9206/tags/stats?user_id=1"

# 3. Student generates quiz for weak tags
curl -X POST "http://localhost:9206/quiz/generate/by-tags" \
  -H "Content-Type: application/json" \
  -d '{"tags": ["loops", "arrays"], "use_and": false, "user_id": 1, "count": 10}'
```

### Instructor Analytics Flow

```bash
# 1. View all question analytics
curl "http://localhost:9206/instructor/questions/analytics" \
  -H "Authorization: Bearer <admin_token>"

# 2. Review flagged question
curl "http://localhost:9206/instructor/questions/5/analytics" \
  -H "Authorization: Bearer <admin_token>"

# 3. Adjust difficulty based on data
curl -X PUT "http://localhost:9206/instructor/questions/5/difficulty" \
  -H "Authorization: Bearer <admin_token>" \
  -H "Content-Type: application/json" \
  -d '{"difficulty": 4}'

# 4. View test results
curl "http://localhost:9206/instructor/tests/1/results" \
  -H "Authorization: Bearer <admin_token>"

# 5. Export results
curl "http://localhost:9206/instructor/tests/1/export" \
  -H "Authorization: Bearer <admin_token>" \
  -o results.csv
```

---

## Testing

All Phase 4 features are covered by comprehensive unit tests in `api/phase4_test.go`:

- **Quiz Attempt History**: 7 tests
- **Timed Test Taking**: 9 tests
- **Tag-Based Review**: 6 tests
- **Test Management**: 3 tests
- **Question Analytics**: 8 tests

**Total: 33 tests**

Run tests with:
```bash
go test ./api -v
```

---

## Future Enhancements

1. **Real-time Timer Sync**: WebSocket support for real-time timer synchronization
2. **Proctoring Features**: Camera monitoring, screen recording, tab switching detection
3. **Advanced Analytics**: Heat maps, learning curve visualization, predictive performance
4. **Adaptive Difficulty**: Dynamic difficulty adjustment during quizzes
5. **Team Tests**: Support for group assessments
6. **Question Banks**: Randomized question selection from banks
7. **Partial Credit**: Support for partially correct answers
8. **Rich Media**: Support for images, videos, code snippets in questions
