package models

import (
	"testing"
	"time"
)

// Test 1: User model creation and fields
func TestUserModel(t *testing.T) {
	user := User{
		ID:          1,
		Username:    "testuser",
		Email:       "test@example.com",
		IsAdmin:     false,
		TotalPoints: 100,
		CreatedAt:   time.Now(),
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	if user.TotalPoints != 100 {
		t.Errorf("Expected total points 100, got %d", user.TotalPoints)
	}
}

// Test 2: Topic model creation and fields
func TestTopicModel(t *testing.T) {
	topic := Topic{
		ID:          1,
		Name:        "Go Programming",
		Description: "Learn Go basics",
		CreatedAt:   time.Now(),
	}

	if topic.Name != "Go Programming" {
		t.Errorf("Expected name 'Go Programming', got '%s'", topic.Name)
	}
	if topic.Description != "Learn Go basics" {
		t.Errorf("Expected description 'Learn Go basics', got '%s'", topic.Description)
	}
}

// Test 3: Question model with tags
func TestQuestionModelWithTags(t *testing.T) {
	question := Question{
		ID:              1,
		TopicID:         1,
		QuestionText:    "What is a goroutine?",
		Difficulty:      3,
		Tags:            []string{"concurrency", "basics"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}

	if len(question.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(question.Tags))
	}
	if question.Difficulty != 3 {
		t.Errorf("Expected difficulty 3, got %d", question.Difficulty)
	}
	if question.TopicQuestionID != "1.1" {
		t.Errorf("Expected topic question ID '1.1', got '%s'", question.TopicQuestionID)
	}
}

// Test 4: Answer model with correctness flag
func TestAnswerModel(t *testing.T) {
	answer := Answer{
		ID:         1,
		QuestionID: 1,
		AnswerText: "A lightweight thread managed by Go runtime",
		IsCorrect:  true,
	}

	if !answer.IsCorrect {
		t.Error("Expected answer to be correct")
	}
	if answer.AnswerText == "" {
		t.Error("Answer text should not be empty")
	}
}

// Test 5: QuizAttempt model with results tracking
func TestQuizAttemptModel(t *testing.T) {
	now := time.Now()
	completed := now.Add(30 * time.Minute)

	quizAttempt := QuizAttempt{
		ID:              1,
		UserID:          1,
		StartTime:       now,
		CompletedTime:   &completed,
		TotalQuestions:  10,
		CorrectAnswers:  7,
		QuestionIDs:     []string{"1.1", "1.2", "2.1"},
		QuestionResults: map[string]bool{"1.1": true, "1.2": false, "2.1": true},
	}

	if quizAttempt.TotalQuestions != 10 {
		t.Errorf("Expected 10 total questions, got %d", quizAttempt.TotalQuestions)
	}
	if quizAttempt.CorrectAnswers != 7 {
		t.Errorf("Expected 7 correct answers, got %d", quizAttempt.CorrectAnswers)
	}
	if len(quizAttempt.QuestionResults) != 3 {
		t.Errorf("Expected 3 question results, got %d", len(quizAttempt.QuestionResults))
	}
}
