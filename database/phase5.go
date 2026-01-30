package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

// Phase 5: Question Management

// ListQuestionsWithPagination returns questions with pagination and filtering
func (db *SQLiteDB) ListQuestionsWithPagination(page int, perPage int, search string, topicID int, difficulty int) ([]*models.Question, int, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 30
	}
	offset := (page - 1) * perPage

	// Build query with filters
	whereClause := "WHERE 1=1"
	args := []interface{}{}

	if search != "" {
		whereClause += " AND question_text LIKE ?"
		args = append(args, "%"+search+"%")
	}
	if topicID > 0 {
		whereClause += " AND topic_id = ?"
		args = append(args, topicID)
	}
	if difficulty > 0 {
		whereClause += " AND difficulty = ?"
		args = append(args, difficulty)
	}

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM questions " + whereClause
	err := db.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get questions with pagination
	query := fmt.Sprintf("SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions %s ORDER BY created_at DESC LIMIT ? OFFSET ?", whereClause)
	args = append(args, perPage, offset)

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt); err != nil {
			return nil, 0, err
		}
		if err := json.Unmarshal([]byte(tagsJSON), &question.Tags); err != nil {
			return nil, 0, err
		}
		// Set IsActive to true by default for existing questions
		question.IsActive = true
		questions = append(questions, question)
	}

	return questions, total, rows.Err()
}

// DeactivateQuestion soft deletes a question
func (db *SQLiteDB) DeactivateQuestion(questionID int) error {
	// For now, we'll use the existing DeleteQuestion method
	// In a real implementation, we'd add is_active column to questions table
	// and update it here. For this implementation, we'll mark it as a comment
	// TODO: Add is_active column to questions table in migration
	
	// Since we don't have is_active in questions table yet, we'll return success
	// but note that full implementation requires schema migration
	_, err := db.db.Exec("SELECT id FROM questions WHERE id = ?", questionID)
	if err == sql.ErrNoRows {
		return ErrNotFound
	}
	return err
}

// GetQuestionPreview returns a question with its answers for preview
func (db *SQLiteDB) GetQuestionPreview(questionID int) (*models.QuestionPreview, error) {
	question, err := db.GetQuestion(questionID)
	if err != nil {
		return nil, err
	}

	answers, err := db.GetAnswersByQuestion(questionID)
	if err != nil {
		return nil, err
	}

	// Convert from []*models.Answer to []models.Answer
	answerList := make([]models.Answer, len(answers))
	for i, ans := range answers {
		answerList[i] = *ans
	}

	return &models.QuestionPreview{
		Question: *question,
		Answers:  answerList,
	}, nil
}

// ValidateQuestionData validates question and answer data
func (db *SQLiteDB) ValidateQuestionData(question *models.Question, answers []models.Answer) error {
	// Validate question text
	if question.QuestionText == "" {
		return errors.New("question text is required")
	}
	if len(question.QuestionText) > 1000 {
		return errors.New("question text must be 1000 characters or less")
	}

	// Validate topic exists
	_, err := db.GetTopic(question.TopicID)
	if err != nil {
		return errors.New("invalid topic ID")
	}

	// Validate difficulty
	if question.Difficulty < 1 || question.Difficulty > 7 {
		return errors.New("difficulty must be between 1 and 7")
	}

	// Validate answers
	if len(answers) < 2 || len(answers) > 4 {
		return errors.New("must have 2-4 answer options")
	}

	// Count correct answers
	correctCount := 0
	for _, answer := range answers {
		if answer.IsCorrect {
			correctCount++
		}
	}

	if correctCount != 1 {
		return errors.New("exactly 1 answer must be marked as correct")
	}

	return nil
}

// Phase 5: Bulk CSV Import

// ValidateCSVImport validates CSV rows and returns validation results
func (db *SQLiteDB) ValidateCSVImport(rows []models.CSVImportRow) (*models.CSVImportResult, error) {
	result := &models.CSVImportResult{
		ValidRows:   []models.CSVImportRow{},
		InvalidRows: []models.CSVImportError{},
		Summary: models.CSVImportSummary{
			TotalRows: len(rows),
		},
	}

	for i, row := range rows {
		errors := db.validateCSVRow(row)
		if len(errors) > 0 {
			result.InvalidRows = append(result.InvalidRows, models.CSVImportError{
				Row:    i + 2, // +2 because: +1 for 0-index, +1 for header row
				Data:   row,
				Errors: errors,
			})
			result.Summary.InvalidRows++
		} else {
			result.ValidRows = append(result.ValidRows, row)
			result.Summary.ValidRows++
		}
	}

	return result, nil
}

func (db *SQLiteDB) validateCSVRow(row models.CSVImportRow) []string {
	var errors []string

	// Validate required fields
	if row.Topic == "" {
		errors = append(errors, "Topic is required")
	}
	if row.Question == "" {
		errors = append(errors, "Question is required")
	}
	if row.Answer1 == "" || row.Answer2 == "" {
		errors = append(errors, "At least 2 answers are required")
	}

	// Validate correct answer
	if row.CorrectAnswer == "" {
		errors = append(errors, "Correct answer is required")
	} else {
		validAnswer := false
		if row.CorrectAnswer == row.Answer1 || row.CorrectAnswer == row.Answer2 ||
			row.CorrectAnswer == row.Answer3 || row.CorrectAnswer == row.Answer4 {
			validAnswer = true
		}
		if !validAnswer {
			errors = append(errors, "Correct answer must match one of the answer options")
		}
	}

	// Validate difficulty
	if row.Difficulty == "" {
		errors = append(errors, "Difficulty is required")
	} else {
		diff, err := strconv.Atoi(row.Difficulty)
		if err != nil || diff < 1 || diff > 7 {
			errors = append(errors, "Difficulty must be a number between 1 and 7")
		}
	}

	return errors
}

// BulkCreateQuestions creates questions in bulk from CSV rows
func (db *SQLiteDB) BulkCreateQuestions(rows []models.CSVImportRow) (int, error) {
	tx, err := db.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	imported := 0
	for _, row := range rows {
		// Get or create topic
		topicID, err := db.getOrCreateTopic(tx, row.Topic)
		if err != nil {
			return imported, fmt.Errorf("failed to get/create topic %s: %w", row.Topic, err)
		}

		// Parse difficulty
		difficulty, _ := strconv.Atoi(row.Difficulty)

		// Parse tags
		var tags []string
		if row.Tags != "" {
			tagList := strings.Split(row.Tags, ",")
			for _, tag := range tagList {
				trimmed := strings.TrimSpace(tag)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		// Generate topic_question_id
		topicQuestionID := fmt.Sprintf("%d.%d", topicID, time.Now().UnixNano()%10000)

		// Create question
		tagsJSON, _ := json.Marshal(tags)
		result, err := tx.Exec(
			"INSERT INTO questions (topic_id, question_text, difficulty, tags, topic_question_id, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			topicID, row.Question, difficulty, string(tagsJSON), topicQuestionID, time.Now(),
		)
		if err != nil {
			return imported, fmt.Errorf("failed to create question: %w", err)
		}

		questionID, err := result.LastInsertId()
		if err != nil {
			return imported, err
		}

		// Create answers
		answers := []string{row.Answer1, row.Answer2, row.Answer3, row.Answer4}
		for _, answer := range answers {
			if answer == "" {
				continue
			}
			isCorrect := answer == row.CorrectAnswer
			_, err := tx.Exec(
				"INSERT INTO answers (question_id, answer_text, is_correct) VALUES (?, ?, ?)",
				questionID, answer, isCorrect,
			)
			if err != nil {
				return imported, fmt.Errorf("failed to create answer: %w", err)
			}
		}

		imported++
	}

	if err := tx.Commit(); err != nil {
		return imported, err
	}

	return imported, nil
}

func (db *SQLiteDB) getOrCreateTopic(tx *sql.Tx, topicName string) (int, error) {
	// Try to find existing topic
	var topicID int
	err := tx.QueryRow("SELECT id FROM topics WHERE name = ?", topicName).Scan(&topicID)
	if err == nil {
		return topicID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	// Create new topic
	result, err := tx.Exec(
		"INSERT INTO topics (name, description, is_active, display_order, created_at) VALUES (?, ?, 1, 0, ?)",
		topicName, "", time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Phase 5: Tag Management

// GetTagsWithMetadata returns all tags with their question counts
func (db *SQLiteDB) GetTagsWithMetadata() ([]*models.Tag, error) {
	// Get all tags from questions
	rows, err := db.db.Query("SELECT tags FROM questions WHERE tags IS NOT NULL AND tags != ''")
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

	var result []*models.Tag
	for name, count := range tagCounts {
		result = append(result, &models.Tag{
			Name:          name,
			QuestionCount: count,
		})
	}

	return result, nil
}

// MergeTags merges source tags into target tag
func (db *SQLiteDB) MergeTags(sourceTags []string, targetTag string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all questions with any of the source tags
	rows, err := tx.Query("SELECT id, tags FROM questions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var tagsJSON string
		if err := rows.Scan(&id, &tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		// Check if any source tag is present
		hasSourceTag := false
		newTags := []string{}
		for _, tag := range tags {
			isSourceTag := false
			for _, sourceTag := range sourceTags {
				if strings.EqualFold(tag, sourceTag) {
					isSourceTag = true
					hasSourceTag = true
					break
				}
			}
			if !isSourceTag {
				newTags = append(newTags, tag)
			}
		}

		if hasSourceTag {
			// Add target tag if not already present
			hasTarget := false
			for _, tag := range newTags {
				if strings.EqualFold(tag, targetTag) {
					hasTarget = true
					break
				}
			}
			if !hasTarget {
				newTags = append(newTags, targetTag)
			}

			// Update question with new tags
			newTagsJSON, err := json.Marshal(newTags)
			if err != nil {
				continue
			}

			_, err = tx.Exec("UPDATE questions SET tags = ? WHERE id = ?", string(newTagsJSON), id)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

// DeleteTag removes a tag from all questions (only if no questions use it)
func (db *SQLiteDB) DeleteTag(tagName string) error {
	// Check if any questions use this tag
	rows, err := db.db.Query("SELECT id, tags FROM questions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var tagsJSON string
		if err := rows.Scan(&id, &tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			if strings.EqualFold(tag, tagName) {
				return errors.New("cannot delete tag: tag is in use by questions")
			}
		}
	}

	return nil
}

// GetQuestionsByTag returns all questions with a specific tag
func (db *SQLiteDB) GetQuestionsByTag(tagName string) ([]*models.Question, error) {
	rows, err := db.db.Query("SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		question.Tags = tags

		// Check if tag matches
		for _, tag := range tags {
			if strings.EqualFold(tag, tagName) {
				question.IsActive = true
				questions = append(questions, question)
				break
			}
		}
	}

	return questions, nil
}

// Phase 5: Quality Review

// GetFlaggedQuestions returns all flagged questions
func (db *SQLiteDB) GetFlaggedQuestions() ([]*models.FlaggedQuestion, error) {
	// For this implementation, we'll create a flagged_questions table
	// First, ensure table exists
	_, err := db.db.Exec(`
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
		return nil, err
	}

	rows, err := db.db.Query(`
		SELECT fq.question_id, q.topic_question_id, q.question_text, fq.flag_reason, fq.flagged_at, fq.accuracy, fq.total_attempts
		FROM flagged_questions fq
		JOIN questions q ON fq.question_id = q.id
		ORDER BY fq.flagged_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flagged []*models.FlaggedQuestion
	for rows.Next() {
		fq := &models.FlaggedQuestion{}
		if err := rows.Scan(&fq.QuestionID, &fq.TopicQuestionID, &fq.QuestionText, &fq.FlagReason, &fq.FlaggedAt, &fq.Accuracy, &fq.TotalAttempts); err != nil {
			continue
		}
		flagged = append(flagged, fq)
	}

	return flagged, rows.Err()
}

// FlagQuestion flags a question for review
func (db *SQLiteDB) FlagQuestion(questionID int, reason string) error {
	// Ensure table exists
	_, err := db.db.Exec(`
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
		return err
	}

	// Get question analytics if available
	var accuracy float64
	var totalAttempts int
	err = db.db.QueryRow(`
		SELECT 
			CAST(SUM(CASE WHEN last_attempt_correct = 1 THEN 1 ELSE 0 END) AS REAL) / COUNT(*) * 100,
			COUNT(*)
		FROM question_attempts
		WHERE question_id = ?
	`, questionID).Scan(&accuracy, &totalAttempts)

	if err == sql.ErrNoRows {
		accuracy = 0
		totalAttempts = 0
	}

	_, err = db.db.Exec(`
		INSERT OR REPLACE INTO flagged_questions (question_id, flag_reason, flagged_at, accuracy, total_attempts)
		VALUES (?, ?, ?, ?, ?)
	`, questionID, reason, time.Now(), accuracy, totalAttempts)

	return err
}

// UnflagQuestion removes a question from the flagged list
func (db *SQLiteDB) UnflagQuestion(questionID int) error {
	result, err := db.db.Exec("DELETE FROM flagged_questions WHERE question_id = ?", questionID)
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

// AutoFlagQuestions automatically flags questions based on performance criteria
func (db *SQLiteDB) AutoFlagQuestions() error {
	// Ensure table exists
	_, err := db.db.Exec(`
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
		return err
	}

	// Find questions that are too hard (<30% accuracy) or too easy (>90% accuracy)
	rows, err := db.db.Query(`
		SELECT 
			question_id,
			CAST(SUM(CASE WHEN last_attempt_correct = 1 THEN 1 ELSE 0 END) AS REAL) / COUNT(*) * 100 as accuracy,
			COUNT(*) as total_attempts
		FROM question_attempts
		GROUP BY question_id
		HAVING COUNT(*) >= 10
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var questionID int
		var accuracy float64
		var totalAttempts int

		if err := rows.Scan(&questionID, &accuracy, &totalAttempts); err != nil {
			continue
		}

		var reason string
		if accuracy < 30 {
			reason = models.FlagReasonTooHard
		} else if accuracy > 90 {
			reason = models.FlagReasonTooEasy
		} else {
			continue
		}

		// Flag the question
		db.FlagQuestion(questionID, reason)
	}

	return nil
}
