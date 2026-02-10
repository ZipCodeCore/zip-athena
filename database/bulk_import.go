package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

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
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			// Fallback to an empty JSON array if tags cannot be marshaled
			tagsJSON = []byte("[]")
		}
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

