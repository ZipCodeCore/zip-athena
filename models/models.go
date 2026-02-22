package models

import "time"

// User represents a student or admin in the system
type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Password    string    `json:"password,omitempty"` // Only used for auth, not returned in API
	IsAdmin     bool      `json:"is_admin"`
	Category    string    `json:"category"` // "java" or "data" - determines accessible topics and resources
	TotalPoints int       `json:"total_points"`
	CreatedAt   time.Time `json:"created_at"`
}

// Topic represents a subject area with multiple questions
type Topic struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`      // "java", "python", "sql" - determines which user categories can access
	IsActive     bool      `json:"is_active"`     // For soft delete (Phase 3)
	DisplayOrder int       `json:"display_order"` // For reordering (Phase 3)
	CreatedAt    time.Time `json:"created_at"`
}

// Question represents a quiz question
type Question struct {
	ID              int       `json:"id"`
	TopicID         int       `json:"topic_id"`
	QuestionText    string    `json:"question_text"`
	Difficulty      int       `json:"difficulty"` // 1-7, 1 being easiest
	Tags            []string  `json:"tags"`
	TopicQuestionID string    `json:"topic_question_id"` // e.g., "3.17"
	IsActive        bool      `json:"is_active"`         // For soft delete (Phase 5)
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
	ID                int             `json:"id"`
	UserID            int             `json:"user_id"`
	QuizID            *int            `json:"quiz_id,omitempty"`
	ParentAttemptID   *int            `json:"parent_attempt_id,omitempty"` // Reference to original attempt if this is a retake
	AttemptNumber     int             `json:"attempt_number"`              // 1 for first attempt, 2 for first retake, etc.
	IsRetake          bool            `json:"is_retake"`                   // True if this is a retake attempt
	StartTime         time.Time       `json:"start_time"`
	CompletedTime     *time.Time      `json:"completed_time,omitempty"`
	TotalQuestions    int             `json:"total_questions"`
	CorrectAnswers    int             `json:"correct_answers"`
	DifficultyAverage float64         `json:"difficulty_average"`  // Average difficulty of all questions in quiz (1-7 scale)
	Description       string          `json:"quiz_description"`    // Human-readable description (e.g., "Python Basics Quiz")
	QuestionIDs       []string        `json:"question_ids"`        // List of topic_question_id
	QuestionResults   map[string]bool `json:"question_results"`   // topic_question_id -> correct/incorrect
}

// QuestionAttempt tracks individual question attempts by a student
type QuestionAttempt struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	QuestionID         int       `json:"question_id"`
	TopicQuestionID    string    `json:"topic_question_id"`
	CorrectAttempts    int       `json:"correct_attempts"`
	IncorrectAttempts  int       `json:"incorrect_attempts"`
	LastAttemptTime    time.Time `json:"last_attempt_time"`
	LastAttemptCorrect bool      `json:"last_attempt_correct"`
}

// Quiz represents a structured quiz with a topic and number (e.g., java.1, python.2)
type Quiz struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"` // e.g., "java.1", "python.2"
	TopicID     int       `json:"topic_id"`
	Category    string    `json:"category"` // "java", "python", "sql", etc.
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AssignedQuiz represents a quiz assigned to a specific student
type AssignedQuiz struct {
	ID         int       `json:"id"`
	QuizID     int       `json:"quiz_id"`
	UserID     int       `json:"user_id"`
	AssignedAt time.Time `json:"assigned_at"`
	AssignedBy int       `json:"assigned_by"`
}

// AssignedQuizInfo combines quiz info with assignment details for student view
type AssignedQuizInfo struct {
	AssignedQuiz
	QuizName        string     `json:"quiz_name"`
	QuizDescription string     `json:"quiz_description"`
	QuizCategory    string     `json:"quiz_category"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	Score           *float64   `json:"score,omitempty"`
	CorrectAnswers  *int       `json:"correct_answers,omitempty"`
	TotalQuestions  *int       `json:"total_questions,omitempty"`
}

// AssignedQuizStudentResult represents a single student's result on an assigned quiz
type AssignedQuizStudentResult struct {
	UserID         int        `json:"user_id"`
	Username       string     `json:"username"`
	AssignedAt     time.Time  `json:"assigned_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	Score          *float64   `json:"score,omitempty"`
	CorrectAnswers *int       `json:"correct_answers,omitempty"`
	TotalQuestions *int       `json:"total_questions,omitempty"`
}

// AssignedQuizStatus represents the instructor status view for an assigned quiz
type AssignedQuizStatus struct {
	QuizID         int                         `json:"quiz_id"`
	QuizName       string                      `json:"quiz_name"`
	QuizCategory   string                      `json:"quiz_category"`
	QuizDescription string                     `json:"quiz_description"`
	StudentResults []AssignedQuizStudentResult `json:"student_results"`
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

// Test represents a formal assessment (Phase 3+4)
type Test struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	TopicID     int        `json:"topic_id"`
	Description string     `json:"description"`
	TimeLimit   int        `json:"time_limit"`           // in minutes, 0 for unlimited
	StartDate   *time.Time `json:"start_date,omitempty"` // Phase 4: availability window
	EndDate     *time.Time `json:"end_date,omitempty"`   // Phase 4: availability window
	IsPublished bool       `json:"is_published"`         // Phase 4: publish status
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
	UserID           int                `json:"user_id"`
	Username         string             `json:"username"`
	TotalQuizzes     int                `json:"total_quizzes"`
	CompletedQuizzes int                `json:"completed_quizzes"`
	AverageScore     float64            `json:"average_score"`
	TotalPoints      int                `json:"total_points"`
	QuestionAccuracy float64            `json:"question_accuracy"`
	TopicBreakdown   map[string]float64 `json:"topic_breakdown"` // topic -> accuracy
	RecentActivity   []QuizAttempt      `json:"recent_activity"`
	FSRSMetrics      FSRSMetrics        `json:"fsrs_metrics"`
}

// FSRSMetrics represents FSRS-related statistics (Phase 3)
type FSRSMetrics struct {
	TotalCards    int     `json:"total_cards"`
	DueCards      int     `json:"due_cards"`
	NewCards      int     `json:"new_cards"`
	ReviewCards   int     `json:"review_cards"`
	AvgStability  float64 `json:"avg_stability"`
	AvgDifficulty float64 `json:"avg_difficulty"`
}

// CohortAnalytics represents analytics for all students (Phase 3)
type CohortAnalytics struct {
	TotalStudents      int                 `json:"total_students"`
	ActiveStudents     int                 `json:"active_students"` // students active in last 7 days
	AverageCohortScore float64             `json:"average_cohort_score"`
	TotalQuizzesTaken  int                 `json:"total_quizzes_taken"`
	TopPerformers      []StudentProgress   `json:"top_performers"`    // top 5 students
	BottomPerformers   []StudentProgress   `json:"bottom_performers"` // bottom 5 students
	TopicDifficulty    map[string]float64  `json:"topic_difficulty"`  // topic -> avg difficulty
	QuestionStats      []QuestionStatistic `json:"question_stats"`
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
	Quiz      Quiz                  `json:"quiz"`
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

// WeakAreas represents a student's weak topics and tags (Phase 3)
type WeakAreas struct {
	UserID     int             `json:"user_id"`
	Username   string          `json:"username"`
	WeakTopics []TopicWeakness `json:"weak_topics"` // Topics with < 60% accuracy
	WeakTags   []TagWeakness   `json:"weak_tags"`   // Tags with < 60% accuracy
}

// TopicWeakness represents weakness in a specific topic (Phase 3)
type TopicWeakness struct {
	TopicID       int     `json:"topic_id"`
	TopicName     string  `json:"topic_name"`
	Accuracy      float64 `json:"accuracy"`
	TotalAttempts int     `json:"total_attempts"`
}

// TagWeakness represents weakness in a specific tag/concept (Phase 3)
type TagWeakness struct {
	Tag           string  `json:"tag"`
	Accuracy      float64 `json:"accuracy"`
	TotalAttempts int     `json:"total_attempts"`
}

// TestAttempt represents a student's attempt at taking a timed test (Phase 4)
type TestAttempt struct {
	ID              int             `json:"id"`
	UserID          int             `json:"user_id"`
	TestID          int             `json:"test_id"`
	StartTime       time.Time       `json:"start_time"`
	SubmitTime      *time.Time      `json:"submit_time,omitempty"`
	TimeRemaining   int             `json:"time_remaining"` // seconds remaining when submitted
	IsAutoSubmit    bool            `json:"is_auto_submit"` // true if auto-submitted when timer expired
	TotalQuestions  int             `json:"total_questions"`
	CorrectAnswers  int             `json:"correct_answers"`
	Score           float64         `json:"score"`            // percentage score
	Answers         map[string]int  `json:"answers"`          // question_id -> answer_id
	QuestionResults map[string]bool `json:"question_results"` // question_id -> correct/incorrect
}

// TagStats represents statistics for a specific tag (Phase 4)
type TagStats struct {
	Tag             string  `json:"tag"`
	QuestionCount   int     `json:"question_count"`
	TotalAttempts   int     `json:"total_attempts"`
	CorrectAttempts int     `json:"correct_attempts"`
	Accuracy        float64 `json:"accuracy"`
}

// QuestionAnalytics represents performance analytics for a question (Phase 4)
type QuestionAnalytics struct {
	QuestionID            int         `json:"question_id"`
	TopicQuestionID       string      `json:"topic_question_id"`
	QuestionText          string      `json:"question_text"`
	TopicID               int         `json:"topic_id"`
	TopicName             string      `json:"topic_name"`
	Tags                  []string    `json:"tags"`
	CurrentDifficulty     int         `json:"current_difficulty"`
	TotalAttempts         int         `json:"total_attempts"`
	CorrectAttempts       int         `json:"correct_attempts"`
	Accuracy              float64     `json:"accuracy"`
	RecommendedDifficulty int         `json:"recommended_difficulty"`
	FlagForReview         bool        `json:"flag_for_review"`     // true if too easy/hard
	AnswerDistribution    map[int]int `json:"answer_distribution"` // answer_id -> count
}

// TestResult represents results for a test (Phase 4)
type TestResult struct {
	TestID         int                 `json:"test_id"`
	TestName       string              `json:"test_name"`
	StudentResults []StudentTestResult `json:"student_results"`
}

// StudentTestResult represents a single student's test result (Phase 4)
type StudentTestResult struct {
	UserID         int       `json:"user_id"`
	Username       string    `json:"username"`
	Score          float64   `json:"score"`
	CorrectAnswers int       `json:"correct_answers"`
	TotalQuestions int       `json:"total_questions"`
	SubmitTime     time.Time `json:"submit_time"`
	TimeTaken      int       `json:"time_taken"` // in seconds
}

// LearningResource represents a learning guide document
type LearningResource struct {
	ID            int       `json:"id"`
	Language      string    `json:"language"`       // java, python
	Phase         int       `json:"phase"`          // 1-5
	SectionNumber string    `json:"section_number"` // e.g., "3", "5", "10"
	SectionTitle  string    `json:"section_title"`
	FilePath      string    `json:"file_path"`
	DisplayOrder  int       `json:"display_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

// QuestionLearningLink links questions to learning resources
type QuestionLearningLink struct {
	ID                 int       `json:"id"`
	QuestionID         int       `json:"question_id"`
	LearningResourceID int       `json:"learning_resource_id"`
	RelevanceScore     float64   `json:"relevance_score"`
	CreatedAt          time.Time `json:"created_at"`
}

// UserLearningProgress tracks student's learning material viewing
type UserLearningProgress struct {
	ID                   int       `json:"id"`
	UserID               int       `json:"user_id"`
	LearningResourceID   int       `json:"learning_resource_id"`
	FirstViewedAt        time.Time `json:"first_viewed_at"`
	LastViewedAt         time.Time `json:"last_viewed_at"`
	ViewCount            int       `json:"view_count"`
	CompletionPercentage float64   `json:"completion_percentage"`
}

// PhaseOutline represents a phase in the curriculum with its resources
type PhaseOutline struct {
	Phase       int                `json:"phase"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Resources   []LearningResource `json:"resources"`
}

// LearningResourceWithProgress combines a resource with user's progress
type LearningResourceWithProgress struct {
	LearningResource
	ViewCount            int        `json:"view_count"`
	LastViewedAt         *time.Time `json:"last_viewed_at,omitempty"`
	CompletionPercentage float64    `json:"completion_percentage"`
}

// Phase 5: Content Management & Quality

// QuestionPreview represents a question in preview mode (Phase 5)
type QuestionPreview struct {
	Question Question `json:"question"`
	Answers  []Answer `json:"answers"`
}

// CSVImportRow represents a single row from CSV import (Phase 5)
type CSVImportRow struct {
	Topic         string `json:"topic"`
	Question      string `json:"question"`
	Answer1       string `json:"answer1"`
	Answer2       string `json:"answer2"`
	Answer3       string `json:"answer3"`
	Answer4       string `json:"answer4"`
	CorrectAnswer string `json:"correct_answer"`
	Difficulty    string `json:"difficulty"`
	Tags          string `json:"tags"`
}

// CSVImportResult represents the result of CSV import validation (Phase 5)
type CSVImportResult struct {
	ValidRows   []CSVImportRow   `json:"valid_rows"`
	InvalidRows []CSVImportError `json:"invalid_rows"`
	Summary     CSVImportSummary `json:"summary"`
}

// CSVImportError represents an error in a CSV row (Phase 5)
type CSVImportError struct {
	Row    int          `json:"row"`
	Data   CSVImportRow `json:"data"`
	Errors []string     `json:"errors"`
}

// CSVImportSummary represents summary of import operation (Phase 5)
type CSVImportSummary struct {
	TotalRows    int `json:"total_rows"`
	ValidRows    int `json:"valid_rows"`
	InvalidRows  int `json:"invalid_rows"`
	ImportedRows int `json:"imported_rows"`
}

// Tag represents a tag with metadata (Phase 5)
type Tag struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	QuestionCount int    `json:"question_count"`
}

// FlaggedQuestion represents a question flagged for quality review (Phase 5)
type FlaggedQuestion struct {
	QuestionID      int       `json:"question_id"`
	TopicQuestionID string    `json:"topic_question_id"`
	QuestionText    string    `json:"question_text"`
	FlagReason      string    `json:"flag_reason"`
	FlaggedAt       time.Time `json:"flagged_at"`
	Accuracy        float64   `json:"accuracy,omitempty"`
	TotalAttempts   int       `json:"total_attempts,omitempty"`
}

// FlagReason constants (Phase 5)
const (
	FlagReasonTooHard         = "TOO_HARD"         // <30% accuracy
	FlagReasonTooEasy         = "TOO_EASY"         // >90% accuracy
	FlagReasonAmbiguous       = "AMBIGUOUS"        // Mixed feedback
	FlagReasonPoorDistractors = "POOR_DISTRACTORS" // Skewed distribution
	FlagReasonManual          = "MANUAL"           // Instructor flagged
)
