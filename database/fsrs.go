package database

import (
"database/sql"

"github.com/ZipCodeCore/zip-athena/models"
)

func (db *SQLiteDB) CreateUserQuestionState(state *models.UserQuestionState) error {
	result, err := db.db.Exec(
		"INSERT INTO user_question_states (user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		state.UserID, state.QuestionID, state.Due, state.Stability, state.Difficulty, state.ElapsedDays, state.ScheduledDays, state.Reps, state.Lapses, state.State, state.LastReview,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	state.ID = int(id)
	return nil
}

func (db *SQLiteDB) GetUserQuestionState(userID, questionID int) (*models.UserQuestionState, error) {
	state := &models.UserQuestionState{}
	err := db.db.QueryRow(
		"SELECT id, user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review FROM user_question_states WHERE user_id = ? AND question_id = ?",
		userID, questionID,
	).Scan(&state.ID, &state.UserID, &state.QuestionID, &state.Due, &state.Stability, &state.Difficulty, &state.ElapsedDays, &state.ScheduledDays, &state.Reps, &state.Lapses, &state.State, &state.LastReview)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return state, nil
}

func (db *SQLiteDB) GetUserQuestionStatesByUser(userID int) ([]*models.UserQuestionState, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review FROM user_question_states WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.UserQuestionState
	for rows.Next() {
		state := &models.UserQuestionState{}
		if err := rows.Scan(&state.ID, &state.UserID, &state.QuestionID, &state.Due, &state.Stability, &state.Difficulty, &state.ElapsedDays, &state.ScheduledDays, &state.Reps, &state.Lapses, &state.State, &state.LastReview); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (db *SQLiteDB) GetDueUserQuestionStates(userID int) ([]*models.UserQuestionState, error) {
	rows, err := db.db.Query(
		"SELECT id, user_id, question_id, due, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, state, last_review FROM user_question_states WHERE user_id = ? AND due <= datetime('now')",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var states []*models.UserQuestionState
	for rows.Next() {
		state := &models.UserQuestionState{}
		if err := rows.Scan(&state.ID, &state.UserID, &state.QuestionID, &state.Due, &state.Stability, &state.Difficulty, &state.ElapsedDays, &state.ScheduledDays, &state.Reps, &state.Lapses, &state.State, &state.LastReview); err != nil {
			return nil, err
		}
		states = append(states, state)
	}
	return states, rows.Err()
}

func (db *SQLiteDB) UpdateUserQuestionState(state *models.UserQuestionState) error {
	result, err := db.db.Exec(
		"UPDATE user_question_states SET due = ?, stability = ?, difficulty = ?, elapsed_days = ?, scheduled_days = ?, reps = ?, lapses = ?, state = ?, last_review = ? WHERE id = ?",
		state.Due, state.Stability, state.Difficulty, state.ElapsedDays, state.ScheduledDays, state.Reps, state.Lapses, state.State, state.LastReview, state.ID,
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

