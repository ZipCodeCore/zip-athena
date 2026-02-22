package database

import (
	"database/sql"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

// AssignQuizToUser assigns a quiz to a specific user
func (db *SQLiteDB) AssignQuizToUser(quizID, userID, assignedBy int) error {
	_, err := db.db.Exec(
		"INSERT OR IGNORE INTO assigned_quizzes (quiz_id, user_id, assigned_by, assigned_at) VALUES (?, ?, ?, ?)",
		quizID, userID, assignedBy, time.Now(),
	)
	return err
}

// UnassignQuizFromUser removes a quiz assignment from a user
func (db *SQLiteDB) UnassignQuizFromUser(quizID, userID int) error {
	result, err := db.db.Exec(
		"DELETE FROM assigned_quizzes WHERE quiz_id = ? AND user_id = ?",
		quizID, userID,
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

// GetAssignedQuizzesForUser returns all quizzes assigned to a user, with completion info
func (db *SQLiteDB) GetAssignedQuizzesForUser(userID int) ([]*models.AssignedQuizInfo, error) {
	query := `
		SELECT aq.id, aq.quiz_id, aq.user_id, aq.assigned_at, aq.assigned_by,
		       q.name, COALESCE(q.description, ''), COALESCE(q.category, ''),
		       qa.completed_time, qa.correct_answers, qa.total_questions
		FROM assigned_quizzes aq
		JOIN quizzes q ON aq.quiz_id = q.id
		LEFT JOIN quiz_attempts qa ON qa.quiz_id = q.id
		    AND qa.user_id = aq.user_id
		    AND qa.completed_time IS NOT NULL
		    AND qa.id = (
		        SELECT MAX(qa2.id) FROM quiz_attempts qa2
		        WHERE qa2.quiz_id = q.id AND qa2.user_id = aq.user_id AND qa2.completed_time IS NOT NULL
		    )
		WHERE aq.user_id = ?
		ORDER BY aq.assigned_at DESC
	`
	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.AssignedQuizInfo
	for rows.Next() {
		info := &models.AssignedQuizInfo{}
		var completedTime sql.NullTime
		var correctAnswers, totalQuestions sql.NullInt64

		err := rows.Scan(
			&info.ID, &info.QuizID, &info.UserID, &info.AssignedAt, &info.AssignedBy,
			&info.QuizName, &info.QuizDescription, &info.QuizCategory,
			&completedTime, &correctAnswers, &totalQuestions,
		)
		if err != nil {
			return nil, err
		}

		if completedTime.Valid {
			info.CompletedAt = &completedTime.Time
			if correctAnswers.Valid {
				ca := int(correctAnswers.Int64)
				info.CorrectAnswers = &ca
			}
			if totalQuestions.Valid {
				tq := int(totalQuestions.Int64)
				info.TotalQuestions = &tq
				if tq > 0 && correctAnswers.Valid {
					score := float64(correctAnswers.Int64) / float64(tq) * 100.0
					info.Score = &score
				}
			}
		}

		result = append(result, info)
	}
	return result, rows.Err()
}

// GetAssignedQuizUserIDs returns all user IDs assigned to a quiz
func (db *SQLiteDB) GetAssignedQuizUserIDs(quizID int) ([]int, error) {
	rows, err := db.db.Query(
		"SELECT user_id FROM assigned_quizzes WHERE quiz_id = ?",
		quizID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, rows.Err()
}

// GetAssignedQuizStatus returns the status of a quiz with all student results
func (db *SQLiteDB) GetAssignedQuizStatus(quizID int) (*models.AssignedQuizStatus, error) {
	// Get quiz info
	quiz, err := db.GetQuiz(quizID)
	if err != nil {
		return nil, err
	}

	status := &models.AssignedQuizStatus{
		QuizID:          quiz.ID,
		QuizName:        quiz.Name,
		QuizCategory:    quiz.Category,
		QuizDescription: quiz.Description,
		StudentResults:  []models.AssignedQuizStudentResult{},
	}

	query := `
		SELECT u.id, u.username, aq.assigned_at,
		       qa.completed_time, qa.correct_answers, qa.total_questions
		FROM assigned_quizzes aq
		JOIN users u ON aq.user_id = u.id
		LEFT JOIN quiz_attempts qa ON qa.quiz_id = aq.quiz_id
		    AND qa.user_id = aq.user_id
		    AND qa.completed_time IS NOT NULL
		    AND qa.id = (
		        SELECT MAX(qa2.id) FROM quiz_attempts qa2
		        WHERE qa2.quiz_id = aq.quiz_id AND qa2.user_id = aq.user_id AND qa2.completed_time IS NOT NULL
		    )
		WHERE aq.quiz_id = ?
		ORDER BY u.username
	`
	rows, err := db.db.Query(query, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := models.AssignedQuizStudentResult{}
		var completedTime sql.NullTime
		var correctAnswers, totalQuestions sql.NullInt64

		err := rows.Scan(
			&r.UserID, &r.Username, &r.AssignedAt,
			&completedTime, &correctAnswers, &totalQuestions,
		)
		if err != nil {
			return nil, err
		}

		if completedTime.Valid {
			r.CompletedAt = &completedTime.Time
			if correctAnswers.Valid {
				ca := int(correctAnswers.Int64)
				r.CorrectAnswers = &ca
			}
			if totalQuestions.Valid {
				tq := int(totalQuestions.Int64)
				r.TotalQuestions = &tq
				if tq > 0 && correctAnswers.Valid {
					score := float64(correctAnswers.Int64) / float64(tq) * 100.0
					r.Score = &score
				}
			}
		}

		status.StudentResults = append(status.StudentResults, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return status, nil
}

// ListAssignedQuizzes returns all quiz assignments
func (db *SQLiteDB) ListAssignedQuizzes() ([]*models.AssignedQuiz, error) {
	rows, err := db.db.Query(
		"SELECT id, quiz_id, user_id, assigned_at, assigned_by FROM assigned_quizzes ORDER BY assigned_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []*models.AssignedQuiz
	for rows.Next() {
		a := &models.AssignedQuiz{}
		if err := rows.Scan(&a.ID, &a.QuizID, &a.UserID, &a.AssignedAt, &a.AssignedBy); err != nil {
			return nil, err
		}
		assignments = append(assignments, a)
	}
	return assignments, rows.Err()
}
