package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

func TestGenerateQuizNoDuplicateQuestions(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a test user
	user := &models.User{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "password",
		Category: "java",
		IsAdmin:  false,
	}
	if err := db.CreateUser(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a test topic
	topic := &models.Topic{
		Name:        "Test Topic",
		Description: "Topic for testing duplicate questions",
		Category:    "java",
	}
	if err := db.CreateTopic(topic); err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}

	// Create questions with some having the same TopicQuestionID (duplicates)
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 1",
			Difficulty:      5,
			TopicQuestionID: "1.1",
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 2",
			Difficulty:      5,
			TopicQuestionID: "1.2",
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 3 (duplicate of 1.1)",
			Difficulty:      5,
			TopicQuestionID: "1.1", // Duplicate TopicQuestionID
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 4",
			Difficulty:      5,
			TopicQuestionID: "1.3",
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 5",
			Difficulty:      5,
			TopicQuestionID: "1.4",
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 6 (duplicate of 1.2)",
			Difficulty:      5,
			TopicQuestionID: "1.2", // Duplicate TopicQuestionID
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Question 7",
			Difficulty:      5,
			TopicQuestionID: "1.5",
		},
	}

	for _, q := range questions {
		if err := db.CreateQuestion(q); err != nil {
			t.Fatalf("Failed to create question: %v", err)
		}
		// Create dummy answers for each question
		for i := 0; i < 4; i++ {
			answer := &models.Answer{
				QuestionID: q.ID,
				AnswerText: "Test answer",
				IsCorrect:  i == 0,
			}
			if err := db.CreateAnswer(answer); err != nil {
				t.Fatalf("Failed to create answer: %v", err)
			}
		}
	}

	// Generate a quiz
	reqBody := GenerateQuizRequest{
		UserID:       user.ID,
		TopicID:      topic.ID,
		NumQuestions: 5,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/quiz/generate", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.GenerateQuizHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp GenerateQuizResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify no duplicate questions based on TopicQuestionID
	seenQuestionIDs := make(map[string]bool)
	for _, quizQuestion := range resp.Questions {
		topicQuestionID := quizQuestion.Question.TopicQuestionID
		if seenQuestionIDs[topicQuestionID] {
			t.Errorf("Duplicate question found in quiz: TopicQuestionID %s appears more than once", topicQuestionID)
		}
		seenQuestionIDs[topicQuestionID] = true
	}

	// Verify we got the requested number of unique questions (or less if not enough unique ones)
	if len(resp.Questions) > reqBody.NumQuestions {
		t.Errorf("Expected at most %d questions, got %d", reqBody.NumQuestions, len(resp.Questions))
	}

	// Verify the quiz attempt has unique question IDs
	attempt, err := db.GetQuizAttempt(resp.QuizAttemptID)
	if err != nil {
		t.Fatalf("Failed to get quiz attempt: %v", err)
	}

	seenAttemptQuestionIDs := make(map[string]bool)
	for _, qid := range attempt.QuestionIDs {
		if seenAttemptQuestionIDs[qid] {
			t.Errorf("Duplicate question ID found in quiz attempt: %s appears more than once", qid)
		}
		seenAttemptQuestionIDs[qid] = true
	}

	t.Logf("Successfully generated quiz with %d unique questions (requested %d)", len(resp.Questions), reqBody.NumQuestions)
}

func TestGenerateQuizAllUniqueQuestions(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a test user
	user := &models.User{
		Username: "testuser2",
		Email:    "test2@test.com",
		Password: "password",
		Category: "java",
		IsAdmin:  false,
	}
	if err := db.CreateUser(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a test topic
	topic := &models.Topic{
		Name:        "Test Topic 2",
		Description: "Topic with all unique questions",
		Category:    "java",
	}
	if err := db.CreateTopic(topic); err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}

	// Create 10 unique questions
	for i := 1; i <= 10; i++ {
		question := &models.Question{
			TopicID:         topic.ID,
			QuestionText:    "Question " + string(rune(i)),
			Difficulty:      5,
			TopicQuestionID: "2." + string(rune(i)),
		}
		if err := db.CreateQuestion(question); err != nil {
			t.Fatalf("Failed to create question: %v", err)
		}
		// Create dummy answers
		for j := 0; j < 4; j++ {
			answer := &models.Answer{
				QuestionID: question.ID,
				AnswerText: "Test answer",
				IsCorrect:  j == 0,
			}
			if err := db.CreateAnswer(answer); err != nil {
				t.Fatalf("Failed to create answer: %v", err)
			}
		}
	}

	// Generate a quiz with all questions
	reqBody := GenerateQuizRequest{
		UserID:       user.ID,
		TopicID:      topic.ID,
		NumQuestions: 10,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/quiz/generate", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.GenerateQuizHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp GenerateQuizResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Should get exactly 10 unique questions
	if len(resp.Questions) != 10 {
		t.Errorf("Expected 10 questions, got %d", len(resp.Questions))
	}

	// Verify all are unique
	seenQuestionIDs := make(map[string]bool)
	for _, quizQuestion := range resp.Questions {
		topicQuestionID := quizQuestion.Question.TopicQuestionID
		if seenQuestionIDs[topicQuestionID] {
			t.Errorf("Duplicate question found: %s", topicQuestionID)
		}
		seenQuestionIDs[topicQuestionID] = true
	}

	t.Logf("Successfully verified all %d questions are unique", len(resp.Questions))
}
