# Phase 1 MVP - Quick Start Guide

## Overview

Phase 1 MVP implements the core student experience for Athena, including:
- User authentication (registration and login)
- SQLite3 database for persistent storage
- Quiz generation from topic-based questions
- Interactive quiz-taking with multiple choice questions
- Progress tracking and dashboard
- Bootstrap-based responsive UI

## Prerequisites

- Go 1.16 or higher
- SQLite3 (included via go-sqlite3 driver)
- Modern web browser

## Getting Started

### 1. Build the Application

```bash
go build -o athena .
```

### 2. Seed the Database with Sample Data

```bash
go run cmd/seed/main.go
```

This will create:
- 3 topics: Java Fundamentals, Python Basics, SQL Database
- 11 sample questions across all topics with multiple choice answers
- SQLite database file: `athena.db`

### 3. Start the Server

```bash
./athena
```

The server will start on `http://localhost:9206`

Alternatively, you can specify a custom database path:
```bash
DB_PATH=/path/to/custom.db ./athena
```

### 4. Access the Application

Open your browser and navigate to:
```
http://localhost:9206/index.html
```

## Using the Application

### Registration
1. On the home page, fill out the registration form with:
   - Username
   - Email
   - Password
2. Click "Register"
3. Upon success, you can log in with your credentials

### Login
1. Enter your username and password
2. Click "Login"
3. You'll be redirected to the dashboard

### Dashboard
The dashboard shows:
- **Total Points**: Your cumulative score
- **Quizzes Taken**: Number of completed quizzes
- **Questions Answered**: Total questions you've attempted
- **Correct Answers**: Total correct answers
- **Available Topics**: Click on any topic to start a quiz
- **Recent Quiz Attempts**: History of your quiz sessions

### Taking a Quiz
1. Click "Start Quiz" on any topic card
2. Answer each multiple choice question
3. Click "Submit Answer" to check your response
4. See immediate feedback (correct/incorrect)
5. Click "Next Question" to continue
6. View your final score at the end

## API Endpoints

### Authentication
- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login existing user

### Topics
- `GET /topics` - List all topics
- `POST /topics` - Create a new topic (admin)

### Questions
- `GET /questions?topic_id={id}` - Get questions for a topic
- `POST /questions` - Create a new question (admin)

### Quiz
- `POST /quiz/generate` - Generate a new quiz
- `POST /quiz/submit-answer` - Submit an answer
- `POST /quiz/complete` - Complete a quiz
- `GET /progress?user_id={id}` - Get user progress

### Users
- `GET /users?id={id}` - Get user by ID
- `POST /users` - Create user

## Running Tests

```bash
go test -v ./...
```

All 20 existing tests should pass.

## Database Schema

The SQLite database includes the following tables:
- `users` - User accounts with authentication
- `topics` - Quiz topics
- `questions` - Quiz questions with difficulty and tags
- `answers` - Answer options for questions
- `quiz_attempts` - User quiz sessions
- `question_attempts` - Individual question attempt history

## Features Implemented

✅ User authentication with password hashing
✅ SQLite3 persistent database
✅ Basic quiz generation (random selection)
✅ Multiple choice question answering
✅ Progress tracking per user
✅ Student dashboard with statistics
✅ Bootstrap-based responsive UI
✅ Clean and logical interface
✅ Topic-based quiz organization

## What's Next (Phase 2)

- FSRS (Free Spaced Repetition Scheduler) algorithm
- Adaptive quiz generation based on performance
- Advanced progress analytics
- Tag-based insights
- Instructor dashboard

## Troubleshooting

**Database locked error**: Make sure only one instance of the server is running

**Port already in use**: Change the port in `main.go` or kill the existing process

**CSS not loading**: Ensure you have internet access for Bootstrap CDN, or download Bootstrap locally

## File Structure

```
zip-athena/
├── main.go                 # Application entry point
├── models/                 # Data models
├── database/              # Database implementations
│   ├── database.go        # Database interface
│   ├── sqlite.go          # SQLite implementation
├── api/                   # API handlers
│   ├── handlers.go        # Core handlers
│   ├── auth.go           # Authentication handlers
│   └── quiz.go           # Quiz-related handlers
├── cmd/seed/             # Database seeding utility
│   └── main.go
├── static/               # Frontend files
│   ├── index.html        # Login/Registration page
│   ├── dashboard.html    # Student dashboard
│   ├── quiz.html        # Quiz taking interface
│   ├── css/
│   │   └── style.css
│   └── js/
│       ├── auth.js
│       ├── dashboard.js
│       └── quiz.js
└── athena.db            # SQLite database (generated)
```
