package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/fsrs"
	"github.com/ZipCodeCore/zip-athena/models"
)

// Test 1: List users as admin
func TestListUsersAsAdmin(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin user
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create regular user
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	req := httptest.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("X-User-ID", "1") // Admin user ID
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.ListUsersHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var users []*models.User
	json.NewDecoder(rec.Body).Decode(&users)

	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}
}

// Test 2: Non-admin cannot access admin endpoints
func TestNonAdminCannotAccessAdminEndpoints(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create regular user
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	req := httptest.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("X-User-ID", "1") // Non-admin user ID
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.ListUsersHandler)(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 Forbidden, got %d", rec.Code)
	}
}

// Test 3: Create quiz
func TestCreateQuiz(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin user
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create topic
	topic := &models.Topic{
		Name:        "Java",
		Description: "Java programming",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	quiz := models.Quiz{
		Name:        "java.1",
		TopicID:     topic.ID,
		Description: "First Java quiz",
	}

	body, _ := json.Marshal(quiz)
	req := httptest.NewRequest("POST", "/admin/quizzes", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.CreateQuizHandler)(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rec.Code)
	}

	var createdQuiz models.Quiz
	json.NewDecoder(rec.Body).Decode(&createdQuiz)

	if createdQuiz.Name != "java.1" {
		t.Errorf("Expected quiz name 'java.1', got %s", createdQuiz.Name)
	}
}

// Test 4: Update quiz
func TestUpdateQuiz(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create quiz
	quiz := &models.Quiz{
		Name:        "java.1",
		TopicID:     1,
		Description: "Original description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Update quiz
	quiz.Description = "Updated description"
	body, _ := json.Marshal(quiz)

	req := httptest.NewRequest("PUT", "/admin/quizzes?id=1", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateQuizHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var updated models.Quiz
	json.NewDecoder(rec.Body).Decode(&updated)

	if updated.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", updated.Description)
	}
}

// Test 5: Delete quiz
func TestDeleteQuiz(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create quiz
	quiz := &models.Quiz{
		Name:        "java.1",
		TopicID:     1,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	req := httptest.NewRequest("DELETE", "/admin/quizzes?id=1", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.DeleteQuizHandler)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}

	// Verify deletion
	_, err := db.GetQuiz(1)
	if err != database.ErrNotFound {
		t.Error("Quiz should have been deleted")
	}
}

// Test 6: List quizzes
func TestListQuizzes(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create quizzes
	for i := 1; i <= 3; i++ {
		quiz := &models.Quiz{
			Name:        "java." + strconv.Itoa(i),
			TopicID:     1,
			Description: "Test quiz",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		db.CreateQuiz(quiz)
	}

	req := httptest.NewRequest("GET", "/admin/quizzes", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.ListQuizzesHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var quizzes []*models.Quiz
	json.NewDecoder(rec.Body).Decode(&quizzes)

	if len(quizzes) != 3 {
		t.Errorf("Expected 3 quizzes, got %d", len(quizzes))
	}
}

// Test 7: Review question with FSRS
func TestReviewQuestionFSRS(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create user
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	// Create question
	question := &models.Question{
		TopicID:         1,
		QuestionText:    "Test question",
		Difficulty:      2,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Review with good rating
	reviewReq := ReviewQuestionRequest{
		UserID:     user.ID,
		QuestionID: question.ID,
		Rating:     fsrs.RatingGood,
	}

	body, _ := json.Marshal(reviewReq)
	req := httptest.NewRequest("POST", "/fsrs/review", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ReviewQuestionHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var resp ReviewQuestionResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.State == nil {
		t.Error("Expected state in response")
	}

	if resp.State.Stability == 0 {
		t.Error("Expected non-zero stability after review")
	}
}

// Test 8: Invalid rating returns error
func TestReviewQuestionInvalidRating(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	reviewReq := ReviewQuestionRequest{
		UserID:     1,
		QuestionID: 1,
		Rating:     99, // Invalid rating
	}

	body, _ := json.Marshal(reviewReq)
	req := httptest.NewRequest("POST", "/fsrs/review", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ReviewQuestionHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rec.Code)
	}
}

// Test 9: Get due questions
func TestGetDueQuestions(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create user
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	// Create questions
	for i := 1; i <= 3; i++ {
		question := &models.Question{
			TopicID:         1,
			QuestionText:    "Question " + strconv.Itoa(i),
			Difficulty:      2,
			TopicQuestionID: "1." + strconv.Itoa(i),
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(question)

		// Create answer for each question
		answer := &models.Answer{
			QuestionID: i,
			AnswerText: "Answer",
			IsCorrect:  true,
		}
		db.CreateAnswer(answer)

		// Create due state
		state := &models.UserQuestionState{
			UserID:     user.ID,
			QuestionID: i,
			Due:        time.Now().Add(-1 * time.Hour), // Due in the past
			Stability:  1.0,
			Difficulty: 5.0,
			State:      1,
		}
		db.CreateUserQuestionState(state)
	}

	req := httptest.NewRequest("GET", "/fsrs/due?user_id=1&limit=3", nil)
	rec := httptest.NewRecorder()

	handler.GetDueQuestionsHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var resp GetDueQuestionsResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp.Count != 3 {
		t.Errorf("Expected 3 due questions, got %d", resp.Count)
	}
}

// Test 10: Adaptive quiz generation prioritizes due questions
func TestAdaptiveQuizGeneration(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create user
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	// Create topic
	topic := &models.Topic{
		Name:        "Java",
		Description: "Java programming",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create questions: some due, some new
	for i := 1; i <= 5; i++ {
		question := &models.Question{
			TopicID:         topic.ID,
			QuestionText:    "Question " + strconv.Itoa(i),
			Difficulty:      2,
			TopicQuestionID: "1." + strconv.Itoa(i),
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(question)

		// Create answer
		answer := &models.Answer{
			QuestionID: i,
			AnswerText: "Answer",
			IsCorrect:  true,
		}
		db.CreateAnswer(answer)

		// Make first 2 questions due
		if i <= 2 {
			state := &models.UserQuestionState{
				UserID:     user.ID,
				QuestionID: i,
				Due:        time.Now().Add(-1 * time.Hour),
				Stability:  1.0,
				Difficulty: 5.0,
				State:      2,
			}
			db.CreateUserQuestionState(state)
		}
	}

	genReq := GenerateQuizRequest{
		UserID:       user.ID,
		TopicID:      topic.ID,
		NumQuestions: 3,
	}

	body, _ := json.Marshal(genReq)
	req := httptest.NewRequest("POST", "/fsrs/quiz/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.GenerateAdaptiveQuizHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var resp GenerateQuizResponse
	json.NewDecoder(rec.Body).Decode(&resp)

	if len(resp.Questions) != 3 {
		t.Errorf("Expected 3 questions, got %d", len(resp.Questions))
	}

	// First 2 questions should be the due ones (IDs 1 and 2)
	if resp.Questions[0].Question.ID > 2 || resp.Questions[1].Question.ID > 2 {
		t.Error("Expected due questions to be prioritized")
	}
}

// Test 11: Update user
func TestUpdateUser(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create user to update
	user := &models.User{
		Username:    "user1",
		Email:       "user1@test.com",
		TotalPoints: 10,
		CreatedAt:   time.Now(),
	}
	db.CreateUser(user)

	// Update user
	user.TotalPoints = 50
	body, _ := json.Marshal(user)

	req := httptest.NewRequest("PUT", "/admin/users?id=2", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateUserHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify update
	updated, _ := db.GetUser(2)
	if updated.TotalPoints != 50 {
		t.Errorf("Expected 50 points, got %d", updated.TotalPoints)
	}
}

// Test 12: Delete user
func TestDeleteUser(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create user to delete
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	req := httptest.NewRequest("DELETE", "/admin/users?id=2", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.DeleteUserHandler)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}

	// Verify deletion
	_, err := db.GetUser(2)
	if err != database.ErrNotFound {
		t.Error("User should have been deleted")
	}
}

// Test 13: List all questions
func TestListAllQuestions(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create questions
	for i := 1; i <= 5; i++ {
		question := &models.Question{
			TopicID:         1,
			QuestionText:    "Question " + strconv.Itoa(i),
			Difficulty:      2,
			Tags:            []string{"test"},
			TopicQuestionID: "1." + strconv.Itoa(i),
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(question)
	}

	req := httptest.NewRequest("GET", "/admin/questions", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.ListQuestionsHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	var questions []*models.Question
	json.NewDecoder(rec.Body).Decode(&questions)

	if len(questions) != 5 {
		t.Errorf("Expected 5 questions, got %d", len(questions))
	}
}

// Test 14: Update question
func TestUpdateQuestion(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create question
	question := &models.Question{
		TopicID:         1,
		QuestionText:    "Original question",
		Difficulty:      2,
		Tags:            []string{"test"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Update question
	question.QuestionText = "Updated question"
	body, _ := json.Marshal(question)

	req := httptest.NewRequest("PUT", "/admin/questions?id=1", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateQuestionHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify update
	updated, _ := db.GetQuestion(1)
	if updated.QuestionText != "Updated question" {
		t.Errorf("Expected updated text, got %s", updated.QuestionText)
	}
}

// Test 15: Delete question
func TestDeleteQuestion(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create question
	question := &models.Question{
		TopicID:         1,
		QuestionText:    "Test question",
		Difficulty:      2,
		Tags:            []string{"test"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	req := httptest.NewRequest("DELETE", "/admin/questions?id=1", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.DeleteQuestionHandler)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}

	// Verify deletion
	_, err := db.GetQuestion(1)
	if err != database.ErrNotFound {
		t.Error("Question should have been deleted")
	}
}

// Test: List all topics via admin endpoint
func TestListAdminTopicsHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin user
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create topics with different categories
	db.CreateTopic(&models.Topic{Name: "Java Fundamentals", Description: "Java basics", Category: "java", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "Python Basics", Description: "Python basics", Category: "python", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "SQL Database", Description: "SQL basics", Category: "sql", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "Inactive Topic", Description: "Should not appear", Category: "java", IsActive: false})

	req := httptest.NewRequest("GET", "/admin/topics", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.ListAdminTopicsHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var topics []*models.Topic
	if err := json.NewDecoder(rec.Body).Decode(&topics); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Admin should see all 4 topics including the inactive one
	if len(topics) != 4 {
		t.Errorf("Expected 4 topics (including inactive), got %d", len(topics))
	}

	// Verify the inactive topic IS returned (admin sees all topics)
	foundInactive := false
	for _, topic := range topics {
		if !topic.IsActive {
			foundInactive = true
		}
	}
	if !foundInactive {
		t.Error("Admin list should include inactive topics")
	}
}

// Test: List admin topics requires admin authentication
func TestListAdminTopicsHandlerRequiresAdmin(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Non-admin user
	user := &models.User{
		Username:  "student",
		Email:     "student@test.com",
		IsAdmin:   false,
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	req := httptest.NewRequest("GET", "/admin/topics", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.ListAdminTopicsHandler)(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("Expected status 403, got %d", rec.Code)
	}
}

// Test 16: Update topic
func TestUpdateTopic(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create topic
	topic := &models.Topic{
		Name:        "Java",
		Description: "Original description",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Update topic
	topic.Description = "Updated description"
	body, _ := json.Marshal(topic)

	req := httptest.NewRequest("PUT", "/admin/topics?id=1", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateTopicHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// Verify update
	updated, _ := db.GetTopic(1)
	if updated.Description != "Updated description" {
		t.Errorf("Expected updated description, got %s", updated.Description)
	}
}

// Test 17: Delete topic
func TestDeleteTopic(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create topic
	topic := &models.Topic{
		Name:        "Java",
		Description: "Test topic",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	req := httptest.NewRequest("DELETE", "/admin/topics?id=1", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.DeleteTopicHandler)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}

	// Verify soft deletion (topic should still exist but be inactive)
	topic, err := db.GetTopic(1)
	if err != nil {
		t.Error("Topic should still exist after soft delete, got error:", err)
	}
	if topic.IsActive {
		t.Error("Topic should be inactive after soft delete")
	}
}

// TestCreateAdminTopicHandler verifies that POST /admin/topics creates a topic and
// that the new topic immediately appears in the admin list.
func TestCreateAdminTopicHandler(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	admin := &models.User{Username: "admin", Email: "admin@test.com", IsAdmin: true, CreatedAt: time.Now()}
	db.CreateUser(admin)

	newTopic := models.Topic{Name: "Shell Usage", Description: "bash and zsh", Category: "all"}
	body, _ := json.Marshal(newTopic)

	req := httptest.NewRequest("POST", "/admin/topics", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.CreateAdminTopicHandler)(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// Topic must appear in the admin list immediately after creation
	topics, err := db.ListTopics()
	if err != nil {
		t.Fatalf("ListTopics failed: %v", err)
	}
	if len(topics) != 1 {
		t.Errorf("Expected 1 topic in list, got %d", len(topics))
	}
	if topics[0].Name != "Shell Usage" {
		t.Errorf("Expected topic name 'Shell Usage', got '%s'", topics[0].Name)
	}
}

// TestUpdateTopicRemainsInAdminList verifies that after updating a topic via
// PUT /admin/topics the topic is still returned by GET /admin/topics.
// This is a regression test for the "topic disappears after editing" bug.
func TestUpdateTopicRemainsInAdminList(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	admin := &models.User{Username: "admin", Email: "admin@test.com", IsAdmin: true, CreatedAt: time.Now()}
	db.CreateUser(admin)

	// Create topic with explicit IsActive=true (mirrors SQLite behaviour)
	topic := &models.Topic{Name: "Java", Description: "desc", Category: "java", IsActive: true, DisplayOrder: 3, CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Send a partial update – name/description/category and is_active explicitly set
	update := map[string]interface{}{
		"name":        "Java Updated",
		"description": "updated desc",
		"category":    "all",
		"is_active":   true,
	}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/admin/topics?id=1", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateTopicHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Topic must still be present in the admin list
	topics, err := db.ListTopics()
	if err != nil {
		t.Fatalf("ListTopics failed: %v", err)
	}
	found := false
	for _, tp := range topics {
		if tp.Name == "Java Updated" {
			found = true
			if !tp.IsActive {
				t.Error("Topic should still be active after update")
			}
			if tp.DisplayOrder != 3 {
				t.Errorf("DisplayOrder should be preserved (3), got %d", tp.DisplayOrder)
			}
			break
		}
	}
	if !found {
		t.Error("Updated topic not found in admin list – topic disappeared after edit")
	}
}

// TestUpdateInstructorTopicRemainsInList verifies that the instructor update
// handler preserves is_active so the topic does not vanish after being edited.
func TestUpdateInstructorTopicRemainsInList(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	admin := &models.User{Username: "admin", Email: "admin@test.com", IsAdmin: true, CreatedAt: time.Now()}
	db.CreateUser(admin)

	topic := &models.Topic{Name: "Python Basics", Description: "intro", Category: "python", IsActive: true, DisplayOrder: 2, CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Partial update – no is_active or display_order in request body
	update := map[string]interface{}{
		"name":        "Python Basics Updated",
		"description": "updated intro",
		"category":    "all",
	}
	body, _ := json.Marshal(update)

	req := httptest.NewRequest("PUT", "/instructor/topics?id=1", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateInstructorTopicHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Topic must still be present in the list (is_active must not have been zeroed)
	topics, err := db.ListTopics()
	if err != nil {
		t.Fatalf("ListTopics failed: %v", err)
	}
	found := false
	for _, tp := range topics {
		if tp.Name == "Python Basics Updated" {
			found = true
			if !tp.IsActive {
				t.Error("Topic should still be active after instructor update")
			}
			break
		}
	}
	if !found {
		t.Error("Updated topic not found in list after instructor update – topic disappeared")
	}
}

// Test 18: Update answer
func TestUpdateAnswer(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create answer
	answer := &models.Answer{
		QuestionID: 1,
		AnswerText: "Original answer",
		IsCorrect:  false,
	}
	db.CreateAnswer(answer)

	// Update answer
	answer.AnswerText = "Updated answer"
	answer.IsCorrect = true
	body, _ := json.Marshal(answer)

	req := httptest.NewRequest("PUT", "/admin/answers?id=1", bytes.NewBuffer(body))
	req.Header.Set("X-User-ID", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateAnswerHandler)(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
}

// Test 19: Delete answer
func TestDeleteAnswer(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin
	admin := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(admin)

	// Create answer
	answer := &models.Answer{
		QuestionID: 1,
		AnswerText: "Test answer",
		IsCorrect:  true,
	}
	db.CreateAnswer(answer)

	req := httptest.NewRequest("DELETE", "/admin/answers?id=1", nil)
	req.Header.Set("X-User-ID", "1")
	rec := httptest.NewRecorder()

	handler.AdminMiddleware(handler.DeleteAnswerHandler)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", rec.Code)
	}
}

// Test 20: Multiple FSRS reviews track progress correctly
func TestMultipleFSRSReviews(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create user
	user := &models.User{
		Username:  "user1",
		Email:     "user1@test.com",
		CreatedAt: time.Now(),
	}
	db.CreateUser(user)

	// Create question
	question := &models.Question{
		TopicID:         1,
		QuestionText:    "Test question",
		Difficulty:      2,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// First review - Good
	reviewReq := ReviewQuestionRequest{
		UserID:     user.ID,
		QuestionID: question.ID,
		Rating:     fsrs.RatingGood,
	}

	body, _ := json.Marshal(reviewReq)
	req := httptest.NewRequest("POST", "/fsrs/review", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.ReviewQuestionHandler(rec, req)

	var resp1 ReviewQuestionResponse
	json.NewDecoder(rec.Body).Decode(&resp1)

	// Second review - Good
	body2, _ := json.Marshal(reviewReq)
	req2 := httptest.NewRequest("POST", "/fsrs/review", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()

	handler.ReviewQuestionHandler(rec2, req2)

	var resp2 ReviewQuestionResponse
	json.NewDecoder(rec2.Body).Decode(&resp2)

	// Reps should increase
	if resp2.State.Reps != 2 {
		t.Errorf("Expected 2 reps, got %d", resp2.State.Reps)
	}

	// State should be tracked correctly
	if resp2.State.UserID != user.ID || resp2.State.QuestionID != question.ID {
		t.Error("State should track correct user and question IDs")
	}

	// Stability should be positive
	if resp2.State.Stability <= 0 {
		t.Errorf("Stability should be positive, got %v", resp2.State.Stability)
	}
}
