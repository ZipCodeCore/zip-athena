# Athena - Student Quiz Application

Athena is a specialized testing application for students who are learning how to develop software. It is a web-based three-tier application with an API business logic middle-end and a relational database backend.

## Current Status: Phase 1 MVP Complete ✅

Phase 1 implements the core student experience including user authentication, quiz taking, progress tracking, and a Bootstrap-based UI. See [PHASE1_README.md](PHASE1_README.md) for detailed usage instructions.

## Quick Start (Phase 1)

```bash
# Seed the database with sample questions
go run cmd/seed/main.go

# Start the server
go run main.go

# Open browser to http://localhost:9206
```

## Phase 0 - Core Implementation

This phase includes the foundational structure with 20 comprehensive tests covering:
- Data models
- Database operations
- API handlers
- Integration tests

## Architecture

- **Backend**: Go API server
- **Database**: SQLite3 (Phase 1) - ready for PostgreSQL/MySQL migration
- **Frontend**: HTML/JavaScript with Bootstrap CSS framework

## Project Structure

```
zip-athena/
├── models/          # Data models (User, Topic, Question, Answer, etc.)
├── database/        # Database layer with interface and in-memory implementation
├── api/             # HTTP API handlers
└── main.go          # Application entry point
```

## Running the Application

```bash
# Build the application
go build -o athena .

# Run the server
./athena
```

The server runs on port 9206.

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests for a specific package
go test -v ./models
go test -v ./database
go test -v ./api
```

## Test Coverage (20 Tests)

### Model Tests (5 tests)
1. **TestUserModel** - Validates User model structure and fields
2. **TestTopicModel** - Validates Topic model structure
3. **TestQuestionModelWithTags** - Tests Question model with tags and difficulty
4. **TestAnswerModel** - Tests Answer model with correctness tracking
5. **TestQuizAttemptModel** - Tests QuizAttempt with results tracking

### Database Tests (10 tests)
6. **TestCreateAndGetUser** - Tests user creation and retrieval
7. **TestGetNonExistentUser** - Tests error handling for missing users
8. **TestUpdateUserPoints** - Tests updating user points
9. **TestCreateAndListTopics** - Tests topic creation and listing
10. **TestGetQuestionsByTopic** - Tests filtering questions by topic
11. **TestCreateAndGetAnswers** - Tests answer creation and retrieval
12. **TestCreateAndGetQuizAttempt** - Tests quiz attempt tracking
13. **TestGetQuizAttemptsByUser** - Tests filtering quiz attempts by user
14. **TestCreateAndGetQuestionAttempt** - Tests individual question attempt tracking
15. **TestConcurrentUserCreation** - Tests thread-safety with concurrent operations

### API Handler Tests (5 tests)
16. **TestCreateUserHandler** - Tests POST /users endpoint
17. **TestGetUserHandler** - Tests GET /users endpoint
18. **TestCreateTopicHandler** - Tests POST /topics endpoint
19. **TestListTopicsHandler** - Tests GET /topics endpoint
20. **TestCreateQuizAttemptHandler** - Tests POST /quiz-attempts endpoint

## API Endpoints

- `POST /users` - Create a new user
- `GET /users?id={id}` - Get user by ID
- `POST /topics` - Create a new topic
- `GET /topics` - List all topics
- `POST /questions` - Create a new question
- `GET /questions?topic_id={id}` - Get questions by topic
- `POST /quiz-attempts` - Create a new quiz attempt
- `GET /quiz-attempts?user_id={id}` - Get quiz attempts by user

## Key Features

- Thread-safe in-memory database
- RESTful API design
- Comprehensive test coverage
- Space repetition system ready (data model supports it)
- Difficulty tracking (1-7 scale)
- Tag-based question categorization
- Quiz attempt tracking with detailed results
- Question attempt history for spaced repetition algorithms

## Next Steps

- Implement frontend JavaScript application
- Add Sqlite3 database backend
- Implement spaced repetition algorithm
- Add authentication and authorization
- Implement test attempt functionality
- Add admin dashboard
- Implement student progress tracking
