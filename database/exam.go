package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

// ---- SQLiteDB implementation of ExamRepository ----

func (db *SQLiteDB) CreateExam(exam *models.Exam) error {
	categoriesJSON, _ := json.Marshal(exam.Categories)
	difficultyJSON, _ := json.Marshal(exam.DifficultyLevels)
	questionIDsJSON, _ := json.Marshal(exam.QuestionIDs)

	result, err := db.db.Exec(
		`INSERT INTO exams (name, categories, difficulty_levels, question_ids, question_count, access_token, created_by, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		exam.Name, string(categoriesJSON), string(difficultyJSON), string(questionIDsJSON),
		exam.QuestionCount, exam.AccessToken, exam.CreatedBy, exam.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	exam.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetExam(id int) (*models.Exam, error) {
	return db.scanExam(db.db.QueryRow(
		`SELECT id, name, categories, difficulty_levels, question_ids, question_count, access_token, created_by, created_at
		 FROM exams WHERE id = ?`, id))
}

func (db *SQLiteDB) GetExamByToken(token string) (*models.Exam, error) {
	return db.scanExam(db.db.QueryRow(
		`SELECT id, name, categories, difficulty_levels, question_ids, question_count, access_token, created_by, created_at
		 FROM exams WHERE access_token = ?`, token))
}

func (db *SQLiteDB) scanExam(row *sql.Row) (*models.Exam, error) {
	exam := &models.Exam{}
	var catJSON, diffJSON, qidsJSON string
	err := row.Scan(
		&exam.ID, &exam.Name, &catJSON, &diffJSON, &qidsJSON,
		&exam.QuestionCount, &exam.AccessToken, &exam.CreatedBy, &exam.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(catJSON), &exam.Categories)
	json.Unmarshal([]byte(diffJSON), &exam.DifficultyLevels)
	json.Unmarshal([]byte(qidsJSON), &exam.QuestionIDs)
	return exam, nil
}

func (db *SQLiteDB) ListExams() ([]*models.Exam, error) {
	rows, err := db.db.Query(
		`SELECT id, name, categories, difficulty_levels, question_ids, question_count, access_token, created_by, created_at
		 FROM exams ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exams []*models.Exam
	for rows.Next() {
		exam := &models.Exam{}
		var catJSON, diffJSON, qidsJSON string
		if err := rows.Scan(
			&exam.ID, &exam.Name, &catJSON, &diffJSON, &qidsJSON,
			&exam.QuestionCount, &exam.AccessToken, &exam.CreatedBy, &exam.CreatedAt,
		); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(catJSON), &exam.Categories)
		json.Unmarshal([]byte(diffJSON), &exam.DifficultyLevels)
		json.Unmarshal([]byte(qidsJSON), &exam.QuestionIDs)
		exams = append(exams, exam)
	}
	if exams == nil {
		exams = []*models.Exam{}
	}
	return exams, rows.Err()
}

func (db *SQLiteDB) DeleteExam(id int) error {
	result, err := db.db.Exec(`DELETE FROM exams WHERE id = ?`, id)
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

// GetQuestionsByCategories returns up to `limit` questions from the given
// categories and difficulty levels, chosen in random order from all matches.
// DifficultyLevels mapping: 1 = difficulty 1-2, 2 = difficulty 3-4, 3 = difficulty 5-7.
func (db *SQLiteDB) GetQuestionsByCategories(categories []string, difficultyLevels []int, limit int) ([]*models.Question, error) {
	diffRanges := difficultyRanges(difficultyLevels)

	// Build a query that joins topics for category filtering and filters by difficulty
	query := `
		SELECT DISTINCT q.id, q.topic_id, q.question_text, q.difficulty, q.tags, q.topic_question_id, q.created_at
		FROM questions q
		JOIN topics t ON q.topic_id = t.id
		WHERE t.is_active = 1
	`
	args := []interface{}{}

	if len(categories) > 0 {
		placeholders := make([]string, len(categories))
		for i, c := range categories {
			placeholders[i] = "?"
			args = append(args, c)
		}
		query += " AND (t.category IN (" + joinStrings(placeholders, ",") + ") OR t.category = 'all')"
	}

	if len(diffRanges) > 0 {
		conds := make([]string, len(diffRanges))
		for i, r := range diffRanges {
			conds[i] = "(q.difficulty >= ? AND q.difficulty <= ?)"
			args = append(args, r[0], r[1])
		}
		query += " AND (" + joinStrings(conds, " OR ") + ")"
	}

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []*models.Question
	for rows.Next() {
		q := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&q.ID, &q.TopicID, &q.QuestionText, &q.Difficulty, &tagsJSON, &q.TopicQuestionID, &q.CreatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsJSON), &q.Tags)
		all = append(all, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Shuffle and cap to limit
	rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

// GetQuestionsByTopicIDs returns up to `limit` questions from the given topic IDs
// and difficulty levels, chosen in random order.
func (db *SQLiteDB) GetQuestionsByTopicIDs(topicIDs []int, difficultyLevels []int, limit int) ([]*models.Question, error) {
	if len(topicIDs) == 0 {
		return nil, nil
	}

	diffRanges := difficultyRanges(difficultyLevels)

	placeholders := make([]string, len(topicIDs))
	args := []interface{}{}
	for i, id := range topicIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := `
		SELECT DISTINCT q.id, q.topic_id, q.question_text, q.difficulty, q.tags, q.topic_question_id, q.created_at
		FROM questions q
		JOIN topics t ON q.topic_id = t.id
		WHERE t.is_active = 1
		  AND q.topic_id IN (` + joinStrings(placeholders, ",") + `)`

	if len(diffRanges) > 0 {
		conds := make([]string, len(diffRanges))
		for i, r := range diffRanges {
			conds[i] = "(q.difficulty >= ? AND q.difficulty <= ?)"
			args = append(args, r[0], r[1])
		}
		query += " AND (" + joinStrings(conds, " OR ") + ")"
	}

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []*models.Question
	for rows.Next() {
		q := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&q.ID, &q.TopicID, &q.QuestionText, &q.Difficulty, &tagsJSON, &q.TopicQuestionID, &q.CreatedAt); err != nil {
			return nil, err
		}
		if tagsJSON != "" {
			if err := json.Unmarshal([]byte(tagsJSON), &q.Tags); err != nil {
				q.Tags = []string{}
			}
		}
		all = append(all, q)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}

func (db *SQLiteDB) AssignExamToUser(examID, userID, assignedBy int) error {
	user, err := db.GetUser(userID)
	if err != nil {
		return err
	}
	if user.IsAdmin {
		return fmt.Errorf("cannot assign exam to admin user")
	}
	if !user.IsActive {
		return fmt.Errorf("user is inactive")
	}
	_, err = db.db.Exec(
		`INSERT OR IGNORE INTO exam_assignments (exam_id, user_id, assigned_by, assigned_at) VALUES (?, ?, ?, ?)`,
		examID, userID, assignedBy, time.Now(),
	)
	return err
}

func (db *SQLiteDB) UnassignExamFromUser(examID, userID int) error {
	result, err := db.db.Exec(
		`DELETE FROM exam_assignments WHERE exam_id = ? AND user_id = ?`,
		examID, userID,
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

func (db *SQLiteDB) GetAssignedExamsForUser(userID int) ([]*models.ExamInfo, error) {
	query := `
		SELECT ea.id, ea.exam_id, ea.user_id, ea.assigned_at, ea.assigned_by,
		       e.name, e.question_count, e.access_token,
		       att.id, att.completed_at, att.correct_answers, att.total_questions, att.score
		FROM exam_assignments ea
		JOIN exams e ON ea.exam_id = e.id
		LEFT JOIN exam_attempts att ON att.exam_id = e.id
		    AND att.user_id = ea.user_id
		    AND att.completed_at IS NOT NULL
		    AND att.id = (
		        SELECT MAX(a2.id) FROM exam_attempts a2
		        WHERE a2.exam_id = e.id AND a2.user_id = ea.user_id AND a2.completed_at IS NOT NULL
		    )
		WHERE ea.user_id = ?
		ORDER BY ea.assigned_at DESC
	`
	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.ExamInfo
	for rows.Next() {
		info := &models.ExamInfo{}
		var attemptID sql.NullInt64
		var completedAt sql.NullTime
		var correctAnswers, totalQuestions sql.NullInt64
		var score sql.NullFloat64

		if err := rows.Scan(
			&info.ID, &info.ExamID, &info.UserID, &info.AssignedAt, &info.AssignedBy,
			&info.ExamName, &info.QuestionCount, &info.AccessToken,
			&attemptID, &completedAt, &correctAnswers, &totalQuestions, &score,
		); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			info.CompletedAt = &completedAt.Time
			if attemptID.Valid {
				aid := int(attemptID.Int64)
				info.AttemptID = &aid
			}
			if correctAnswers.Valid {
				ca := int(correctAnswers.Int64)
				info.CorrectAnswers = &ca
			}
			if totalQuestions.Valid {
				tq := int(totalQuestions.Int64)
				info.TotalQuestions = &tq
			}
			if score.Valid {
				s := score.Float64
				info.Score = &s
			}
		}
		result = append(result, info)
	}
	if result == nil {
		result = []*models.ExamInfo{}
	}
	return result, rows.Err()
}

func (db *SQLiteDB) CreateExamAttempt(attempt *models.ExamAttempt) error {
	answersJSON, _ := json.Marshal(attempt.Answers)
	resultsJSON, _ := json.Marshal(attempt.QuestionResults)

	result, err := db.db.Exec(
		`INSERT INTO exam_attempts (exam_id, user_id, started_at, total_questions, correct_answers, score, answers, question_results)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		attempt.ExamID, attempt.UserID, attempt.StartedAt,
		attempt.TotalQuestions, attempt.CorrectAnswers, attempt.Score,
		string(answersJSON), string(resultsJSON),
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

func (db *SQLiteDB) GetExamAttempt(id int) (*models.ExamAttempt, error) {
	attempt := &models.ExamAttempt{}
	var completedAt sql.NullTime
	var answersJSON, resultsJSON string

	err := db.db.QueryRow(
		`SELECT id, exam_id, user_id, started_at, completed_at, total_questions, correct_answers, score, answers, question_results
		 FROM exam_attempts WHERE id = ?`, id,
	).Scan(
		&attempt.ID, &attempt.ExamID, &attempt.UserID, &attempt.StartedAt, &completedAt,
		&attempt.TotalQuestions, &attempt.CorrectAnswers, &attempt.Score,
		&answersJSON, &resultsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if completedAt.Valid {
		attempt.CompletedAt = &completedAt.Time
	}
	json.Unmarshal([]byte(answersJSON), &attempt.Answers)
	json.Unmarshal([]byte(resultsJSON), &attempt.QuestionResults)
	return attempt, nil
}

// SubmitExamAttempt grades the exam and records the results.
func (db *SQLiteDB) SubmitExamAttempt(attemptID int, answers map[string]int) error {
	attempt, err := db.GetExamAttempt(attemptID)
	if err != nil {
		return err
	}

	exam, err := db.GetExam(attempt.ExamID)
	if err != nil {
		return err
	}

	attempt.Answers = answers
	attempt.QuestionResults = make(map[string]bool)
	attempt.CorrectAnswers = 0
	attempt.TotalQuestions = len(exam.QuestionIDs)

	for _, qid := range exam.QuestionIDs {
		key := itoa(qid)
		answerID, ok := answers[key]
		if !ok {
			attempt.QuestionResults[key] = false
			continue
		}
		answer, err := db.GetAnswer(answerID)
		if err == nil && answer.QuestionID == qid && answer.IsCorrect {
			attempt.CorrectAnswers++
			attempt.QuestionResults[key] = true
		} else {
			attempt.QuestionResults[key] = false
		}
	}

	if attempt.TotalQuestions > 0 {
		attempt.Score = float64(attempt.CorrectAnswers) / float64(attempt.TotalQuestions) * 100
	}
	now := time.Now()
	attempt.CompletedAt = &now

	answersJSON, _ := json.Marshal(attempt.Answers)
	resultsJSON, _ := json.Marshal(attempt.QuestionResults)

	_, err = db.db.Exec(
		`UPDATE exam_attempts SET completed_at = ?, total_questions = ?, correct_answers = ?, score = ?, answers = ?, question_results = ?
		 WHERE id = ?`,
		now, attempt.TotalQuestions, attempt.CorrectAnswers, attempt.Score,
		string(answersJSON), string(resultsJSON), attemptID,
	)
	return err
}

func (db *SQLiteDB) GetExamStatus(examID int) (*models.ExamStatus, error) {
	exam, err := db.GetExam(examID)
	if err != nil {
		return nil, err
	}

	status := &models.ExamStatus{
		ExamID:         exam.ID,
		ExamName:       exam.Name,
		QuestionCount:  exam.QuestionCount,
		StudentResults: []models.ExamStudentResult{},
	}

	query := `
		SELECT u.id, u.username, ea.assigned_at,
		       att.id, att.completed_at, att.correct_answers, att.total_questions, att.score
		FROM exam_assignments ea
		JOIN users u ON ea.user_id = u.id
		LEFT JOIN exam_attempts att ON att.exam_id = ea.exam_id
		    AND att.user_id = ea.user_id
		    AND att.completed_at IS NOT NULL
		    AND att.id = (
		        SELECT MAX(a2.id) FROM exam_attempts a2
		        WHERE a2.exam_id = ea.exam_id AND a2.user_id = ea.user_id AND a2.completed_at IS NOT NULL
		    )
		WHERE ea.exam_id = ?
		ORDER BY u.username
	`
	rows, err := db.db.Query(query, examID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		r := models.ExamStudentResult{}
		var attemptID sql.NullInt64
		var completedAt sql.NullTime
		var correctAnswers, totalQuestions sql.NullInt64
		var score sql.NullFloat64

		if err := rows.Scan(
			&r.UserID, &r.Username, &r.AssignedAt,
			&attemptID, &completedAt, &correctAnswers, &totalQuestions, &score,
		); err != nil {
			return nil, err
		}
		if completedAt.Valid {
			r.CompletedAt = &completedAt.Time
			if attemptID.Valid {
				aid := int(attemptID.Int64)
				r.AttemptID = &aid
			}
			if correctAnswers.Valid {
				ca := int(correctAnswers.Int64)
				r.CorrectAnswers = &ca
			}
			if totalQuestions.Valid {
				tq := int(totalQuestions.Int64)
				r.TotalQuestions = &tq
			}
			if score.Valid {
				s := score.Float64
				r.Score = &s
			}
		}
		status.StudentResults = append(status.StudentResults, r)
	}
	return status, rows.Err()
}

// ---- helpers ----

func joinStrings(ss []string, sep string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

func itoa(n int) string {
	return json.Number(fmt.Sprintf("%d", n)).String()
}

// difficultyRanges converts exam difficulty levels (1,2,3) to [min,max] DB ranges.
func difficultyRanges(levels []int) [][2]int {
	mapping := map[int][2]int{
		1: {1, 2},
		2: {3, 4},
		3: {5, 7},
	}
	var result [][2]int
	seen := map[int]bool{}
	for _, l := range levels {
		if r, ok := mapping[l]; ok && !seen[l] {
			result = append(result, r)
			seen[l] = true
		}
	}
	return result
}
