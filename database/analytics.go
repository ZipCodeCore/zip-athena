package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

func (db *SQLiteDB) GetStudentProgress(userID int) (*models.StudentProgress, error) {
	// Get user info
	user, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	progress := &models.StudentProgress{
		UserID:         userID,
		Username:       user.Username,
		TotalPoints:    user.TotalPoints,
		TopicBreakdown: make(map[string]float64),
		RecentActivity: make([]models.QuizAttempt, 0),
	}

	// Get quiz attempt stats
	rows, err := db.db.Query(`
		SELECT COUNT(*), 
		       SUM(CASE WHEN completed_time IS NOT NULL THEN 1 ELSE 0 END),
		       AVG(CASE WHEN completed_time IS NOT NULL AND total_questions > 0 
		                THEN CAST(correct_answers AS REAL) / total_questions * 100 
		                ELSE 0 END)
		FROM quiz_attempts WHERE user_id = ?
	`, userID)
	if err == nil && rows.Next() {
		rows.Scan(&progress.TotalQuizzes, &progress.CompletedQuizzes, &progress.AverageScore)
		rows.Close()
	}

	// Get recent activity (last 5 completed quizzes)
	activityRows, err := db.db.Query(`
		SELECT id, user_id, start_time, completed_time, total_questions, correct_answers, question_ids, question_results
		FROM quiz_attempts 
		WHERE user_id = ? AND completed_time IS NOT NULL
		ORDER BY start_time DESC LIMIT 5
	`, userID)
	if err == nil {
		defer activityRows.Close()
		for activityRows.Next() {
			attempt := models.QuizAttempt{}
			var questionIDsJSON, questionResultsJSON string
			err := activityRows.Scan(&attempt.ID, &attempt.UserID, &attempt.StartTime, &attempt.CompletedTime,
				&attempt.TotalQuestions, &attempt.CorrectAnswers, &questionIDsJSON, &questionResultsJSON)
			if err == nil {
				json.Unmarshal([]byte(questionIDsJSON), &attempt.QuestionIDs)
				json.Unmarshal([]byte(questionResultsJSON), &attempt.QuestionResults)
				progress.RecentActivity = append(progress.RecentActivity, attempt)
			}
		}
	}

	// Get question accuracy
	qaRows, err := db.db.Query(`
		SELECT SUM(correct_attempts + incorrect_attempts), SUM(correct_attempts)
		FROM question_attempts WHERE user_id = ?
	`, userID)
	if err == nil && qaRows.Next() {
		var totalAttempts, correctAttempts int
		qaRows.Scan(&totalAttempts, &correctAttempts)
		qaRows.Close()
		if totalAttempts > 0 {
			progress.QuestionAccuracy = float64(correctAttempts) / float64(totalAttempts) * 100
		}
	}

	// Get topic breakdown
	topicRows, err := db.db.Query(`
		SELECT t.name, 
		       CAST(SUM(qa.correct_attempts) AS REAL) / SUM(qa.correct_attempts + qa.incorrect_attempts) * 100
		FROM question_attempts qa
		INNER JOIN questions q ON qa.question_id = q.id
		INNER JOIN topics t ON q.topic_id = t.id
		WHERE qa.user_id = ?
		GROUP BY t.name
		HAVING SUM(qa.correct_attempts + qa.incorrect_attempts) > 0
	`, userID)
	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var topicName string
			var accuracy float64
			if topicRows.Scan(&topicName, &accuracy) == nil {
				progress.TopicBreakdown[topicName] = accuracy
			}
		}
	}

	// Get FSRS metrics
	fsrsRows, err := db.db.Query(`
		SELECT COUNT(*),
		       SUM(CASE WHEN state = 0 THEN 1 ELSE 0 END),
		       SUM(CASE WHEN state = 2 THEN 1 ELSE 0 END),
		       SUM(CASE WHEN due <= datetime('now') THEN 1 ELSE 0 END),
		       AVG(stability),
		       AVG(difficulty)
		FROM user_question_states WHERE user_id = ?
	`, userID)
	if err == nil && fsrsRows.Next() {
		fsrsRows.Scan(&progress.FSRSMetrics.TotalCards, &progress.FSRSMetrics.NewCards,
			&progress.FSRSMetrics.ReviewCards, &progress.FSRSMetrics.DueCards,
			&progress.FSRSMetrics.AvgStability, &progress.FSRSMetrics.AvgDifficulty)
		fsrsRows.Close()
	}

	return progress, nil
}

func (db *SQLiteDB) GetCohortAnalytics() (*models.CohortAnalytics, error) {
	analytics := &models.CohortAnalytics{
		TopicDifficulty:  make(map[string]float64),
		QuestionStats:    make([]models.QuestionStatistic, 0),
		TopPerformers:    make([]models.StudentProgress, 0),
		BottomPerformers: make([]models.StudentProgress, 0),
	}

	// Count total students (non-admin)
	row := db.db.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = 0")
	row.Scan(&analytics.TotalStudents)

	// Count active students (last 7 days)
	activeRow := db.db.QueryRow(`
		SELECT COUNT(DISTINCT user_id) FROM quiz_attempts 
		WHERE start_time >= datetime('now', '-7 days')
	`)
	activeRow.Scan(&analytics.ActiveStudents)

	// Get total quizzes and average score
	statsRow := db.db.QueryRow(`
		SELECT COUNT(*),
		       AVG(CASE WHEN total_questions > 0 
		                THEN CAST(correct_answers AS REAL) / total_questions * 100 
		                ELSE 0 END)
		FROM quiz_attempts WHERE completed_time IS NOT NULL
	`)
	statsRow.Scan(&analytics.TotalQuizzesTaken, &analytics.AverageCohortScore)

	// Get top performers
	topRows, err := db.db.Query(`
		SELECT user_id, AVG(CAST(correct_answers AS REAL) / total_questions * 100) as avg_score
		FROM quiz_attempts 
		WHERE completed_time IS NOT NULL AND total_questions > 0
		GROUP BY user_id
		ORDER BY avg_score DESC LIMIT 5
	`)
	if err == nil {
		defer topRows.Close()
		for topRows.Next() {
			var userID int
			var avgScore float64
			if topRows.Scan(&userID, &avgScore) == nil {
				if progress, err := db.GetStudentProgress(userID); err == nil {
					analytics.TopPerformers = append(analytics.TopPerformers, *progress)
				}
			}
		}
	}

	// Get bottom performers
	bottomRows, err := db.db.Query(`
		SELECT user_id, AVG(CAST(correct_answers AS REAL) / total_questions * 100) as avg_score
		FROM quiz_attempts 
		WHERE completed_time IS NOT NULL AND total_questions > 0
		GROUP BY user_id
		ORDER BY avg_score ASC LIMIT 5
	`)
	if err == nil {
		defer bottomRows.Close()
		for bottomRows.Next() {
			var userID int
			var avgScore float64
			if bottomRows.Scan(&userID, &avgScore) == nil {
				if progress, err := db.GetStudentProgress(userID); err == nil {
					analytics.BottomPerformers = append(analytics.BottomPerformers, *progress)
				}
			}
		}
	}

	// Get topic difficulty
	topicRows, err := db.db.Query(`
		SELECT t.name, AVG(uqs.difficulty)
		FROM user_question_states uqs
		INNER JOIN questions q ON uqs.question_id = q.id
		INNER JOIN topics t ON q.topic_id = t.id
		GROUP BY t.name
	`)
	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var topicName string
			var avgDifficulty float64
			if topicRows.Scan(&topicName, &avgDifficulty) == nil {
				analytics.TopicDifficulty[topicName] = avgDifficulty
			}
		}
	}

	// Get question statistics
	questionRows, err := db.db.Query(`
		SELECT q.id, q.topic_question_id, q.question_text,
		       SUM(qa.correct_attempts + qa.incorrect_attempts) as total_attempts,
		       SUM(qa.correct_attempts) as correct_attempts,
		       CAST(SUM(qa.correct_attempts) AS REAL) / SUM(qa.correct_attempts + qa.incorrect_attempts) * 100 as accuracy,
		       COALESCE(AVG(uqs.difficulty), 0) as avg_difficulty
		FROM questions q
		LEFT JOIN question_attempts qa ON q.id = qa.question_id
		LEFT JOIN user_question_states uqs ON q.id = uqs.question_id
		GROUP BY q.id
		HAVING total_attempts > 0
		ORDER BY accuracy ASC
		LIMIT 20
	`)
	if err == nil {
		defer questionRows.Close()
		for questionRows.Next() {
			stat := models.QuestionStatistic{}
			questionRows.Scan(&stat.QuestionID, &stat.TopicQuestionID, &stat.QuestionText,
				&stat.TotalAttempts, &stat.CorrectAttempts, &stat.Accuracy, &stat.AvgDifficulty)
			analytics.QuestionStats = append(analytics.QuestionStats, stat)
		}
	}

	return analytics, nil
}

func (db *SQLiteDB) ListStudentsWithMetrics(searchQuery string, sortBy string, filterCohort string) ([]*models.StudentProgress, error) {
	// Get all non-admin users
	query := "SELECT id FROM users WHERE is_admin = 0 AND is_active = 1"
	args := []interface{}{}

	// Add search filter if provided
	if searchQuery != "" {
		query += " AND (username LIKE ? OR email LIKE ?)"
		searchPattern := "%" + searchQuery + "%"
		args = append(args, searchPattern, searchPattern)
	}

	rows, err := db.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]*models.StudentProgress, 0)
	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			continue
		}

		// Get full progress for each student
		progress, err := db.GetStudentProgress(userID)
		if err != nil {
			continue
		}
		students = append(students, progress)
	}

	// Sort students based on sortBy parameter
	// For now, we return unsorted - sorting can be done in the handler
	// TODO: Implement sorting by accuracy, completion rate, last activity

	return students, nil
}

func (db *SQLiteDB) GetStudentQuizHistory(userID int) ([]*models.QuizAttempt, error) {
	// Verify user exists
	_, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	rows, err := db.db.Query(`
		SELECT id, user_id, quiz_id, parent_attempt_id, attempt_number, is_retake,
		       start_time, completed_time, total_questions, correct_answers
		FROM quiz_attempts
		WHERE user_id = ?
		ORDER BY start_time DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attempts := make([]*models.QuizAttempt, 0)
	for rows.Next() {
		attempt := &models.QuizAttempt{}
		var quizID, parentAttemptID *int
		var completedTime *time.Time

		err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&quizID,
			&parentAttemptID,
			&attempt.AttemptNumber,
			&attempt.IsRetake,
			&attempt.StartTime,
			&completedTime,
			&attempt.TotalQuestions,
			&attempt.CorrectAnswers,
		)
		if err != nil {
			continue
		}

		attempt.QuizID = quizID
		attempt.ParentAttemptID = parentAttemptID
		attempt.CompletedTime = completedTime
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}

func (db *SQLiteDB) GetStudentWeakAreas(userID int) (*models.WeakAreas, error) {
	user, err := db.GetUser(userID)
	if err != nil {
		return nil, err
	}

	weakAreas := &models.WeakAreas{
		UserID:     userID,
		Username:   user.Username,
		WeakTopics: make([]models.TopicWeakness, 0),
		WeakTags:   make([]models.TagWeakness, 0),
	}

	// Get weak topics (< 60% accuracy)
	topicRows, err := db.db.Query(`
		SELECT t.id, t.name,
		       SUM(qa.correct_attempts) as correct,
		       SUM(qa.correct_attempts + qa.incorrect_attempts) as total,
		       CAST(SUM(qa.correct_attempts) AS REAL) / SUM(qa.correct_attempts + qa.incorrect_attempts) * 100 as accuracy
		FROM topics t
		JOIN questions q ON t.id = q.topic_id
		JOIN question_attempts qa ON q.id = qa.question_id
		WHERE qa.user_id = ?
		GROUP BY t.id, t.name
		HAVING total > 0 AND accuracy < 60
		ORDER BY accuracy ASC
	`, userID)

	if err == nil {
		defer topicRows.Close()
		for topicRows.Next() {
			var topicID int
			var topicName string
			var correct, total int
			var accuracy float64

			if topicRows.Scan(&topicID, &topicName, &correct, &total, &accuracy) == nil {
				weakAreas.WeakTopics = append(weakAreas.WeakTopics, models.TopicWeakness{
					TopicID:       topicID,
					TopicName:     topicName,
					Accuracy:      accuracy,
					TotalAttempts: total,
				})
			}
		}
	}

	// Get weak tags (< 60% accuracy)
	// This requires extracting tags from questions and calculating accuracy per tag
	// For simplicity, we'll query question attempts and aggregate by tags
	qaRows, err := db.db.Query(`
		SELECT qa.question_id, qa.correct_attempts, qa.incorrect_attempts
		FROM question_attempts qa
		WHERE qa.user_id = ?
	`, userID)

	if err == nil {
		defer qaRows.Close()

		// Map to store tag statistics: tag -> [correct, total]
		tagStats := make(map[string][2]int)

		for qaRows.Next() {
			var questionID, correct, incorrect int
			if qaRows.Scan(&questionID, &correct, &incorrect) == nil {
				// Get question tags
				question, err := db.GetQuestion(questionID)
				if err != nil {
					continue
				}

				total := correct + incorrect
				for _, tag := range question.Tags {
					stats := tagStats[tag]
					stats[0] += correct
					stats[1] += total
					tagStats[tag] = stats
				}
			}
		}

		// Calculate accuracy for each tag and filter weak ones
		for tag, stats := range tagStats {
			correct := stats[0]
			total := stats[1]

			if total > 0 {
				accuracy := float64(correct) / float64(total) * 100
				if accuracy < 60 {
					weakAreas.WeakTags = append(weakAreas.WeakTags, models.TagWeakness{
						Tag:           tag,
						Accuracy:      accuracy,
						TotalAttempts: total,
					})
				}
			}
		}
	}

	return weakAreas, nil
}

func (db *SQLiteDB) GetQuestionAnalytics() ([]*models.QuestionAnalytics, error) {
	query := `
		SELECT q.id, q.topic_question_id, q.question_text, q.topic_id, t.name,
			   q.tags, q.difficulty,
			   COALESCE(SUM(qa.correct_attempts + qa.incorrect_attempts), 0) as total_attempts,
			   COALESCE(SUM(qa.correct_attempts), 0) as correct_attempts
		FROM questions q
		LEFT JOIN topics t ON q.topic_id = t.id
		LEFT JOIN question_attempts qa ON q.id = qa.question_id
		GROUP BY q.id
	`

	rows, err := db.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analytics []*models.QuestionAnalytics
	for rows.Next() {
		a := &models.QuestionAnalytics{
			AnswerDistribution: make(map[int]int),
		}
		var tagsJSON string

		err := rows.Scan(
			&a.QuestionID, &a.TopicQuestionID, &a.QuestionText, &a.TopicID, &a.TopicName,
			&tagsJSON, &a.CurrentDifficulty, &a.TotalAttempts, &a.CorrectAttempts,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(tagsJSON), &a.Tags)

		// Calculate accuracy
		if a.TotalAttempts > 0 {
			a.Accuracy = float64(a.CorrectAttempts) / float64(a.TotalAttempts) * 100
		}

		// Calculate recommended difficulty
		if a.Accuracy > 90 {
			a.RecommendedDifficulty = a.CurrentDifficulty + 1
			if a.RecommendedDifficulty > 7 {
				a.RecommendedDifficulty = 7
			}
			a.FlagForReview = true
		} else if a.Accuracy < 30 {
			a.RecommendedDifficulty = a.CurrentDifficulty - 1
			if a.RecommendedDifficulty < 1 {
				a.RecommendedDifficulty = 1
			}
			a.FlagForReview = true
		} else {
			a.RecommendedDifficulty = a.CurrentDifficulty
		}

		analytics = append(analytics, a)
	}

	return analytics, rows.Err()
}

func (db *SQLiteDB) GetQuestionDetailedAnalytics(questionID int) (*models.QuestionAnalytics, error) {
	query := `
		SELECT q.id, q.topic_question_id, q.question_text, q.topic_id, t.name,
			   q.tags, q.difficulty,
			   COALESCE(SUM(qa.correct_attempts + qa.incorrect_attempts), 0) as total_attempts,
			   COALESCE(SUM(qa.correct_attempts), 0) as correct_attempts
		FROM questions q
		LEFT JOIN topics t ON q.topic_id = t.id
		LEFT JOIN question_attempts qa ON q.id = qa.question_id
		WHERE q.id = ?
		GROUP BY q.id
	`

	a := &models.QuestionAnalytics{
		AnswerDistribution: make(map[int]int),
	}
	var tagsJSON string

	err := db.db.QueryRow(query, questionID).Scan(
		&a.QuestionID, &a.TopicQuestionID, &a.QuestionText, &a.TopicID, &a.TopicName,
		&tagsJSON, &a.CurrentDifficulty, &a.TotalAttempts, &a.CorrectAttempts,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(tagsJSON), &a.Tags)

	// Calculate accuracy
	if a.TotalAttempts > 0 {
		a.Accuracy = float64(a.CorrectAttempts) / float64(a.TotalAttempts) * 100
	}

	// Calculate recommended difficulty
	if a.Accuracy > 90 {
		a.RecommendedDifficulty = a.CurrentDifficulty + 1
		if a.RecommendedDifficulty > 7 {
			a.RecommendedDifficulty = 7
		}
		a.FlagForReview = true
	} else if a.Accuracy < 30 {
		a.RecommendedDifficulty = a.CurrentDifficulty - 1
		if a.RecommendedDifficulty < 1 {
			a.RecommendedDifficulty = 1
		}
		a.FlagForReview = true
	} else {
		a.RecommendedDifficulty = a.CurrentDifficulty
	}

	// Get answer distribution
	a.AnswerDistribution, _ = db.GetAnswerDistribution(questionID)

	return a, nil
}

func (db *SQLiteDB) GetAnswerDistribution(questionID int) (map[int]int, error) {
	// This would require tracking answer selections in quiz/test attempts
	// For now, return empty distribution
	return make(map[int]int), nil
}

func (db *SQLiteDB) UpdateQuestionDifficultyFromAnalytics(questionID int, newDifficulty int) error {
	if newDifficulty < 1 || newDifficulty > 7 {
		return fmt.Errorf("difficulty must be between 1 and 7")
	}

	_, err := db.db.Exec("UPDATE questions SET difficulty = ? WHERE id = ?", newDifficulty, questionID)
	return err
}
