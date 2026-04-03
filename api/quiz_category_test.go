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

	javaTracingTopic := &models.Topic{
		Name:        "Java Code Tracing",
		Description: "Trace Java code",
		Category:    "java",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(javaTracingTopic)

	pythonTracingTopic := &models.Topic{
		Name:        "Python Code Tracing",
		Description: "Trace Python code",
		Category:    "python",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(pythonTracingTopic)

	sqlTopic := &models.Topic{
		Name:        "SQL",
		Description: "Learn SQL",
		Category:    "all",
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

	javaTracingQuestion := &models.Question{
		TopicID:         javaTracingTopic.ID,
		QuestionText:    "What is the output of this Java snippet?",
		Difficulty:      2,
		Tags:            []string{"tracing"},
		TopicQuestionID: "java-tracing.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(javaTracingQuestion)

	pythonTracingQuestion := &models.Question{
		TopicID:         pythonTracingTopic.ID,
		QuestionText:    "What is the output of this Python snippet?",
		Difficulty:      2,
		Tags:            []string{"tracing"},
		TopicQuestionID: "python-tracing.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(pythonTracingQuestion)

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

	// Test 3: Java user can generate Java tracing quiz
	t.Run("JavaUserCanAccessJavaTracingTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       javaUser.ID,
			TopicID:      javaTracingTopic.ID,
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

	// Test 4: Java user cannot generate Python quiz
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

	// Test 5: Java user cannot generate Python tracing quiz
	t.Run("JavaUserCannotAccessPythonTracingTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       javaUser.ID,
			TopicID:      pythonTracingTopic.ID,
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

	// Test 6: Data user can generate Python quiz
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

	// Test 7: Data user can generate Python tracing quiz
	t.Run("DataUserCanAccessPythonTracingTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       dataUser.ID,
			TopicID:      pythonTracingTopic.ID,
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

	// Test 8: Data user can generate SQL quiz
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

	// Test 9: Data user cannot generate Java quiz
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

	// Test 10: Data user cannot generate Java tracing quiz
	t.Run("DataUserCannotAccessJavaTracingTopic", func(t *testing.T) {
		reqBody := GenerateQuizRequest{
			UserID:       dataUser.ID,
			TopicID:      javaTracingTopic.ID,
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

	// Test 11: Admin can generate Java quiz
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

	// Test 12: Admin can generate Python quiz
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

	// Test 13: Admin can generate SQL quiz
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
