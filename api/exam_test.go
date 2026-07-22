package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// setupExamDB creates a test database with questions for exam tests
func setupExamDB(t *testing.T) (*database.InMemoryDB, *models.User, *models.User) {
	t.Helper()
	db := database.NewInMemoryDB()

	admin := &models.User{Username: "examadmin", Email: "examadmin@test.com", IsAdmin: true}
	db.CreateUser(admin)

	student := &models.User{Username: "examstudent", Email: "examstudent@test.com", IsAdmin: false, Category: "java"}
	db.CreateUser(student)

	// Create a topic and questions
	topic := &models.Topic{Name: "Java Basics", Category: "java", IsActive: true, CreatedAt: time.Now()}
	db.CreateTopic(topic)

	for i := 1; i <= 5; i++ {
		q := &models.Question{
			TopicID:         topic.ID,
			QuestionText:    fmt.Sprintf("Question %d?", i),
			Difficulty:      i%3 + 1, // cycles through 2,3,1
			Tags:            []string{"java"},
			TopicQuestionID: fmt.Sprintf("%d.%d", topic.ID, i),
			IsActive:        true,
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(q)
		db.CreateAnswer(&models.Answer{QuestionID: q.ID, AnswerText: "Correct", IsCorrect: true})
		db.CreateAnswer(&models.Answer{QuestionID: q.ID, AnswerText: "Wrong", IsCorrect: false})
	}

	return db, admin, student
}

func TestGenerateExamHandler(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	body, _ := json.Marshal(map[string]interface{}{
		"name":              "Java Midterm",
		"categories":        []string{"java"},
		"difficulty_levels": []int{1, 2, 3},
		"question_count":    3,
	})

	req := httptest.NewRequest(http.MethodPost, "/instructor/exams/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GenerateExamHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if result["success"] != true {
		t.Errorf("Expected success=true, got %v", result["success"])
	}
	if result["exam_name"].(string) != "Java Midterm" {
		t.Errorf("Expected exam_name 'Java Midterm', got %v", result["exam_name"])
	}
	if result["access_token"] == "" || result["access_token"] == nil {
		t.Error("Expected non-empty access_token")
	}
	if result["exam_url"] == "" || result["exam_url"] == nil {
		t.Error("Expected non-empty exam_url")
	}
}

func TestGenerateExamHandler_MissingName(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	body, _ := json.Marshal(map[string]interface{}{
		"categories":        []string{"java"},
		"difficulty_levels": []int{1},
	})

	req := httptest.NewRequest(http.MethodPost, "/instructor/exams/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GenerateExamHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestListExamsHandler(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	// Create an exam first
	db.CreateExam(&models.Exam{
		Name: "Test Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "testtoken123",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	})

	req := httptest.NewRequest(http.MethodGet, "/instructor/exams", nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.ListExamsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var exams []*models.Exam
	json.NewDecoder(w.Body).Decode(&exams)
	if len(exams) != 1 {
		t.Errorf("Expected 1 exam, got %d", len(exams))
	}
}

func TestAssignExamToStudentsHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "Assigned Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "assigntoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)

	body, _ := json.Marshal(map[string]interface{}{
		"user_ids": []int{student.ID},
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/instructor/exams/assign?exam_id=%d", exam.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.AssignExamToStudentsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["assigned"].(float64) != 1 {
		t.Errorf("Expected assigned=1, got %v", result["assigned"])
	}
}

func TestAssignExamToStudentsHandler_InactiveStudentRejected(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	student.IsActive = false
	if err := db.UpdateUser(student); err != nil {
		t.Fatalf("failed to deactivate student: %v", err)
	}

	exam := &models.Exam{
		Name: "Assigned Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "inactive-token",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)

	body, _ := json.Marshal(map[string]interface{}{
		"user_ids": []int{student.ID},
	})

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/instructor/exams/assign?exam_id=%d", exam.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.AssignExamToStudentsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["assigned"].(float64) != 0 {
		t.Errorf("Expected assigned=0, got %v", result["assigned"])
	}
}

func TestGetAssignedExamsHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "My Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "mytoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/assigned-exams?user_id=%d", student.ID), nil)
	w := httptest.NewRecorder()

	handler.GetAssignedExamsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var exams []*models.ExamInfo
	json.NewDecoder(w.Body).Decode(&exams)
	if len(exams) != 1 {
		t.Errorf("Expected 1 assigned exam, got %d", len(exams))
	}
	if exams[0].ExamID != exam.ID {
		t.Errorf("Expected exam ID %d, got %d", exam.ID, exams[0].ExamID)
	}
}

func TestStartExamHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	// Get IDs of questions in the DB
	questions, _ := db.ListQuestions()
	qIDs := make([]int, 0, 3)
	for i, q := range questions {
		if i >= 3 {
			break
		}
		qIDs = append(qIDs, q.ID)
	}

	exam := &models.Exam{
		Name: "Start Exam", Categories: []string{"java"}, DifficultyLevels: []int{1, 2, 3},
		QuestionIDs: qIDs, QuestionCount: len(qIDs), AccessToken: "starttoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/exam/start?token=%s&user_id=%d", exam.AccessToken, student.ID), nil)
	w := httptest.NewRecorder()

	handler.StartExamHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["attempt_id"] == nil || result["attempt_id"].(float64) == 0 {
		t.Error("Expected attempt_id to be set")
	}
	if result["questions"] == nil {
		t.Error("Expected questions in response")
	}
	questionsResult, ok := result["questions"].([]interface{})
	if !ok || len(questionsResult) == 0 {
		t.Fatalf("Expected at least one question in response, got %v", result["questions"])
	}
	firstQuestion, ok := questionsResult[0].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected first question to be an object, got %T", questionsResult[0])
	}
	if firstQuestion["category"] != "Java Basics" {
		t.Errorf("Expected question category %q, got %v", "Java Basics", firstQuestion["category"])
	}
}

func TestStartExamHandler_NotAssigned(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "Unassigned Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "unassignedtoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	// Do NOT assign to student

	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/exam/start?token=%s&user_id=%d", exam.AccessToken, student.ID), nil)
	w := httptest.NewRecorder()

	handler.StartExamHandler(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", w.Code)
	}
}

func TestStartExamHandler_AllowsAdminPreview(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "Preview Exam", Categories: []string{"java"}, DifficultyLevels: []int{1, 2, 3},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "preview-token",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodPost,
		fmt.Sprintf("/exam/start?token=%s&user_id=%d", exam.AccessToken, admin.ID), nil)
	w := httptest.NewRecorder()

	handler.StartExamHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestSubmitExamHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	questions, _ := db.ListQuestions()
	qIDs := []int{questions[0].ID}

	exam := &models.Exam{
		Name: "Submit Exam", Categories: []string{"java"}, DifficultyLevels: []int{1, 2, 3},
		QuestionIDs: qIDs, QuestionCount: 1, AccessToken: "submittoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	// Create an attempt
	attempt := &models.ExamAttempt{
		ExamID: exam.ID, UserID: student.ID, StartedAt: time.Now(),
		TotalQuestions: 1, Answers: map[string]int{}, QuestionResults: map[string]bool{},
	}
	db.CreateExamAttempt(attempt)

	// Get correct answer for question
	answers, _ := db.GetAnswersByQuestion(questions[0].ID)
	var correctAnswerID int
	for _, a := range answers {
		if a.IsCorrect {
			correctAnswerID = a.ID
			break
		}
	}

	body, _ := json.Marshal(map[string]interface{}{
		"attempt_id": attempt.ID,
		"answers":    map[string]int{fmt.Sprintf("%d", questions[0].ID): correctAnswerID},
	})

	req := httptest.NewRequest(http.MethodPost, "/exam/submit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SubmitExamHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["score"].(float64) != 100.0 {
		t.Errorf("Expected score 100, got %v", result["score"])
	}
}

func TestGetExamStatusHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "Status Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "statustoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodGet,
		fmt.Sprintf("/instructor/exams/status?exam_id=%d", exam.ID), nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GetExamStatusHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var status models.ExamStatus
	json.NewDecoder(w.Body).Decode(&status)
	if status.ExamID != exam.ID {
		t.Errorf("Expected exam ID %d, got %d", exam.ID, status.ExamID)
	}
	if len(status.StudentResults) != 1 {
		t.Errorf("Expected 1 student result, got %d", len(status.StudentResults))
	}
}

func TestGetAvailableExamCategoriesHandler(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/instructor/exams/categories", nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GetAvailableExamCategoriesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var categories []string
	json.NewDecoder(w.Body).Decode(&categories)
	found := false
	for _, c := range categories {
		if c == "java" {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected 'java' in categories, got %v", categories)
	}
}

func TestGetAvailableExamTopicsHandler(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	req := httptest.NewRequest(http.MethodGet, "/instructor/exams/topics", nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GetAvailableExamTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var topics []map[string]interface{}
	json.NewDecoder(w.Body).Decode(&topics)
	if len(topics) == 0 {
		t.Fatal("Expected at least one topic")
	}
	found := false
	for _, topic := range topics {
		if topic["name"] == "Java Basics" {
			found = true
			if topic["category"] != "java" {
				t.Errorf("Expected category 'java', got %v", topic["category"])
			}
		}
	}
	if !found {
		t.Errorf("Expected 'Java Basics' in topics, got %v", topics)
	}
}

func TestGenerateExamHandler_WithTopicIDs(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	// Get the topic ID that was created in setupExamDB
	topics, err := db.ListAllTopics()
	if err != nil || len(topics) == 0 {
		t.Fatal("No topics in DB")
	}
	topicID := topics[0].ID

	body, _ := json.Marshal(map[string]interface{}{
		"name":              "Topic-based Exam",
		"topic_ids":         []int{topicID},
		"difficulty_levels": []int{1, 2, 3},
		"question_count":    3,
	})

	req := httptest.NewRequest(http.MethodPost, "/instructor/exams/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.GenerateExamHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)
	if result["success"] != true {
		t.Errorf("Expected success=true, got %v", result["success"])
	}
}

func TestDeleteExamHandler(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "To Delete", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "deletetoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)

	req := httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/instructor/exams?exam_id=%d", exam.ID), nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.DeleteExamHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify deleted
	if _, err := db.GetExam(exam.ID); err == nil {
		t.Error("Expected exam to be deleted")
	}
}

func TestDeleteExamHandler_NotFound(t *testing.T) {
	db, admin, _ := setupExamDB(t)
	handler := NewHandler(db)

	req := httptest.NewRequest(http.MethodDelete, "/instructor/exams?exam_id=99999", nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.DeleteExamHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestUnassignExamFromUserHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "Unassign Exam", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "unassigntoken2",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	req := httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/instructor/exams/unassign?exam_id=%d&user_id=%d", exam.ID, student.ID), nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.UnassignExamFromUserHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	// Verify unassigned
	exams, _ := db.GetAssignedExamsForUser(student.ID)
	for _, e := range exams {
		if e.ExamID == exam.ID {
			t.Error("Expected exam to be unassigned from student")
		}
	}
}

func TestUnassignExamFromUserHandler_NotFound(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	exam := &models.Exam{
		Name: "Not Assigned", Categories: []string{"java"}, DifficultyLevels: []int{1},
		QuestionIDs: []int{1}, QuestionCount: 1, AccessToken: "notassignedtoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	// Do NOT assign

	req := httptest.NewRequest(http.MethodDelete,
		fmt.Sprintf("/instructor/exams/unassign?exam_id=%d&user_id=%d", exam.ID, student.ID), nil)
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", admin.ID))
	w := httptest.NewRecorder()

	handler.UnassignExamFromUserHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestGetExamSubmissionHandler(t *testing.T) {
	db, admin, student := setupExamDB(t)
	handler := NewHandler(db)

	questions, _ := db.ListQuestions()
	qIDs := []int{questions[0].ID}

	exam := &models.Exam{
		Name: "Submission Exam", Categories: []string{"java"}, DifficultyLevels: []int{1, 2, 3},
		QuestionIDs: qIDs, QuestionCount: 1, AccessToken: "submissiontoken",
		CreatedBy: admin.ID, CreatedAt: time.Now(),
	}
	db.CreateExam(exam)
	db.AssignExamToUser(exam.ID, student.ID, admin.ID)

	// Create and submit an attempt
	attempt := &models.ExamAttempt{
		ExamID: exam.ID, UserID: student.ID, StartedAt: time.Now(),
		TotalQuestions: 1, Answers: map[string]int{}, QuestionResults: map[string]bool{},
	}
	db.CreateExamAttempt(attempt)

	answers, _ := db.GetAnswersByQuestion(questions[0].ID)
	var correctAnswerID int
	for _, a := range answers {
		if a.IsCorrect {
			correctAnswerID = a.ID
			break
		}
	}
	db.SubmitExamAttempt(attempt.ID, map[string]int{fmt.Sprintf("%d", questions[0].ID): correctAnswerID})

	t.Run("student can view own submission", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/exam/submission?attempt_id=%d&user_id=%d", attempt.ID, student.ID), nil)
		w := httptest.NewRecorder()
		handler.GetExamSubmissionHandler(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
		var result map[string]interface{}
		json.NewDecoder(w.Body).Decode(&result)
		if result["exam_name"] != "Submission Exam" {
			t.Errorf("Expected exam_name 'Submission Exam', got %v", result["exam_name"])
		}
		qs, ok := result["questions"].([]interface{})
		if !ok || len(qs) == 0 {
			t.Fatal("Expected questions in response")
		}
	})

	t.Run("admin can view any submission", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/exam/submission?attempt_id=%d&user_id=%d", attempt.ID, admin.ID), nil)
		w := httptest.NewRecorder()
		handler.GetExamSubmissionHandler(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	t.Run("other student cannot view", func(t *testing.T) {
		other := &models.User{Username: "other", Email: "other@test.com", IsAdmin: false}
		db.CreateUser(other)
		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/exam/submission?attempt_id=%d&user_id=%d", attempt.ID, other.ID), nil)
		w := httptest.NewRecorder()
		handler.GetExamSubmissionHandler(w, req)
		if w.Code != http.StatusForbidden {
			t.Errorf("Expected 403, got %d", w.Code)
		}
	})

	t.Run("attempt_id includes attempt_id in status", func(t *testing.T) {
		status, err := db.GetExamStatus(exam.ID)
		if err != nil {
			t.Fatal(err)
		}
		if len(status.StudentResults) == 0 {
			t.Fatal("Expected at least one student result")
		}
		if status.StudentResults[0].AttemptID == nil {
			t.Error("Expected AttemptID to be set in ExamStudentResult after submission")
		}
	})
}
