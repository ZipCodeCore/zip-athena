package database

import (
	"database/sql"
	"fmt"

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

// initSchema creates the database tables
func (db *SQLiteDB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0,
		category TEXT DEFAULT 'java',
		total_points INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		category TEXT DEFAULT 'java',
		is_active BOOLEAN DEFAULT 1,
		display_order INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		topic_id INTEGER NOT NULL,
		question_text TEXT NOT NULL,
		difficulty INTEGER NOT NULL,
		tags TEXT,
		topic_question_id TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (topic_id) REFERENCES topics(id)
	);
	-- Note: is_active column for Phase 5 soft deletes to be added via migration

	CREATE TABLE IF NOT EXISTS answers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question_id INTEGER NOT NULL,
		answer_text TEXT NOT NULL,
		is_correct BOOLEAN NOT NULL,
		FOREIGN KEY (question_id) REFERENCES questions(id)
	);

	CREATE TABLE IF NOT EXISTS quiz_attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		quiz_id INTEGER,
		parent_attempt_id INTEGER,
		attempt_number INTEGER DEFAULT 1,
		is_retake BOOLEAN DEFAULT 0,
		start_time DATETIME NOT NULL,
		completed_time DATETIME,
		total_questions INTEGER DEFAULT 0,
		correct_answers INTEGER DEFAULT 0,
		difficulty_average REAL DEFAULT 0,
		description TEXT DEFAULT '',
		question_ids TEXT,
		question_results TEXT,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (quiz_id) REFERENCES quizzes(id),
		FOREIGN KEY (parent_attempt_id) REFERENCES quiz_attempts(id)
	);

	CREATE TABLE IF NOT EXISTS question_attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		question_id INTEGER NOT NULL,
		topic_question_id TEXT NOT NULL,
		correct_attempts INTEGER DEFAULT 0,
		incorrect_attempts INTEGER DEFAULT 0,
		last_attempt_time DATETIME NOT NULL,
		last_attempt_correct BOOLEAN NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (question_id) REFERENCES questions(id)
	);

	CREATE TABLE IF NOT EXISTS quizzes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		topic_id INTEGER NOT NULL,
		description TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (topic_id) REFERENCES topics(id)
	);

	CREATE TABLE IF NOT EXISTS user_question_states (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		question_id INTEGER NOT NULL,
		due DATETIME NOT NULL,
		stability REAL NOT NULL,
		difficulty REAL NOT NULL,
		elapsed_days INTEGER NOT NULL,
		scheduled_days INTEGER NOT NULL,
		reps INTEGER NOT NULL,
		lapses INTEGER NOT NULL,
		state INTEGER NOT NULL,
		last_review DATETIME,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (question_id) REFERENCES questions(id),
		UNIQUE(user_id, question_id)
	);

	CREATE TABLE IF NOT EXISTS tests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		topic_id INTEGER NOT NULL,
		description TEXT,
		time_limit INTEGER DEFAULT 0,
		start_date DATETIME,
		end_date DATETIME,
		is_published INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (topic_id) REFERENCES topics(id)
	);

	CREATE TABLE IF NOT EXISTS test_attempts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		test_id INTEGER NOT NULL,
		start_time DATETIME NOT NULL,
		submit_time DATETIME,
		time_remaining INTEGER DEFAULT 0,
		is_auto_submit INTEGER DEFAULT 0,
		total_questions INTEGER DEFAULT 0,
		correct_answers INTEGER DEFAULT 0,
		score REAL DEFAULT 0.0,
		answers TEXT DEFAULT '{}',
		question_results TEXT DEFAULT '{}',
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS quiz_questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		quiz_id INTEGER NOT NULL,
		question_id INTEGER NOT NULL,
		question_order INTEGER NOT NULL,
		FOREIGN KEY (quiz_id) REFERENCES quizzes(id) ON DELETE CASCADE,
		FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
		UNIQUE(quiz_id, question_id)
	);

	CREATE TABLE IF NOT EXISTS test_questions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		test_id INTEGER NOT NULL,
		question_id INTEGER NOT NULL,
		question_order INTEGER NOT NULL,
		FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE,
		FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
		UNIQUE(test_id, question_id)
	);

	CREATE TABLE IF NOT EXISTS learning_resources (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		language TEXT NOT NULL CHECK(language IN ('java', 'python')),
		phase INTEGER NOT NULL CHECK(phase BETWEEN 1 AND 5),
		section_number TEXT NOT NULL,
		section_title TEXT NOT NULL,
		file_path TEXT NOT NULL,
		display_order INTEGER NOT NULL,
		is_active BOOLEAN DEFAULT 1,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(language, section_number)
	);

	CREATE INDEX IF NOT EXISTS idx_learning_resources_language ON learning_resources(language);
	CREATE INDEX IF NOT EXISTS idx_learning_resources_phase ON learning_resources(phase);
	CREATE INDEX IF NOT EXISTS idx_learning_resources_active ON learning_resources(is_active);

	CREATE TABLE IF NOT EXISTS user_learning_progress (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		learning_resource_id INTEGER NOT NULL,
		first_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		last_viewed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		view_count INTEGER DEFAULT 1,
		completion_percentage REAL DEFAULT 0 CHECK(completion_percentage >= 0 AND completion_percentage <= 100),
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (learning_resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
		UNIQUE(user_id, learning_resource_id)
	);

	CREATE INDEX IF NOT EXISTS idx_user_learning_progress_user ON user_learning_progress(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_learning_progress_resource ON user_learning_progress(learning_resource_id);

	CREATE TABLE IF NOT EXISTS question_learning_links (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		question_id INTEGER NOT NULL,
		learning_resource_id INTEGER NOT NULL,
		relevance_score REAL DEFAULT 1.0 CHECK(relevance_score >= 0 AND relevance_score <= 1),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
		FOREIGN KEY (learning_resource_id) REFERENCES learning_resources(id) ON DELETE CASCADE,
		UNIQUE(question_id, learning_resource_id)
	);

	CREATE INDEX IF NOT EXISTS idx_question_learning_links_question ON question_learning_links(question_id);
	CREATE INDEX IF NOT EXISTS idx_question_learning_links_resource ON question_learning_links(learning_resource_id);
	`

	_, err := db.db.Exec(schema)
	if err != nil {
		return err
	}

	// Run migrations for existing databases
	return db.runMigrations()
}

// runMigrations applies database migrations for schema changes
func (db *SQLiteDB) runMigrations() error {
	// Migration: Add is_active and display_order to topics table if they don't exist
	// Check if columns exist
	var hasIsActive, hasDisplayOrder bool

	// Check for is_active column
	rows, err := db.db.Query("PRAGMA table_info(topics)")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dfltValue interface{}
			if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err == nil {
				if name == "is_active" {
					hasIsActive = true
				}
				if name == "display_order" {
					hasDisplayOrder = true
				}
			}
		}
	}

	// Add missing columns
	if !hasIsActive {
		_, err := db.db.Exec("ALTER TABLE topics ADD COLUMN is_active BOOLEAN DEFAULT 1")
		if err != nil {
			return fmt.Errorf("failed to add is_active column: %w", err)
		}
	}

	if !hasDisplayOrder {
		_, err := db.db.Exec("ALTER TABLE topics ADD COLUMN display_order INTEGER DEFAULT 0")
		if err != nil {
			return fmt.Errorf("failed to add display_order column: %w", err)
		}
	}

	// Phase 5: Create flagged_questions table if it doesn't exist
	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS flagged_questions (
			question_id INTEGER PRIMARY KEY,
			flag_reason TEXT NOT NULL,
			flagged_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			accuracy REAL,
			total_attempts INTEGER,
			FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create flagged_questions table: %w", err)
	}

	// Migration: Activate phases 2 and 3 in learning_resources
	// Only runs if learning_resources table exists (safe no-op otherwise)
	_, err = db.db.Exec(`
		UPDATE learning_resources SET is_active = 1 WHERE phase IN (2, 3) AND is_active = 0
	`)
	if err != nil {
		return fmt.Errorf("failed to activate phases 2 and 3: %w", err)
	}

	// Migration: Fix file_path values that don't match actual filenames on disk
	filePathFixes := map[string]string{
		"java/8-classes-and-objects.md":    "java/8-oop-classes-and-objects.md",
		"java/9-inheritance.md":            "java/9-oop-inheritance.md",
		"java/10-polymorphism.md":          "java/10-oop-polymorphism.md",
		"java/15-collections.md":           "java/15-collections-framework.md",
		"python/11-modules-packages.md":    "python/11-modules-and-packages.md",
		"python/14-classes-and-objects.md": "python/14-oop-classes-and-objects.md",
		"python/15-inheritance.md":         "python/15-oop-inheritance.md",
		"python/16-polymorphism.md":        "python/16-oop-polymorphism.md",
	}
	for oldPath, newPath := range filePathFixes {
		_, err = db.db.Exec(`UPDATE learning_resources SET file_path = ? WHERE file_path = ?`, newPath, oldPath)
		if err != nil {
			return fmt.Errorf("failed to fix file_path %s -> %s: %w", oldPath, newPath, err)
		}
	}

	// Migration: Add description column to quiz_attempts table if it doesn't exist
	var hasDescription bool
	qaRows, err := db.db.Query("PRAGMA table_info(quiz_attempts)")
	if err == nil {
		defer qaRows.Close()
		for qaRows.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dfltValue interface{}
			if err := qaRows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err == nil {
				if name == "description" {
					hasDescription = true
				}
			}
		}
	}
	if !hasDescription {
		_, err := db.db.Exec("ALTER TABLE quiz_attempts ADD COLUMN description TEXT DEFAULT ''")
		if err != nil {
			return fmt.Errorf("failed to add description column to quiz_attempts: %w", err)
		}
	}

	return nil
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
