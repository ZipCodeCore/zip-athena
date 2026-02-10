package database

import (
"database/sql"

"github.com/ZipCodeCore/zip-athena/models"
)

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

