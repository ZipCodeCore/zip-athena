# Import/Export Format Documentation

This document describes the JSON formats used for importing and exporting Questions, Quizzes, and Tests in the Athena Learning Platform.

## Table of Contents

1. [Questions Format](#questions-format)
2. [Quizzes Format](#quizzes-format)
3. [Tests Format](#tests-format)
4. [Usage Examples](#usage-examples)
5. [API Endpoints](#api-endpoints)

---

## Questions Format

### Question with Answers Structure

Questions are exported/imported with their associated answers in the following format:

```json
{
  "question": {
    "topic_id": 1,
    "question_text": "What is the output of System.out.println(5 + 2)?",
    "difficulty": 2,
    "tags": ["java", "basics", "operators"],
    "topic_question_id": "1.5"
  },
  "answers": [
    {
      "answer_text": "7",
      "is_correct": true
    },
    {
      "answer_text": "52",
      "is_correct": false
    },
    {
      "answer_text": "Error",
      "is_correct": false
    },
    {
      "answer_text": "5 + 2",
      "is_correct": false
    }
  ]
}
```

### Field Descriptions

**Question Object:**
- `topic_id` (integer, required): The ID of the topic this question belongs to
- `question_text` (string, required): The text of the question
- `difficulty` (integer, required): Question difficulty level (1-7, where 1 is easiest)
- `tags` (array of strings, optional): Tags for categorizing the question
- `topic_question_id` (string, required): Unique identifier within the topic (e.g., "1.5" for topic 1, question 5)

**Answer Object:**
- `answer_text` (string, required): The text of the answer option
- `is_correct` (boolean, required): Whether this answer is correct (at least one answer must be correct)

### Bulk Questions Export/Import

When exporting or importing multiple questions, use an array of Question with Answers objects:

```json
[
  {
    "question": { ... },
    "answers": [ ... ]
  },
  {
    "question": { ... },
    "answers": [ ... ]
  }
]
```

---

## Quizzes Format

### Quiz Export Structure

When exporting a quiz, the complete quiz with all associated questions and answers is returned:

```json
{
  "quiz": {
    "id": 1,
    "name": "java.1",
    "topic_id": 1,
    "description": "Introduction to Java - Basic Concepts",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "questions": [
    {
      "question": {
        "id": 5,
        "topic_id": 1,
        "question_text": "What is Java?",
        "difficulty": 1,
        "tags": ["java", "basics"],
        "topic_question_id": "1.1",
        "created_at": "2024-01-10T09:00:00Z"
      },
      "answers": [
        {
          "id": 17,
          "question_id": 5,
          "answer_text": "A programming language",
          "is_correct": true
        },
        {
          "id": 18,
          "question_id": 5,
          "answer_text": "A type of coffee",
          "is_correct": false
        }
      ]
    }
  ]
}
```

### Quiz Import Structure

When importing a quiz, provide the quiz metadata and questions (IDs are auto-generated):

```json
{
  "name": "python.1",
  "topic_id": 2,
  "description": "Introduction to Python Programming",
  "questions": [
    {
      "question": {
        "topic_id": 2,
        "question_text": "What is Python primarily used for?",
        "difficulty": 1,
        "tags": ["python", "basics", "intro"],
        "topic_question_id": "2.1"
      },
      "answers": [
        {
          "answer_text": "Web development and data science",
          "is_correct": true
        },
        {
          "answer_text": "Only mobile app development",
          "is_correct": false
        },
        {
          "answer_text": "Database management only",
          "is_correct": false
        }
      ]
    },
    {
      "question": {
        "topic_id": 2,
        "question_text": "Which of these is a valid Python variable name?",
        "difficulty": 2,
        "tags": ["python", "variables", "syntax"],
        "topic_question_id": "2.2"
      },
      "answers": [
        {
          "answer_text": "my_variable",
          "is_correct": true
        },
        {
          "answer_text": "2variable",
          "is_correct": false
        },
        {
          "answer_text": "my-variable",
          "is_correct": false
        }
      ]
    }
  ]
}
```

### Field Descriptions (Quiz Import)

- `name` (string, required): Name of the quiz (e.g., "python.1", "java.advanced")
- `topic_id` (integer, required): The ID of the topic this quiz belongs to
- `description` (string, optional): Description of the quiz content and objectives
- `questions` (array, required): Array of question-with-answers objects (see Questions Format)

**Note:** Questions are automatically linked to the quiz in the order they appear in the array (order 0, 1, 2, etc.)

---

## Tests Format

Tests use a similar structure to Quizzes but include additional fields for formal assessments.

### Test Structure

```json
{
  "id": 1,
  "name": "Java Midterm Exam",
  "topic_id": 1,
  "description": "Comprehensive midterm covering Java basics through OOP",
  "time_limit": 90,
  "created_at": "2024-01-20T14:00:00Z",
  "updated_at": "2024-01-20T14:00:00Z"
}
```

### Field Descriptions

- `id` (integer, auto-generated): Unique identifier for the test
- `name` (string, required): Name of the test
- `topic_id` (integer, required): The ID of the topic this test covers
- `description` (string, optional): Description of the test scope and content
- `time_limit` (integer, optional): Time limit in minutes (0 or omitted = no time limit)
- `created_at` (timestamp, auto-generated): When the test was created
- `updated_at` (timestamp, auto-generated): When the test was last modified

**Note:** Like quizzes, tests can have questions linked to them via the many-to-many relationship endpoints.

---

## Usage Examples

### Example 1: Exporting a Quiz

**Request:**
```http
GET /instructor/quizzes/export?id=1
X-User-ID: 1
```

**Response:**
```json
{
  "quiz": {
    "id": 1,
    "name": "java.1",
    "topic_id": 1,
    "description": "Introduction to Java",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "questions": [...]
}
```

### Example 2: Importing a Quiz

**Request:**
```http
POST /instructor/quizzes/import
Content-Type: application/json
X-User-ID: 1

{
  "name": "sql.basic",
  "topic_id": 3,
  "description": "SQL Fundamentals",
  "questions": [
    {
      "question": {
        "topic_id": 3,
        "question_text": "What does SQL stand for?",
        "difficulty": 1,
        "tags": ["sql", "basics"],
        "topic_question_id": "3.1"
      },
      "answers": [
        {
          "answer_text": "Structured Query Language",
          "is_correct": true
        },
        {
          "answer_text": "Simple Question Language",
          "is_correct": false
        }
      ]
    }
  ]
}
```

**Response:**
```json
{
  "message": "Quiz imported successfully",
  "quiz_id": 7,
  "quiz_name": "sql.basic",
  "question_count": 1
}
```

### Example 3: Exporting All Questions

**Request:**
```http
GET /instructor/questions/export
X-User-ID: 1
```

**Response:**
```json
[
  {
    "question": {
      "id": 1,
      "topic_id": 1,
      "question_text": "What is Java?",
      "difficulty": 1,
      "tags": ["java", "basics"],
      "topic_question_id": "1.1",
      "created_at": "2024-01-10T09:00:00Z"
    },
    "answers": [...]
  },
  {
    "question": {...},
    "answers": [...]
  }
]
```

### Example 4: Importing Questions

**Request:**
```http
POST /instructor/questions/import
Content-Type: application/json
X-User-ID: 1

[
  {
    "question": {
      "topic_id": 2,
      "question_text": "What is a Python list?",
      "difficulty": 2,
      "tags": ["python", "data-structures"],
      "topic_question_id": "2.15"
    },
    "answers": [
      {
        "answer_text": "An ordered, mutable collection",
        "is_correct": true
      },
      {
        "answer_text": "An unordered set",
        "is_correct": false
      }
    ]
  }
]
```

**Response:**
```json
{
  "message": "Questions imported successfully",
  "count": 1
}
```

### Example 5: Adding a Question to a Quiz

**Request:**
```http
POST /admin/quiz-questions
Content-Type: application/json
X-User-ID: 1

{
  "quiz_id": 1,
  "question_id": 15,
  "order": 3
}
```

**Response:**
```json
{
  "message": "Question added to quiz successfully"
}
```

### Example 6: Adding a Question to a Test

**Request:**
```http
POST /admin/test-questions
Content-Type: application/json
X-User-ID: 1

{
  "test_id": 1,
  "question_id": 20,
  "order": 5
}
```

**Response:**
```json
{
  "message": "Question added to test successfully"
}
```

---

## API Endpoints

### Import/Export Endpoints

All import/export endpoints require admin authentication via the `X-User-ID` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/instructor/quizzes/export?id={quiz_id}` | Export a single quiz with all questions and answers |
| POST | `/instructor/quizzes/import` | Import a quiz with questions and answers |
| GET | `/instructor/questions/export` | Export all questions with answers |
| POST | `/instructor/questions/import` | Import multiple questions with answers |

### Quiz-Question Relationship Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/quiz-questions?quiz_id={id}` | Get all questions for a quiz (in order) |
| POST | `/admin/quiz-questions` | Add a question to a quiz |
| DELETE | `/admin/quiz-questions` | Remove a question from a quiz |

### Test-Question Relationship Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/test-questions?test_id={id}` | Get all questions for a test (in order) |
| POST | `/admin/test-questions` | Add a question to a test |
| DELETE | `/admin/test-questions` | Remove a question from a test |

### Test Management Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/admin/tests` | List all tests |
| POST | `/admin/tests` | Create a new test |
| PUT | `/admin/tests?id={test_id}` | Update a test |
| DELETE | `/admin/tests?id={test_id}` | Delete a test |

---

## Best Practices

1. **Topic ID Validation**: Always ensure the `topic_id` exists before importing questions or quizzes
2. **Unique topic_question_id**: Each question should have a unique `topic_question_id` within its topic
3. **Answer Correctness**: Always include at least one correct answer for each question
4. **Difficulty Scaling**: Use difficulty levels consistently (1=easiest, 7=hardest)
5. **Tags**: Use consistent tag naming conventions for better organization
6. **Order Management**: When adding questions to quizzes/tests, use sequential order values (0, 1, 2, ...)
7. **Backup**: Always export quizzes before making significant changes
8. **Validation**: Validate JSON structure before importing to avoid errors

---

## Notes

- All timestamps are in ISO 8601 format (UTC)
- IDs are auto-generated and should not be included in import requests
- Import operations create new records; they do not update existing ones
- The many-to-many relationship allows questions to be reused across multiple quizzes and tests
- Questions maintain their order within quizzes and tests via the `order` field
- Export files can be saved and re-imported to duplicate content or migrate between instances
