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
	QuizID             *int      `json:"quiz_id,omitempty"`
	ParentAttemptID    *int      `json:"parent_attempt_id,omitempty"` // Reference to original attempt if this is a retake
	AttemptNumber      int       `json:"attempt_number"` // 1 for first attempt, 2 for first retake, etc.
	IsRetake           bool      `json:"is_retake"` // True if this is a retake attempt
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

// Test represents a formal assessment (Phase 3)
type Test struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	TopicID     int       `json:"topic_id"`
	Description string    `json:"description"`
	TimeLimit   int       `json:"time_limit"` // in minutes, 0 for unlimited
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// QuizQuestion represents many-to-many relationship between Quiz and Question (Phase 3)
type QuizQuestion struct {
	ID         int `json:"id"`
	QuizID     int `json:"quiz_id"`
	QuestionID int `json:"question_id"`
	Order      int `json:"order"` // Order of question in quiz
}

// TestQuestion represents many-to-many relationship between Test and Question (Phase 3)
type TestQuestion struct {
	ID         int `json:"id"`
	TestID     int `json:"test_id"`
	QuestionID int `json:"question_id"`
	Order      int `json:"order"` // Order of question in test
}

// StudentProgress represents analytics data for a student (Phase 3)
type StudentProgress struct {
	UserID              int                `json:"user_id"`
	Username            string             `json:"username"`
	TotalQuizzes        int                `json:"total_quizzes"`
	CompletedQuizzes    int                `json:"completed_quizzes"`
	AverageScore        float64            `json:"average_score"`
	TotalPoints         int                `json:"total_points"`
	QuestionAccuracy    float64            `json:"question_accuracy"`
	TopicBreakdown      map[string]float64 `json:"topic_breakdown"`      // topic -> accuracy
	RecentActivity      []QuizAttempt      `json:"recent_activity"`
	FSRSMetrics         FSRSMetrics        `json:"fsrs_metrics"`
}

// FSRSMetrics represents FSRS-related statistics (Phase 3)
type FSRSMetrics struct {
	TotalCards     int     `json:"total_cards"`
	DueCards       int     `json:"due_cards"`
	NewCards       int     `json:"new_cards"`
	ReviewCards    int     `json:"review_cards"`
	AvgStability   float64 `json:"avg_stability"`
	AvgDifficulty  float64 `json:"avg_difficulty"`
}

// CohortAnalytics represents analytics for all students (Phase 3)
type CohortAnalytics struct {
	TotalStudents       int                    `json:"total_students"`
	ActiveStudents      int                    `json:"active_students"`      // students active in last 7 days
	AverageCohortScore  float64                `json:"average_cohort_score"`
	TotalQuizzesTaken   int                    `json:"total_quizzes_taken"`
	TopPerformers       []StudentProgress      `json:"top_performers"`       // top 5 students
	BottomPerformers    []StudentProgress      `json:"bottom_performers"`    // bottom 5 students
	TopicDifficulty     map[string]float64     `json:"topic_difficulty"`     // topic -> avg difficulty
	QuestionStats       []QuestionStatistic    `json:"question_stats"`
}

// QuestionStatistic represents statistics for a specific question (Phase 3)
type QuestionStatistic struct {
	QuestionID      int     `json:"question_id"`
	TopicQuestionID string  `json:"topic_question_id"`
	QuestionText    string  `json:"question_text"`
	TotalAttempts   int     `json:"total_attempts"`
	CorrectAttempts int     `json:"correct_attempts"`
	Accuracy        float64 `json:"accuracy"`
	AvgDifficulty   float64 `json:"avg_difficulty"`
}

// QuizExport represents the export format for quizzes (Phase 3)
type QuizExport struct {
	Quiz      Quiz                `json:"quiz"`
	Questions []QuestionWithAnswers `json:"questions"`
}

// QuestionWithAnswers represents a question with its answers for export (Phase 3)
type QuestionWithAnswers struct {
	Question Question `json:"question"`
	Answers  []Answer `json:"answers"`
}

// QuizImport represents the import format for quizzes (Phase 3)
type QuizImport struct {
	Name        string                `json:"name"`
	TopicID     int                   `json:"topic_id"`
	Description string                `json:"description"`
	Questions   []QuestionWithAnswers `json:"questions"`
}
