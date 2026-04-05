# Learning Guide: Using Athena to Learn Software Development

## Introduction

Athena is not just a quiz application—it's a complete, production-quality codebase designed to teach modern software development practices. This guide helps you learn by exploring real code that solves real problems.

## Who This Guide is For

- **Bootcamp students** learning Go, web APIs, or databases
- **Self-taught developers** wanting to see how real apps are structured
- **Junior developers** studying best practices and patterns
- **Anyone** curious about building three-tier web applications

## Understanding Student Categories

Athena uses a **category** system to tailor topic visibility to each student's curriculum:

| User Category | Topics Visible |
|---------------|----------------|
| `java`        | Java Fundamentals, Java Code Tracing, SQL Database |
| `data`        | Python Basics, Python Code Tracing, SQL Database |

This is enforced by `api/topic_access.go:userCanAccessTopic()`. Admins bypass this filter and see all topics.

Code for this lives in:
- `constants/constants.go` — defines `CategoryJava`, `CategoryData`, `TopicCategoryJava`, `TopicCategoryPython`, `TopicCategoryAll`
- `api/topic_access.go` — the access-control predicate used by quiz and dashboard endpoints
- `database/migrations/005_topic_category_all.up.sql` — how the category system evolved

## Learning Objectives

By studying Athena, you'll learn:

1. **Three-tier architecture** - How frontend, backend, and database work together
2. **RESTful API design** - Building clean, intuitive HTTP endpoints
3. **Database design** - Schema design, migrations, and SQL
4. **Security** - Authentication, authorization, and data protection
5. **Testing** - Unit tests, integration tests, and TDD
6. **Go programming** - Practical Go with real use cases
7. **Algorithms** - Implementing spaced repetition (FSRS)
8. **Code organization** - Package structure and separation of concerns

## Available Quiz Topics

| Topic | Category | Description |
|-------|----------|-------------|
| Java Fundamentals | `java` | Core Java concepts — variables, loops, classes |
| Java Code Tracing | `java` | Read Java snippets and predict printed output |
| Python Basics | `data` | Introduction to Python — syntax, lists, functions |
| Python Code Tracing | `data` | Read Python snippets and predict printed output |
| SQL Database | `all` | Structured Query Language — queries, joins, aggregation |

**Code Tracing topics** are especially useful for building mental models of how programs execute. Questions contain code blocks formatted in Markdown (` ```java` / ` ```python`), rendered as formatted code in the quiz UI.

## Learning Pathways

Choose a pathway based on your goals:

### Pathway 1: "I'm New to Go"

**Goal**: Learn Go by reading working code

**Steps**:
1. Start with `main.go` - The entry point
   - See how configuration is loaded
   - Understand server initialization
   - Learn about graceful shutdown

2. Explore `models/models.go` - Data structures
   - Study Go structs and tags
   - Understand JSON serialization
   - See validation logic

3. Read `api/auth.go` - Simple API handlers
   - Learn HTTP handler patterns
   - Understand request/response flow
   - See error handling

4. Study `database/sqlite.go` - Database operations
   - Learn Go's `database/sql` package
   - Understand prepared statements
   - See connection management

**Exercises**:
- Trace a request from browser → handler → database
- Add a new struct field to `User`
- Write a simple handler that returns JSON
- Query the database from Go code

**Time Estimate**: 2-3 hours

---

### Pathway 2: "I Want to Understand Web APIs"

**Goal**: Learn RESTful API design and implementation

**Steps**:
1. Review `API_EXAMPLES.md` - See the API surface
   - Understand REST conventions
   - Learn endpoint naming
   - Study HTTP methods and status codes

2. Read `router/router.go` - Routing setup
   - See how URLs map to handlers
   - Understand middleware chains
   - Learn route organization

3. Study `api/handlers.go` - Core business logic
   - Trace quiz generation flow
   - See how JSON is parsed/returned
   - Understand validation patterns

4. Explore `api/middleware.go` - Cross-cutting concerns
   - Learn middleware patterns
   - Study authentication flow
   - See rate limiting implementation

**Exercises**:
- Design a new endpoint (e.g., GET /stats)
- Implement the handler
- Add tests for your endpoint
- Document it in API_EXAMPLES.md

**Time Estimate**: 3-4 hours

---

### Pathway 3: "I Need to Learn Databases"

**Goal**: Understand database design, SQL, and migrations

**Steps**:
1. Study `database/migrations/001_initial_schema.up.sql`
   - See table designs
   - Understand foreign keys
   - Learn about indexes

2. Read `database/database.go` - Interface design
   - Understand abstraction layers
   - See method signatures
   - Learn interface patterns in Go

3. Explore `database/sqlite.go` - SQL implementation
   - Study parameterized queries
   - Learn transaction handling
   - See error handling with databases

4. Review `DATABASE_MIGRATIONS.md` - Schema evolution
   - Understand migration strategies
   - Learn about versioning
   - See rollback procedures

**Exercises**:
- Design a new table (e.g., user_badges)
- Write a migration to create it
- Add interface methods for CRUD
- Implement in sqlite.go
- Write tests for your operations

**Time Estimate**: 4-5 hours

---

### Pathway 4: "I Want to Master Testing"

**Goal**: Learn testing strategies and TDD

**Steps**:
1. Start with `fsrs/fsrs_test.go` - Unit tests
   - See how pure functions are tested
   - Understand table-driven tests
   - Learn assertion patterns

2. Study `api/handlers_test.go` - Integration tests
   - See how to test HTTP handlers
   - Learn about test fixtures
   - Understand mocking strategies

3. Read `api/phase3_test.go` - Feature tests
   - See complete workflow testing
   - Understand test data setup
   - Learn about test isolation

4. Explore test helpers in test files
   - Study helper functions
   - See test utilities
   - Learn about setup/teardown

**Exercises**:
- Write a unit test for a new function
- Add an integration test for an endpoint
- Create a test that verifies a complete user workflow
- Measure and improve test coverage

**Time Estimate**: 3-4 hours

---

### Pathway 5: "I Want to Understand Security"

**Goal**: Learn security best practices in web apps

**Steps**:
1. Study `api/auth.go` - Authentication
   - See bcrypt password hashing
   - Understand JWT token generation
   - Learn token validation

2. Read `api/jwt.go` - Token management
   - Study JWT structure
   - See expiration handling
   - Learn about claims

3. Explore `api/middleware.go` - Security middleware
   - See security headers
   - Understand CORS
   - Learn about rate limiting

4. Review `api/admin.go` - Authorization
   - See role-based access control
   - Understand permission checks
   - Learn about user isolation

5. Study `api/topic_access.go` - Category-based access control
   - Understand how student category (`java` vs `data`) gates topic visibility
   - See how `constants/constants.go` defines the category values
   - Trace the path from JWT claims → user record → topic filter

**Exercises**:
- Trace a login request end-to-end
- Implement a "reset password" feature
- Add a new admin-only endpoint
- Write tests for authorization logic
- Explain what happens when a `data` student requests a `java`-only topic

**Time Estimate**: 3-4 hours

---

### Pathway 6: "I Want to Learn Algorithms"

**Goal**: Understand and implement algorithms in real code

**Steps**:
1. Read about FSRS:
   - [FSRS on GitHub](https://github.com/open-spaced-repetition/fsrs4anki/wiki/Free-Spaced-Repetition-Scheduler)
   - Understand the problem it solves
   - Learn the core concepts

2. Study `fsrs/fsrs.go` - Implementation
   - See algorithm in Go code
   - Understand parameter tuning
   - Learn about state machines

3. Trace through `api/fsrs_handlers.go` - Usage
   - See how algorithm integrates with API
   - Understand data flow
   - Learn about state persistence

4. Review `fsrs/fsrs_test.go` - Verification
   - See how algorithm is tested
   - Understand edge cases
   - Learn about validation

**Exercises**:
- Trace a question review through FSRS
- Modify algorithm parameters and test effects
- Implement a simpler spaced repetition algorithm
- Compare results with FSRS

**Time Estimate**: 4-6 hours

---

## Deep Dive Topics

### Understanding Request Flow

Let's trace what happens when you submit a quiz answer:

1. **Browser** → `static/js/quiz.js:submitAnswer()`
   ```javascript
   fetch('/quiz/submit-answer', {
     method: 'POST',
     body: JSON.stringify({ question_id, answer_id })
   })
   ```

2. **Server** → Middleware chain processes request
   - Security headers added
   - CORS check passes
   - Rate limit checked
   - JWT token validated (optional)

3. **Router** → `router/router.go` matches route
   ```go
   r.Post("/quiz/submit-answer", h.SubmitAnswer)
   ```

4. **Handler** → `api/handlers.go:SubmitAnswer()`
   - Parse JSON from request body
   - Validate question/answer IDs
   - Call database to check correctness

5. **Database** → `database/sqlite.go`
   - Execute SQL query
   - Return answer details
   - Update statistics

6. **Response** → Back through handler
   - Build JSON response
   - Set HTTP status code
   - Return to browser

7. **Browser** → Update UI
   - Show correct/incorrect
   - Update score display
   - Enable "Next Question" button

**Exercise**: Add logging at each step to trace a real request.

---

### Understanding Data Flow

How does a quiz get generated?

1. **User clicks "Start Quiz"** on dashboard
2. **Frontend** → Calls `POST /quiz/generate`
3. **Handler** → `api/handlers.go:GenerateAdaptiveQuiz()`
   ```go
   // Get questions for topic
   questions := db.GetQuestions(topicID)
   
   // Check FSRS state for each
   for each question {
     state := db.GetUserQuestionState(userID, questionID)
     prioritize based on due date
   }
   
   // Select top N questions
   selectedQuestions := prioritizeAndSelect(questions, limit)
   
   // Create quiz attempt
   attemptID := db.CreateQuizAttempt(...)
   ```

4. **Database** → Stores attempt, returns ID
5. **Response** → Returns question data to browser
6. **Frontend** → Displays first question

**Exercise**: Modify the prioritization logic and observe effects.

---

### Understanding State Management

How does FSRS track learning progress?

1. **Initial state** when question first seen:
   ```go
   state := UserQuestionState{
     Stability: 0,
     Difficulty: 0,
     State: StateNew,
     Reps: 0,
   }
   ```

2. **After reviewing** with rating 3 (Good):
   ```go
   newState := fsrs.Review(state, rating, now)
   // Stability increases
   // Difficulty adjusts
   // Due date calculated
   // State → StateLearning
   ```

3. **Database updated** with new state:
   ```go
   db.UpdateUserQuestionState(userID, questionID, newState)
   ```

4. **Future quizzes** check due date:
   ```go
   if now > state.Due {
     // Include in quiz (needs review)
   }
   ```

**Exercise**: Track your own learning on a question through multiple reviews.

---

## Code Reading Techniques

### 1. Top-Down Reading

Start with high-level flow, drill into details:

```
main.go
  ↓
router/router.go (what endpoints exist?)
  ↓
api/handlers.go (what does each do?)
  ↓
database/sqlite.go (how is data stored?)
```

### 2. Bottom-Up Reading

Start with data, work up to UI:

```
database schema (what data exists?)
  ↓
models/models.go (how is it represented?)
  ↓
api/handlers.go (how is it manipulated?)
  ↓
static/js/*.js (how is it displayed?)
```

### 3. Feature-Focused Reading

Pick a feature, trace it completely:

```
Quiz Retake:
  → static/js/dashboard.js (button click)
  → api/quiz.go:RetakeQuiz (handler)
  → database/sqlite.go:CreateRetakeAttempt (storage)
  → api/quiz_retake_test.go (tests)

Topic Category Access:
  → constants/constants.go (category constants)
  → api/topic_access.go:userCanAccessTopic (predicate)
  → api/handlers.go (where predicate is applied)
  → api/quiz_category_test.go (tests)
```

### 4. Test-Driven Reading

Read tests to understand behavior:

```
api/phase3_test.go:TestStudentProgress
  ↓ (what's being tested?)
api/instructor.go:GetStudentProgress
  ↓ (how does it work?)
database/sqlite.go:GetStudentProgress
  ↓ (what data is queried?)
```

---

## Common Patterns to Study

### 1. Handler Pattern

All API handlers follow this structure:

```go
func (h *Handler) SomeEndpoint(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req RequestType
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", 400)
        return
    }
    
    // 2. Validate input
    if req.Field == "" {
        http.Error(w, "Field required", 400)
        return
    }
    
    // 3. Business logic / database call
    result, err := h.db.SomeOperation(req.Field)
    if err != nil {
        log.Printf("Error: %v", err)
        http.Error(w, "Internal error", 500)
        return
    }
    
    // 4. Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

**Study**: Find 3 handlers that follow this pattern.

### 2. Interface Pattern

Database abstraction:

```go
// Interface defines contract
type Database interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

// Implementation provides concrete behavior
type SQLiteDB struct { /* ... */ }
func (db *SQLiteDB) GetUser(id int) (*User, error) { /* ... */ }
```

**Study**: Why use interfaces? What are the benefits?

### 3. Middleware Pattern

Composable request processing:

```go
func Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Before request
        doSomething()
        
        // Call next handler
        next(w, r)
        
        // After request
        doSomethingElse()
    }
}
```

**Study**: Trace a request through the middleware chain.

### 4. Table-Driven Test Pattern

Efficient, readable tests:

```go
tests := []struct {
    name     string
    input    int
    expected string
}{
    {"case 1", 1, "one"},
    {"case 2", 2, "two"},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got := function(tt.input)
        if got != tt.expected {
            t.Errorf("got %v, want %v", got, tt.expected)
        }
    })
}
```

**Study**: Find examples in `fsrs_test.go`.

---

## Exercises by Skill Level

### Beginner

1. **Add a field to User model**
   - Add `bio string` to `User` struct
   - Update database schema
   - Display on dashboard

2. **Create a simple endpoint**
   - Add `GET /health` that returns `{"status": "ok"}`
   - Add to router
   - Test in browser

3. **Read and understand a test**
   - Pick any test in `api/handlers_test.go`
   - Understand what it's testing
   - Run it and see it pass

### Intermediate

1. **Implement a new feature**
   - Add "favorite questions" functionality
   - Users can star questions
   - Display favorites list

2. **Add analytics**
   - Create endpoint that returns accuracy by tag
   - Group question attempts by tag
   - Calculate success rates

3. **Improve existing feature**
   - Add question difficulty to quiz results
   - Calculate average difficulty attempted
   - Display on dashboard

### Advanced

1. **Implement complex feature**
   - Add "study groups" where users can share quizzes
   - Design schema
   - Build API
   - Add frontend
   - Write tests

2. **Optimize performance**
   - Profile the application
   - Identify slow queries
   - Add indexes
   - Measure improvement

3. **Add new algorithm**
   - Implement Leitner system (alternative to FSRS)
   - Create side-by-side comparison
   - Analyze which works better for different users

---

## Study Groups and Discussion

### Discussion Questions

1. **Architecture**: Why is the database abstracted behind an interface?
2. **Security**: How would you add two-factor authentication?
3. **Scalability**: What would break first under high load?
4. **Testing**: How would you test the frontend JavaScript?
5. **Design**: Should FSRS parameters be configurable per user?

### Code Review Practice

1. Pick a pull request or feature
2. Review the code changes
3. Ask yourself:
   - Is it correct?
   - Is it secure?
   - Is it tested?
   - Is it maintainable?
   - Are there edge cases?

### Pair Programming Ideas

1. Implement a feature together
2. One person writes test, other writes code
3. Refactor existing code to improve clarity
4. Debug an issue together

---

## Additional Resources

### Go Learning
- [Tour of Go](https://go.dev/tour/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://go.dev/doc/effective_go)

### Web APIs
- [REST API Tutorial](https://restfulapi.net/)
- [HTTP Status Codes](https://httpstatuses.com/)
- [JSON API](https://jsonapi.org/)

### Databases
- [SQLite Tutorial](https://www.sqlitetutorial.net/)
- [Database Design](https://www.databasedesign.info/)
- [SQL Learning](https://sqlbolt.com/)

### Testing
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)
- [Table-Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Testing Best Practices](https://github.com/golang/go/wiki/TestComments)

### Algorithms
- [FSRS Documentation](https://github.com/open-spaced-repetition/fsrs4anki/wiki)
- [Spaced Repetition](https://en.wikipedia.org/wiki/Spaced_repetition)
- [Algorithm Visualization](https://visualgo.net/)

---

## Your Learning Journey

### Week 1: Foundation
- Get app running
- Understand architecture
- Read core files
- Run tests

### Week 2: Deep Dive
- Pick a pathway
- Complete exercises
- Ask questions
- Write code

### Week 3: Build
- Implement new feature
- Write tests
- Get code review
- Iterate

### Week 4: Master
- Refactor existing code
- Optimize performance
- Document learnings
- Help others

---

## Getting Help

- Read `ARCHITECTURE.md` for structural questions
- See `API_EXAMPLES.md` for endpoint details
- Check `DEVELOPMENT_GUIDE.md` for contribution guidelines
- Ask questions in GitHub Discussions
- Join study groups or pair programming sessions

## Success Metrics

You'll know you've learned Athena when you can:

✅ Explain the three-tier architecture
✅ Trace a request from browser to database
✅ Understand the FSRS algorithm
✅ Write a new API endpoint with tests
✅ Design and implement a new feature
✅ Review and understand any file in the codebase
✅ Debug issues using logs and tests
✅ Explain security measures implemented
✅ Teach these concepts to others

Happy learning! 🎓
