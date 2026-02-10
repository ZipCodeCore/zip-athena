package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB implements the Database interface using SQLite3
type SQLiteDB struct {
	db *sql.DB
}

// DB returns the underlying sql.DB connection
func (db *SQLiteDB) DB() *sql.DB {
	return db.db
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB := &SQLiteDB{db: db}
	if err := sqlDB.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return sqlDB, nil
}

// Close closes the database connection
func (db *SQLiteDB) Close() error {
	return db.db.Close()
}

// initSchema creates the database tables using migrations
func (db *SQLiteDB) initSchema() error {
	// Find migrations directory relative to this file or working directory
	migrationsPath := findMigrationsPath()
	return db.RunMigrations(migrationsPath)
}

// findMigrationsPath locates the migrations directory
func findMigrationsPath() string {
	// Try relative to working directory first
	paths := []string{
		"database/migrations",
		"../database/migrations",
		"../../database/migrations",
	}
	
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			absPath, err := filepath.Abs(path)
			if err != nil {
				// If we can't get absolute path, use the relative path
				return path
			}
			return absPath
		}
	}
	
	// Default to relative path
	return "database/migrations"
}

// User operations
// Topic operations
// Question operations
// Answer operations
// QuizAttempt operations
// QuestionAttempt operations
// Quiz operations
// UserQuestionState operations
// Admin operations - Users
// Admin operations - Questions
// Admin operations - Answers
// Admin operations - Topics
// GetAnswer retrieves an answer by ID
// Phase 3: Test operations
// Phase 3: Many-to-many Quiz-Question relationships
// Phase 3: Many-to-many Test-Question relationships
// Phase 3: Analytics operations
// ListStudentsWithMetrics returns a list of all students with their performance metrics
// GetStudentQuizHistory returns all quiz attempts for a student
// GetStudentWeakAreas identifies topics and tags where student is struggling (< 60% accuracy)
// ReorderTopics updates the display_order for topics
// Phase 4: SQLite implementations to append to sqlite.go

// Phase 4: Quiz History Methods

// GetQuizHistoryByUser retrieves quiz attempts for a user with pagination
// GetQuizHistoryByTopic retrieves quiz attempts for a user filtered by topic
// GetQuizAttemptWithDetails retrieves a quiz attempt with full question details
// Phase 4: Timed Test Taking Methods

// CreateTestAttempt creates a new test attempt
// GetTestAttempt retrieves a test attempt by ID
// UpdateTestAttemptAnswer updates an answer in a test attempt
// SubmitTestAttempt submits a test attempt and calculates the score
// GetTestResults retrieves the results of a test attempt
// GetAvailableTests retrieves available tests for a student
// GetTestAttemptsByUser retrieves all test attempts for a user for a specific test
// Phase 4: Tag-Based Review Methods

// GetAllTags retrieves all unique tags from all questions
// GetTagsWithCounts retrieves all tags with their question counts
// GetTagStats retrieves accuracy statistics for each tag for a user
// GetQuestionsByTags retrieves questions that match the given tags
// Phase 4: Test Management (Instructor) Methods

// PublishTest sets a test as published
// GetTestResultsByTest retrieves all test results for a specific test
// ExportTestResults exports test results as CSV
// Phase 4: Question Performance Analytics Methods

// GetQuestionAnalytics retrieves analytics for all questions
// GetQuestionDetailedAnalytics retrieves detailed analytics for a single question
// GetAnswerDistribution retrieves distribution of answers chosen for a question
// UpdateQuestionDifficultyFromAnalytics updates a question's difficulty
