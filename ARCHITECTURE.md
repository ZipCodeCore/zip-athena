# Athena Architecture Guide

## Overview

Athena is a three-tier web application for student learning through adaptive quizzing. This document explains how the application is structured and how its components work together.

## High-Level Architecture

```
┌─────────────────┐
│   Web Browser   │
│  (Static HTML/  │
│   JavaScript)   │
└────────┬────────┘
         │ HTTP/JSON
         │
┌────────▼────────┐
│   Go API Server │
│   (Chi Router)  │
│   + Middleware  │
└────────┬────────┘
         │ SQL
         │
┌────────▼────────┐
│  SQLite Database│
│   (athena.db)   │
└─────────────────┘
```

## Application Layers

### 1. Frontend Layer (Static Files)

**Location**: `/static/`

The frontend is built with vanilla JavaScript and Bootstrap 5, providing a responsive web interface.

**Key Components**:
- `index.html` - Login and registration page
- `dashboard.html` - Student dashboard showing progress and available quizzes
- `quiz.html` - Interactive quiz-taking interface
- `admin.html` - Admin content management dashboard
- `instructor.html` - Instructor analytics and student progress views
- `learning.html` - Learning resources and curriculum browser

**Frontend Architecture**:
- No build step required - works directly in the browser
- Uses Bootstrap for responsive design
- API communication via `fetch()` for REST calls
- JWT token stored in `localStorage` for authentication
- Mobile-responsive design for accessibility

### 2. API/Business Logic Layer (Go)

**Location**: `/api/`, `/router/`, `/main.go`

The backend is a RESTful API server built with Go and the Chi router.

#### Request Flow

```
HTTP Request
    ↓
Security Headers Middleware
    ↓
CORS Middleware
    ↓
Rate Limiting Middleware (60 req/min)
    ↓
Request Size Middleware (10MB max)
    ↓
Content Type Middleware
    ↓
JWT Optional Middleware
    ↓
Route Handler (api/*.go)
    ↓
Database Layer (database/*.go)
    ↓
HTTP Response (JSON)
```

#### Key API Modules

**`api/handlers.go`** - Core quiz generation and submission logic
- Quiz generation with difficulty balancing
- Answer submission and validation
- Progress tracking

**`api/auth.go`** - User authentication
- Registration with bcrypt password hashing
- Login with JWT token generation
- Token validation

**`api/admin.go`** - Admin content management
- CRUD operations for users, questions, topics, quizzes
- Bulk CSV import with validation
- Question flagging and review

**`api/instructor.go`** - Instructor analytics
- Student progress tracking
- Cohort analytics
- Quiz/question import/export

**`api/phase4.go`** - Advanced student features
- Quiz history and detailed review
- Timed test-taking
- Learning resource access

**`api/phase5.go`** - Content quality management
- Flagged question review
- Tag management and analytics
- Question preview

**`api/fsrs_handlers.go`** - Spaced repetition algorithm
- FSRS-4.5 adaptive scheduling
- Review submission with ratings
- Due question prioritization

#### Middleware Stack

**`api/middleware.go`** contains:
- **Security Headers**: CSP, X-Frame-Options, etc.
- **CORS**: Configurable cross-origin resource sharing
- **Rate Limiting**: Token bucket algorithm
- **Request Size Limiting**: Prevents large request attacks
- **Content Type Validation**: Ensures proper content types
- **JWT Optional**: Attaches user context if token present

### 3. Data Layer (SQLite)

**Location**: `/database/`, `/models/`

The data layer uses SQLite3 for persistence with a clean abstraction interface.

#### Database Architecture

**Interface Design** (`database/database.go`):
- Defines all database operations as interface methods
- Allows for easy testing and potential DB engine swaps
- Clean separation between business logic and data access

**SQLite Implementation** (`database/sqlite.go`):
- Connection pooling with `sql.DB`
- Prepared statements for performance
- Transaction support
- Thread-safe concurrent operations

**Migrations** (`database/migrations/`):
- Schema versioning with golang-migrate
- Up/down migrations for rollback capability
- Automatic application on startup

#### Data Models

**`models/models.go`** defines core entities:

```go
User              // Student/instructor/admin accounts
Topic             // Subject areas (Java, Python, etc.)
Question          // Quiz questions with difficulty and tags
Answer            // Multiple choice options
QuizAttempt       // User quiz sessions with scores
QuestionAttempt   // Individual question history
Quiz              // Named quiz collections
UserQuestionState // FSRS spaced repetition state
Test              // Formal timed assessments
TestAttempt       // Test-taking sessions
LearningResource  // Curriculum content links
FlaggedQuestion   // Quality review queue
```

#### Database Schema Highlights

**Users & Authentication**:
- bcrypt password hashing
- Role-based access (is_admin flag)
- Category preference for personalization

**Questions & Content**:
- Tagged for organization
- Difficulty levels (1-4)
- Topic associations
- Learning resource links

**Quiz System**:
- Flexible quiz-question relationships
- Parent-child attempt tracking for retakes
- Detailed attempt metadata

**FSRS Spaced Repetition**:
- Per-user, per-question state tracking
- Stability and difficulty metrics
- Due date scheduling
- Review history

## Core Algorithms

### 1. FSRS Spaced Repetition

**Location**: `/fsrs/fsrs.go`

Athena uses the FSRS-4.5 (Free Spaced Repetition Scheduler) algorithm for optimal learning.

**How It Works**:
1. Each user-question pair has a `UserQuestionState`
2. When a question is answered, user rates difficulty (1-4):
   - 1 = Again (failed)
   - 2 = Hard
   - 3 = Good
   - 4 = Easy
3. Algorithm calculates:
   - **Stability**: How long knowledge persists
   - **Difficulty**: Question's inherent challenge for user
   - **Next Review**: When to show question again
4. States track learning progress:
   - New → Learning → Review → Relearning (if failed)

**Key Concepts**:
- Questions due for review are prioritized
- Successful reviews increase interval
- Failed reviews reset to shorter intervals
- Optimal retention through spaced intervals

### 2. Adaptive Quiz Generation

**Location**: `api/handlers.go` - `GenerateAdaptiveQuiz()`

Quizzes are generated intelligently based on FSRS state:

```
Priority Order:
1. Due Questions (past due date)
2. New Questions (never reviewed)
3. Future Questions (not yet due)
```

This ensures students review material at optimal times.

### 3. Quiz Retake System

**Location**: `api/quiz.go` - `RetakeQuiz()`

Students can retake quizzes with scores < 100%:
- Exact same questions reproduced
- Tracks as child attempt of original
- FSRS integrates incorrect answers
- Attempt number incremented

## Configuration

**Location**: `/config/config.go`

Environment-based configuration with sensible defaults:

```go
Server:
  - PORT (default: 9206)
  - Read/Write/Idle timeouts

Database:
  - DB_TYPE (sqlite)
  - DB_PATH (athena.db)

Security:
  - JWT_SECRET
  - Token expiration

CORS:
  - Allowed origins/methods/headers

Logging:
  - Level (info, debug, error)
  - Output format
```

## Security Features

1. **Authentication**:
   - bcrypt password hashing (cost 10)
   - JWT tokens for stateless auth
   - Secure token storage

2. **Authorization**:
   - Role-based access control
   - Admin-only endpoints protected
   - User isolation (can't access other users' data)

3. **Input Validation**:
   - Request size limits (10MB)
   - Content-type validation
   - SQL injection prevention via prepared statements
   - XSS protection via JSON encoding

4. **Rate Limiting**:
   - 60 requests per minute per IP
   - Burst allowance of 20 requests
   - Prevents abuse and DoS

5. **Security Headers**:
   - Content Security Policy
   - X-Frame-Options (prevent clickjacking)
   - X-Content-Type-Options
   - Referrer-Policy

## Performance Considerations

1. **Database**:
   - Indexes on frequently queried columns
   - Connection pooling
   - Prepared statements cached
   - Efficient joins and aggregations

2. **API**:
   - Minimal JSON payloads
   - Pagination for large lists
   - Streaming responses for exports

3. **Frontend**:
   - No build step = fast development
   - CDN-hosted Bootstrap (cached by browsers)
   - Lazy loading of resources

## Testing Strategy

**Location**: `*_test.go` files throughout codebase

Athena has 100+ tests across multiple layers:

1. **Unit Tests**:
   - FSRS algorithm (20 tests)
   - Data models (validation, serialization)

2. **Integration Tests**:
   - API endpoints (60+ tests)
   - Database operations
   - Authentication flows

3. **Feature Tests**:
   - Complete user workflows
   - Quiz generation and submission
   - Admin operations
   - Instructor analytics

**Running Tests**:
```bash
go test -v ./...           # All tests
go test -v ./fsrs          # FSRS algorithm
go test -v ./api           # API handlers
go test -v ./database      # Database layer
```

## Deployment Architecture

**Development**:
```bash
go run main.go
```

**Production** (containerized):
```dockerfile
FROM golang:1.x
COPY . /app
RUN go build -o athena
CMD ["/app/athena"]
```

**Database**:
- SQLite file (`athena.db`)
- Migrations auto-applied on startup
- Backups recommended before schema changes

## Extending the Application

### Adding a New API Endpoint

1. Define handler in appropriate `api/*.go` file
2. Add route in `router/router.go`
3. Add tests in `api/*_test.go`
4. Document in `API_EXAMPLES.md`

### Adding a New Database Table

1. Create migration in `database/migrations/`
2. Add model to `models/models.go`
3. Add interface methods to `database/database.go`
4. Implement in `database/sqlite.go`
5. Write tests

### Adding Frontend Features

1. Update HTML in `static/`
2. Add JavaScript logic in `static/js/`
3. Call new API endpoints
4. Test across browsers

## Key Design Decisions

1. **Why Go?**
   - Fast, compiled language
   - Excellent concurrency
   - Great standard library
   - Easy deployment (single binary)

2. **Why SQLite?**
   - Zero configuration
   - Perfect for small to medium deployments
   - File-based (easy backups)
   - ACID compliant

3. **Why No Frontend Framework?**
   - Simplicity for learning
   - No build step complexity
   - Easy to understand for students
   - Fast iteration

4. **Why FSRS?**
   - Research-backed algorithm
   - Superior to traditional Anki/SM-2
   - Open source
   - Proven effective

## Learning from This Codebase

This application demonstrates:

1. **Three-tier architecture** - Clean separation of concerns
2. **RESTful API design** - Standard HTTP methods and status codes
3. **Database abstraction** - Interface-based design for flexibility
4. **Middleware patterns** - Cross-cutting concerns (auth, logging, etc.)
5. **Security best practices** - Hashing, tokens, rate limiting, validation
6. **Testing strategies** - Unit, integration, and feature tests
7. **Configuration management** - Environment-based settings
8. **Error handling** - Structured logging and error responses
9. **Algorithms in practice** - FSRS spaced repetition
10. **Code organization** - Modular packages with clear responsibilities

## Further Reading

- **Getting Started**: `GETTING_STARTED.md` - Set up and run the app
- **Development**: `DEVELOPMENT_GUIDE.md` - Contribute to the codebase
- **Learning**: `LEARNING_GUIDE.md` - Use this app to learn software development
- **API Reference**: `API_EXAMPLES.md` - Complete API documentation
- **Database**: `DATABASE_MIGRATIONS.md` - Schema management
