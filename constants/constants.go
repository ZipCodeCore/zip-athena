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
// Topics are assigned one of three category values:
//   - "java"   – visible to java students (and admins)
//   - "python" – visible to data students (and admins)
//   - "all"    – visible to all students regardless of category
const (
	TopicCategoryJava   = "java"
	TopicCategoryPython = "python"
	TopicCategoryAll    = "all"
)

// User Roles (using IsAdmin bool in User model, but defining for clarity)
const (
	RoleStudent = "student"
	RoleAdmin   = "admin"
)
