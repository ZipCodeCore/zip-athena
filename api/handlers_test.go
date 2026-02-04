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

	// Create a test user
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		IsAdmin:  true, // Admin can see all topics
		Category: "java",
	}
	db.CreateUser(user)

	// Create some topics with correct naming
	db.CreateTopic(&models.Topic{Name: "Java Fundamentals", Description: "Java basics", Category: "java"})
	db.CreateTopic(&models.Topic{Name: "Python Basics", Description: "Python basics", Category: "python"})

	req := httptest.NewRequest(http.MethodGet, "/topics?user_id=1", nil)
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

// Test category-based topic filtering
func TestListTopicsHandlerCategoryFiltering(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create test users with different categories
	javaUser := &models.User{
		Username: "javastudent",
		Email:    "java@example.com",
		Password: "password",
		IsAdmin:  false,
		Category: "java",
	}
	db.CreateUser(javaUser)

	dataUser := &models.User{
		Username: "datastudent",
		Email:    "data@example.com",
		Password: "password",
		IsAdmin:  false,
		Category: "data",
	}
	db.CreateUser(dataUser)

	adminUser := &models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "password",
		IsAdmin:  true,
		Category: "java",
	}
	db.CreateUser(adminUser)

	// Create topics with different categories
	db.CreateTopic(&models.Topic{Name: "Java Fundamentals", Description: "Java basics", Category: "java"})
	db.CreateTopic(&models.Topic{Name: "Python Basics", Description: "Python basics", Category: "python"})
	db.CreateTopic(&models.Topic{Name: "SQL Database", Description: "SQL basics", Category: "sql"})

	// Test 1: Java student should only see Java and SQL topics
	req := httptest.NewRequest(http.MethodGet, "/topics?user_id=1", nil)
	w := httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for Java student, got %d", http.StatusOK, w.Code)
	}

	var javaTopics []*models.Topic
	json.NewDecoder(w.Body).Decode(&javaTopics)

	if len(javaTopics) != 2 {
		t.Errorf("Java student should see 2 topics (Java, SQL), got %d", len(javaTopics))
	}

	// Verify no Python topics are returned
	for _, topic := range javaTopics {
		if topic.Category == "python" {
			t.Errorf("Java student should not see Python topics, but got: %s", topic.Name)
		}
	}

	// Test 2: Data student should only see Python and SQL topics
	req = httptest.NewRequest(http.MethodGet, "/topics?user_id=2", nil)
	w = httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for Data student, got %d", http.StatusOK, w.Code)
	}

	var dataTopics []*models.Topic
	json.NewDecoder(w.Body).Decode(&dataTopics)

	if len(dataTopics) != 2 {
		t.Errorf("Data student should see 2 topics (Python, SQL), got %d", len(dataTopics))
	}

	// Verify no Java topics are returned
	for _, topic := range dataTopics {
		if topic.Category == "java" {
			t.Errorf("Data student should not see Java topics, but got: %s", topic.Name)
		}
	}

	// Test 3: Admin should see all topics
	req = httptest.NewRequest(http.MethodGet, "/topics?user_id=3", nil)
	w = httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for Admin, got %d", http.StatusOK, w.Code)
	}

	var adminTopics []*models.Topic
	json.NewDecoder(w.Body).Decode(&adminTopics)

	if len(adminTopics) != 3 {
		t.Errorf("Admin should see all 3 topics, got %d", len(adminTopics))
	}

	// Test 4: Missing user_id should return error
	req = httptest.NewRequest(http.MethodGet, "/topics", nil)
	w = httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing user_id, got %d", http.StatusBadRequest, w.Code)
	}

	// Test 5: Invalid user_id should return error
	req = httptest.NewRequest(http.MethodGet, "/topics?user_id=invalid", nil)
	w = httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid user_id, got %d", http.StatusBadRequest, w.Code)
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

// Test 21: GetProgressHandler should count questions from all quizzes including in-progress ones
func TestGetProgressHandlerIncludesInProgressQuizzes(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "progressuser",
		Email:    "progress@example.com",
	}
	if err := db.CreateUser(user); err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create a completed quiz attempt
	completedTime := time.Now()
	completedAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		StartTime:       time.Now().Add(-1 * time.Hour),
		CompletedTime:   &completedTime,
		TotalQuestions:  5,
		CorrectAnswers:  3,
		QuestionIDs:     []string{"1.1", "1.2", "1.3", "1.4", "1.5"},
		QuestionResults: map[string]bool{"1.1": true, "1.2": true, "1.3": true, "1.4": false, "1.5": false},
	}
	if err := db.CreateQuizAttempt(completedAttempt); err != nil {
		t.Fatalf("Failed to create completed quiz attempt: %v", err)
	}

	// Create an in-progress quiz attempt (no CompletedTime)
	inProgressAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		StartTime:       time.Now(),
		CompletedTime:   nil, // This is in progress
		TotalQuestions:  4,
		CorrectAnswers:  2,
		QuestionIDs:     []string{"2.1", "2.2", "2.3", "2.4"},
		QuestionResults: map[string]bool{"2.1": true, "2.2": true},
	}
	if err := db.CreateQuizAttempt(inProgressAttempt); err != nil {
		t.Fatalf("Failed to create in-progress quiz attempt: %v", err)
	}

	// Get progress
	req := httptest.NewRequest(http.MethodGet, "/progress?user_id="+strconv.Itoa(user.ID), nil)
	w := httptest.NewRecorder()

	handler.GetProgressHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var progress map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&progress); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify total_questions includes both completed and in-progress quizzes
	// Expected: 5 (from completed quiz) + 2 (answered in in-progress quiz) = 7
	totalQuestions := int(progress["total_questions"].(float64))
	expectedTotalQuestions := 7 // 5 from completed + 2 answered in progress
	if totalQuestions != expectedTotalQuestions {
		t.Errorf("Expected total_questions to be %d (from both completed and in-progress quizzes), got %d", expectedTotalQuestions, totalQuestions)
	}

	// Verify total_correct includes both completed and in-progress quizzes
	// Expected: 3 (completed) + 2 (in-progress) = 5
	totalCorrect := int(progress["total_correct"].(float64))
	expectedTotalCorrect := 5
	if totalCorrect != expectedTotalCorrect {
		t.Errorf("Expected total_correct to be %d (from both completed and in-progress quizzes), got %d", expectedTotalCorrect, totalCorrect)
	}

	// Verify total_quizzes includes both completed and in-progress
	totalQuizzes := int(progress["total_quizzes"].(float64))
	if totalQuizzes != 2 {
		t.Errorf("Expected total_quizzes to be 2, got %d", totalQuizzes)
	}
}

// TestCreateUserWithPasswordHashing tests that passwords are hashed when creating users via admin
func TestCreateUserWithPasswordHashing(t *testing.T) {
	// Create a temporary SQLite database for this test
	tmpDir := t.TempDir()
	dbPath := tmpDir + "/test_user_creation.db"

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create SQLite database: %v", err)
	}
	defer db.Close()

	handler := NewHandler(db)

	// Create a user with a plain text password (like admin panel does)
	testPassword := "testpass123"
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: testPassword,
		IsAdmin:  false,
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateUserHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}

	var created models.User
	json.NewDecoder(w.Body).Decode(&created)

	// Verify the returned user doesn't have password
	if created.Password != "" {
		t.Error("Expected password to be empty in response")
	}

	// Now try to login with the created user (correct password)
	loginReq := LoginRequest{
		Username: "testuser",
		Password: testPassword,
	}

	loginBody, _ := json.Marshal(loginReq)
	loginHTTPReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginBody))
	loginW := httptest.NewRecorder()

	handler.LoginHandler(loginW, loginHTTPReq)

	if loginW.Code != http.StatusOK {
		t.Fatalf("Expected login to succeed with status %d, got %d. Body: %s", http.StatusOK, loginW.Code, loginW.Body.String())
	}

	var authResp AuthResponse
	json.NewDecoder(loginW.Body).Decode(&authResp)

	if authResp.User.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", authResp.User.Username)
	}

	if authResp.Message != "Login successful" {
		t.Errorf("Expected message 'Login successful', got '%s'", authResp.Message)
	}

	// Test that login fails with incorrect password
	wrongLoginReq := LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	wrongLoginBody, _ := json.Marshal(wrongLoginReq)
	wrongLoginHTTPReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(wrongLoginBody))
	wrongLoginW := httptest.NewRecorder()

	handler.LoginHandler(wrongLoginW, wrongLoginHTTPReq)

	if wrongLoginW.Code != http.StatusUnauthorized {
		t.Errorf("Expected login to fail with status %d for wrong password, got %d", http.StatusUnauthorized, wrongLoginW.Code)
	}
}

// TestUpdateUserPasswordHashing tests that passwords are hashed when updating users via admin
func TestUpdateUserPasswordHashing(t *testing.T) {
// Create a temporary SQLite database for this test
tmpDir := t.TempDir()
dbPath := tmpDir + "/test_user_update.db"

db, err := database.NewSQLiteDB(dbPath)
if err != nil {
t.Fatalf("Failed to create SQLite database: %v", err)
}
defer db.Close()

handler := NewHandler(db)

// Create a user first
initialPassword := "initial123"
user := models.User{
Username: "testuser",
Email:    "test@example.com",
Password: initialPassword,
IsAdmin:  false,
}

createBody, _ := json.Marshal(user)
createReq := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(createBody))
createW := httptest.NewRecorder()

handler.CreateUserHandler(createW, createReq)

if createW.Code != http.StatusCreated {
t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, createW.Code, createW.Body.String())
}

var created models.User
json.NewDecoder(createW.Body).Decode(&created)
userID := created.ID

// Verify initial login works
loginReq := LoginRequest{
Username: "testuser",
Password: initialPassword,
}
loginBody, _ := json.Marshal(loginReq)
loginHTTPReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginBody))
loginW := httptest.NewRecorder()

handler.LoginHandler(loginW, loginHTTPReq)

if loginW.Code != http.StatusOK {
t.Fatalf("Expected initial login to succeed with status %d, got %d", http.StatusOK, loginW.Code)
}

// Now update the user's password via admin endpoint
newPassword := "newpass456"
updateUser := models.User{
Username:    "testuser",
Email:       "test@example.com",
Password:    newPassword,
IsAdmin:     false,
TotalPoints: 0,
}

updateBody, _ := json.Marshal(updateUser)
updateReq := httptest.NewRequest(http.MethodPut, "/admin/users?id="+strconv.Itoa(userID), bytes.NewBuffer(updateBody))
updateReq.Header.Set("X-User-ID", "1") // Simulate admin user
updateW := httptest.NewRecorder()

handler.UpdateUserHandler(updateW, updateReq)

if updateW.Code != http.StatusOK {
t.Fatalf("Expected update status %d, got %d. Body: %s", http.StatusOK, updateW.Code, updateW.Body.String())
}

// Verify old password no longer works
oldLoginReq := LoginRequest{
Username: "testuser",
Password: initialPassword,
}
oldLoginBody, _ := json.Marshal(oldLoginReq)
oldLoginHTTPReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(oldLoginBody))
oldLoginW := httptest.NewRecorder()

handler.LoginHandler(oldLoginW, oldLoginHTTPReq)

if oldLoginW.Code != http.StatusUnauthorized {
t.Errorf("Expected old password to fail with status %d, got %d", http.StatusUnauthorized, oldLoginW.Code)
}

// Verify new password works
newLoginReq := LoginRequest{
Username: "testuser",
Password: newPassword,
}
newLoginBody, _ := json.Marshal(newLoginReq)
newLoginHTTPReq := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(newLoginBody))
newLoginW := httptest.NewRecorder()

handler.LoginHandler(newLoginW, newLoginHTTPReq)

if newLoginW.Code != http.StatusOK {
t.Fatalf("Expected login with new password to succeed with status %d, got %d. Body: %s", http.StatusOK, newLoginW.Code, newLoginW.Body.String())
}

var authResp AuthResponse
json.NewDecoder(newLoginW.Body).Decode(&authResp)

if authResp.User.Username != "testuser" {
t.Errorf("Expected username 'testuser', got '%s'", authResp.User.Username)
}

if authResp.Message != "Login successful" {
t.Errorf("Expected message 'Login successful', got '%s'", authResp.Message)
}
}
