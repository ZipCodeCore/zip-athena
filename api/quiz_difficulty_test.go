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

func TestGenerateQuizDifficultyFiltering(t *testing.T) {
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
	db.CreateUser(user)

	// Create a test topic
	topic := &models.Topic{
		Name:        "Test Topic",
		Description: "Test topic for difficulty filtering",
		Category:    "java",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create questions with different difficulty levels
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "Easy question 1",
			Difficulty:      1,
			TopicQuestionID: "1.1",
			Tags:            []string{"test"},
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Easy question 2",
			Difficulty:      1,
			TopicQuestionID: "1.2",
			Tags:            []string{"test"},
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Medium question 1",
			Difficulty:      2,
			TopicQuestionID: "2.1",
			Tags:            []string{"test"},
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Medium question 2",
			Difficulty:      2,
			TopicQuestionID: "2.2",
			Tags:            []string{"test"},
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Hard question 1",
			Difficulty:      3,
			TopicQuestionID: "3.1",
			Tags:            []string{"test"},
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "Hard question 2",
			Difficulty:      3,
			TopicQuestionID: "3.2",
			Tags:            []string{"test"},
			CreatedAt:       time.Now(),
		},
	}

	for _, q := range questions {
		db.CreateQuestion(q)
		// Add at least one answer for each question
		db.CreateAnswer(&models.Answer{
			QuestionID: q.ID,
			AnswerText: "Test answer",
			IsCorrect:  true,
		})
	}

	tests := []struct {
		name               string
		difficultyLevel    int
		expectedMaxDiff    int
		expectedMinQuestions int
	}{
		{
			name:               "Level 1 - Only easiest questions",
			difficultyLevel:    1,
			expectedMaxDiff:    1,
			expectedMinQuestions: 2, // Should have 2 questions with difficulty 1
		},
		{
			name:               "Level 2 - Easy and medium questions",
			difficultyLevel:    2,
			expectedMaxDiff:    2,
			expectedMinQuestions: 4, // Should have 2 easy + 2 medium = 4 questions
		},
		{
			name:               "Level 3 - All difficulty levels",
			difficultyLevel:    3,
			expectedMaxDiff:    3,
			expectedMinQuestions: 6, // Should have all 6 questions
		},
		{
			name:               "No difficulty filter - All questions",
			difficultyLevel:    0,
			expectedMaxDiff:    3,
			expectedMinQuestions: 6, // Should have all 6 questions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := GenerateQuizRequest{
				UserID:          user.ID,
				TopicID:         topic.ID,
				NumQuestions:    10, // Request more than available to get all questions
				DifficultyLevel: tt.difficultyLevel,
			}

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.GenerateQuizHandler(w, req)

			if w.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
			}

			var resp GenerateQuizResponse
			if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// Check that we got questions
			if len(resp.Questions) < tt.expectedMinQuestions {
				t.Errorf("Expected at least %d questions, got %d", tt.expectedMinQuestions, len(resp.Questions))
			}

			// Check that all questions have difficulty <= expected max
			for _, q := range resp.Questions {
				if q.Question.Difficulty > tt.expectedMaxDiff {
					t.Errorf("Found question with difficulty %d, expected max %d for level %d",
						q.Question.Difficulty, tt.expectedMaxDiff, tt.difficultyLevel)
				}
			}
		})
	}
}

func TestGenerateQuizNoDifficultyQuestionsError(t *testing.T) {
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
	db.CreateUser(user)

	// Create a test topic
	topic := &models.Topic{
		Name:        "Test Topic",
		Description: "Test topic with only hard questions",
		Category:    "java",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create only hard questions
	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Very hard question",
		Difficulty:      7,
		TopicQuestionID: "7.1",
		Tags:            []string{"test"},
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)
	db.CreateAnswer(&models.Answer{
		QuestionID: question.ID,
		AnswerText: "Test answer",
		IsCorrect:  true,
	})

	// Try to generate quiz with difficulty level 1 (should fail)
	reqBody := GenerateQuizRequest{
		UserID:          user.ID,
		TopicID:         topic.ID,
		NumQuestions:    10,
		DifficultyLevel: 1,
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.GenerateQuizHandler(w, req)

	// Should return 404 with appropriate error message
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	responseBody := w.Body.String()
	expectedError := "No questions found for the selected difficulty level"
	if responseBody != expectedError+"\n" {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, responseBody)
	}
}
