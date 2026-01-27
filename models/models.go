package models

import "time"

// User represents a student or admin in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"` // Only used for auth, not returned in API
	IsAdmin      bool      `json:"is_admin"`
	TotalPoints  int       `json:"total_points"`
	CreatedAt    time.Time `json:"created_at"`
}

// Topic represents a subject area with multiple questions
type Topic struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Question represents a quiz question
type Question struct {
	ID              int       `json:"id"`
	TopicID         int       `json:"topic_id"`
	QuestionText    string    `json:"question_text"`
	Difficulty      int       `json:"difficulty"` // 1-7, 1 being easiest
	Tags            []string  `json:"tags"`
	TopicQuestionID string    `json:"topic_question_id"` // e.g., "3.17"
	CreatedAt       time.Time `json:"created_at"`
}

// Answer represents a possible answer to a question
type Answer struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"question_id"`
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"is_correct"`
}

// QuizAttempt represents a student's attempt at taking a quiz
type QuizAttempt struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	StartTime          time.Time `json:"start_time"`
	CompletedTime      *time.Time `json:"completed_time,omitempty"`
	TotalQuestions     int       `json:"total_questions"`
	CorrectAnswers     int       `json:"correct_answers"`
	QuestionIDs        []string  `json:"question_ids"` // List of topic_question_id
	QuestionResults    map[string]bool `json:"question_results"` // topic_question_id -> correct/incorrect
}

// QuestionAttempt tracks individual question attempts by a student
type QuestionAttempt struct {
	ID                  int       `json:"id"`
	UserID              int       `json:"user_id"`
	QuestionID          int       `json:"question_id"`
	TopicQuestionID     string    `json:"topic_question_id"`
	CorrectAttempts     int       `json:"correct_attempts"`
	IncorrectAttempts   int       `json:"incorrect_attempts"`
	LastAttemptTime     time.Time `json:"last_attempt_time"`
	LastAttemptCorrect  bool      `json:"last_attempt_correct"`
}

// Quiz represents a structured quiz with a topic and number (e.g., java.1, python.2)
type Quiz struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`        // e.g., "java.1", "python.2"
	TopicID     int       `json:"topic_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserQuestionState tracks FSRS state for each user-question pair
type UserQuestionState struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	QuestionID    int       `json:"question_id"`
	Due           time.Time `json:"due"`
	Stability     float64   `json:"stability"`
	Difficulty    float64   `json:"difficulty"`
	ElapsedDays   int       `json:"elapsed_days"`
	ScheduledDays int       `json:"scheduled_days"`
	Reps          int       `json:"reps"`
	Lapses        int       `json:"lapses"`
	State         int       `json:"state"` // 0=New, 1=Learning, 2=Review, 3=Relearning
	LastReview    time.Time `json:"last_review"`
}
