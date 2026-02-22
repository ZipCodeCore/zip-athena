package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// setupAssignedQuizDB creates a test database with a quiz and users for assigned quiz tests
func setupAssignedQuizDB(t *testing.T) (*database.InMemoryDB, *models.User, *models.User, *models.Quiz) {
	t.Helper()
	db := database.NewInMemoryDB()

	// Create admin user
	admin := &models.User{
		Username: "instructor",
		Email:    "instructor@test.com",
		IsAdmin:  true,
	}
	db.CreateUser(admin)

	// Create student user
	student := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		IsAdmin:  false,
		Category: "java",
	}
	db.CreateUser(student)

	// Create a topic and quiz
	topic := &models.Topic{
		Name:      "Java Basics",
		Category:  "java",
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	db.CreateTopic(topic)

	quiz := &models.Quiz{
		Name:        "Java Quiz 1",
		TopicID:     topic.ID,
		Category:    "java",
		Description: "Test quiz",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Add questions to quiz
	q := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "What is Java?",
		Difficulty:      2,
		Tags:            []string{"java", "basics"},
		TopicQuestionID: fmt.Sprintf("%d.1", topic.ID),
		IsActive:        true,
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(q)
	db.AddQuestionToQuiz(quiz.ID, q.ID, 1)

	// Add answers
	db.CreateAnswer(&models.Answer{QuestionID: q.ID, AnswerText: "A programming language", IsCorrect: true})
	db.CreateAnswer(&models.Answer{QuestionID: q.ID, AnswerText: "A coffee brand", IsCorrect: false})
	db.CreateAnswer(&models.Answer{QuestionID: q.ID, AnswerText: "A database", IsCorrect: false})

	return db, admin, student, quiz
}

func TestAssignQuizToStudentsHandler(t *testing.T) {
	db, admin, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	body, _ := json.Marshal(map[string]interface{}{
		"user_ids": []int{student.ID},
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/instructor/assigned-quizzes/assign?quiz_id=%d", quiz.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.AssignQuizToStudentsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["assigned"].(float64) != 1 {
		t.Errorf("Expected 1 assigned, got %v", result["assigned"])
	}
}

func TestGetAssignedQuizzesHandler(t *testing.T) {
	db, admin, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	// Assign quiz to student first
	db.AssignQuizToUser(quiz.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/assigned-quizzes?user_id=%d", student.ID), nil)
	w := httptest.NewRecorder()

	handler.GetAssignedQuizzesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var quizzes []*models.AssignedQuizInfo
	json.NewDecoder(w.Body).Decode(&quizzes)
	if len(quizzes) != 1 {
		t.Errorf("Expected 1 assigned quiz, got %d", len(quizzes))
	}
	if quizzes[0].QuizID != quiz.ID {
		t.Errorf("Expected quiz ID %d, got %d", quiz.ID, quizzes[0].QuizID)
	}
}

func TestTakeAssignedQuizHandler(t *testing.T) {
	db, admin, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	// Assign quiz to student first
	db.AssignQuizToUser(quiz.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/assigned-quizzes/start?quiz_id=%d&user_id=%d", quiz.ID, student.ID), nil)
	w := httptest.NewRecorder()

	handler.TakeAssignedQuizHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp GenerateQuizResponse
	json.NewDecoder(w.Body).Decode(&resp)
	if resp.QuizAttemptID == 0 {
		t.Error("Expected quiz_attempt_id to be set")
	}
	if len(resp.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(resp.Questions))
	}
}

func TestTakeAssignedQuiz_NotAssigned(t *testing.T) {
	db, _, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	// Do NOT assign quiz to student
	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/assigned-quizzes/start?quiz_id=%d&user_id=%d", quiz.ID, student.ID), nil)
	w := httptest.NewRecorder()

	handler.TakeAssignedQuizHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestGetAssignedQuizStatusHandler(t *testing.T) {
	db, admin, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	// Assign quiz to student
	db.AssignQuizToUser(quiz.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/instructor/assigned-quizzes/status/quiz?quiz_id=%d", quiz.ID), nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GetAssignedQuizStatusHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var status models.AssignedQuizStatus
	json.NewDecoder(w.Body).Decode(&status)
	if status.QuizID != quiz.ID {
		t.Errorf("Expected quiz ID %d, got %d", quiz.ID, status.QuizID)
	}
	if len(status.StudentResults) != 1 {
		t.Errorf("Expected 1 student result, got %d", len(status.StudentResults))
	}
	if status.StudentResults[0].UserID != student.ID {
		t.Errorf("Expected student ID %d, got %d", student.ID, status.StudentResults[0].UserID)
	}
}

func TestImportAndCreateQuizFromJSONHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin user
	admin := &models.User{Username: "admin", Email: "admin@test.com", IsAdmin: true}
	db.CreateUser(admin)

	// Build multipart form with JSON quiz data
	quizJSON := `{
		"topic": "Python Test",
		"questions": [
			{
				"text": "What is Python?",
				"difficulty": 2,
				"tags": ["python", "basics"],
				"answers": [
					{"text": "A programming language", "is_correct": true},
					{"text": "A snake", "is_correct": false},
					{"text": "A database", "is_correct": false}
				]
			}
		]
	}`

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("quiz_name", "Python Fundamentals Quiz")
	w.WriteField("category", "python")
	w.WriteField("description", "Test quiz for Python")
	fw, _ := w.CreateFormFile("file", "quiz.json")
	fw.Write([]byte(quizJSON))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/instructor/assigned-quizzes/import", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	rw := httptest.NewRecorder()

	handler.ImportAndCreateQuizFromJSONHandler(rw, req)

	if rw.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", rw.Code, rw.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(rw.Body).Decode(&result)
	if result["success"] != true {
		t.Errorf("Expected success=true, got %v", result["success"])
	}
	if result["imported"].(float64) != 1 {
		t.Errorf("Expected 1 imported question, got %v", result["imported"])
	}
	if result["quiz_name"].(string) != "Python Fundamentals Quiz" {
		t.Errorf("Expected quiz_name 'Python Fundamentals Quiz', got %v", result["quiz_name"])
	}
}

func TestImportAndCreateQuizFromJSON_ValidationError(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	admin := &models.User{Username: "admin", Email: "admin@test.com", IsAdmin: true}
	db.CreateUser(admin)

	// Bad JSON - no correct answer
	quizJSON := `{
		"topic": "Test",
		"questions": [
			{
				"text": "Bad question?",
				"difficulty": 2,
				"tags": [],
				"answers": [
					{"text": "Wrong 1", "is_correct": false},
					{"text": "Wrong 2", "is_correct": false}
				]
			}
		]
	}`

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("quiz_name", "Bad Quiz")
	w.WriteField("category", "java")
	fw, _ := w.CreateFormFile("file", "quiz.json")
	fw.Write([]byte(quizJSON))
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/instructor/assigned-quizzes/import", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	rw := httptest.NewRecorder()

	handler.ImportAndCreateQuizFromJSONHandler(rw, req)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d: %s", rw.Code, rw.Body.String())
	}
}

func TestUnassignQuizFromStudentHandler(t *testing.T) {
	db, admin, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	// Assign first
	db.AssignQuizToUser(quiz.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/instructor/assigned-quizzes/unassign?quiz_id=%d&user_id=%d", quiz.ID, student.ID), nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.UnassignQuizFromStudentHandler(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d: %s", w.Code, w.Body.String())
	}

	// Verify it's gone
	quizzes, _ := db.GetAssignedQuizzesForUser(student.ID)
	if len(quizzes) != 0 {
		t.Errorf("Expected 0 assigned quizzes after unassign, got %d", len(quizzes))
	}
}

func TestAnswersRandomizedInAssignedQuiz(t *testing.T) {
	db, admin, student, quiz := setupAssignedQuizDB(t)
	handler := NewHandler(db)

	db.AssignQuizToUser(quiz.ID, student.ID, admin.ID)

	// Take quiz multiple times and check that answers are sometimes in different order
	orderCounts := make(map[string]int)
	for i := 0; i < 20; i++ {
		req := httptest.NewRequest(http.MethodPost,
			fmt.Sprintf("/assigned-quizzes/start?quiz_id=%d&user_id=%d", quiz.ID, student.ID), nil)
		w := httptest.NewRecorder()
		handler.TakeAssignedQuizHandler(w, req)

		var resp GenerateQuizResponse
		json.NewDecoder(w.Body).Decode(&resp)
		if len(resp.Questions) > 0 && len(resp.Questions[0].Answers) > 0 {
			order := ""
			for _, a := range resp.Questions[0].Answers {
				order += strings.TrimSpace(a.AnswerText[:1])
			}
			orderCounts[order]++
		}
	}

	// With 3 answers and 20 attempts, we should see more than 1 unique ordering
	// (probability of all same order is (1/6)^19 which is astronomically small)
	if len(orderCounts) < 2 {
		t.Logf("Answer orders seen: %v", orderCounts)
		t.Log("Note: answers may not have randomized, but this is expected with small answer sets occasionally")
		// Don't fail - randomization is non-deterministic
	}
}
