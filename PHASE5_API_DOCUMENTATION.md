# Phase 5 API Documentation

## Authentication
All Phase 5 endpoints require authentication via the `X-User-ID` header and admin privileges.

```
X-User-ID: <user_id>
```

The user must have `is_admin: true` in the database.

---

## Question Management APIs

### List Questions (Paginated)
Get a paginated list of questions with optional filters.

**Endpoint:** `GET /instructor/questions`

**Query Parameters:**
- `page` (optional, default: 1) - Page number
- `per_page` (optional, default: 30, max: 100) - Results per page
- `search` (optional) - Search in question text
- `topic_id` (optional) - Filter by topic ID
- `difficulty` (optional) - Filter by difficulty (1-7)

**Response:**
```json
{
  "questions": [
    {
      "id": 1,
      "topic_id": 1,
      "question_text": "What is a goroutine?",
      "difficulty": 3,
      "tags": ["concurrency", "basics"],
      "topic_question_id": "1.1",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "per_page": 30,
  "total_pages": 4
}
```

### Create Question
Create a new question with answers.

**Endpoint:** `POST /instructor/questions`

**Request Body:**
```json
{
  "question": {
    "topic_id": 1,
    "question_text": "What is a goroutine?",
    "difficulty": 3,
    "tags": ["concurrency", "basics"],
    "topic_question_id": "1.1"
  },
  "answers": [
    {"answer_text": "A lightweight thread", "is_correct": true},
    {"answer_text": "A function", "is_correct": false},
    {"answer_text": "A package", "is_correct": false},
    {"answer_text": "A variable", "is_correct": false}
  ]
}
```

**Validation Rules:**
- Question text: required, max 1000 chars
- Answers: 2-4 required
- Exactly 1 correct answer required
- Topic: must exist
- Difficulty: 1-7 required

**Response:** `201 Created`
```json
{
  "question": {...},
  "answers": [...]
}
```

### Update Question
Update an existing question and its answers.

**Endpoint:** `PUT /instructor/questions?id={question_id}`

**Request Body:** Same as Create Question

**Response:** `200 OK`
```json
{
  "question": {...},
  "answers": [...]
}
```

### Deactivate Question
Soft delete a question (sets `is_active` to false).

**Endpoint:** `DELETE /instructor/questions?id={question_id}`

**Response:** `204 No Content`

### Preview Question
Get question preview as students see it.

**Endpoint:** `GET /instructor/questions/{id}/preview`

**Response:**
```json
{
  "question": {
    "id": 1,
    "topic_id": 1,
    "question_text": "What is a goroutine?",
    "difficulty": 3,
    "tags": ["concurrency", "basics"],
    "topic_question_id": "1.1",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "answers": [
    {"id": 1, "question_id": 1, "answer_text": "A lightweight thread", "is_correct": true},
    {"id": 2, "question_id": 1, "answer_text": "A function", "is_correct": false}
  ]
}
```

---

## Bulk CSV Import APIs

### CSV Format
```csv
Topic,Question,Answer1,Answer2,Answer3,Answer4,CorrectAnswer,Difficulty,Tags
Java,"What is Java?","A language","A drink","An island","","A language",2,"basics,oop"
Python,"What is a list?","Array","Set","Dict","","Array",1,"data-structures"
```

### Upload and Import CSV
Upload a CSV file and immediately import valid questions.

**Endpoint:** `POST /admin/questions/import`

**Request:**
- Content-Type: `multipart/form-data`
- Field: `file` (CSV file)

**Response:**
```json
{
  "valid_rows": [
    {
      "topic": "Java",
      "question": "What is Java?",
      "answer1": "A language",
      "answer2": "A drink",
      "answer3": "An island",
      "answer4": "",
      "correct_answer": "A language",
      "difficulty": "2",
      "tags": "basics,oop"
    }
  ],
  "invalid_rows": [
    {
      "row": 3,
      "data": {...},
      "errors": [
        "Topic is required",
        "Difficulty must be between 1 and 7"
      ]
    }
  ],
  "summary": {
    "total_rows": 100,
    "valid_rows": 95,
    "invalid_rows": 5,
    "imported_rows": 95
  }
}
```

### Validate CSV Without Importing
Validate a CSV file without importing.

**Endpoint:** `POST /admin/questions/import/validate`

**Request:** Same as Upload and Import

**Response:** Same as Upload and Import (but `imported_rows` will be 0)

### Confirm Import
Import previously validated rows.

**Endpoint:** `POST /admin/questions/import/confirm`

**Request Body:**
```json
[
  {
    "topic": "Java",
    "question": "What is Java?",
    "answer1": "A language",
    "answer2": "A drink",
    "correct_answer": "A language",
    "difficulty": "2",
    "tags": "basics"
  }
]
```

**Response:**
```json
{
  "imported": 1
}
```

---

## Tag Management APIs

### List All Tags
Get all tags with metadata.

**Endpoint:** `GET /admin/tags`

**Response:**
```json
[
  {
    "name": "concurrency",
    "description": "",
    "question_count": 15
  },
  {
    "name": "basics",
    "description": "",
    "question_count": 42
  }
]
```

### Merge Tags
Merge multiple tags into one target tag.

**Endpoint:** `POST /admin/tags/merge`

**Request Body:**
```json
{
  "source_tags": ["goroutines", "go-routines"],
  "target_tag": "concurrency"
}
```

**Response:**
```json
{
  "message": "Tags merged successfully"
}
```

### Delete Tag
Delete a tag (only if no questions use it).

**Endpoint:** `DELETE /admin/tags?name={tag_name}`

**Response:** `204 No Content`

**Error (if tag in use):**
```json
{
  "error": "cannot delete tag: tag is in use by questions"
}
```

### Get Questions by Tag
Retrieve all questions with a specific tag.

**Endpoint:** `GET /admin/tags/questions?name={tag_name}`

**Response:**
```json
[
  {
    "id": 1,
    "topic_id": 1,
    "question_text": "What is a goroutine?",
    "difficulty": 3,
    "tags": ["concurrency", "basics"],
    "topic_question_id": "1.1",
    "created_at": "2024-01-01T00:00:00Z"
  }
]
```

---

## Quality Review APIs

### Get Flagged Questions
Retrieve the queue of flagged questions.

**Endpoint:** `GET /admin/questions/flagged`

**Response:**
```json
[
  {
    "question_id": 1,
    "topic_question_id": "1.1",
    "question_text": "What is a goroutine?",
    "flag_reason": "TOO_EASY",
    "flagged_at": "2024-01-01T00:00:00Z",
    "accuracy": 95.5,
    "total_attempts": 200
  }
]
```

**Flag Reasons:**
- `TOO_HARD` - Accuracy < 30%
- `TOO_EASY` - Accuracy > 90%
- `AMBIGUOUS` - Mixed feedback
- `POOR_DISTRACTORS` - Skewed answer distribution
- `MANUAL` - Instructor flagged

### Flag Question
Manually flag a question for review.

**Endpoint:** `POST /admin/questions/flag?id={question_id}`

**Request Body:**
```json
{
  "reason": "MANUAL"
}
```

**Response:**
```json
{
  "message": "Question flagged successfully"
}
```

### Unflag Question
Remove a flag from a question.

**Endpoint:** `DELETE /admin/questions/flag?id={question_id}`

**Response:** `204 No Content`

### Resolve Question
Mark a flagged question as resolved (same as unflag).

**Endpoint:** `PUT /admin/questions/resolve?id={question_id}`

**Response:** `204 No Content`

### Auto-Flag Questions
Run automatic flagging based on performance analytics.

**Endpoint:** `POST /admin/questions/auto-flag`

**Criteria:**
- Flags questions with accuracy < 30% as `TOO_HARD`
- Flags questions with accuracy > 90% as `TOO_EASY`
- Requires at least 10 attempts per question

**Response:**
```json
{
  "message": "Questions auto-flagged successfully"
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request parameters"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized: User ID required"
}
```

### 403 Forbidden
```json
{
  "error": "Forbidden: Admin access required"
}
```

### 404 Not Found
```json
{
  "error": "Record not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## Usage Examples

### Example 1: Create a Question
```bash
curl -X POST http://localhost:9206/instructor/questions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "question": {
      "topic_id": 1,
      "question_text": "What is Go?",
      "difficulty": 2,
      "tags": ["basics"],
      "topic_question_id": "1.1"
    },
    "answers": [
      {"answer_text": "A programming language", "is_correct": true},
      {"answer_text": "A game", "is_correct": false}
    ]
  }'
```

### Example 2: Import Questions from CSV
```bash
curl -X POST http://localhost:9206/admin/questions/import \
  -H "X-User-ID: 1" \
  -F "file=@questions.csv"
```

### Example 3: Merge Tags
```bash
curl -X POST http://localhost:9206/admin/tags/merge \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 1" \
  -d '{
    "source_tags": ["go-routines", "goroutines"],
    "target_tag": "concurrency"
  }'
```

### Example 4: Get Flagged Questions
```bash
curl -X GET http://localhost:9206/admin/questions/flagged \
  -H "X-User-ID: 1"
```

---

## Notes

1. **Pagination**: All list endpoints support pagination. Use `page` and `per_page` parameters.
2. **Filtering**: Question list supports multiple filters that can be combined.
3. **Validation**: All create/update operations validate input before processing.
4. **Transactions**: CSV import uses transactions to ensure data integrity.
5. **Soft Deletes**: Questions are deactivated, not permanently deleted.
6. **Case Insensitive**: Tag operations are case-insensitive.
7. **Auto-Creation**: CSV import automatically creates topics if they don't exist.
8. **Performance**: Bulk import tested with 500+ questions successfully.
