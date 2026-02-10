package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

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

func (db *SQLiteDB) GetTestResults(attemptID int) (*models.TestAttempt, error) {
	return db.GetTestAttempt(attemptID)
}

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

func (db *SQLiteDB) PublishTest(testID int) error {
	_, err := db.db.Exec("UPDATE tests SET is_published = 1, updated_at = ? WHERE id = ?", time.Now(), testID)
	return err
}

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

