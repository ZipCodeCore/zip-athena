# Athena - Student Quiz Application

Athena is a specialized testing application for students who are learning how to develop software. It is a web-based three-tier application with an API business logic middle-end and a relational database backend.

## Current Status: Phase 5 Complete ✅

Athena now includes comprehensive student learning, instructor analytics, and admin content management, with FSRS-4.5 spaced repetition, timed tests, learning resources, bulk import, and quality review.

## Recent Updates

### February 2026 (Recent PRs)
- **Simplified migration flow** with a cleaner CLI and fewer moving parts
- **SQLite layer split into focused modules** for clearer database responsibilities
- **Routing cleanup** with consolidated API registration and consistent middleware
- **Frontend JS refactor** for smaller, better-factored UI code
- **Mobile-responsive UI improvements** across student and admin views
- **Quiz generation fixes** for filtering, difficulty handling, and duplicate prevention
- **Dashboard and admin fixes** (user stats, category display, question editing, labels, available topics)

### January 28, 2026
- **[Bug Fix]** Added comprehensive error logging to all API handlers
  - All database errors now logged to stderr with contextual information
  - Silent error cases now properly logged (quiz generation, FSRS updates, progress tracking)
  - Replaced stdout logging (fmt.Printf) with stderr logging (log.Printf)
  - Improves server monitoring and debugging capabilities

### January 28, 2026 (Earlier)
- **[Feature]** Quiz Retake System (Phase 2.5)
  - Students can retake quizzes with scores less than 100%
  - Exact same questions reproduced in retakes
  - FSRS integration: incorrect answers get shorter review intervals
  - Parent-child attempt tracking with attempt numbers
  - 13 new comprehensive tests added
  - Authorization validation ensures students only retake their own quizzes

### Phase 5 Features
- **Content Management & Quality**
  - Bulk CSV question import with validation
  - Tag management (merge, delete, analytics)
  - Flagged question review and resolution
  - Instructor question management with previews

### Phase 4 Features
- **Student Learning Enhancements**
  - Quiz attempt history and detailed review
  - Timed test-taking with results
  - Tag-based review quizzes
  - Learning resources and curriculum outline

### Phase 3 Features
- **Instructor Analytics**
  - Student progress tracking with detailed performance metrics
  - Cohort analytics for class-wide statistics
  - Individual student assessment views
- **Import/Export**
  - Quiz import/export in JSON format
  - Question bank import/export
  - Bulk question management capabilities
- **Test Management**
  - Create and manage formal tests
  - Question-to-test assignment
  - Test CRUD operations
- **Quiz-Question Relationships**
  - Dynamic quiz-question linking
  - Reusable question pools
  - Flexible quiz composition

### Phase 2 Features
- **FSRS-4.5 Spaced Repetition Algorithm**
  - Adaptive scheduling based on difficulty and performance
  - 20 comprehensive unit tests
  - Difficulty and stability tracking
- **Adaptive Quiz Generation**
  - Prioritizes due questions over new or future questions
  - Integration with spaced repetition system
- **Admin Dashboard**
  - Full CRUD operations for users, topics, questions, quizzes
  - User management and statistics
  - Question and quiz management
- **Enhanced Question Bank**
  - Java: 20 questions
  - Python: 18 questions
  - Go: 15 questions

### Phase 1 Features
- User authentication with bcrypt password hashing
- Quiz generation and taking
- Progress tracking
- Bootstrap-based responsive UI
- SQLite3 database backend

## Quick Start

```bash
# Seed the database with sample questions
go run cmd/seed/main.go

# Start the server
go run main.go

# Open browser to http://localhost:9206
```

## Database Migrations

Athena uses [golang-migrate](https://github.com/golang-migrate/migrate) for database version control and schema management. Migrations are automatically applied when the application starts.

### Manual Migration Management

```bash
# Build the migration tool
go build -o migrate ./cmd/migrate

# Check current migration version
./migrate -db athena.db -cmd version

# Apply all pending migrations
./migrate -db athena.db -cmd up

# Roll back the last migration
./migrate -db athena.db -cmd down -steps 1
```

See [DATABASE_MIGRATIONS.md](DATABASE_MIGRATIONS.md) for complete migration documentation.

## Architecture

- **Backend**: Go API server (Chi router) with structured logging and middleware
- **Middleware**: CORS, rate limiting, request size limits, content-type checks, JWT optional auth
- **Database**: SQLite3 with schema migrations via golang-migrate
- **Frontend**: HTML/JavaScript with Bootstrap, mobile-responsive layout
- **Algorithm**: FSRS-4.5 spaced repetition system

## Project Structure

```
zip-athena/
├── api/             # HTTP API handlers
├── config/          # Environment-based configuration
├── database/        # Database layer with SQLite implementation
│   ├── migrations/  # SQL migration files for schema versioning
│   └── migrate.go   # Migration runner functions
├── fsrs/            # FSRS-4.5 spaced repetition algorithm
├── learning/        # Learning resource content and guides
├── logger/          # Structured logging setup
├── models/          # Data models (User, Topic, Question, Answer, QuizAttempt, etc.)
├── questions/       # Seed question banks
├── router/          # Route registration and middleware
├── static/          # Frontend HTML, CSS, and JavaScript
├── cmd/
│   ├── seed/        # Database seeding utility
│   └── migrate/     # Migration management CLI tool
└── main.go          # Application entry point
```

## Running the Application

```bash
# Build the application
go build -o athena .

# Run the server
./athena

# Or run directly
go run main.go
```

The server runs on port 9206 by default. Configure with `PORT` and other env vars in [config/config.go](config/config.go).

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests for a specific package
go test -v ./models
go test -v ./database
go test -v ./api
go test -v ./fsrs
```

### Test Coverage

100+ tests across database, API handlers, FSRS, and phased features.

## API Endpoints

Most admin and instructor endpoints require JWT authentication. See [API_EXAMPLES.md](API_EXAMPLES.md).

### Authentication
- `POST /auth/register` - User registration
- `POST /auth/login` - User login

### Student Endpoints
- `POST /quiz/generate` - Generate a new quiz
- `POST /quiz/generate/by-tags` - Generate a tag-filtered quiz
- `POST /quiz/submit-answer` - Submit answer for a question
- `POST /quiz/complete` - Complete a quiz
- `POST /quiz/retake` - Retake a quiz (for scores < 100%)
- `GET /progress?user_id={id}` - Get student progress
- `GET /students/quiz-history` - Quiz history list
- `GET /students/quiz-history/{id}` - Quiz attempt details

### Tests (Student)
- `GET /tests/available` - Available tests
- `GET /tests/{id}` - Test details
- `POST /tests/{id}/start` - Start a test
- `PUT /test-attempts/{id}/answer` - Submit test answer
- `POST /test-attempts/{id}/submit` - Submit test
- `GET /test-attempts/{id}/results` - Test results

### FSRS Endpoints
- `POST /fsrs/review` - Review a question with rating (1-4)
- `GET /fsrs/due?user_id={id}` - Get due questions for review
- `POST /fsrs/quiz/generate` - Generate adaptive quiz based on FSRS

### Admin Endpoints (Protected)
- `GET /admin/users` - List all users
- `POST /admin/users` - Create user
- `PUT /admin/users` - Update user
- `DELETE /admin/users` - Delete user
- `POST /admin/questions/import` - CSV import questions
- `POST /admin/questions/import/validate` - Validate CSV
- `POST /admin/questions/import/confirm` - Confirm import
- `GET /admin/questions/flagged` - Flagged question queue
- `POST /admin/questions/auto-flag` - Auto-flag questions
- `GET /admin/questions` - List all questions
- `POST /admin/questions` - Create question
- `PUT /admin/questions` - Update question
- `DELETE /admin/questions` - Delete question
- `GET /admin/quizzes` - List all quizzes
- `POST /admin/quizzes` - Create quiz
- `PUT /admin/quizzes` - Update quiz
- `DELETE /admin/quizzes` - Delete quiz
- Similar endpoints for topics, answers, and tests

### Instructor Endpoints (Protected)
- `GET /instructor/student-progress?user_id={id}` - Get detailed student progress
- `GET /instructor/cohort-analytics` - Get cohort-wide statistics
- `GET /instructor/quizzes/export?id={id}` - Export quiz as JSON
- `POST /instructor/quizzes/import` - Import quiz from JSON
- `GET /instructor/questions/export` - Export questions as JSON
- `POST /instructor/questions/import` - Import questions from JSON
- `GET /instructor/questions` - List questions with filters
- `GET /instructor/questions/{id}/preview` - Preview question

## Key Features

- **Spaced Repetition Learning**: FSRS-4.5 algorithm for optimal review scheduling
- **Quiz Retakes**: Students can retake quizzes to improve understanding
- **Adaptive Quizzing**: Automatically prioritizes questions that need review
- **Comprehensive Error Logging**: All errors logged to stderr for monitoring
- **Admin Dashboard**: Full content management and quality review
- **Instructor Analytics**: Track student and cohort performance
- **Import/Export**: Bulk question and quiz management
- **Bulk CSV Import**: Validation and confirmation workflow
- **Tag-Based Review**: Filter quizzes by tags and track accuracy
- **Timed Tests**: Formal assessments with results export
- **Learning Resources**: Curriculum outline and resource linking
- **Thread-safe Database**: Concurrent operation support
- **RESTful API**: Clean, well-documented endpoints
- **Secure Authentication**: Bcrypt password hashing
- **Progress Tracking**: Detailed statistics and history
- **Tag-based Organization**: Categorize questions by topics and tags
- **Mobile Responsive UI**: Improved layouts for smaller screens

## Admin Access

Default admin account for development:
- **Username**: admin
- **Password**: admin
- **Capabilities**: Full access to all admin and instructor features

## Error Logging

All errors that prevent users from completing operations are logged to stderr with:
- Descriptive error messages
- Contextual information (user IDs, question IDs, etc.)
- Timestamp information via Go's log package
- Helps with debugging and monitoring in production

## Documentation

- `PHASE1_README.md` - Phase 1 implementation details
- `PHASE2_README.md` - Phase 2 requirements and design
- `PHASE2_COMPLETE.md` - Phase 2 completion summary
- `PHASE3_COMPLETE.md` - Phase 3 completion summary
- `PHASE4_COMPLETE.md` - Phase 4 completion summary
- `PHASE5_COMPLETE.md` - Phase 5 completion summary
- `QUIZ_RETAKE_FEATURE.md` - Quiz retake feature documentation
- `IMPLEMENTATION_SUMMARY.md` - Retake implementation details
- `IMPLEMENTATION_PHASES.md` - Complete phased implementation plan
- `IMPORT_EXPORT_FORMAT.md` - Import/export format specification
- `API_EXAMPLES.md` - API usage examples
- `SPECIFICATION.md` - Complete system specification
- `TEST_SUMMARY.md` - Testing documentation

## Technology Stack

- **Language**: Go 1.x
- **Database**: SQLite3 (production-ready for small to medium deployments)
- **Password Hashing**: bcrypt
- **Frontend**: Vanilla JavaScript with Bootstrap 5
- **Algorithm**: FSRS-4.5 (research-backed spaced repetition)

## Future Enhancements

- Real-time progress notifications
- Question difficulty auto-adjustment
- Gamification features (achievements, leaderboards)
- Mobile app support
- Advanced analytics dashboards
- Question version control
- Collaborative quiz creation
