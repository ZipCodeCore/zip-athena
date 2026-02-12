# Athena - Adaptive Learning Platform

Athena is an intelligent quiz application that helps students learn programming through adaptive, spaced repetition. It's also a complete, production-quality codebase designed to teach modern software development practices.

**Dual Purpose**:
1. 📚 **Learning Tool** - Use it to learn Java, Python, and SQL through intelligent quizzing
2. 🎓 **Teaching Resource** - Study the codebase to learn three-tier architecture, RESTful APIs, databases, and Go programming

## Quick Start

Get up and running in 5 minutes:

```bash
# 1. Clone the repository
git clone https://github.com/ZipCodeCore/zip-athena.git
cd zip-athena

# 2. Seed the database with sample questions
go run cmd/seed/main.go

# 3. Start the server
go run main.go

# 4. Open browser to http://localhost:9206
```

## What Makes Athena Special?

### Intelligent Learning with FSRS
- Uses the **FSRS-4.5 spaced repetition algorithm** (research-backed, superior to traditional Anki/SM-2)
- Questions appear at optimal intervals for maximum retention
- Adapts to your performance - struggles more, masters less

### Student-Focused Features
- Quiz retakes to improve understanding
- Detailed progress tracking and analytics
- Timed tests for assessment
- Learning resources linked to questions
- Tag-based review for focused practice

### Built for Learning Software Development
- Clean three-tier architecture (Frontend → API → Database)
- Well-tested (100+ tests)
- Production-ready code patterns
- Comprehensive documentation
- Real-world security practices

## Documentation

### For Students & Learners

- **[Getting Started Guide](GETTING_STARTED.md)** - Set up and run Athena in 5 minutes
- **[Learning Guide](LEARNING_GUIDE.md)** - Use this codebase to learn software development
  - 6 guided learning pathways (Go, APIs, Databases, Testing, Security, Algorithms)
  - Code reading techniques
  - Exercises by skill level
  - Study group discussion questions

### For Developers & Contributors

- **[Architecture Guide](ARCHITECTURE.md)** - Understand how Athena is built
  - High-level architecture overview
  - Component responsibilities
  - Core algorithms (FSRS, adaptive quiz generation)
  - Design decisions and rationale
  
- **[Development Guide](DEVELOPMENT_GUIDE.md)** - Contribute to the codebase
  - Development workflow
  - Adding features/endpoints/tables
  - Testing guidelines
  - Code quality standards

### Technical References

- **[API Examples](docs/development/API_EXAMPLES.md)** - Complete API documentation
- **[Database Migrations](docs/development/DATABASE_MIGRATIONS.md)** - Schema management
- **[Import/Export Format](docs/development/IMPORT_EXPORT_FORMAT.md)** - Data interchange
- **[Test Summary](docs/development/TEST_SUMMARY.md)** - Testing documentation

### Additional Resources

- **[Specification](docs/SPECIFICATION.md)** - Complete system specification
- **[Curriculum](docs/CURRICULUM.md)** - Learning content organization
- **[Development History](docs/archive/)** - Phase implementation details

## Architecture Overview

```
┌─────────────────┐
│   Web Browser   │  ← Vanilla JS + Bootstrap
│  (Static HTML)  │    No build step required
└────────┬────────┘
         │ HTTP/JSON
┌────────▼────────┐
│   Go API Server │  ← Chi router + Middleware
│   (RESTful API) │    JWT auth, CORS, rate limiting
└────────┬────────┘
         │ SQL
┌────────▼────────┐
│ SQLite Database │  ← Migrations, indexes
│   (athena.db)   │    Thread-safe, ACID-compliant
└─────────────────┘
```

## Using Athena

### As a Student

1. **Register/Login** at http://localhost:9206
2. **Browse Topics** on your dashboard (Java, Python, SQL)
3. **Take Quizzes** - Answer questions, get immediate feedback
4. **Track Progress** - See your accuracy, points, and improvement
5. **Review Material** - Access learning resources
6. **Retake Quizzes** - Improve your understanding

### As an Instructor

1. **Login** as instructor
2. **View Student Progress** at http://localhost:9206/instructor.html
3. **Monitor Class Analytics** - See cohort-wide statistics
4. **Identify Struggling Students** - Focus help where needed
5. **Import/Export Content** - Share quiz banks

### As an Admin

1. **Login** with admin credentials (username: `admin`, password: `admin`)
2. **Manage Content** at http://localhost:9206/admin.html
3. **Create/Edit Questions** - Build question banks
4. **Bulk Import** - Upload questions via CSV
5. **Review Flagged Questions** - Quality control

## For Developers

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test -v ./api
go test -v ./fsrs
```

**Test Coverage**: 100+ tests across all layers (database, API, FSRS algorithm)

### Common Development Tasks

```bash
# Format code
go fmt ./...

# Check for issues
go vet ./...

# Build application
go build -o athena .

# Run with custom config
PORT=8080 LOG_LEVEL=debug go run main.go

# Manage database migrations
go build -o migrate ./cmd/migrate
./migrate -db athena.db -cmd version
```

See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) for complete development documentation.

## Key Features

### For Students
- 📚 **Adaptive Quizzing** - FSRS algorithm schedules optimal review intervals
- 🔄 **Quiz Retakes** - Retake quizzes to improve understanding
- 📊 **Progress Tracking** - Detailed statistics and history
- ⏱️ **Timed Tests** - Formal assessments with results
- 📖 **Learning Resources** - Curriculum guides linked to questions
- 🏷️ **Tag-Based Review** - Focus on specific concepts

### For Instructors
- 👥 **Student Analytics** - Track individual student progress
- 📈 **Cohort Analytics** - Class-wide performance statistics
- 🎯 **Identify Struggles** - Find students who need help
- 📦 **Import/Export** - Share quiz banks and questions

### For Admins
- ✏️ **Question Management** - Create, edit, delete questions
- 📥 **Bulk CSV Import** - Upload question banks efficiently
- 🚩 **Quality Review** - Flag and review problematic questions
- 🏷️ **Tag Management** - Organize and analyze by tags
- 👤 **User Management** - Full user CRUD operations

## API Endpoints

Athena provides a comprehensive RESTful API. See [API Examples](docs/development/API_EXAMPLES.md) for complete documentation.

### Example Endpoints

**Authentication**:
- `POST /auth/register` - Register new user
- `POST /auth/login` - Login and get JWT token

**Student Quizzing**:
- `POST /quiz/generate` - Generate adaptive quiz
- `POST /quiz/submit-answer` - Submit an answer
- `POST /quiz/complete` - Complete quiz
- `POST /quiz/retake` - Retake a quiz
- `GET /progress?user_id={id}` - Get progress stats

**FSRS Spaced Repetition**:
- `POST /fsrs/review` - Review question with rating (1-4)
- `GET /fsrs/due?user_id={id}` - Get due questions

**Admin** (requires JWT):
- `GET /admin/users` - List all users
- `POST /admin/questions` - Create question
- `POST /admin/questions/import` - Bulk CSV import

**Instructor** (requires JWT):
- `GET /instructor/student-progress?user_id={id}` - Student details
- `GET /instructor/cohort-analytics` - Class statistics

## Learning Pathways

Want to learn from this codebase? Choose a pathway in the [Learning Guide](LEARNING_GUIDE.md):

1. **"I'm New to Go"** - Learn Go by reading working code (2-3 hours)
2. **"I Want to Understand Web APIs"** - RESTful API design (3-4 hours)
3. **"I Need to Learn Databases"** - Schema design, SQL, migrations (4-5 hours)
4. **"I Want to Master Testing"** - Unit, integration, feature tests (3-4 hours)
5. **"I Want to Understand Security"** - Auth, authorization, best practices (3-4 hours)
6. **"I Want to Learn Algorithms"** - FSRS spaced repetition (4-6 hours)

Each pathway includes:
- Guided steps through the codebase
- Hands-on exercises
- Time estimates
- Learning objectives

## Contributing

We welcome contributions! See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) for:

- Development setup
- Code style guidelines
- Testing requirements
- Pull request process
- Common development tasks

**Quick Start for Contributors**:
```bash
# Fork the repository
# Clone your fork
git clone https://github.com/YOUR_USERNAME/zip-athena.git

# Create a feature branch
git checkout -b feature/your-feature

# Make changes and test
go test -v ./...

# Submit pull request
```

## Project Structure

```
zip-athena/
├── api/              # HTTP handlers and middleware
├── config/           # Configuration management
├── database/         # Data layer with SQLite
│   └── migrations/   # Schema migrations
├── fsrs/             # FSRS algorithm
├── models/           # Data models
├── router/           # HTTP routing
├── static/           # Frontend (HTML/JS/CSS)
├── cmd/              # CLI tools (seed, migrate)
├── docs/             # Documentation
│   ├── development/  # Technical references
│   └── archive/      # Development history
└── main.go          # Application entry point
```

## Technology Stack

- **Language**: Go 1.x - Fast, concurrent, easy deployment
- **Router**: Chi - Lightweight, composable HTTP router
- **Database**: SQLite3 - Zero-config, file-based, ACID-compliant
- **Frontend**: Vanilla JS + Bootstrap 5 - No build step needed
- **Security**: bcrypt password hashing, JWT tokens
- **Algorithm**: FSRS-4.5 - Research-backed spaced repetition

## License

See [LICENSE](LICENSE) file for details.

## Acknowledgments

Built with modern software development practices to serve as both a functional learning platform and an educational resource for understanding production-quality codebases.
