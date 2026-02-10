package database

import (
	"database/sql"

	"github.com/ZipCodeCore/zip-athena/models"
)

// CreateAnswer creates a new answer for a question
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

// GetAnswersByQuestion retrieves all answers for a specific question
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

// UpdateAnswer updates an answer's information
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

// DeleteAnswer deletes an answer by ID
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
