# Development Guide

## Introduction

This guide helps you contribute to Athena, whether you're fixing bugs, adding features, or improving documentation. We welcome contributions from developers of all skill levels!

## Development Setup

### Prerequisites

- **Go 1.16+**: [Download](https://go.dev/dl/)
- **Git**: [Download](https://git-scm.com/)
- **SQLite3**: Usually pre-installed on macOS/Linux
- **Code Editor**: VS Code, GoLand, or your favorite editor

### Recommended VS Code Extensions

- Go (golang.go)
- SQLite Viewer
- REST Client
- GitLens

### Initial Setup

```bash
# 1. Fork the repository on GitHub

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/zip-athena.git
cd zip-athena

# 3. Add upstream remote
git remote add upstream https://github.com/ZipCodeCore/zip-athena.git

# 4. Install Go dependencies
go mod download

# 5. Seed the database
go run cmd/seed/main.go

# 6. Run tests to verify setup
go test -v ./...

# 7. Start the server
go run main.go
```

## Project Structure

```
zip-athena/
├── api/                  # HTTP handlers and middleware
│   ├── handlers.go       # Core quiz logic
│   ├── auth.go          # Authentication
│   ├── admin.go         # Admin endpoints
│   ├── instructor.go    # Instructor analytics
│   ├── fsrs_handlers.go # FSRS endpoints
│   ├── phase4.go        # Advanced student features
│   ├── phase5.go        # Content management
│   ├── middleware.go    # Request middleware
│   └── *_test.go       # API tests
│
├── config/              # Configuration management
│   └── config.go       # Environment-based config
│
├── database/            # Data layer
│   ├── database.go     # Interface definition
│   ├── sqlite.go       # SQLite implementation
│   └── migrations/     # Schema migrations
│
├── fsrs/               # FSRS algorithm
│   ├── fsrs.go        # Algorithm implementation
│   └── fsrs_test.go   # Algorithm tests
│
├── models/             # Data models
│   ├── models.go      # All model definitions
│   └── models_test.go # Model tests
│
├── router/             # HTTP routing
│   └── router.go      # Route registration
│
├── logger/             # Logging utilities
│   └── logger.go
│
├── static/             # Frontend files
│   ├── *.html         # Pages
│   ├── js/            # JavaScript
│   └── css/           # Stylesheets
│
├── cmd/                # Command-line tools
│   ├── seed/          # Database seeding
│   └── migrate/       # Migration management
│
└── main.go            # Application entry point
```

## Development Workflow

### 1. Planning Your Change

Before writing code:

1. **Check existing issues** - Someone may be working on it
2. **Open an issue** - Describe what you want to do
3. **Discuss approach** - Get feedback before investing time
4. **Get approval** - Especially for large changes

### 2. Creating a Branch

```bash
# Update main branch
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/bug-description
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions

### 3. Making Changes

#### Code Style

Follow Go conventions:
- Run `gofmt` on all code
- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small

```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Run linter (if installed)
golangci-lint run
```

#### Commit Messages

Write clear, descriptive commit messages:

```
Add quiz retake feature

- Students can now retake quizzes with scores < 100%
- Retakes use same questions as original attempt
- FSRS schedules are updated based on incorrect answers
- Added 13 new tests for retake functionality

Closes #123
```

Format:
- First line: Brief summary (50 chars or less)
- Blank line
- Detailed description (what, why, how)
- Reference issues/PRs

### 4. Testing Your Changes

**Always write tests** for new code:

```bash
# Run all tests
go test -v ./...

# Run specific package
go test -v ./api

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./api -run TestQuizGeneration
```

#### Test Requirements

1. **Unit tests** for new functions
2. **Integration tests** for API endpoints
3. **Edge case tests** for error handling
4. **Performance tests** for algorithms

#### Writing Good Tests

```go
func TestFeatureName(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()
    
    // Create test data
    user := createTestUser(db)
    
    // Execute
    result, err := function(user.ID)
    
    // Assert
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if result.Field != expected {
        t.Errorf("Expected %v, got %v", expected, result.Field)
    }
}
```

### 5. Running the Application

```bash
# Development mode with auto-reload
air  # If you have air installed

# Or manual restart
go run main.go

# Build and run
go build -o athena .
./athena

# With environment variables
PORT=8080 LOG_LEVEL=debug go run main.go
```

### 6. Submitting a Pull Request

```bash
# 1. Ensure tests pass
go test -v ./...

# 2. Commit your changes
git add .
git commit -m "Your descriptive commit message"

# 3. Push to your fork
git push origin feature/your-feature-name

# 4. Create pull request on GitHub
# - Describe your changes
# - Reference related issues
# - Add screenshots for UI changes
# - Request review from maintainers
```

#### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] All tests pass
- [ ] Added new tests
- [ ] Tested manually

## Checklist
- [ ] Code follows project style
- [ ] Self-reviewed code
- [ ] Commented complex code
- [ ] Updated documentation
- [ ] No new warnings
```

## Common Development Tasks

### Adding a New API Endpoint

1. **Define the handler** in appropriate file (e.g., `api/handlers.go`):

```go
func (h *Handler) NewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Parse request
    var req RequestType
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // Business logic
    result, err := h.db.SomeOperation(req.Field)
    if err != nil {
        log.Printf("Error in NewEndpoint: %v", err)
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    
    // Return response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

2. **Add route** in `router/router.go`:

```go
r.Post("/your/endpoint", h.NewEndpoint)
```

3. **Write tests** in `api/*_test.go`:

```go
func TestNewEndpoint(t *testing.T) {
    // Test implementation
}
```

4. **Document** in `API_EXAMPLES.md`

### Adding a Database Table

1. **Create migration files**:

```bash
# Create new migration
touch database/migrations/004_add_new_table.up.sql
touch database/migrations/004_add_new_table.down.sql
```

2. **Write up migration** (`004_add_new_table.up.sql`):

```sql
CREATE TABLE IF NOT EXISTS new_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    field TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_new_table_user_id ON new_table(user_id);
```

3. **Write down migration** (`004_add_new_table.down.sql`):

```sql
DROP INDEX IF EXISTS idx_new_table_user_id;
DROP TABLE IF EXISTS new_table;
```

4. **Add model** to `models/models.go`:

```go
type NewTable struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    Field     string    `json:"field"`
    CreatedAt time.Time `json:"created_at"`
}
```

5. **Add interface methods** to `database/database.go`:

```go
type Database interface {
    // ... existing methods
    CreateNewTableEntry(*NewTable) error
    GetNewTableEntry(id int) (*NewTable, error)
}
```

6. **Implement in SQLite** (`database/sqlite.go`):

```go
func (db *SQLiteDB) CreateNewTableEntry(entry *NewTable) error {
    query := `INSERT INTO new_table (user_id, field) VALUES (?, ?)`
    result, err := db.db.Exec(query, entry.UserID, entry.Field)
    if err != nil {
        return err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    entry.ID = int(id)
    return nil
}
```

7. **Test migrations**:

```bash
./migrate -db test.db -cmd up
./migrate -db test.db -cmd down
```

### Adding Frontend Features

1. **Update HTML** in `static/*.html`:

```html
<div id="new-feature">
    <button onclick="doSomething()">Click Me</button>
</div>
```

2. **Add JavaScript** in `static/js/*.js`:

```javascript
async function doSomething() {
    const response = await fetch('/api/endpoint', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('token')}`
        },
        body: JSON.stringify({ data: 'value' })
    });
    
    const result = await response.json();
    updateUI(result);
}
```

3. **Test in browser** across different screen sizes

### Debugging

#### Server Logs

```go
// Add logging in handlers
log.Printf("Debug: processing request for user %d", userID)

// Run with debug level
LOG_LEVEL=debug go run main.go
```

#### Database Inspection

```bash
# Open database
sqlite3 athena.db

# List tables
.tables

# Query data
SELECT * FROM users;

# Check schema
.schema users

# Exit
.quit
```

#### HTTP Debugging

Use curl or a REST client:

```bash
# Test endpoint
curl -X POST http://localhost:9206/quiz/generate \
  -H "Content-Type: application/json" \
  -d '{"topic_id": 1, "limit": 5}'

# With auth token
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9206/admin/users
```

## Code Quality

### Running Linters

```bash
# go vet (built-in)
go vet ./...

# golangci-lint (recommended)
golangci-lint run

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Code Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in browser
go tool cover -html=coverage.out

# Check coverage percentage
go test -cover ./...
```

### Performance Profiling

```bash
# CPU profiling
go test -cpuprofile cpu.prof -bench .

# Memory profiling
go test -memprofile mem.prof -bench .

# View profile
go tool pprof cpu.prof
```

## Database Management

### Running Migrations

```bash
# Build migration tool
go build -o migrate ./cmd/migrate

# Check version
./migrate -db athena.db -cmd version

# Apply all pending
./migrate -db athena.db -cmd up

# Rollback one
./migrate -db athena.db -cmd down -steps 1

# Rollback all
./migrate -db athena.db -cmd down
```

### Seeding Data

```bash
# Full seed (drops existing data)
go run cmd/seed/main.go

# Custom seeding
# Edit cmd/seed/main.go to add your test data
```

### Backing Up Database

```bash
# Backup
cp athena.db athena.db.backup

# Restore
cp athena.db.backup athena.db

# Export to SQL
sqlite3 athena.db .dump > backup.sql

# Import from SQL
sqlite3 athena.db < backup.sql
```

## Common Issues

### Port Already in Use

```bash
# Find process
lsof -i :9206

# Kill it
kill -9 <PID>
```

### Database Locked

```bash
# Only one process can write at a time
# Stop all running instances
pkill -f athena

# Restart
go run main.go
```

### Tests Failing

```bash
# Clean test cache
go clean -testcache

# Run with verbose output
go test -v ./...

# Run single test
go test -v -run TestName ./package
```

### Import Issues

```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify modules
go mod verify
```

## Best Practices

### General Guidelines

1. **Keep it simple** - Prefer simple, readable code over clever code
2. **One thing well** - Functions should do one thing
3. **Error handling** - Always handle errors, never ignore them
4. **Testing** - Write tests before fixing bugs
5. **Documentation** - Comment exported functions and complex logic

### Go-Specific

1. **Use gofmt** - Always format code
2. **Handle errors** - Check every error
3. **Use interfaces** - For abstraction and testing
4. **Prefer composition** - Over inheritance
5. **Use context** - For cancellation and timeouts

### API Design

1. **REST conventions** - Use proper HTTP methods and status codes
2. **Consistent naming** - Follow established patterns
3. **Validation** - Validate all inputs
4. **Error responses** - Provide helpful error messages
5. **Documentation** - Document all endpoints

### Database

1. **Use transactions** - For multi-step operations
2. **Parameterized queries** - Prevent SQL injection
3. **Indexes** - Add for frequently queried columns
4. **Migrations** - Never modify data directly
5. **Backups** - Backup before schema changes

### Security

1. **Never commit secrets** - Use environment variables
2. **Hash passwords** - Use bcrypt
3. **Validate input** - Never trust user input
4. **Rate limiting** - Prevent abuse
5. **Least privilege** - Only grant necessary permissions

## Release Process

### Versioning

We use semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Creating a Release

1. Update version in code
2. Update CHANGELOG.md
3. Run all tests
4. Tag release: `git tag v1.2.3`
5. Push tag: `git push origin v1.2.3`
6. Create GitHub release with notes

## Getting Help

### Documentation

- `ARCHITECTURE.md` - System architecture
- `API_EXAMPLES.md` - API reference
- `DATABASE_MIGRATIONS.md` - Migration guide
- `LEARNING_GUIDE.md` - Learn from the codebase

### Community

- GitHub Issues - Bug reports and feature requests
- GitHub Discussions - Questions and ideas
- Pull Requests - Code review and feedback

### Maintainers

Tag maintainers in issues/PRs for:
- Architecture decisions
- Breaking changes
- New features
- Security concerns

## Contributing Guidelines

### Code of Conduct

- Be respectful and inclusive
- Assume good intentions
- Provide constructive feedback
- Focus on ideas, not individuals

### Review Process

1. **Submit PR** - With clear description
2. **Automated checks** - Tests must pass
3. **Code review** - Maintainer reviews code
4. **Address feedback** - Make requested changes
5. **Approval** - Maintainer approves
6. **Merge** - Maintainer merges PR

### Review Criteria

- Code quality and style
- Test coverage
- Documentation
- Performance impact
- Security implications
- Backward compatibility

## Acknowledgments

Thank you for contributing to Athena! Your efforts help make this a better learning platform for students worldwide.

## Resources

- [Go Documentation](https://go.dev/doc/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
- [Semantic Versioning](https://semver.org/)

Happy coding! 🚀
