package database

import (
	"database/sql"
	"encoding/json"

	"github.com/ZipCodeCore/zip-athena/models"
)

func (db *SQLiteDB) CreateQuiz(quiz *models.Quiz) error {
	result, err := db.db.Exec(
		"INSERT INTO quizzes (name, topic_id, category, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		quiz.Name, quiz.TopicID, quiz.Category, quiz.Description, quiz.CreatedAt, quiz.UpdatedAt,
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
		"SELECT id, name, topic_id, COALESCE(category, ''), description, created_at, updated_at FROM quizzes WHERE id = ?",
		id,
	).Scan(&quiz.ID, &quiz.Name, &quiz.TopicID, &quiz.Category, &quiz.Description, &quiz.CreatedAt, &quiz.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return quiz, nil
}

func (db *SQLiteDB) ListQuizzes() ([]*models.Quiz, error) {
	rows, err := db.db.Query("SELECT id, name, topic_id, COALESCE(category, ''), description, created_at, updated_at FROM quizzes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quizzes []*models.Quiz
	for rows.Next() {
		quiz := &models.Quiz{}
		if err := rows.Scan(&quiz.ID, &quiz.Name, &quiz.TopicID, &quiz.Category, &quiz.Description, &quiz.CreatedAt, &quiz.UpdatedAt); err != nil {
			return nil, err
		}
		quizzes = append(quizzes, quiz)
	}
	return quizzes, rows.Err()
}

func (db *SQLiteDB) UpdateQuiz(quiz *models.Quiz) error {
	result, err := db.db.Exec(
		"UPDATE quizzes SET name = ?, topic_id = ?, category = ?, description = ?, updated_at = ? WHERE id = ?",
		quiz.Name, quiz.TopicID, quiz.Category, quiz.Description, quiz.UpdatedAt, quiz.ID,
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
		"INSERT INTO quiz_attempts (user_id, quiz_id, parent_attempt_id, attempt_number, is_retake, start_time, completed_time, total_questions, correct_answers, difficulty_average, description, question_ids, question_results) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		attempt.UserID, attempt.QuizID, attempt.ParentAttemptID, attempt.AttemptNumber, attempt.IsRetake, attempt.StartTime, attempt.CompletedTime, attempt.TotalQuestions, attempt.CorrectAnswers, attempt.DifficultyAverage, attempt.Description, string(questionIDsJSON), string(questionResultsJSON),
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
	var quizID, parentAttemptID sql.NullInt64

	err := db.db.QueryRow(
		"SELECT id, user_id, quiz_id, parent_attempt_id, attempt_number, is_retake, start_time, completed_time, total_questions, correct_answers, difficulty_average, description, question_ids, question_results FROM quiz_attempts WHERE id = ?",
		id,
	).Scan(&attempt.ID, &attempt.UserID, &quizID, &parentAttemptID, &attempt.AttemptNumber, &attempt.IsRetake, &attempt.StartTime, &completedTime, &attempt.TotalQuestions, &attempt.CorrectAnswers, &attempt.DifficultyAverage, &attempt.Description, &questionIDsJSON, &questionResultsJSON)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	if completedTime.Valid {
		attempt.CompletedTime = &completedTime.Time
	}

	if quizID.Valid {
		qid := int(quizID.Int64)
		attempt.QuizID = &qid
	}

	if parentAttemptID.Valid {
		paid := int(parentAttemptID.Int64)
		attempt.ParentAttemptID = &paid
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
		"SELECT id, user_id, quiz_id, parent_attempt_id, attempt_number, is_retake, start_time, completed_time, total_questions, correct_answers, difficulty_average, description, question_ids, question_results FROM quiz_attempts WHERE user_id = ?",
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
		var quizID, parentAttemptID sql.NullInt64

		if err := rows.Scan(&attempt.ID, &attempt.UserID, &quizID, &parentAttemptID, &attempt.AttemptNumber, &attempt.IsRetake, &attempt.StartTime, &completedTime, &attempt.TotalQuestions, &attempt.CorrectAnswers, &attempt.DifficultyAverage, &attempt.Description, &questionIDsJSON, &questionResultsJSON); err != nil {
			return nil, err
		}

		if completedTime.Valid {
			attempt.CompletedTime = &completedTime.Time
		}

		if quizID.Valid {
			qid := int(quizID.Int64)
			attempt.QuizID = &qid
		}

		if parentAttemptID.Valid {
			paid := int(parentAttemptID.Int64)
			attempt.ParentAttemptID = &paid
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

func (db *SQLiteDB) GetAllQuizAttempts() ([]*models.QuizAttempt, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, quiz_id, parent_attempt_id, attempt_number, is_retake, start_time, completed_time, total_questions, correct_answers, difficulty_average, description, question_ids, question_results FROM quiz_attempts",
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
		var quizID, parentAttemptID sql.NullInt64

		if err := rows.Scan(&attempt.ID, &attempt.UserID, &quizID, &parentAttemptID, &attempt.AttemptNumber, &attempt.IsRetake, &attempt.StartTime, &completedTime, &attempt.TotalQuestions, &attempt.CorrectAnswers, &attempt.DifficultyAverage, &attempt.Description, &questionIDsJSON, &questionResultsJSON); err != nil {
			return nil, err
		}

		if completedTime.Valid {
			attempt.CompletedTime = &completedTime.Time
		}

		if quizID.Valid {
			qid := int(quizID.Int64)
			attempt.QuizID = &qid
		}

		if parentAttemptID.Valid {
			paid := int(parentAttemptID.Int64)
			attempt.ParentAttemptID = &paid
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
		"UPDATE quiz_attempts SET quiz_id = ?, parent_attempt_id = ?, attempt_number = ?, is_retake = ?, completed_time = ?, total_questions = ?, correct_answers = ?, difficulty_average = ?, description = ?, question_ids = ?, question_results = ? WHERE id = ?",
		attempt.QuizID, attempt.ParentAttemptID, attempt.AttemptNumber, attempt.IsRetake, attempt.CompletedTime, attempt.TotalQuestions, attempt.CorrectAnswers, attempt.DifficultyAverage, attempt.Description, string(questionIDsJSON), string(questionResultsJSON), attempt.ID,
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

func (db *SQLiteDB) GetQuizAttemptWithDetails(attemptID int) (*models.QuizAttempt, error) {
	// First get the quiz attempt
	attempt, err := db.GetQuizAttempt(attemptID)
	if err != nil {
		return nil, err
	}

	// Return the attempt - the API handler will need to populate questions separately
	return attempt, nil
}

