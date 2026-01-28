package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

func TestRetakeQuizHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create questions
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a goroutine?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.1",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a channel?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.2",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is defer?",
			Difficulty:      2,
			Tags:            []string{"basics"},
			TopicQuestionID: "go.3",
			CreatedAt:       time.Now(),
		},
	}

	for _, q := range questions {
		db.CreateQuestion(q)

		// Create answers for each question
		db.CreateAnswer(&models.Answer{
			QuestionID: q.ID,
			AnswerText: "Correct answer",
			IsCorrect:  true,
		})
		db.CreateAnswer(&models.Answer{
			QuestionID: q.ID,
			AnswerText: "Wrong answer",
			IsCorrect:  false,
		})
	}

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create a completed quiz attempt with less than 100% score
	completedTime := time.Now()
	parentAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-30 * time.Minute),
		CompletedTime:   &completedTime,
		TotalQuestions:  3,
		CorrectAnswers:  2, // 2 out of 3 correct (66.67%)
		QuestionIDs:     []string{"go.1", "go.2", "go.3"},
		QuestionResults: map[string]bool{
			"go.1": true,
			"go.2": true,
			"go.3": false, // This one was incorrect
		},
	}
	db.CreateQuizAttempt(parentAttempt)

	// Test retake generation
	reqBody := RetakeQuizRequest{
		ParentAttemptID: parentAttempt.ID,
		UserID:          user.ID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.RetakeQuizHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var response GenerateQuizResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify the retake quiz has same questions
	if len(response.Questions) != 3 {
		t.Errorf("Expected 3 questions, got %d", len(response.Questions))
	}

	// Get the retake attempt from database
	retakeAttempt, err := db.GetQuizAttempt(response.QuizAttemptID)
	if err != nil {
		t.Fatalf("Failed to get retake attempt: %v", err)
	}

	// Verify retake properties
	if retakeAttempt.IsRetake != true {
		t.Errorf("Expected IsRetake to be true, got false")
	}

	if retakeAttempt.AttemptNumber != 2 {
		t.Errorf("Expected AttemptNumber to be 2, got %d", retakeAttempt.AttemptNumber)
	}

	if retakeAttempt.ParentAttemptID == nil || *retakeAttempt.ParentAttemptID != parentAttempt.ID {
		t.Errorf("Expected ParentAttemptID to be %d, got %v", parentAttempt.ID, retakeAttempt.ParentAttemptID)
	}

	// Verify same questions
	if len(retakeAttempt.QuestionIDs) != len(parentAttempt.QuestionIDs) {
		t.Errorf("Expected %d questions, got %d", len(parentAttempt.QuestionIDs), len(retakeAttempt.QuestionIDs))
	}

	for i, qid := range parentAttempt.QuestionIDs {
		if retakeAttempt.QuestionIDs[i] != qid {
			t.Errorf("Question mismatch at index %d: expected %s, got %s", i, qid, retakeAttempt.QuestionIDs[i])
		}
	}
}

func TestRetakeQuizWithPerfectScore(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create a completed quiz attempt with 100% score
	completedTime := time.Now()
	parentAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-30 * time.Minute),
		CompletedTime:   &completedTime,
		TotalQuestions:  3,
		CorrectAnswers:  3, // 3 out of 3 correct (100%)
		QuestionIDs:     []string{"go.1", "go.2", "go.3"},
		QuestionResults: map[string]bool{
			"go.1": true,
			"go.2": true,
			"go.3": true,
		},
	}
	db.CreateQuizAttempt(parentAttempt)

	// Test retake generation - should fail
	reqBody := RetakeQuizRequest{
		ParentAttemptID: parentAttempt.ID,
		UserID:          user.ID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.RetakeQuizHandler(w, req)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRetakeQuizIncompleteAttempt(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create an incomplete quiz attempt
	parentAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-30 * time.Minute),
		CompletedTime:   nil, // Not completed
		TotalQuestions:  3,
		CorrectAnswers:  0,
		QuestionIDs:     []string{"go.1", "go.2", "go.3"},
		QuestionResults: map[string]bool{},
	}
	db.CreateQuizAttempt(parentAttempt)

	// Test retake generation - should fail
	reqBody := RetakeQuizRequest{
		ParentAttemptID: parentAttempt.ID,
		UserID:          user.ID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.RetakeQuizHandler(w, req)

	// Should return 400 Bad Request
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRetakeQuizUnauthorizedUser(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create two users
	user1 := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user1)

	user2 := &models.User{
		Username: "student2",
		Email:    "student2@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user2)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create a completed quiz attempt for user1
	completedTime := time.Now()
	parentAttempt := &models.QuizAttempt{
		UserID:          user1.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-30 * time.Minute),
		CompletedTime:   &completedTime,
		TotalQuestions:  3,
		CorrectAnswers:  2,
		QuestionIDs:     []string{"go.1", "go.2", "go.3"},
		QuestionResults: map[string]bool{
			"go.1": true,
			"go.2": true,
			"go.3": false,
		},
	}
	db.CreateQuizAttempt(parentAttempt)

	// Test retake generation by user2 - should fail
	reqBody := RetakeQuizRequest{
		ParentAttemptID: parentAttempt.ID,
		UserID:          user2.ID, // Different user
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.RetakeQuizHandler(w, req)

	// Should return 401 Unauthorized
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestMultipleRetakes(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create questions
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a goroutine?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.1",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a channel?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.2",
			CreatedAt:       time.Now(),
		},
	}

	for _, q := range questions {
		db.CreateQuestion(q)
		db.CreateAnswer(&models.Answer{
			QuestionID: q.ID,
			AnswerText: "Correct answer",
			IsCorrect:  true,
		})
	}

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create initial attempt with partial score
	completedTime1 := time.Now()
	attempt1 := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-60 * time.Minute),
		CompletedTime:   &completedTime1,
		TotalQuestions:  2,
		CorrectAnswers:  1,
		QuestionIDs:     []string{"go.1", "go.2"},
		QuestionResults: map[string]bool{
			"go.1": true,
			"go.2": false,
		},
	}
	db.CreateQuizAttempt(attempt1)

	// First retake
	reqBody1 := RetakeQuizRequest{
		ParentAttemptID: attempt1.ID,
		UserID:          user.ID,
	}
	body1, _ := json.Marshal(reqBody1)
	req1 := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body1))
	w1 := httptest.NewRecorder()
	handler.RetakeQuizHandler(w1, req1)

	var response1 GenerateQuizResponse
	json.NewDecoder(w1.Body).Decode(&response1)

	// Mark first retake as completed with partial score
	attempt2, _ := db.GetQuizAttempt(response1.QuizAttemptID)
	completedTime2 := time.Now()
	attempt2.CompletedTime = &completedTime2
	attempt2.CorrectAnswers = 1
	attempt2.QuestionResults = map[string]bool{
		"go.1": true,
		"go.2": false,
	}
	db.UpdateQuizAttempt(attempt2)

	// Verify first retake properties
	if attempt2.AttemptNumber != 2 {
		t.Errorf("First retake: Expected AttemptNumber 2, got %d", attempt2.AttemptNumber)
	}
	if *attempt2.ParentAttemptID != attempt1.ID {
		t.Errorf("First retake: Expected ParentAttemptID %d, got %d", attempt1.ID, *attempt2.ParentAttemptID)
	}

	// Second retake
	reqBody2 := RetakeQuizRequest{
		ParentAttemptID: attempt2.ID,
		UserID:          user.ID,
	}
	body2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body2))
	w2 := httptest.NewRecorder()
	handler.RetakeQuizHandler(w2, req2)

	var response2 GenerateQuizResponse
	json.NewDecoder(w2.Body).Decode(&response2)

	// Verify second retake properties
	attempt3, _ := db.GetQuizAttempt(response2.QuizAttemptID)
	if attempt3.AttemptNumber != 3 {
		t.Errorf("Second retake: Expected AttemptNumber 3, got %d", attempt3.AttemptNumber)
	}
	if *attempt3.ParentAttemptID != attempt2.ID {
		t.Errorf("Second retake: Expected ParentAttemptID %d, got %d", attempt2.ID, *attempt3.ParentAttemptID)
	}

	// Verify all attempts have same questions
	if len(attempt1.QuestionIDs) != len(attempt2.QuestionIDs) || len(attempt1.QuestionIDs) != len(attempt3.QuestionIDs) {
		t.Error("All attempts should have same number of questions")
	}
}
