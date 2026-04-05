# Getting Started with Athena

## What is Athena?

Athena is a web-based student quiz application that helps you learn programming through adaptive, intelligent quizzing. It uses a scientifically-proven spaced repetition algorithm to optimize your learning and retention.

**Key Features**:
- 📚 Adaptive quizzing with FSRS spaced repetition
- 🎯 Topics in Java, Python, SQL, and Code Tracing
- 📊 Track your progress and see detailed analytics
- ⏱️ Timed tests for assessment
- 🔄 Retake quizzes to improve understanding
- 📖 Learning resources linked to questions

## Student Categories

Every student account is assigned a **category** that controls which quiz topics they can see:

| Category | Topics Visible |
|----------|----------------|
| `java`   | Java Fundamentals, Java Code Tracing, SQL Database |
| `data`   | Python Basics, Python Code Tracing, SQL Database |

Admins can see all topics regardless of category. Contact your instructor to have your account assigned the correct category.

## Prerequisites

Before you begin, make sure you have:

- **Go 1.16 or higher** - [Download here](https://go.dev/dl/)
- **Git** - [Download here](https://git-scm.com/downloads)
- **A modern web browser** (Chrome, Firefox, Safari, Edge)
- **Basic command line familiarity**

You don't need Node.js, npm, or any other build tools!

## Quick Start (5 Minutes)

### 1. Clone the Repository

```bash
git clone https://github.com/ZipCodeCore/zip-athena.git
cd zip-athena
```

### 2. Seed the Database

This creates sample questions, topics, and an admin user:

```bash
go run cmd/seed/main.go
```

**What this creates**:
- 100+ questions across Java, Python, SQL, and Code Tracing topics
- Admin user (username: `admin`, password: `admin`)
- Java student (username: `user`, password: `user`, category: `java`)
- Data/Python student (username: `data`, password: `data`, category: `data`)
- Sample topics, quizzes, and quiz history
- Database file: `athena.db`

### 3. Start the Server

```bash
go run main.go
```

You should see:
```
INFO[0000] Starting Athena API server
INFO[0000] Database initialized successfully
INFO[0000] Athena API server listening on :9206
```

### 4. Open in Your Browser

Navigate to: **http://localhost:9206/index.html**

### 5. Try It Out!

> **Note**: Self-registration is currently disabled. Use one of the pre-seeded accounts below or contact your instructor to have an account created for you.

**Option A: Login as a Java Student**
1. Username: `user`
2. Password: `user`
3. You will see Java and SQL topics on your dashboard

**Option B: Login as a Data/Python Student**
1. Username: `data`
2. Password: `data`
3. You will see Python and SQL topics on your dashboard

**Option C: Login as Admin**
1. Username: `admin`
2. Password: `admin`
3. Go to http://localhost:9206/admin.html to manage content

## Understanding Your First Quiz

When you start a quiz, Athena:

1. **Selects questions** based on:
   - Questions you need to review (due for spaced repetition)
   - New questions you haven't seen
   - Your selected topic

2. **Shows one question at a time** with:
   - Multiple choice answers
   - Immediate feedback when you submit
   - Explanation (when available)

3. **Tracks your performance**:
   - Correct/incorrect count
   - Overall score
   - Updates your FSRS schedule for future reviews

4. **Adapts to you**:
   - Questions you struggle with appear more frequently
   - Questions you master appear less often
   - Optimal spacing for long-term retention

## Key Concepts

### Spaced Repetition (FSRS)

Athena uses FSRS-4.5, a research-backed algorithm that schedules reviews at optimal intervals:

- **Just answered correctly?** → See it again in a few days
- **Struggled with it?** → See it again soon
- **Mastered it?** → See it again in weeks/months

This is more effective than cramming or random review.

### Quiz Retakes

Got less than 100% on a quiz? You can retake it!

- Same questions, fresh attempt
- Helps reinforce what you missed
- Updates your FSRS schedule
- Tracks improvement over time

### Topics and Tags

Questions are organized by:
- **Topics**: Broad areas such as "Java Fundamentals", "Python Basics", "SQL Database", "Java Code Tracing", and "Python Code Tracing"
- **Tags**: Specific concepts like "loops", "variables", "inheritance"

Which topics you see depends on your account's **category** (`java` or `data`). Use tags to focus on specific areas you want to improve.

## Common Tasks

### Taking a Quiz

1. **Dashboard**: Click on a topic card
2. **Answer**: Select your answer and click "Submit Answer"
3. **Review**: See if you were correct
4. **Continue**: Click "Next Question"
5. **Complete**: Finish all questions to see your score

### Viewing Your Progress

1. Go to your **Dashboard** (http://localhost:9206/dashboard.html)
2. See your stats:
   - Total points earned
   - Quizzes taken
   - Accuracy percentage
   - Recent quiz history

### Reviewing Learning Resources

1. Go to **Learning Resources** (http://localhost:9206/learning.html)
2. Browse topics organized by language (Java, Python)
3. Click to read detailed guides on concepts
4. Resources are linked to related questions

### Taking a Timed Test

1. Go to **Tests** from your dashboard
2. Click "Start Test" on an available test
3. Timer counts down as you work
4. Submit before time runs out
5. View detailed results after submission

## Admin Features

Login as admin to access:

### Managing Questions

1. Go to http://localhost:9206/admin.html
2. Click "Questions" tab
3. Create, edit, or delete questions
4. Bulk import via CSV

### Managing Users

1. "Users" tab in admin panel
2. View all registered users
3. See user statistics
4. Delete users if needed

### Creating Quizzes

1. "Quizzes" tab
2. Create named quiz collections
3. Associate with topics
4. Students can access via dashboard

## Instructor Features

Login as an instructor to:

### View Student Progress

1. Go to http://localhost:9206/instructor.html
2. See list of all students
3. Click a student to see:
   - Quiz history
   - Accuracy by topic
   - Struggling areas
   - Recent activity

### Cohort Analytics

1. "Analytics" tab
2. See class-wide statistics
3. Identify struggling students
4. Track improvement over time

### Import/Export Content

1. Export questions as JSON
2. Share quiz banks
3. Import from other instructors
4. Backup and restore content

## Troubleshooting

### Port Already in Use

If you see "address already in use":

```bash
# Find and kill the process on port 9206
lsof -ti:9206 | xargs kill -9

# Or use a different port
PORT=8080 go run main.go
```

### Database Locked

Only one server instance can access `athena.db` at a time:

```bash
# Stop all running instances
pkill -f "go run main.go"

# Restart
go run main.go
```

### CSS Not Loading

Make sure you have internet access for Bootstrap CDN, or:

1. Download Bootstrap locally
2. Update HTML files to reference local files

### Can't Login

If you forgot your password:

```bash
# Re-run seed to reset admin password
go run cmd/seed/main.go
```

This recreates the admin user with password `admin`.

## Next Steps

Now that you're up and running:

1. **Explore the Code**: See [ARCHITECTURE.md](ARCHITECTURE.md) to understand how it's built
2. **Learn from It**: See [LEARNING_GUIDE.md](LEARNING_GUIDE.md) for educational pathways
3. **Contribute**: See [DEVELOPMENT_GUIDE.md](DEVELOPMENT_GUIDE.md) to add features
4. **Read the API**: See [API_EXAMPLES.md](API_EXAMPLES.md) for endpoint documentation

## Configuration

Customize Athena with environment variables:

```bash
# Database
export DB_PATH=/path/to/custom.db

# Server
export PORT=8080
export READ_TIMEOUT=30s
export WRITE_TIMEOUT=30s

# Security
export JWT_SECRET=your-secret-key

# Logging
export LOG_LEVEL=debug

# Start server
go run main.go
```

See `config/config.go` for all available options.

## Development Workflow

### Making Changes

```bash
# 1. Edit code
vim api/handlers.go

# 2. Run tests
go test -v ./api

# 3. Restart server
# Press Ctrl+C in the running server terminal
go run main.go

# 4. Test in browser
# Refresh browser and verify changes
```

### Running Tests

```bash
# All tests
go test -v ./...

# Specific package
go test -v ./api
go test -v ./fsrs
go test -v ./database

# With coverage
go test -cover ./...
```

### Database Migrations

If you modify the schema:

```bash
# Build migration tool
go build -o migrate ./cmd/migrate

# Check current version
./migrate -db athena.db -cmd version

# Apply pending migrations
./migrate -db athena.db -cmd up

# Rollback last migration
./migrate -db athena.db -cmd down -steps 1
```

See [DATABASE_MIGRATIONS.md](DATABASE_MIGRATIONS.md) for details.

## Getting Help

- **Architecture questions?** → See `ARCHITECTURE.md`
- **API questions?** → See `API_EXAMPLES.md`
- **Learning pathways?** → See `LEARNING_GUIDE.md`
- **Contributing?** → See `DEVELOPMENT_GUIDE.md`
- **Bugs/Issues?** → [Open an issue](https://github.com/ZipCodeCore/zip-athena/issues)

## Resources

- **FSRS Algorithm**: [Open Spaced Repetition](https://github.com/open-spaced-repetition)
- **Go Documentation**: [go.dev](https://go.dev/doc/)
- **Chi Router**: [go-chi](https://go-chi.io/)
- **Bootstrap**: [getbootstrap.com](https://getbootstrap.com/)

## Success Checklist

You're ready to go when you can:

- ✅ Start the server without errors
- ✅ Register a new user account
- ✅ Login successfully
- ✅ Generate and complete a quiz
- ✅ See your progress on the dashboard
- ✅ Access admin panel (with admin credentials)
- ✅ Run tests successfully (`go test ./...`)

Welcome to Athena! Happy learning! 🚀
