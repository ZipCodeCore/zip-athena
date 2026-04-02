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
	db.CreateTopic(&models.Topic{Name: "Java Fundamentals", Description: "Java basics", Category: "java", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "Python Basics", Description: "Python basics", Category: "python", IsActive: true})

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
	db.CreateTopic(&models.Topic{Name: "Java Fundamentals", Description: "Java basics", Category: "java", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "Java Code Tracing", Description: "Trace Java code", Category: "java-tracing", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "Python Basics", Description: "Python basics", Category: "python", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "Python Code Tracing", Description: "Trace Python code", Category: "python-tracing", IsActive: true})
	db.CreateTopic(&models.Topic{Name: "SQL Database", Description: "SQL basics", Category: "sql", IsActive: true})

	// Test 1: Java student should only see Java and SQL topics
	req := httptest.NewRequest(http.MethodGet, "/topics?user_id=1", nil)
	w := httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d for Java student, got %d", http.StatusOK, w.Code)
	}

	var javaTopics []*models.Topic
	json.NewDecoder(w.Body).Decode(&javaTopics)

	if len(javaTopics) != 3 {
		t.Errorf("Java student should see 3 topics (Java, Java tracing, SQL), got %d", len(javaTopics))
	}

	// Verify no Python topics are returned
	for _, topic := range javaTopics {
		if topic.Category == "python" || topic.Category == "python-tracing" {
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

	if len(dataTopics) != 3 {
		t.Errorf("Data student should see 3 topics (Python, Python tracing, SQL), got %d", len(dataTopics))
	}

	// Verify no Java topics are returned
	for _, topic := range dataTopics {
		if topic.Category == "java" || topic.Category == "java-tracing" {
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

	if len(adminTopics) != 5 {
		t.Errorf("Admin should see all 5 topics, got %d", len(adminTopics))
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

	// Test 6: Non-existent user_id should return 404
	req = httptest.NewRequest(http.MethodGet, "/topics?user_id=999", nil)
	w = httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent user_id, got %d", http.StatusNotFound, w.Code)
	}
}

// Test: Inactive topics should not be returned
func TestListTopicsHandlerInactiveFiltering(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a test user
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
		IsAdmin:  true,
		Category: "java",
	}
	db.CreateUser(user)

	// Create active and inactive topics
	activeTopic := &models.Topic{
		Name:        "Active Java Topic",
		Description: "Active topic",
		Category:    "java",
		IsActive:    true,
	}
	db.CreateTopic(activeTopic)

	inactiveTopic := &models.Topic{
		Name:        "Inactive Python Topic",
		Description: "Inactive topic",
		Category:    "python",
		IsActive:    false,
	}
	db.CreateTopic(inactiveTopic)

	// List topics - should only return active ones
	req := httptest.NewRequest(http.MethodGet, "/topics?user_id=1", nil)
	w := httptest.NewRecorder()
	handler.ListTopicsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var topics []*models.Topic
	json.NewDecoder(w.Body).Decode(&topics)

	// Should only return 1 active topic
	if len(topics) != 1 {
		t.Errorf("Expected 1 active topic, got %d", len(topics))
	}

	// Verify it's the active one
	if len(topics) > 0 && !topics[0].IsActive {
		t.Errorf("Expected topic to be active")
	}

	// Verify inactive topic is not in the list
	for _, topic := range topics {
		if topic.Name == "Inactive Python Topic" {
			t.Errorf("Inactive topic should not be returned")
		}
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

// Test: GetProgressHandler should only return stats for the specific user, not aggregate
func TestGetProgressHandlerUserSpecificStats(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create two users
	user1 := &models.User{
		Username: "user1",
		Email:    "user1@example.com",
		Category: "data",
	}
	if err := db.CreateUser(user1); err != nil {
		t.Fatalf("Failed to create user1: %v", err)
	}

	user2 := &models.User{
		Username: "user2",
		Email:    "user2@example.com",
		Category: "java",
	}
	if err := db.CreateUser(user2); err != nil {
		t.Fatalf("Failed to create user2: %v", err)
	}

	// Create quiz attempts for user1
	completedTime := time.Now()
	user1Attempt := &models.QuizAttempt{
		UserID:          user1.ID,
		StartTime:       time.Now().Add(-1 * time.Hour),
		CompletedTime:   &completedTime,
		TotalQuestions:  10,
		CorrectAnswers:  7,
		QuestionIDs:     []string{"1.1", "1.2", "1.3", "1.4", "1.5", "1.6", "1.7", "1.8", "1.9", "1.10"},
		QuestionResults: make(map[string]bool),
	}
	if err := db.CreateQuizAttempt(user1Attempt); err != nil {
		t.Fatalf("Failed to create user1 quiz attempt: %v", err)
	}

	// Create many quiz attempts for user2 (to make sure they don't leak into user1's stats)
	for i := 0; i < 5; i++ {
		completedTime2 := time.Now()
		user2Attempt := &models.QuizAttempt{
			UserID:          user2.ID,
			StartTime:       time.Now().Add(-1 * time.Hour),
			CompletedTime:   &completedTime2,
			TotalQuestions:  20,
			CorrectAnswers:  15,
			QuestionIDs:     []string{},
			QuestionResults: make(map[string]bool),
		}
		if err := db.CreateQuizAttempt(user2Attempt); err != nil {
			t.Fatalf("Failed to create user2 quiz attempt: %v", err)
		}
	}

	// Get progress for user1
	req := httptest.NewRequest(http.MethodGet, "/progress?user_id="+strconv.Itoa(user1.ID), nil)
	w := httptest.NewRecorder()
	handler.GetProgressHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var progress map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&progress); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Verify user1's stats - should NOT include user2's data
	totalQuizzes := int(progress["total_quizzes"].(float64))
	if totalQuizzes != 1 {
		t.Errorf("Expected user1 to have 1 quiz, got %d (should not include user2's 5 quizzes)", totalQuizzes)
	}

	totalQuestions := int(progress["total_questions"].(float64))
	if totalQuestions != 10 {
		t.Errorf("Expected user1 to have answered 10 questions, got %d (should not include user2's 100 questions)", totalQuestions)
	}

	totalCorrect := int(progress["total_correct"].(float64))
	if totalCorrect != 7 {
		t.Errorf("Expected user1 to have 7 correct answers, got %d (should not include user2's 75 correct)", totalCorrect)
	}

	// Verify quiz_attempts array only contains user1's attempts
	quizAttempts := progress["quiz_attempts"].([]interface{})
	if len(quizAttempts) != 1 {
		t.Errorf("Expected user1 to have 1 quiz attempt in list, got %d", len(quizAttempts))
	}
}

// TestGetProgressHandlerAdminViewsAllAttempts tests that admin users can see all quiz attempts
func TestGetProgressHandlerAdminViewsAllAttempts(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create an admin user
	adminUser := &models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Category: "admin",
		IsAdmin:  true,
	}
	if err := db.CreateUser(adminUser); err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	// Create a regular user
	regularUser := &models.User{
		Username: "student",
		Email:    "student@example.com",
		Category: "java",
		IsAdmin:  false,
	}
	if err := db.CreateUser(regularUser); err != nil {
		t.Fatalf("Failed to create regular user: %v", err)
	}

	// Create quiz attempts for admin
	completedTime := time.Now()
	adminAttempt := &models.QuizAttempt{
		UserID:          adminUser.ID,
		StartTime:       time.Now().Add(-2 * time.Hour),
		CompletedTime:   &completedTime,
		TotalQuestions:  5,
		CorrectAnswers:  4,
		QuestionIDs:     []string{"1.1", "1.2", "1.3", "1.4", "1.5"},
		QuestionResults: make(map[string]bool),
	}
	if err := db.CreateQuizAttempt(adminAttempt); err != nil {
		t.Fatalf("Failed to create admin quiz attempt: %v", err)
	}

	// Create quiz attempts for regular user
	completedTime2 := time.Now()
	regularAttempt := &models.QuizAttempt{
		UserID:          regularUser.ID,
		StartTime:       time.Now().Add(-1 * time.Hour),
		CompletedTime:   &completedTime2,
		TotalQuestions:  10,
		CorrectAnswers:  7,
		QuestionIDs:     []string{"2.1", "2.2", "2.3", "2.4", "2.5", "2.6", "2.7", "2.8", "2.9", "2.10"},
		QuestionResults: make(map[string]bool),
	}
	if err := db.CreateQuizAttempt(regularAttempt); err != nil {
		t.Fatalf("Failed to create regular user quiz attempt: %v", err)
	}

	// Get progress for admin user - should see ALL attempts
	req := httptest.NewRequest(http.MethodGet, "/progress?user_id="+strconv.Itoa(adminUser.ID), nil)
	w := httptest.NewRecorder()
	handler.GetProgressHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var adminProgress map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&adminProgress); err != nil {
		t.Fatalf("Failed to decode admin response: %v", err)
	}

	// Admin should see all quiz attempts (both their own and regular user's)
	adminQuizAttempts := adminProgress["quiz_attempts"].([]interface{})
	if len(adminQuizAttempts) != 2 {
		t.Errorf("Expected admin to see 2 quiz attempts (all attempts), got %d", len(adminQuizAttempts))
	}

	// Get progress for regular user - should see only their own attempts
	req2 := httptest.NewRequest(http.MethodGet, "/progress?user_id="+strconv.Itoa(regularUser.ID), nil)
	w2 := httptest.NewRecorder()
	handler.GetProgressHandler(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w2.Code)
	}

	var regularProgress map[string]interface{}
	if err := json.NewDecoder(w2.Body).Decode(&regularProgress); err != nil {
		t.Fatalf("Failed to decode regular user response: %v", err)
	}

	// Regular user should see only their own attempts
	regularQuizAttempts := regularProgress["quiz_attempts"].([]interface{})
	if len(regularQuizAttempts) != 1 {
		t.Errorf("Expected regular user to see 1 quiz attempt (only their own), got %d", len(regularQuizAttempts))
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

// Test: Verify dashboard stats are per-user, not aggregate (GitHub Issue #74)
// This test explicitly validates that each user sees only their own stats
func TestDashboardStatsArePerUser(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Define test data structure
	type userData struct {
		user             *models.User
		quizCount        int
		questionsPerQuiz int
		correctPerQuiz   int
	}

	testUsers := []userData{
		{
			user:             &models.User{Username: "alice", Email: "alice@example.com", Category: "java"},
			quizCount:        5,
			questionsPerQuiz: 10,
			correctPerQuiz:   8,
		},
		{
			user:             &models.User{Username: "bob", Email: "bob@example.com", Category: "python"},
			quizCount:        3,
			questionsPerQuiz: 10,
			correctPerQuiz:   9,
		},
		{
			user:             &models.User{Username: "charlie", Email: "charlie@example.com", Category: "data"},
			quizCount:        7,
			questionsPerQuiz: 10,
			correctPerQuiz:   5,
		},
	}

	// Create users and validate IDs are assigned
	for i, userData := range testUsers {
		if err := db.CreateUser(userData.user); err != nil {
			t.Fatalf("Failed to create user %s: %v", userData.user.Username, err)
		}
		if userData.user.ID == 0 {
			t.Fatalf("User %s was not assigned a valid ID", userData.user.Username)
		}
		// Ensure each user has a unique ID
		for j := 0; j < i; j++ {
			if testUsers[j].user.ID == userData.user.ID {
				t.Fatalf("Users %s and %s have duplicate IDs", testUsers[j].user.Username, userData.user.Username)
			}
		}
	}

	// Create quiz attempts for each user based on test data
	for _, userData := range testUsers {
		for i := 0; i < userData.quizCount; i++ {
			completedTime := time.Now()
			attempt := &models.QuizAttempt{
				UserID:          userData.user.ID,
				StartTime:       time.Now().Add(-1 * time.Hour),
				CompletedTime:   &completedTime,
				TotalQuestions:  userData.questionsPerQuiz,
				CorrectAnswers:  userData.correctPerQuiz,
				QuestionIDs:     []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"},
				QuestionResults: make(map[string]bool),
			}
			if err := db.CreateQuizAttempt(attempt); err != nil {
				t.Fatalf("Failed to create quiz attempt for %s: %v", userData.user.Username, err)
			}
		}
	}

	// Verify each user gets only their own stats, not aggregate
	for _, userData := range testUsers {
		req := httptest.NewRequest(http.MethodGet, "/progress?user_id="+strconv.Itoa(userData.user.ID), nil)
		w := httptest.NewRecorder()
		handler.GetProgressHandler(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("[%s] Expected status %d, got %d", userData.user.Username, http.StatusOK, w.Code)
			continue
		}

		var progress map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&progress); err != nil {
			t.Fatalf("[%s] Failed to decode response: %v", userData.user.Username, err)
		}

		// Calculate expected values
		expectedQuizzes := userData.quizCount
		expectedQuestions := userData.quizCount * userData.questionsPerQuiz
		expectedCorrect := userData.quizCount * userData.correctPerQuiz

		// Verify stats are per-user, not aggregate
		totalQuizzes := int(progress["total_quizzes"].(float64))
		if totalQuizzes != expectedQuizzes {
			t.Errorf("[%s] Expected %d quizzes, got %d (must not show aggregate of all users)",
				userData.user.Username, expectedQuizzes, totalQuizzes)
		}

		totalQuestions := int(progress["total_questions"].(float64))
		if totalQuestions != expectedQuestions {
			t.Errorf("[%s] Expected %d questions, got %d (must not show aggregate of all users)",
				userData.user.Username, expectedQuestions, totalQuestions)
		}

		totalCorrect := int(progress["total_correct"].(float64))
		if totalCorrect != expectedCorrect {
			t.Errorf("[%s] Expected %d correct answers, got %d (must not show aggregate of all users)",
				userData.user.Username, expectedCorrect, totalCorrect)
		}
	}

	// Verify aggregate stats across all users would be different from individual stats
	// This ensures we're not accidentally returning the same value for everyone
	var aggregateQuizzes, aggregateQuestions, aggregateCorrect int
	for _, userData := range testUsers {
		aggregateQuizzes += userData.quizCount
		aggregateQuestions += userData.quizCount * userData.questionsPerQuiz
		aggregateCorrect += userData.quizCount * userData.correctPerQuiz
	}

	// Ensure no individual user has the aggregate values (would indicate a bug)
	for _, userData := range testUsers {
		individualQuizzes := userData.quizCount
		individualQuestions := userData.quizCount * userData.questionsPerQuiz
		individualCorrect := userData.quizCount * userData.correctPerQuiz

		if individualQuizzes == aggregateQuizzes || individualQuestions == aggregateQuestions || individualCorrect == aggregateCorrect {
			t.Errorf("Test data for %s matches aggregate values - this would mask bugs. Please use different test data.", userData.user.Username)
		}
	}
}
