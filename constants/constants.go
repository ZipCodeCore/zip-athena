package constants

// FSRS States
const (
	StateNew        = 0
	StateLearning   = 1
	StateReview     = 2
	StateRelearning = 3
)

// Difficulty Levels (1-7 scale)
const (
	DifficultyTrivial  = 1
	DifficultyVeryEasy = 2
	DifficultyEasy     = 3
	DifficultyMedium   = 4
	DifficultyHard     = 5
	DifficultyVeryHard = 6
	DifficultyExpert   = 7
)

// Difficulty bounds
const (
	MinDifficulty = DifficultyTrivial
	MaxDifficulty = DifficultyExpert
)

// User Categories
const (
	CategoryJava = "java"
	CategoryData = "data"
)

// Topic Categories
const (
	TopicCategoryJava          = "java"
	TopicCategoryPython        = "python"
	TopicCategorySQL           = "sql"
	TopicCategoryJavaTracing   = "java-tracing"
	TopicCategoryPythonTracing = "python-tracing"
)

// User Roles (using IsAdmin bool in User model, but defining for clarity)
const (
	RoleStudent = "student"
	RoleAdmin   = "admin"
)
