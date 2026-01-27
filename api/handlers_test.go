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

// Test 16: Create user via API
func TestCreateUserHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		IsAdmin:  false,
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateUserHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var created models.User
	json.NewDecoder(w.Body).Decode(&created)

	if created.Username != user.Username {
		t.Errorf("Expected username '%s', got '%s'", user.Username, created.Username)
	}
	if created.ID == 0 {
		t.Error("Expected user ID to be set")
	}
}

// Test 17: Get user via API
func TestGetUserHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user first
	user := &models.User{
		Username: "getuser",
		Email:    "get@example.com",
	}
	db.CreateUser(user)

	req := httptest.NewRequest(http.MethodGet, "/users?id=1", nil)
	w := httptest.NewRecorder()

	handler.GetUserHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var retrieved models.User
	json.NewDecoder(w.Body).Decode(&retrieved)

	if retrieved.Username != user.Username {
		t.Errorf("Expected username '%s', got '%s'", user.Username, retrieved.Username)
	}
}

// Test 18: Create topic via API
func TestCreateTopicHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	topic := models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
	}

	body, _ := json.Marshal(topic)
	req := httptest.NewRequest(http.MethodPost, "/topics", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateTopicHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var created models.Topic
	json.NewDecoder(w.Body).Decode(&created)

	if created.Name != topic.Name {
		t.Errorf("Expected name '%s', got '%s'", topic.Name, created.Name)
	}
}

// Test 19: List topics via API
func TestListTopicsHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create some topics
	db.CreateTopic(&models.Topic{Name: "JavaScript", Description: "JS basics"})
	db.CreateTopic(&models.Topic{Name: "Python", Description: "Python basics"})

	req := httptest.NewRequest(http.MethodGet, "/topics", nil)
	w := httptest.NewRecorder()

	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var topics []*models.Topic
	json.NewDecoder(w.Body).Decode(&topics)

	if len(topics) != 2 {
		t.Errorf("Expected 2 topics, got %d", len(topics))
	}
}

// Test 20: Create quiz attempt via API
func TestCreateQuizAttemptHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	attempt := models.QuizAttempt{
		UserID:          1,
		TotalQuestions:  10,
		CorrectAnswers:  7,
		QuestionIDs:     []string{"1.1", "1.2"},
		QuestionResults: map[string]bool{"1.1": true, "1.2": false},
	}

	body, _ := json.Marshal(attempt)
	req := httptest.NewRequest(http.MethodPost, "/quiz-attempts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateQuizAttemptHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var created models.QuizAttempt
	json.NewDecoder(w.Body).Decode(&created)

	if created.CorrectAnswers != attempt.CorrectAnswers {
		t.Errorf("Expected %d correct answers, got %d", attempt.CorrectAnswers, created.CorrectAnswers)
	}
}
