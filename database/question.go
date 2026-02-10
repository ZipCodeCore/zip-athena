package database

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/ZipCodeCore/zip-athena/models"
)

// CreateQuestion creates a new question in the database
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

// GetQuestion retrieves a question by ID
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

// GetQuestionsByTopic retrieves all questions for a specific topic
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

// ListQuestions retrieves all questions
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

// UpdateQuestion updates a question's information
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

// DeleteQuestion deletes a question by ID
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

// DeactivateQuestion soft deletes a question (currently delegates to hard delete)
// Once an is_active column is added to the questions table, this should be
// updated to perform a true soft delete instead of a hard delete.
func (db *SQLiteDB) DeactivateQuestion(questionID int) error {
	// For now, delegate to the existing DeleteQuestion implementation
	return db.DeleteQuestion(questionID)
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
	query := "SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions " + whereClause + " ORDER BY created_at DESC LIMIT ? OFFSET ?"
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
