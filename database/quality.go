package database

import (
	"database/sql"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

func (db *SQLiteDB) GetFlaggedQuestions() ([]*models.FlaggedQuestion, error) {
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

func (db *SQLiteDB) FlagQuestion(questionID int, reason string) error {
	// Get question analytics if available
	var accuracy float64
	var totalAttempts int
	err := db.db.QueryRow(`
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

func (db *SQLiteDB) AutoFlagQuestions() error {
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

