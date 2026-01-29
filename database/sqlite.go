package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
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
		total_points INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS topics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
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

	return nil
}

// User operations
func (db *SQLiteDB) CreateUser(user *models.User) error {
	result, err := db.db.Exec(
		"INSERT INTO users (username, email, password, is_admin, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.IsAdmin, user.TotalPoints, user.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetUser(id int) (*models.User, error) {
	user := &models.User{}
	err := db.db.QueryRow(
		"SELECT id, username, email, is_admin, total_points, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.TotalPoints, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *SQLiteDB) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := db.db.QueryRow(
		"SELECT id, username, email, password, is_admin, total_points, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.TotalPoints, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *SQLiteDB) UpdateUserPoints(id int, points int) error {
	result, err := db.db.Exec("UPDATE users SET total_points = ? WHERE id = ?", points, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Topic operations
func (db *SQLiteDB) CreateTopic(topic *models.Topic) error {
	// Default to active for new topics
	if topic.ID == 0 {
		topic.IsActive = true
	}

	// Get max display_order and increment
	if topic.DisplayOrder == 0 {
		var maxOrder int
		err := db.db.QueryRow("SELECT COALESCE(MAX(display_order), 0) FROM topics").Scan(&maxOrder)
		if err == nil {
			topic.DisplayOrder = maxOrder + 1
		}
	}

	result, err := db.db.Exec(
		"INSERT INTO topics (name, description, is_active, display_order, created_at) VALUES (?, ?, ?, ?, ?)",
		topic.Name, topic.Description, topic.IsActive, topic.DisplayOrder, topic.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	topic.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetTopic(id int) (*models.Topic, error) {
	topic := &models.Topic{}
	err := db.db.QueryRow(
		"SELECT id, name, description, is_active, display_order, created_at FROM topics WHERE id = ?",
		id,
	).Scan(&topic.ID, &topic.Name, &topic.Description, &topic.IsActive, &topic.DisplayOrder, &topic.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func (db *SQLiteDB) ListTopics() ([]*models.Topic, error) {
	rows, err := db.db.Query("SELECT id, name, description, is_active, display_order, created_at FROM topics ORDER BY display_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*models.Topic
	for rows.Next() {
		topic := &models.Topic{}
		if err := rows.Scan(&topic.ID, &topic.Name, &topic.Description, &topic.IsActive, &topic.DisplayOrder, &topic.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, rows.Err()
}

// Question operations
func (db *SQLiteDB) CreateQuestion(question *models.Question) error {
	tagsJSON, err := json.Marshal(question.Tags)
	if err != nil {
		return err
	}

	result, err := db.db.Exec(
		"INSERT INTO questions (topic_id, question_text, difficulty, tags, topic_question_id, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		question.TopicID, question.QuestionText, question.Difficulty, string(tagsJSON), question.TopicQuestionID, question.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	question.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetQuestion(id int) (*models.Question, error) {
	question := &models.Question{}
	var tagsJSON string
	err := db.db.QueryRow(
		"SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions WHERE id = ?",
		id,
	).Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(tagsJSON), &question.Tags); err != nil {
		return nil, err
	}
	return question, nil
}

func (db *SQLiteDB) GetQuestionsByTopic(topicID int) ([]*models.Question, error) {
	rows, err := db.db.Query(
		"SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions WHERE topic_id = ?",
		topicID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsJSON), &question.Tags); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

// Answer operations
func (db *SQLiteDB) CreateAnswer(answer *models.Answer) error {
	result, err := db.db.Exec(
		"INSERT INTO answers (question_id, answer_text, is_correct) VALUES (?, ?, ?)",
		answer.QuestionID, answer.AnswerText, answer.IsCorrect,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	answer.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetAnswersByQuestion(questionID int) ([]*models.Answer, error) {
	rows, err := db.db.Query(
		"SELECT id, question_id, answer_text, is_correct FROM answers WHERE question_id = ?",
		questionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*models.Answer
	for rows.Next() {
		answer := &models.Answer{}
		if err := rows.Scan(&answer.ID, &answer.QuestionID, &answer.AnswerText, &answer.IsCorrect); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}
	return answers, rows.Err()
}

// QuizAttempt operations
func (db *SQLiteDB) CreateQuizAttempt(attempt *models.QuizAttempt) error {
	questionIDsJSON, err := json.Marshal(attempt.QuestionIDs)
	if err != nil {
		return err
	}
	questionResultsJSON, err := json.Marshal(attempt.QuestionResults)
	if err != nil {
		return err
	}

	result, err := db.db.Exec(
		"INSERT INTO quiz_attempts (user_id, start_time, completed_time, total_questions, correct_answers, question_ids, question_results) VALUES (?, ?, ?, ?, ?, ?, ?)",
		attempt.UserID, attempt.StartTime, attempt.CompletedTime, attempt.TotalQuestions, attempt.CorrectAnswers, string(questionIDsJSON), string(questionResultsJSON),
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	attempt.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetQuizAttempt(id int) (*models.QuizAttempt, error) {
	attempt := &models.QuizAttempt{}
	var questionIDsJSON, questionResultsJSON string
	var completedTime sql.NullTime

	err := db.db.QueryRow(
		"SELECT id, user_id, start_time, completed_time, total_questions, correct_answers, question_ids, question_results FROM quiz_attempts WHERE id = ?",
		id,
	).Scan(&attempt.ID, &attempt.UserID, &attempt.StartTime, &completedTime, &attempt.TotalQuestions, &attempt.CorrectAnswers, &questionIDsJSON, &questionResultsJSON)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if completedTime.Valid {
		attempt.CompletedTime = &completedTime.Time
	}

	if err := json.Unmarshal([]byte(questionIDsJSON), &attempt.QuestionIDs); err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(questionResultsJSON), &attempt.QuestionResults); err != nil {
		return nil, err
	}

	return attempt, nil
}

func (db *SQLiteDB) GetQuizAttemptsByUser(userID int) ([]*models.QuizAttempt, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, start_time, completed_time, total_questions, correct_answers, question_ids, question_results FROM quiz_attempts WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.QuizAttempt
	for rows.Next() {
		attempt := &models.QuizAttempt{}
		var questionIDsJSON, questionResultsJSON string
		var completedTime sql.NullTime

		if err := rows.Scan(&attempt.ID, &attempt.UserID, &attempt.StartTime, &completedTime, &attempt.TotalQuestions, &attempt.CorrectAnswers, &questionIDsJSON, &questionResultsJSON); err != nil {
			return nil, err
		}

		if completedTime.Valid {
			attempt.CompletedTime = &completedTime.Time
		}

		if err := json.Unmarshal([]byte(questionIDsJSON), &attempt.QuestionIDs); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(questionResultsJSON), &attempt.QuestionResults); err != nil {
			return nil, err
		}

		attempts = append(attempts, attempt)
	}
	return attempts, rows.Err()
}

func (db *SQLiteDB) UpdateQuizAttempt(attempt *models.QuizAttempt) error {
	questionIDsJSON, err := json.Marshal(attempt.QuestionIDs)
	if err != nil {
		return err
	}
	questionResultsJSON, err := json.Marshal(attempt.QuestionResults)
	if err != nil {
		return err
	}

	result, err := db.db.Exec(
		"UPDATE quiz_attempts SET quiz_id = ?, parent_attempt_id = ?, attempt_number = ?, is_retake = ?, completed_time = ?, total_questions = ?, correct_answers = ?, question_ids = ?, question_results = ? WHERE id = ?",
		attempt.QuizID, attempt.ParentAttemptID, attempt.AttemptNumber, attempt.IsRetake, attempt.CompletedTime, attempt.TotalQuestions, attempt.CorrectAnswers, string(questionIDsJSON), string(questionResultsJSON), attempt.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// QuestionAttempt operations
func (db *SQLiteDB) CreateQuestionAttempt(attempt *models.QuestionAttempt) error {
	result, err := db.db.Exec(
		"INSERT INTO question_attempts (user_id, question_id, topic_question_id, correct_attempts, incorrect_attempts, last_attempt_time, last_attempt_correct) VALUES (?, ?, ?, ?, ?, ?, ?)",
		attempt.UserID, attempt.QuestionID, attempt.TopicQuestionID, attempt.CorrectAttempts, attempt.IncorrectAttempts, attempt.LastAttemptTime, attempt.LastAttemptCorrect,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	attempt.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetQuestionAttempt(id int) (*models.QuestionAttempt, error) {
	attempt := &models.QuestionAttempt{}
	err := db.db.QueryRow(
		"SELECT id, user_id, question_id, topic_question_id, correct_attempts, incorrect_attempts, last_attempt_time, last_attempt_correct FROM question_attempts WHERE id = ?",
		id,
	).Scan(&attempt.ID, &attempt.UserID, &attempt.QuestionID, &attempt.TopicQuestionID, &attempt.CorrectAttempts, &attempt.IncorrectAttempts, &attempt.LastAttemptTime, &attempt.LastAttemptCorrect)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return attempt, nil
}

func (db *SQLiteDB) GetQuestionAttemptsByUser(userID int) ([]*models.QuestionAttempt, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, question_id, topic_question_id, correct_attempts, incorrect_attempts, last_attempt_time, last_attempt_correct FROM question_attempts WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.QuestionAttempt
	for rows.Next() {
		attempt := &models.QuestionAttempt{}
		if err := rows.Scan(&attempt.ID, &attempt.UserID, &attempt.QuestionID, &attempt.TopicQuestionID, &attempt.CorrectAttempts, &attempt.IncorrectAttempts, &attempt.LastAttemptTime, &attempt.LastAttemptCorrect); err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}
	return attempts, rows.Err()
}

func (db *SQLiteDB) UpdateQuestionAttempt(attempt *models.QuestionAttempt) error {
	result, err := db.db.Exec(
		"UPDATE question_attempts SET correct_attempts = ?, incorrect_attempts = ?, last_attempt_time = ?, last_attempt_correct = ? WHERE id = ?",
		attempt.CorrectAttempts, attempt.IncorrectAttempts, attempt.LastAttemptTime, attempt.LastAttemptCorrect, attempt.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Quiz operations
func (db *SQLiteDB) CreateQuiz(quiz *models.Quiz) error {
	result, err := db.db.Exec(
		"INSERT INTO quizzes (name, topic_id, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		quiz.Name, quiz.TopicID, quiz.Description, quiz.CreatedAt, quiz.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	quiz.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetQuiz(id int) (*models.Quiz, error) {
	quiz := &models.Quiz{}
	err := db.db.QueryRow(
		"SELECT id, name, topic_id, description, created_at, updated_at FROM quizzes WHERE id = ?",
		id,
	).Scan(&quiz.ID, &quiz.Name, &quiz.TopicID, &quiz.Description, &quiz.CreatedAt, &quiz.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (db *SQLiteDB) ListQuizzes() ([]*models.Quiz, error) {
	rows, err := db.db.Query("SELECT id, name, topic_id, description, created_at, updated_at FROM quizzes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*models.Quiz
	for rows.Next() {
		quiz := &models.Quiz{}
		if err := rows.Scan(&quiz.ID, &quiz.Name, &quiz.TopicID, &quiz.Description, &quiz.CreatedAt, &quiz.UpdatedAt); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func (db *SQLiteDB) UpdateQuiz(quiz *models.Quiz) error {
	result, err := db.db.Exec(
		"UPDATE quizzes SET name = ?, topic_id = ?, description = ?, updated_at = ? WHERE id = ?",
		quiz.Name, quiz.TopicID, quiz.Description, quiz.UpdatedAt, quiz.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) DeleteQuiz(id int) error {
	result, err := db.db.Exec("DELETE FROM quizzes WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// UserQuestionState operations
func (db *SQLiteDB) CreateUserQuestionState(state *models.UserQuestionState) error {
	result, err := db.db.Exec(
		"INSERT INTO user_question_states (user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		state.UserID, state.QuestionID, state.Due, state.Stability, state.Difficulty, state.ElapsedDays, state.ScheduledDays, state.Reps, state.Lapses, state.State, state.LastReview,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	state.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetUserQuestionState(userID, questionID int) (*models.UserQuestionState, error) {
	state := &models.UserQuestionState{}
	err := db.db.QueryRow(
		"SELECT id, user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review FROM user_question_states WHERE user_id = ? AND question_id = ?",
		userID, questionID,
	).Scan(&state.ID, &state.UserID, &state.QuestionID, &state.Due, &state.Stability, &state.Difficulty, &state.ElapsedDays, &state.ScheduledDays, &state.Reps, &state.Lapses, &state.State, &state.LastReview)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (db *SQLiteDB) GetUserQuestionStatesByUser(userID int) ([]*models.UserQuestionState, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review FROM user_question_states WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.UserQuestionState
	for rows.Next() {
		state := &models.UserQuestionState{}
		if err := rows.Scan(&state.ID, &state.UserID, &state.QuestionID, &state.Due, &state.Stability, &state.Difficulty, &state.ElapsedDays, &state.ScheduledDays, &state.Reps, &state.Lapses, &state.State, &state.LastReview); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (db *SQLiteDB) GetDueUserQuestionStates(userID int) ([]*models.UserQuestionState, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review FROM user_question_states WHERE user_id = ? AND due <= datetime('now')",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.UserQuestionState
	for rows.Next() {
		state := &models.UserQuestionState{}
		if err := rows.Scan(&state.ID, &state.UserID, &state.QuestionID, &state.Due, &state.Stability, &state.Difficulty, &state.ElapsedDays, &state.ScheduledDays, &state.Reps, &state.Lapses, &state.State, &state.LastReview); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (db *SQLiteDB) UpdateUserQuestionState(state *models.UserQuestionState) error {
	result, err := db.db.Exec(
		"UPDATE user_question_states SET due = ?, stability = ?, difficulty = ?, elapsed_days = ?, scheduled_days = ?, reps = ?, lapses = ?, state = ?, last_review = ? WHERE id = ?",
		state.Due, state.Stability, state.Difficulty, state.ElapsedDays, state.ScheduledDays, state.Reps, state.Lapses, state.State, state.LastReview, state.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Admin operations - Users
func (db *SQLiteDB) ListUsers() ([]*models.User, error) {
	rows, err := db.db.Query("SELECT id, username, email, is_admin, total_points, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.TotalPoints, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (db *SQLiteDB) UpdateUser(user *models.User) error {
	result, err := db.db.Exec(
		"UPDATE users SET username = ?, email = ?, is_admin = ?, total_points = ? WHERE id = ?",
		user.Username, user.Email, user.IsAdmin, user.TotalPoints, user.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) DeleteUser(id int) error {
	result, err := db.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Admin operations - Questions
func (db *SQLiteDB) ListQuestions() ([]*models.Question, error) {
	rows, err := db.db.Query(
		"SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(tagsJSON), &question.Tags); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}
	return questions, rows.Err()
}

func (db *SQLiteDB) UpdateQuestion(question *models.Question) error {
	tagsJSON, err := json.Marshal(question.Tags)
	if err != nil {
		return err
	}

	result, err := db.db.Exec(
		"UPDATE questions SET topic_id = ?, question_text = ?, difficulty = ?, tags = ?, topic_question_id = ? WHERE id = ?",
		question.TopicID, question.QuestionText, question.Difficulty, string(tagsJSON), question.TopicQuestionID, question.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) DeleteQuestion(id int) error {
	result, err := db.db.Exec("DELETE FROM questions WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Admin operations - Answers
func (db *SQLiteDB) UpdateAnswer(answer *models.Answer) error {
	result, err := db.db.Exec(
		"UPDATE answers SET question_id = ?, answer_text = ?, is_correct = ? WHERE id = ?",
		answer.QuestionID, answer.AnswerText, answer.IsCorrect, answer.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) DeleteAnswer(id int) error {
	result, err := db.db.Exec("DELETE FROM answers WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Admin operations - Topics
func (db *SQLiteDB) UpdateTopic(topic *models.Topic) error {
	result, err := db.db.Exec(
		"UPDATE topics SET name = ?, description = ?, is_active = ?, display_order = ? WHERE id = ?",
		topic.Name, topic.Description, topic.IsActive, topic.DisplayOrder, topic.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) DeleteTopic(id int) error {
	// Soft delete: set is_active to false instead of deleting
	result, err := db.db.Exec("UPDATE topics SET is_active = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// GetAnswer retrieves an answer by ID
func (db *SQLiteDB) GetAnswer(id int) (*models.Answer, error) {
	answer := &models.Answer{}
	err := db.db.QueryRow(
		"SELECT id, question_id, answer_text, is_correct FROM answers WHERE id = ?",
		id,
	).Scan(&answer.ID, &answer.QuestionID, &answer.AnswerText, &answer.IsCorrect)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return answer, nil
}

// Phase 3: Test operations
func (db *SQLiteDB) CreateTest(test *models.Test) error {
	result, err := db.db.Exec(
		"INSERT INTO tests (name, topic_id, description, time_limit, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		test.Name, test.TopicID, test.Description, test.TimeLimit, test.CreatedAt, test.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	test.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetTest(id int) (*models.Test, error) {
	test := &models.Test{}
	err := db.db.QueryRow(
		"SELECT id, name, topic_id, description, time_limit, created_at, updated_at FROM tests WHERE id = ?",
		id,
	).Scan(&test.ID, &test.Name, &test.TopicID, &test.Description, &test.TimeLimit, &test.CreatedAt, &test.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return test, nil
}

func (db *SQLiteDB) ListTests() ([]*models.Test, error) {
	rows, err := db.db.Query("SELECT id, name, topic_id, description, time_limit, created_at, updated_at FROM tests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tests := make([]*models.Test, 0)
	for rows.Next() {
		test := &models.Test{}
		err := rows.Scan(&test.ID, &test.Name, &test.TopicID, &test.Description, &test.TimeLimit, &test.CreatedAt, &test.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	return tests, nil
}

func (db *SQLiteDB) UpdateTest(test *models.Test) error {
	result, err := db.db.Exec(
		"UPDATE tests SET name = ?, topic_id = ?, description = ?, time_limit = ?, updated_at = ? WHERE id = ?",
		test.Name, test.TopicID, test.Description, test.TimeLimit, test.UpdatedAt, test.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) DeleteTest(id int) error {
	result, err := db.db.Exec("DELETE FROM tests WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// Phase 3: Many-to-many Quiz-Question relationships
func (db *SQLiteDB) AddQuestionToQuiz(quizID, questionID, order int) error {
	_, err := db.db.Exec(
		"INSERT INTO quiz_questions (quiz_id, question_id, question_order) VALUES (?, ?, ?)",
		quizID, questionID, order,
	)
	return err
}

func (db *SQLiteDB) RemoveQuestionFromQuiz(quizID, questionID int) error {
	result, err := db.db.Exec(
		"DELETE FROM quiz_questions WHERE quiz_id = ? AND question_id = ?",
		quizID, questionID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) GetQuizQuestions(quizID int) ([]*models.Question, error) {
	rows, err := db.db.Query(`
		SELECT q.id, q.topic_id, q.question_text, q.difficulty, q.tags, q.topic_question_id, q.created_at
		FROM questions q
		INNER JOIN quiz_questions qq ON q.id = qq.question_id
		WHERE qq.quiz_id = ?
		ORDER BY qq.question_order
	`, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]*models.Question, 0)
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt)
		if err != nil {
			return nil, err
		}
		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &question.Tags); err != nil {
				question.Tags = []string{}
			}
		} else {
			question.Tags = []string{}
		}
		questions = append(questions, question)
	}
	return questions, nil
}

// Phase 3: Many-to-many Test-Question relationships
func (db *SQLiteDB) AddQuestionToTest(testID, questionID, order int) error {
	_, err := db.db.Exec(
		"INSERT INTO test_questions (test_id, question_id, question_order) VALUES (?, ?, ?)",
		testID, questionID, order,
	)
	return err
}

func (db *SQLiteDB) RemoveQuestionFromTest(testID, questionID int) error {
	result, err := db.db.Exec(
		"DELETE FROM test_questions WHERE test_id = ? AND question_id = ?",
		testID, questionID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (db *SQLiteDB) GetTestQuestions(testID int) ([]*models.Question, error) {
	rows, err := db.db.Query(`
		SELECT q.id, q.topic_id, q.question_text, q.difficulty, q.tags, q.topic_question_id, q.created_at
		FROM questions q
		INNER JOIN test_questions tq ON q.id = tq.question_id
		WHERE tq.test_id = ?
		ORDER BY tq.question_order
	`, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questions := make([]*models.Question, 0)
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt)
		if err != nil {
			return nil, err
		}
		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &question.Tags); err != nil {
				question.Tags = []string{}
			}
		} else {
			question.Tags = []string{}
		}
		questions = append(questions, question)
	}
	return questions, nil
}

// Phase 3: Analytics operations
func (db *SQLiteDB) GetStudentProgress(userID int) (*models.StudentProgress, error) {
	// Get user info
	user, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	progress := &models.StudentProgress{
		UserID:         userID,
		Username:       user.Username,
		TotalPoints:    user.TotalPoints,
		TopicBreakdown: make(map[string]float64),
		RecentActivity: make([]models.QuizAttempt, 0),
	}

	// Get quiz attempt stats
	rows, err := db.db.Query(`
		SELECT COUNT(*), 
		       SUM(CASE WHEN completed_time IS NOT NULL THEN 1 ELSE 0 END),
		       AVG(CASE WHEN completed_time IS NOT NULL AND total_questions > 0 
		                THEN CAST(correct_answers AS REAL) / total_questions * 100 
		                ELSE 0 END)
		FROM quiz_attempts WHERE user_id = ?
	`, userID)
	if err == nil && rows.Next() {
		rows.Scan(&progress.TotalQuizzes, &progress.CompletedQuizzes, &progress.AverageScore)
		rows.Close()
	}

	// Get recent activity (last 5 completed quizzes)
	activityRows, err := db.db.Query(`
		SELECT id, user_id, start_time, completed_time, total_questions, correct_answers, question_ids, question_results
		FROM quiz_attempts 
		WHERE user_id = ? AND completed_time IS NOT NULL
		ORDER BY start_time DESC LIMIT 5
	`, userID)
	if err == nil {
		defer activityRows.Close()
		for activityRows.Next() {
			attempt := models.QuizAttempt{}
			var questionIDsJSON, questionResultsJSON string
			err := activityRows.Scan(&attempt.ID, &attempt.UserID, &attempt.StartTime, &attempt.CompletedTime,
				&attempt.TotalQuestions, &attempt.CorrectAnswers, &questionIDsJSON, &questionResultsJSON)
			if err == nil {
				json.Unmarshal([]byte(questionIDsJSON), &attempt.QuestionIDs)
				json.Unmarshal([]byte(questionResultsJSON), &attempt.QuestionResults)
				progress.RecentActivity = append(progress.RecentActivity, attempt)
			}
		}
	}

	// Get question accuracy
	qaRows, err := db.db.Query(`
		SELECT SUM(correct_attempts + incorrect_attempts), SUM(correct_attempts)
		FROM question_attempts WHERE user_id = ?
	`, userID)
	if err == nil && qaRows.Next() {
		var totalAttempts, correctAttempts int
		qaRows.Scan(&totalAttempts, &correctAttempts)
		qaRows.Close()
		if totalAttempts > 0 {
			progress.QuestionAccuracy = float64(correctAttempts) / float64(totalAttempts) * 100
		}
	}

	// Get topic breakdown
	topicRows, err := db.db.Query(`
		SELECT t.name, 
		       CAST(SUM(qa.correct_attempts) AS REAL) / SUM(qa.correct_attempts + qa.incorrect_attempts) * 100
		FROM question_attempts qa
		INNER JOIN questions q ON qa.question_id = q.id
		INNER JOIN topics t ON q.topic_id = t.id
		WHERE qa.user_id = ?
		GROUP BY t.name
		HAVING SUM(qa.correct_attempts + qa.incorrect_attempts) > 0
	`, userID)
	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var topicName string
			var accuracy float64
			if topicRows.Scan(&topicName, &accuracy) == nil {
				progress.TopicBreakdown[topicName] = accuracy
			}
		}
	}

	// Get FSRS metrics
	fsrsRows, err := db.db.Query(`
		SELECT COUNT(*),
		       SUM(CASE WHEN state = 0 THEN 1 ELSE 0 END),
		       SUM(CASE WHEN state = 2 THEN 1 ELSE 0 END),
		       SUM(CASE WHEN due <= datetime('now') THEN 1 ELSE 0 END),
		       AVG(stability),
		       AVG(difficulty)
		FROM user_question_states WHERE user_id = ?
	`, userID)
	if err == nil && fsrsRows.Next() {
		fsrsRows.Scan(&progress.FSRSMetrics.TotalCards, &progress.FSRSMetrics.NewCards,
			&progress.FSRSMetrics.ReviewCards, &progress.FSRSMetrics.DueCards,
			&progress.FSRSMetrics.AvgStability, &progress.FSRSMetrics.AvgDifficulty)
		fsrsRows.Close()
	}

	return progress, nil
}

func (db *SQLiteDB) GetCohortAnalytics() (*models.CohortAnalytics, error) {
	analytics := &models.CohortAnalytics{
		TopicDifficulty:  make(map[string]float64),
		QuestionStats:    make([]models.QuestionStatistic, 0),
		TopPerformers:    make([]models.StudentProgress, 0),
		BottomPerformers: make([]models.StudentProgress, 0),
	}

	// Count total students (non-admin)
	row := db.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = 0")
	row.Scan(&analytics.TotalStudents)

	// Count active students (last 7 days)
	activeRow := db.db.QueryRow(`
		SELECT COUNT(DISTINCT user_id) FROM quiz_attempts 
		WHERE start_time >= datetime('now', '-7 days')
	`)
	activeRow.Scan(&analytics.ActiveStudents)

	// Get total quizzes and average score
	statsRow := db.db.QueryRow(`
		SELECT COUNT(*),
		       AVG(CASE WHEN total_questions > 0 
		                THEN CAST(correct_answers AS REAL) / total_questions * 100 
		                ELSE 0 END)
		FROM quiz_attempts WHERE completed_time IS NOT NULL
	`)
	statsRow.Scan(&analytics.TotalQuizzesTaken, &analytics.AverageCohortScore)

	// Get top performers
	topRows, err := db.db.Query(`
		SELECT user_id, AVG(CAST(correct_answers AS REAL) / total_questions * 100) as avg_score
		FROM quiz_attempts 
		WHERE completed_time IS NOT NULL AND total_questions > 0
		GROUP BY user_id
		ORDER BY avg_score DESC LIMIT 5
	`)
	if err == nil {
		defer topRows.Close()
		for topRows.Next() {
			var userID int
			var avgScore float64
			if topRows.Scan(&userID, &avgScore) == nil {
				if progress, err := db.GetStudentProgress(userID); err == nil {
					analytics.TopPerformers = append(analytics.TopPerformers, *progress)
				}
			}
		}
	}

	// Get bottom performers
	bottomRows, err := db.db.Query(`
		SELECT user_id, AVG(CAST(correct_answers AS REAL) / total_questions * 100) as avg_score
		FROM quiz_attempts 
		WHERE completed_time IS NOT NULL AND total_questions > 0
		GROUP BY user_id
		ORDER BY avg_score ASC LIMIT 5
	`)
	if err == nil {
		defer bottomRows.Close()
		for bottomRows.Next() {
			var userID int
			var avgScore float64
			if bottomRows.Scan(&userID, &avgScore) == nil {
				if progress, err := db.GetStudentProgress(userID); err == nil {
					analytics.BottomPerformers = append(analytics.BottomPerformers, *progress)
				}
			}
		}
	}

	// Get topic difficulty
	topicRows, err := db.db.Query(`
		SELECT t.name, AVG(uqs.difficulty)
		FROM user_question_states uqs
		INNER JOIN questions q ON uqs.question_id = q.id
		INNER JOIN topics t ON q.topic_id = t.id
		GROUP BY t.name
	`)
	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var topicName string
			var avgDifficulty float64
			if topicRows.Scan(&topicName, &avgDifficulty) == nil {
				analytics.TopicDifficulty[topicName] = avgDifficulty
			}
		}
	}

	// Get question statistics
	questionRows, err := db.db.Query(`
		SELECT q.id, q.topic_question_id, q.question_text,
		       SUM(qa.correct_attempts + qa.incorrect_attempts) as total_attempts,
		       SUM(qa.correct_attempts) as correct_attempts,
		       CAST(SUM(qa.correct_attempts) AS REAL) / SUM(qa.correct_attempts + qa.incorrect_attempts) * 100 as accuracy,
		       COALESCE(AVG(uqs.difficulty), 0) as avg_difficulty
		FROM questions q
		LEFT JOIN question_attempts qa ON q.id = qa.question_id
		LEFT JOIN user_question_states uqs ON q.id = uqs.question_id
		GROUP BY q.id
		HAVING total_attempts > 0
		ORDER BY accuracy ASC
		LIMIT 20
	`)
	if err == nil {
		defer questionRows.Close()
		for questionRows.Next() {
			stat := models.QuestionStatistic{}
			questionRows.Scan(&stat.QuestionID, &stat.TopicQuestionID, &stat.QuestionText,
				&stat.TotalAttempts, &stat.CorrectAttempts, &stat.Accuracy, &stat.AvgDifficulty)
			analytics.QuestionStats = append(analytics.QuestionStats, stat)
		}
	}

	return analytics, nil
}

// ListStudentsWithMetrics returns a list of all students with their performance metrics
func (db *SQLiteDB) ListStudentsWithMetrics(searchQuery string, sortBy string, filterCohort string) ([]*models.StudentProgress, error) {
	// Get all non-admin users
	query := "SELECT id FROM users WHERE is_admin = 0"
	args := []interface{}{}

	// Add search filter if provided
	if searchQuery != "" {
		query += " AND (username LIKE ? OR email LIKE ?)"
		searchPattern := "%" + searchQuery + "%"
		args = append(args, searchPattern, searchPattern)
	}

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]*models.StudentProgress, 0)
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			continue
		}

		// Get full progress for each student
		progress, err := db.GetStudentProgress(userID)
		if err != nil {
			continue
		}
		students = append(students, progress)
	}

	// Sort students based on sortBy parameter
	// For now, we return unsorted - sorting can be done in the handler
	// TODO: Implement sorting by accuracy, completion rate, last activity

	return students, nil
}

// GetStudentQuizHistory returns all quiz attempts for a student
func (db *SQLiteDB) GetStudentQuizHistory(userID int) ([]*models.QuizAttempt, error) {
	// Verify user exists
	_, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query(`
		SELECT id, user_id, quiz_id, parent_attempt_id, attempt_number, is_retake,
		       start_time, completed_time, total_questions, correct_answers
		FROM quiz_attempts
		WHERE user_id = ?
		ORDER BY start_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attempts := make([]*models.QuizAttempt, 0)
	for rows.Next() {
		attempt := &models.QuizAttempt{}
		var quizID, parentAttemptID *int
		var completedTime *time.Time

		err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&quizID,
			&parentAttemptID,
			&attempt.AttemptNumber,
			&attempt.IsRetake,
			&attempt.StartTime,
			&completedTime,
			&attempt.TotalQuestions,
			&attempt.CorrectAnswers,
		)
		if err != nil {
			continue
		}

		attempt.QuizID = quizID
		attempt.ParentAttemptID = parentAttemptID
		attempt.CompletedTime = completedTime
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

// GetStudentWeakAreas identifies topics and tags where student is struggling (< 60% accuracy)
func (db *SQLiteDB) GetStudentWeakAreas(userID int) (*models.WeakAreas, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	weakAreas := &models.WeakAreas{
		UserID:     userID,
		Username:   user.Username,
		WeakTopics: make([]models.TopicWeakness, 0),
		WeakTags:   make([]models.TagWeakness, 0),
	}

	// Get weak topics (< 60% accuracy)
	topicRows, err := db.db.Query(`
		SELECT t.id, t.name,
		       SUM(qa.correct_attempts) as correct,
		       SUM(qa.correct_attempts + qa.incorrect_attempts) as total,
		       CAST(SUM(qa.correct_attempts) AS REAL) / SUM(qa.correct_attempts + qa.incorrect_attempts) * 100 as accuracy
		FROM topics t
		JOIN questions q ON t.id = q.topic_id
		JOIN question_attempts qa ON q.id = qa.question_id
		WHERE qa.user_id = ?
		GROUP BY t.id, t.name
		HAVING total > 0 AND accuracy < 60
		ORDER BY accuracy ASC
	`, userID)

	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var topicID int
			var topicName string
			var correct, total int
			var accuracy float64

			if topicRows.Scan(&topicID, &topicName, &correct, &total, &accuracy) == nil {
				weakAreas.WeakTopics = append(weakAreas.WeakTopics, models.TopicWeakness{
					TopicID:       topicID,
					TopicName:     topicName,
					Accuracy:      accuracy,
					TotalAttempts: total,
				})
			}
		}
	}

	// Get weak tags (< 60% accuracy)
	// This requires extracting tags from questions and calculating accuracy per tag
	// For simplicity, we'll query question attempts and aggregate by tags
	qaRows, err := db.db.Query(`
		SELECT qa.question_id, qa.correct_attempts, qa.incorrect_attempts
		FROM question_attempts qa
		WHERE qa.user_id = ?
	`, userID)

	if err == nil {
		defer qaRows.Close()

		// Map to store tag statistics: tag -> [correct, total]
		tagStats := make(map[string][2]int)

		for qaRows.Next() {
			var questionID, correct, incorrect int
			if qaRows.Scan(&questionID, &correct, &incorrect) == nil {
				// Get question tags
				question, err := db.GetQuestion(questionID)
				if err != nil {
					continue
				}

				total := correct + incorrect
				for _, tag := range question.Tags {
					stats := tagStats[tag]
					stats[0] += correct
					stats[1] += total
					tagStats[tag] = stats
				}
			}
		}

		// Calculate accuracy for each tag and filter weak ones
		for tag, stats := range tagStats {
			correct := stats[0]
			total := stats[1]

			if total > 0 {
				accuracy := float64(correct) / float64(total) * 100
				if accuracy < 60 {
					weakAreas.WeakTags = append(weakAreas.WeakTags, models.TagWeakness{
						Tag:           tag,
						Accuracy:      accuracy,
						TotalAttempts: total,
					})
				}
			}
		}
	}

	return weakAreas, nil
}

// ReorderTopics updates the display_order for topics
func (db *SQLiteDB) ReorderTopics(topicIDs []int) error {
	// First validate all topics exist
	for _, topicID := range topicIDs {
		var exists int
		err := db.db.QueryRow("SELECT COUNT(*) FROM topics WHERE id = ?", topicID).Scan(&exists)
		if err != nil {
			return err
		}
		if exists == 0 {
			return ErrNotFound
		}
	}

	// Update display_order for each topic based on its position in the array
	for i, topicID := range topicIDs {
		_, err := db.db.Exec("UPDATE topics SET display_order = ? WHERE id = ?", i, topicID)
		if err != nil {
			return err
		}
	}
	return nil
}

// Phase 4: SQLite implementations to append to sqlite.go

// Phase 4: Quiz History Methods

// GetQuizHistoryByUser retrieves quiz attempts for a user with pagination
func (db *SQLiteDB) GetQuizHistoryByUser(userID int, page int, perPage int) ([]*models.QuizAttempt, error) {
	offset := (page - 1) * perPage
	query := `
		SELECT id, user_id, quiz_id, parent_attempt_id, attempt_number, is_retake,
			   start_time, completed_time, total_questions, correct_answers,
			   question_ids, question_results
		FROM quiz_attempts
		WHERE user_id = ? AND completed_time IS NOT NULL
		ORDER BY completed_time DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.db.Query(query, userID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.QuizAttempt
	for rows.Next() {
		attempt := &models.QuizAttempt{}
		var quizID, parentAttemptID sql.NullInt64
		var completedTime sql.NullTime
		var questionIDs, questionResults string

		err := rows.Scan(
			&attempt.ID, &attempt.UserID, &quizID, &parentAttemptID,
			&attempt.AttemptNumber, &attempt.IsRetake, &attempt.StartTime, &completedTime,
			&attempt.TotalQuestions, &attempt.CorrectAnswers, &questionIDs, &questionResults,
		)
		if err != nil {
			return nil, err
		}

		if quizID.Valid {
			id := int(quizID.Int64)
			attempt.QuizID = &id
		}
		if parentAttemptID.Valid {
			id := int(parentAttemptID.Int64)
			attempt.ParentAttemptID = &id
		}
		if completedTime.Valid {
			attempt.CompletedTime = &completedTime.Time
		}

		// Parse JSON fields
		json.Unmarshal([]byte(questionIDs), &attempt.QuestionIDs)
		json.Unmarshal([]byte(questionResults), &attempt.QuestionResults)

		attempts = append(attempts, attempt)
	}

	return attempts, rows.Err()
}

// GetQuizHistoryByTopic retrieves quiz attempts for a user filtered by topic
func (db *SQLiteDB) GetQuizHistoryByTopic(userID int, topicID int, page int, perPage int) ([]*models.QuizAttempt, error) {
	offset := (page - 1) * perPage
	query := `
		SELECT DISTINCT qa.id, qa.user_id, qa.quiz_id, qa.parent_attempt_id, qa.attempt_number, qa.is_retake,
			   qa.start_time, qa.completed_time, qa.total_questions, qa.correct_answers,
			   qa.question_ids, qa.question_results
		FROM quiz_attempts qa
		JOIN questions q ON qa.question_ids LIKE '%' || q.topic_question_id || '%'
		WHERE qa.user_id = ? AND qa.completed_time IS NOT NULL AND q.topic_id = ?
		ORDER BY qa.completed_time DESC
		LIMIT ? OFFSET ?
	`

	rows, err := db.db.Query(query, userID, topicID, perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.QuizAttempt
	for rows.Next() {
		attempt := &models.QuizAttempt{}
		var quizID, parentAttemptID sql.NullInt64
		var completedTime sql.NullTime
		var questionIDs, questionResults string

		err := rows.Scan(
			&attempt.ID, &attempt.UserID, &quizID, &parentAttemptID,
			&attempt.AttemptNumber, &attempt.IsRetake, &attempt.StartTime, &completedTime,
			&attempt.TotalQuestions, &attempt.CorrectAnswers, &questionIDs, &questionResults,
		)
		if err != nil {
			return nil, err
		}

		if quizID.Valid {
			id := int(quizID.Int64)
			attempt.QuizID = &id
		}
		if parentAttemptID.Valid {
			id := int(parentAttemptID.Int64)
			attempt.ParentAttemptID = &id
		}
		if completedTime.Valid {
			attempt.CompletedTime = &completedTime.Time
		}

		json.Unmarshal([]byte(questionIDs), &attempt.QuestionIDs)
		json.Unmarshal([]byte(questionResults), &attempt.QuestionResults)

		attempts = append(attempts, attempt)
	}

	return attempts, rows.Err()
}

// GetQuizAttemptWithDetails retrieves a quiz attempt with full question details
func (db *SQLiteDB) GetQuizAttemptWithDetails(attemptID int) (*models.QuizAttempt, error) {
	// First get the quiz attempt
	attempt, err := db.GetQuizAttempt(attemptID)
	if err != nil {
		return nil, err
	}

	// Return the attempt - the API handler will need to populate questions separately
	return attempt, nil
}

// Phase 4: Timed Test Taking Methods

// CreateTestAttempt creates a new test attempt
func (db *SQLiteDB) CreateTestAttempt(attempt *models.TestAttempt) error {
	answersJSON, _ := json.Marshal(attempt.Answers)
	questionResultsJSON, _ := json.Marshal(attempt.QuestionResults)

	query := `
		INSERT INTO test_attempts (user_id, test_id, start_time, submit_time, time_remaining,
			is_auto_submit, total_questions, correct_answers, score, answers, question_results)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.db.Exec(query, attempt.UserID, attempt.TestID, attempt.StartTime,
		nil, 0, false, 0, 0, 0.0, string(answersJSON), string(questionResultsJSON))
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	attempt.ID = int(id)
	return nil
}

// GetTestAttempt retrieves a test attempt by ID
func (db *SQLiteDB) GetTestAttempt(id int) (*models.TestAttempt, error) {
	query := `
		SELECT id, user_id, test_id, start_time, submit_time, time_remaining,
			   is_auto_submit, total_questions, correct_answers, score, answers, question_results
		FROM test_attempts
		WHERE id = ?
	`

	attempt := &models.TestAttempt{}
	var submitTime sql.NullTime
	var answersJSON, questionResultsJSON string

	err := db.db.QueryRow(query, id).Scan(
		&attempt.ID, &attempt.UserID, &attempt.TestID, &attempt.StartTime, &submitTime,
		&attempt.TimeRemaining, &attempt.IsAutoSubmit, &attempt.TotalQuestions,
		&attempt.CorrectAnswers, &attempt.Score, &answersJSON, &questionResultsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if submitTime.Valid {
		attempt.SubmitTime = &submitTime.Time
	}

	json.Unmarshal([]byte(answersJSON), &attempt.Answers)
	json.Unmarshal([]byte(questionResultsJSON), &attempt.QuestionResults)

	return attempt, nil
}

// UpdateTestAttemptAnswer updates an answer in a test attempt
func (db *SQLiteDB) UpdateTestAttemptAnswer(attemptID int, questionID string, answerID int) error {
	// Get current attempt
	attempt, err := db.GetTestAttempt(attemptID)
	if err != nil {
		return err
	}

	// Update answers map
	if attempt.Answers == nil {
		attempt.Answers = make(map[string]int)
	}
	attempt.Answers[questionID] = answerID

	// Save back to database
	answersJSON, _ := json.Marshal(attempt.Answers)
	_, err = db.db.Exec("UPDATE test_attempts SET answers = ? WHERE id = ?", string(answersJSON), attemptID)
	return err
}

// SubmitTestAttempt submits a test attempt and calculates the score
func (db *SQLiteDB) SubmitTestAttempt(attemptID int, isAutoSubmit bool) error {
	attempt, err := db.GetTestAttempt(attemptID)
	if err != nil {
		return err
	}

	now := time.Now()
	attempt.SubmitTime = &now
	attempt.IsAutoSubmit = isAutoSubmit

	// Get test details
	test, err := db.GetTest(attempt.TestID)
	if err != nil {
		return err
	}

	// Calculate time remaining
	timeLimitSeconds := test.TimeLimit * 60
	elapsed := int(now.Sub(attempt.StartTime).Seconds())
	attempt.TimeRemaining = timeLimitSeconds - elapsed
	if attempt.TimeRemaining < 0 {
		attempt.TimeRemaining = 0
	}

	// Get test questions
	testQuestions, err := db.GetTestQuestions(attempt.TestID)
	if err != nil {
		return err
	}

	attempt.TotalQuestions = len(testQuestions)
	attempt.CorrectAnswers = 0
	attempt.QuestionResults = make(map[string]bool)

	// Check each answer
	for _, q := range testQuestions {
		questionKey := fmt.Sprintf("%d", q.ID)
		answerID, hasAnswer := attempt.Answers[questionKey]

		if hasAnswer {
			answer, err := db.GetAnswer(answerID)
			if err == nil && answer.QuestionID == q.ID && answer.IsCorrect {
				attempt.CorrectAnswers++
				attempt.QuestionResults[questionKey] = true
			} else {
				attempt.QuestionResults[questionKey] = false
			}
		} else {
			attempt.QuestionResults[questionKey] = false
		}
	}

	// Calculate score
	if attempt.TotalQuestions > 0 {
		attempt.Score = float64(attempt.CorrectAnswers) / float64(attempt.TotalQuestions) * 100
	}

	// Update database
	questionResultsJSON, _ := json.Marshal(attempt.QuestionResults)
	query := `
		UPDATE test_attempts
		SET submit_time = ?, time_remaining = ?, is_auto_submit = ?,
			total_questions = ?, correct_answers = ?, score = ?, question_results = ?
		WHERE id = ?
	`
	_, err = db.db.Exec(query, attempt.SubmitTime, attempt.TimeRemaining, attempt.IsAutoSubmit,
		attempt.TotalQuestions, attempt.CorrectAnswers, attempt.Score,
		string(questionResultsJSON), attemptID)
	return err
}

// GetTestResults retrieves the results of a test attempt
func (db *SQLiteDB) GetTestResults(attemptID int) (*models.TestAttempt, error) {
	return db.GetTestAttempt(attemptID)
}

// GetAvailableTests retrieves available tests for a student
func (db *SQLiteDB) GetAvailableTests(userID int) ([]*models.Test, error) {
	now := time.Now()
	query := `
		SELECT id, name, topic_id, description, time_limit, start_date, end_date,
			   is_published, created_at, updated_at
		FROM tests
		WHERE is_published = 1
		  AND (start_date IS NULL OR start_date <= ?)
		  AND (end_date IS NULL OR end_date >= ?)
	`

	rows, err := db.db.Query(query, now, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*models.Test
	for rows.Next() {
		test := &models.Test{}
		var startDate, endDate sql.NullTime

		err := rows.Scan(
			&test.ID, &test.Name, &test.TopicID, &test.Description, &test.TimeLimit,
			&startDate, &endDate, &test.IsPublished, &test.CreatedAt, &test.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if startDate.Valid {
			test.StartDate = &startDate.Time
		}
		if endDate.Valid {
			test.EndDate = &endDate.Time
		}

		tests = append(tests, test)
	}

	return tests, rows.Err()
}

// GetTestAttemptsByUser retrieves all test attempts for a user for a specific test
func (db *SQLiteDB) GetTestAttemptsByUser(userID int, testID int) ([]*models.TestAttempt, error) {
	query := `
		SELECT id, user_id, test_id, start_time, submit_time, time_remaining,
			   is_auto_submit, total_questions, correct_answers, score, answers, question_results
		FROM test_attempts
		WHERE user_id = ? AND test_id = ?
		ORDER BY start_time DESC
	`

	rows, err := db.db.Query(query, userID, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attempts []*models.TestAttempt
	for rows.Next() {
		attempt := &models.TestAttempt{}
		var submitTime sql.NullTime
		var answersJSON, questionResultsJSON string

		err := rows.Scan(
			&attempt.ID, &attempt.UserID, &attempt.TestID, &attempt.StartTime, &submitTime,
			&attempt.TimeRemaining, &attempt.IsAutoSubmit, &attempt.TotalQuestions,
			&attempt.CorrectAnswers, &attempt.Score, &answersJSON, &questionResultsJSON,
		)
		if err != nil {
			return nil, err
		}

		if submitTime.Valid {
			attempt.SubmitTime = &submitTime.Time
		}

		json.Unmarshal([]byte(answersJSON), &attempt.Answers)
		json.Unmarshal([]byte(questionResultsJSON), &attempt.QuestionResults)

		attempts = append(attempts, attempt)
	}

	return attempts, rows.Err()
}

// Phase 4: Tag-Based Review Methods

// GetAllTags retrieves all unique tags from all questions
func (db *SQLiteDB) GetAllTags() ([]string, error) {
	rows, err := db.db.Query("SELECT DISTINCT tags FROM questions WHERE tags IS NOT NULL AND tags != '[]'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

// GetTagsWithCounts retrieves all tags with their question counts
func (db *SQLiteDB) GetTagsWithCounts() (map[string]int, error) {
	rows, err := db.db.Query("SELECT tags FROM questions WHERE tags IS NOT NULL AND tags != '[]'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			tagCounts[tag]++
		}
	}

	return tagCounts, rows.Err()
}

// GetTagStats retrieves accuracy statistics for each tag for a user
func (db *SQLiteDB) GetTagStats(userID int) ([]*models.TagStats, error) {
	query := `
		SELECT q.tags, qa.correct_attempts, qa.incorrect_attempts
		FROM question_attempts qa
		JOIN questions q ON qa.question_id = q.id
		WHERE qa.user_id = ? AND q.tags IS NOT NULL AND q.tags != '[]'
	`

	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagStats := make(map[string]*models.TagStats)
	for rows.Next() {
		var tagsJSON string
		var correctAttempts, incorrectAttempts int

		if err := rows.Scan(&tagsJSON, &correctAttempts, &incorrectAttempts); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			if _, exists := tagStats[tag]; !exists {
				tagStats[tag] = &models.TagStats{Tag: tag}
			}

			stats := tagStats[tag]
			stats.CorrectAttempts += correctAttempts
			stats.TotalAttempts += correctAttempts + incorrectAttempts
		}
	}

	// Get question counts
	tagCounts, _ := db.GetTagsWithCounts()

	// Calculate accuracy
	result := make([]*models.TagStats, 0, len(tagStats))
	for tag, stats := range tagStats {
		stats.QuestionCount = tagCounts[tag]
		if stats.TotalAttempts > 0 {
			stats.Accuracy = float64(stats.CorrectAttempts) / float64(stats.TotalAttempts) * 100
		}
		result = append(result, stats)
	}

	return result, rows.Err()
}

// GetQuestionsByTags retrieves questions that match the given tags
func (db *SQLiteDB) GetQuestionsByTags(tags []string, useAnd bool) ([]*models.Question, error) {
	rows, err := db.db.Query("SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions WHERE tags IS NOT NULL AND tags != '[]'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		q := &models.Question{}
		var tagsJSON string

		err := rows.Scan(&q.ID, &q.TopicID, &q.QuestionText, &q.Difficulty, &tagsJSON, &q.TopicQuestionID, &q.CreatedAt)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(tagsJSON), &q.Tags)

		// Check if question matches tags
		if useAnd {
			hasAllTags := true
			for _, tag := range tags {
				found := false
				for _, qTag := range q.Tags {
					if qTag == tag {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if hasAllTags {
				questions = append(questions, q)
			}
		} else {
			hasAnyTag := false
			for _, tag := range tags {
				for _, qTag := range q.Tags {
					if qTag == tag {
						hasAnyTag = true
						break
					}
				}
				if hasAnyTag {
					break
				}
			}
			if hasAnyTag {
				questions = append(questions, q)
			}
		}
	}

	return questions, rows.Err()
}

// Phase 4: Test Management (Instructor) Methods

// PublishTest sets a test as published
func (db *SQLiteDB) PublishTest(testID int) error {
	_, err := db.db.Exec("UPDATE tests SET is_published = 1, updated_at = ? WHERE id = ?", time.Now(), testID)
	return err
}

// GetTestResultsByTest retrieves all test results for a specific test
func (db *SQLiteDB) GetTestResultsByTest(testID int) (*models.TestResult, error) {
	test, err := db.GetTest(testID)
	if err != nil {
		return nil, err
	}

	result := &models.TestResult{
		TestID:         testID,
		TestName:       test.Name,
		StudentResults: make([]models.StudentTestResult, 0),
	}

	query := `
		SELECT ta.user_id, u.username, ta.score, ta.correct_answers, ta.total_questions,
			   ta.submit_time, ta.time_remaining
		FROM test_attempts ta
		JOIN users u ON ta.user_id = u.id
		WHERE ta.test_id = ? AND ta.submit_time IS NOT NULL
		ORDER BY ta.submit_time DESC
	`

	rows, err := db.db.Query(query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		sr := models.StudentTestResult{}
		var submitTime time.Time
		var timeRemaining int

		err := rows.Scan(&sr.UserID, &sr.Username, &sr.Score, &sr.CorrectAnswers,
			&sr.TotalQuestions, &submitTime, &timeRemaining)
		if err != nil {
			continue
		}

		sr.SubmitTime = submitTime
		if test.TimeLimit > 0 {
			sr.TimeTaken = (test.TimeLimit * 60) - timeRemaining
		}

		result.StudentResults = append(result.StudentResults, sr)
	}

	return result, rows.Err()
}

// ExportTestResults exports test results as CSV
func (db *SQLiteDB) ExportTestResults(testID int) ([]byte, error) {
	results, err := db.GetTestResultsByTest(testID)
	if err != nil {
		return nil, err
	}

	var csv strings.Builder
	csv.WriteString("Username,Score,Correct Answers,Total Questions,Submit Time,Time Taken (seconds)\n")

	for _, sr := range results.StudentResults {
		csv.WriteString(fmt.Sprintf("%s,%.2f,%d,%d,%s,%d\n",
			sr.Username,
			sr.Score,
			sr.CorrectAnswers,
			sr.TotalQuestions,
			sr.SubmitTime.Format(time.RFC3339),
			sr.TimeTaken,
		))
	}

	return []byte(csv.String()), nil
}

// Phase 4: Question Performance Analytics Methods

// GetQuestionAnalytics retrieves analytics for all questions
func (db *SQLiteDB) GetQuestionAnalytics() ([]*models.QuestionAnalytics, error) {
	query := `
		SELECT q.id, q.topic_question_id, q.question_text, q.topic_id, t.name,
			   q.tags, q.difficulty,
			   COALESCE(SUM(qa.correct_attempts + qa.incorrect_attempts), 0) as total_attempts,
			   COALESCE(SUM(qa.correct_attempts), 0) as correct_attempts
		FROM questions q
		LEFT JOIN topics t ON q.topic_id = t.id
		LEFT JOIN question_attempts qa ON q.id = qa.question_id
		GROUP BY q.id
	`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*models.QuestionAnalytics
	for rows.Next() {
		a := &models.QuestionAnalytics{
			AnswerDistribution: make(map[int]int),
		}
		var tagsJSON string

		err := rows.Scan(
			&a.QuestionID, &a.TopicQuestionID, &a.QuestionText, &a.TopicID, &a.TopicName,
			&tagsJSON, &a.CurrentDifficulty, &a.TotalAttempts, &a.CorrectAttempts,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(tagsJSON), &a.Tags)

		// Calculate accuracy
		if a.TotalAttempts > 0 {
			a.Accuracy = float64(a.CorrectAttempts) / float64(a.TotalAttempts) * 100
		}

		// Calculate recommended difficulty
		if a.Accuracy > 90 {
			a.RecommendedDifficulty = a.CurrentDifficulty + 1
			if a.RecommendedDifficulty > 7 {
				a.RecommendedDifficulty = 7
			}
			a.FlagForReview = true
		} else if a.Accuracy < 30 {
			a.RecommendedDifficulty = a.CurrentDifficulty - 1
			if a.RecommendedDifficulty < 1 {
				a.RecommendedDifficulty = 1
			}
			a.FlagForReview = true
		} else {
			a.RecommendedDifficulty = a.CurrentDifficulty
		}

		analytics = append(analytics, a)
	}

	return analytics, rows.Err()
}

// GetQuestionDetailedAnalytics retrieves detailed analytics for a single question
func (db *SQLiteDB) GetQuestionDetailedAnalytics(questionID int) (*models.QuestionAnalytics, error) {
	query := `
		SELECT q.id, q.topic_question_id, q.question_text, q.topic_id, t.name,
			   q.tags, q.difficulty,
			   COALESCE(SUM(qa.correct_attempts + qa.incorrect_attempts), 0) as total_attempts,
			   COALESCE(SUM(qa.correct_attempts), 0) as correct_attempts
		FROM questions q
		LEFT JOIN topics t ON q.topic_id = t.id
		LEFT JOIN question_attempts qa ON q.id = qa.question_id
		WHERE q.id = ?
		GROUP BY q.id
	`

	a := &models.QuestionAnalytics{
		AnswerDistribution: make(map[int]int),
	}
	var tagsJSON string

	err := db.db.QueryRow(query, questionID).Scan(
		&a.QuestionID, &a.TopicQuestionID, &a.QuestionText, &a.TopicID, &a.TopicName,
		&tagsJSON, &a.CurrentDifficulty, &a.TotalAttempts, &a.CorrectAttempts,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(tagsJSON), &a.Tags)

	// Calculate accuracy
	if a.TotalAttempts > 0 {
		a.Accuracy = float64(a.CorrectAttempts) / float64(a.TotalAttempts) * 100
	}

	// Calculate recommended difficulty
	if a.Accuracy > 90 {
		a.RecommendedDifficulty = a.CurrentDifficulty + 1
		if a.RecommendedDifficulty > 7 {
			a.RecommendedDifficulty = 7
		}
		a.FlagForReview = true
	} else if a.Accuracy < 30 {
		a.RecommendedDifficulty = a.CurrentDifficulty - 1
		if a.RecommendedDifficulty < 1 {
			a.RecommendedDifficulty = 1
		}
		a.FlagForReview = true
	} else {
		a.RecommendedDifficulty = a.CurrentDifficulty
	}

	// Get answer distribution
	a.AnswerDistribution, _ = db.GetAnswerDistribution(questionID)

	return a, nil
}

// GetAnswerDistribution retrieves distribution of answers chosen for a question
func (db *SQLiteDB) GetAnswerDistribution(questionID int) (map[int]int, error) {
	// This would require tracking answer selections in quiz/test attempts
	// For now, return empty distribution
	return make(map[int]int), nil
}

// UpdateQuestionDifficultyFromAnalytics updates a question's difficulty
func (db *SQLiteDB) UpdateQuestionDifficultyFromAnalytics(questionID int, newDifficulty int) error {
	if newDifficulty < 1 || newDifficulty > 7 {
		return fmt.Errorf("difficulty must be between 1 and 7")
	}

	_, err := db.db.Exec("UPDATE questions SET difficulty = ? WHERE id = ?", newDifficulty, questionID)
	return err
}
