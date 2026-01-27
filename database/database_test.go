package database

import (
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

// Test 6: Create and retrieve user
func TestCreateAndGetUser(t *testing.T) {
	db := NewInMemoryDB()

	user := &models.User{
		Username:    "john_doe",
		Email:       "john@example.com",
		IsAdmin:     false,
		TotalPoints: 0,
		CreatedAt:   time.Now(),
	}

	err := db.CreateUser(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.ID != 1 {
		t.Errorf("Expected user ID to be 1, got %d", user.ID)
	}

	retrieved, err := db.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}

	if retrieved.Username != user.Username {
		t.Errorf("Expected username '%s', got '%s'", user.Username, retrieved.Username)
	}
}

// Test 7: Get non-existent user
func TestGetNonExistentUser(t *testing.T) {
	db := NewInMemoryDB()

	_, err := db.GetUser(999)
	if err != ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

// Test 8: Update user points
func TestUpdateUserPoints(t *testing.T) {
	db := NewInMemoryDB()

	user := &models.User{
		Username:    "jane_doe",
		Email:       "jane@example.com",
		TotalPoints: 50,
		CreatedAt:   time.Now(),
	}

	db.CreateUser(user)

	err := db.UpdateUserPoints(user.ID, 150)
	if err != nil {
		t.Fatalf("Failed to update user points: %v", err)
	}

	updated, _ := db.GetUser(user.ID)
	if updated.TotalPoints != 150 {
		t.Errorf("Expected points to be 150, got %d", updated.TotalPoints)
	}
}

// Test 9: Create and list topics
func TestCreateAndListTopics(t *testing.T) {
	db := NewInMemoryDB()

	topic1 := &models.Topic{
		Name:        "JavaScript",
		Description: "Learn JavaScript fundamentals",
		CreatedAt:   time.Now(),
	}
	topic2 := &models.Topic{
		Name:        "Python",
		Description: "Learn Python basics",
		CreatedAt:   time.Now(),
	}

	db.CreateTopic(topic1)
	db.CreateTopic(topic2)

	topics, err := db.ListTopics()
	if err != nil {
		t.Fatalf("Failed to list topics: %v", err)
	}

	if len(topics) != 2 {
		t.Errorf("Expected 2 topics, got %d", len(topics))
	}
}

// Test 10: Get questions by topic
func TestGetQuestionsByTopic(t *testing.T) {
	db := NewInMemoryDB()

	question1 := &models.Question{
		TopicID:         1,
		QuestionText:    "What is a variable?",
		Difficulty:      1,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	question2 := &models.Question{
		TopicID:         1,
		QuestionText:    "What is a function?",
		Difficulty:      2,
		TopicQuestionID: "1.2",
		CreatedAt:       time.Now(),
	}
	question3 := &models.Question{
		TopicID:         2,
		QuestionText:    "What is a class?",
		Difficulty:      3,
		TopicQuestionID: "2.1",
		CreatedAt:       time.Now(),
	}

	db.CreateQuestion(question1)
	db.CreateQuestion(question2)
	db.CreateQuestion(question3)

	questions, err := db.GetQuestionsByTopic(1)
	if err != nil {
		t.Fatalf("Failed to get questions: %v", err)
	}

	if len(questions) != 2 {
		t.Errorf("Expected 2 questions for topic 1, got %d", len(questions))
	}
}

// Test 11: Create and get answers for question
func TestCreateAndGetAnswers(t *testing.T) {
	db := NewInMemoryDB()

	answer1 := &models.Answer{
		QuestionID: 1,
		AnswerText: "Correct answer",
		IsCorrect:  true,
	}
	answer2 := &models.Answer{
		QuestionID: 1,
		AnswerText: "Wrong answer",
		IsCorrect:  false,
	}

	db.CreateAnswer(answer1)
	db.CreateAnswer(answer2)

	answers, err := db.GetAnswersByQuestion(1)
	if err != nil {
		t.Fatalf("Failed to get answers: %v", err)
	}

	if len(answers) != 2 {
		t.Errorf("Expected 2 answers, got %d", len(answers))
	}

	correctCount := 0
	for _, ans := range answers {
		if ans.IsCorrect {
			correctCount++
		}
	}
	if correctCount != 1 {
		t.Errorf("Expected 1 correct answer, got %d", correctCount)
	}
}

// Test 12: Create and get quiz attempts
func TestCreateAndGetQuizAttempt(t *testing.T) {
	db := NewInMemoryDB()

	attempt := &models.QuizAttempt{
		UserID:          1,
		StartTime:       time.Now(),
		TotalQuestions:  5,
		CorrectAnswers:  3,
		QuestionIDs:     []string{"1.1", "1.2"},
		QuestionResults: map[string]bool{"1.1": true, "1.2": false},
	}

	err := db.CreateQuizAttempt(attempt)
	if err != nil {
		t.Fatalf("Failed to create quiz attempt: %v", err)
	}

	retrieved, err := db.GetQuizAttempt(attempt.ID)
	if err != nil {
		t.Fatalf("Failed to get quiz attempt: %v", err)
	}

	if retrieved.CorrectAnswers != 3 {
		t.Errorf("Expected 3 correct answers, got %d", retrieved.CorrectAnswers)
	}
}

// Test 13: Get quiz attempts by user
func TestGetQuizAttemptsByUser(t *testing.T) {
	db := NewInMemoryDB()

	attempt1 := &models.QuizAttempt{
		UserID:         1,
		StartTime:      time.Now(),
		TotalQuestions: 5,
		CorrectAnswers: 3,
	}
	attempt2 := &models.QuizAttempt{
		UserID:         1,
		StartTime:      time.Now(),
		TotalQuestions: 10,
		CorrectAnswers: 7,
	}
	attempt3 := &models.QuizAttempt{
		UserID:         2,
		StartTime:      time.Now(),
		TotalQuestions: 5,
		CorrectAnswers: 4,
	}

	db.CreateQuizAttempt(attempt1)
	db.CreateQuizAttempt(attempt2)
	db.CreateQuizAttempt(attempt3)

	attempts, err := db.GetQuizAttemptsByUser(1)
	if err != nil {
		t.Fatalf("Failed to get quiz attempts: %v", err)
	}

	if len(attempts) != 2 {
		t.Errorf("Expected 2 quiz attempts for user 1, got %d", len(attempts))
	}
}

// Test 14: Create and get question attempts
func TestCreateAndGetQuestionAttempt(t *testing.T) {
	db := NewInMemoryDB()

	attempt := &models.QuestionAttempt{
		UserID:             1,
		QuestionID:         1,
		TopicQuestionID:    "1.1",
		CorrectAttempts:    2,
		IncorrectAttempts:  1,
		LastAttemptTime:    time.Now(),
		LastAttemptCorrect: true,
	}

	err := db.CreateQuestionAttempt(attempt)
	if err != nil {
		t.Fatalf("Failed to create question attempt: %v", err)
	}

	retrieved, err := db.GetQuestionAttempt(attempt.ID)
	if err != nil {
		t.Fatalf("Failed to get question attempt: %v", err)
	}

	if retrieved.CorrectAttempts != 2 {
		t.Errorf("Expected 2 correct attempts, got %d", retrieved.CorrectAttempts)
	}
}

// Test 15: Concurrent user creation
func TestConcurrentUserCreation(t *testing.T) {
	db := NewInMemoryDB()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(idx int) {
			user := &models.User{
				Username:  "concurrentuser",
				Email:     "user@example.com",
				CreatedAt: time.Now(),
			}
			db.CreateUser(user)
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Check that all users were created with unique IDs
	if db.nextUserID != 11 {
		t.Errorf("Expected next user ID to be 11, got %d", db.nextUserID)
	}
}
