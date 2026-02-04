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

func TestGenerateQuizCategoryFiltering(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create users with different categories
	javaUser := &models.User{
		Username: "javastudent",
		Email:    "java@test.com",
		Password: "password",
		Category: "java",
		IsAdmin:  false,
	}
	db.CreateUser(javaUser)

	dataUser := &models.User{
		Username: "datastudent",
		Email:    "data@test.com",
		Password: "password",
		Category: "data",
		IsAdmin:  false,
	}
	db.CreateUser(dataUser)

	adminUser := &models.User{
		Username: "adminuser",
		Email:    "admin@test.com",
		Password: "password",
		Category: "java",
		IsAdmin:  true,
	}
	db.CreateUser(adminUser)

	// Create topics
	javaTopic := &models.Topic{
		Name:        "Java Programming",
		Description: "Learn Java",
		Category:    "java",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(javaTopic)

	pythonTopic := &models.Topic{
		Name:        "Python Programming",
		Description: "Learn Python",
		Category:    "python",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(pythonTopic)

	sqlTopic := &models.Topic{
		Name:        "SQL",
		Description: "Learn SQL",
		Category:    "sql",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(sqlTopic)

	// Create questions for each topic
	javaQuestion := &models.Question{
		TopicID:         javaTopic.ID,
		QuestionText:    "What is a Java class?",
		Difficulty:      2,
		Tags:            []string{"oop"},
		TopicQuestionID: "java.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(javaQuestion)

	pythonQuestion := &models.Question{
		TopicID:         pythonTopic.ID,
		QuestionText:    "What is a Python list?",
		Difficulty:      2,
		Tags:            []string{"basics"},
		TopicQuestionID: "python.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(pythonQuestion)

	sqlQuestion := &models.Question{
		TopicID:         sqlTopic.ID,
		QuestionText:    "What is SELECT?",
		Difficulty:      1,
		Tags:            []string{"basics"},
		TopicQuestionID: "sql.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(sqlQuestion)

	// Test 1: Java user can generate Java quiz
	t.Run("JavaUserCanAccessJavaTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       javaUser.ID,
			TopicID:      javaTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 2: Java user can generate SQL quiz
	t.Run("JavaUserCanAccessSQLTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       javaUser.ID,
			TopicID:      sqlTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 3: Java user cannot generate Python quiz
	t.Run("JavaUserCannotAccessPythonTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       javaUser.ID,
			TopicID:      pythonTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 4: Data user can generate Python quiz
	t.Run("DataUserCanAccessPythonTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       dataUser.ID,
			TopicID:      pythonTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 5: Data user can generate SQL quiz
	t.Run("DataUserCanAccessSQLTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       dataUser.ID,
			TopicID:      sqlTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 6: Data user cannot generate Java quiz
	t.Run("DataUserCannotAccessJavaTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       dataUser.ID,
			TopicID:      javaTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusForbidden {
			t.Errorf("Expected status 403 Forbidden, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 7: Admin can generate Java quiz
	t.Run("AdminCanAccessJavaTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       adminUser.ID,
			TopicID:      javaTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 8: Admin can generate Python quiz
	t.Run("AdminCanAccessPythonTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       adminUser.ID,
			TopicID:      pythonTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 9: Admin can generate SQL quiz
	t.Run("AdminCanAccessSQLTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       adminUser.ID,
			TopicID:      sqlTopic.ID,
			NumQuestions: 1,
		}
		reqJSON, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/generate", bytes.NewBuffer(reqJSON))
		w := httptest.NewRecorder()

		handler.GenerateQuizHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})
}
