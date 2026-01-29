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

// Phase 4: Quiz Attempt History Tests

func TestGetQuizHistoryByUser(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

// Create a user
user := &models.User{Username: "student1", Email: "student1@test.com", Password: "pass123"}
db.CreateUser(user)

// Create a topic
topic := &models.Topic{Name: "Java", Description: "Java Programming"}
db.CreateTopic(topic)

// Create completed quiz attempts
now := time.Now()
for i := 0; i < 5; i++ {
attempt := &models.QuizAttempt{
UserID:          user.ID,
StartTime:       now.Add(-time.Duration(i) * time.Hour),
CompletedTime:   &now,
TotalQuestions:  10,
CorrectAnswers:  7 + i,
QuestionIDs:     []string{"1.1", "1.2"},
QuestionResults: map[string]bool{"1.1": true, "1.2": false},
AttemptNumber:   1,
}
db.CreateQuizAttempt(attempt)
}

// Test quiz history retrieval
req := httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&page=1&per_page=3", nil)
w := httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var attempts []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&attempts)

if len(attempts) != 3 {
t.Errorf("Expected 3 attempts with pagination, got %d", len(attempts))
}
}

func TestGetQuizHistoryByTopic(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

// Create user and topics
user := &models.User{Username: "student1", Email: "student1@test.com"}
db.CreateUser(user)

topic1 := &models.Topic{Name: "Java"}
topic2 := &models.Topic{Name: "Python"}
db.CreateTopic(topic1)
db.CreateTopic(topic2)

// Create questions for each topic
q1 := &models.Question{TopicID: topic1.ID, TopicQuestionID: "1.1", QuestionText: "Java Q1", Difficulty: 3}
q2 := &models.Question{TopicID: topic2.ID, TopicQuestionID: "2.1", QuestionText: "Python Q1", Difficulty: 3}
db.CreateQuestion(q1)
db.CreateQuestion(q2)

// Create quiz attempts with different topics
now := time.Now()
attempt1 := &models.QuizAttempt{
UserID:          user.ID,
CompletedTime:   &now,
QuestionIDs:     []string{"1.1"},
QuestionResults: map[string]bool{"1.1": true},
TotalQuestions:  1,
CorrectAnswers:  1,
AttemptNumber:   1,
}
attempt2 := &models.QuizAttempt{
UserID:          user.ID,
CompletedTime:   &now,
QuestionIDs:     []string{"2.1"},
QuestionResults: map[string]bool{"2.1": true},
TotalQuestions:  1,
CorrectAnswers:  1,
AttemptNumber:   1,
}
db.CreateQuizAttempt(attempt1)
db.CreateQuizAttempt(attempt2)

// Test filtering by topic
req := httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&topic_id=1&page=1&per_page=10", nil)
w := httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var attempts []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&attempts)

if len(attempts) != 1 {
t.Errorf("Expected 1 attempt for topic 1, got %d", len(attempts))
}
}

func TestGetQuizAttemptDetails(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
db.CreateUser(user)

now := time.Now()
attempt := &models.QuizAttempt{
UserID:          user.ID,
StartTime:       now,
CompletedTime:   &now,
TotalQuestions:  5,
CorrectAnswers:  3,
QuestionIDs:     []string{"1.1", "1.2", "1.3", "1.4", "1.5"},
QuestionResults: map[string]bool{"1.1": true, "1.2": false, "1.3": true, "1.4": true, "1.5": false},
AttemptNumber:   1,
}
db.CreateQuizAttempt(attempt)

req := httptest.NewRequest(http.MethodGet, "/students/quiz-history/1", nil)
w := httptest.NewRecorder()
handler.GetQuizAttemptDetailsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var result models.QuizAttempt
json.NewDecoder(w.Body).Decode(&result)

if result.TotalQuestions != 5 {
t.Errorf("Expected 5 total questions, got %d", result.TotalQuestions)
}
if result.CorrectAnswers != 3 {
t.Errorf("Expected 3 correct answers, got %d", result.CorrectAnswers)
}
}

func TestQuizHistoryPagination(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
db.CreateUser(user)

// Create 15 attempts
now := time.Now()
for i := 0; i < 15; i++ {
attempt := &models.QuizAttempt{
UserID:          user.ID,
StartTime:       now.Add(-time.Duration(i) * time.Hour),
CompletedTime:   &now,
TotalQuestions:  5,
CorrectAnswers:  3,
QuestionIDs:     []string{"1.1"},
QuestionResults: map[string]bool{"1.1": true},
AttemptNumber:   1,
}
db.CreateQuizAttempt(attempt)
}

// Test first page
req := httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&page=1&per_page=10", nil)
w := httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

var page1 []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&page1)
if len(page1) != 10 {
t.Errorf("Expected 10 attempts on page 1, got %d", len(page1))
}

// Test second page
req = httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&page=2&per_page=10", nil)
w = httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

var page2 []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&page2)
if len(page2) != 5 {
t.Errorf("Expected 5 attempts on page 2, got %d", len(page2))
}
}

func TestQuizHistorySortByDate(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
db.CreateUser(user)

now := time.Now()
older := now.Add(-2 * time.Hour)
oldest := now.Add(-4 * time.Hour)

// Create attempts in random order
attempts := []*models.QuizAttempt{
{UserID: user.ID, StartTime: oldest, CompletedTime: &oldest, TotalQuestions: 5, CorrectAnswers: 2, QuestionIDs: []string{"1.1"}, QuestionResults: map[string]bool{"1.1": true}, AttemptNumber: 1},
{UserID: user.ID, StartTime: now, CompletedTime: &now, TotalQuestions: 5, CorrectAnswers: 4, QuestionIDs: []string{"1.1"}, QuestionResults: map[string]bool{"1.1": true}, AttemptNumber: 1},
{UserID: user.ID, StartTime: older, CompletedTime: &older, TotalQuestions: 5, CorrectAnswers: 3, QuestionIDs: []string{"1.1"}, QuestionResults: map[string]bool{"1.1": true}, AttemptNumber: 1},
}

for _, att := range attempts {
db.CreateQuizAttempt(att)
}

req := httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&page=1&per_page=10", nil)
w := httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

var result []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&result)

// Should be sorted newest first
if result[0].CorrectAnswers != 4 {
t.Errorf("Expected newest attempt first with 4 correct, got %d", result[0].CorrectAnswers)
}
if result[2].CorrectAnswers != 2 {
t.Errorf("Expected oldest attempt last with 2 correct, got %d", result[2].CorrectAnswers)
}
}

func TestQuizHistoryEmptyState(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
db.CreateUser(user)

req := httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&page=1&per_page=10", nil)
w := httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK for empty history, got %d", w.Code)
}

var attempts []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&attempts)

if len(attempts) != 0 {
t.Errorf("Expected 0 attempts for new user, got %d", len(attempts))
}
}

func TestQuizHistoryIncompleteAttemptsExcluded(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
db.CreateUser(user)

now := time.Now()
// Complete attempt
complete := &models.QuizAttempt{
UserID:          user.ID,
StartTime:       now,
CompletedTime:   &now,
TotalQuestions:  5,
CorrectAnswers:  3,
QuestionIDs:     []string{"1.1"},
QuestionResults: map[string]bool{"1.1": true},
AttemptNumber:   1,
}
// Incomplete attempt
incomplete := &models.QuizAttempt{
UserID:          user.ID,
StartTime:       now,
CompletedTime:   nil,
TotalQuestions:  5,
CorrectAnswers:  0,
QuestionIDs:     []string{"1.1"},
QuestionResults: map[string]bool{},
AttemptNumber:   1,
}
db.CreateQuizAttempt(complete)
db.CreateQuizAttempt(incomplete)

req := httptest.NewRequest(http.MethodGet, "/students/quiz-history?user_id=1&page=1&per_page=10", nil)
w := httptest.NewRecorder()
handler.GetQuizHistoryHandler(w, req)

var attempts []*models.QuizAttempt
json.NewDecoder(w.Body).Decode(&attempts)

if len(attempts) != 1 {
t.Errorf("Expected 1 complete attempt, got %d", len(attempts))
}
}

// Phase 4: Timed Test Taking Tests

func TestGetAvailableTests(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

now := time.Now()
future := now.Add(24 * time.Hour)
past := now.Add(-24 * time.Hour)

// Published test within window
test1 := &models.Test{Name: "Test 1", TopicID: topic.ID, TimeLimit: 60, IsPublished: true}
// Unpublished test
test2 := &models.Test{Name: "Test 2", TopicID: topic.ID, TimeLimit: 60, IsPublished: false}
// Published but not started yet
test3 := &models.Test{Name: "Test 3", TopicID: topic.ID, TimeLimit: 60, IsPublished: true, StartDate: &future}
// Published but ended
test4 := &models.Test{Name: "Test 4", TopicID: topic.ID, TimeLimit: 60, IsPublished: true, EndDate: &past}

db.CreateTest(test1)
db.CreateTest(test2)
db.CreateTest(test3)
db.CreateTest(test4)

req := httptest.NewRequest(http.MethodGet, "/tests/available?user_id=1", nil)
w := httptest.NewRecorder()
handler.GetAvailableTestsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var tests []*models.Test
json.NewDecoder(w.Body).Decode(&tests)

if len(tests) != 1 {
t.Errorf("Expected 1 available test, got %d", len(tests))
}
if tests[0].Name != "Test 1" {
t.Errorf("Expected Test 1, got %s", tests[0].Name)
}
}

func TestStartTest(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60, IsPublished: true}
db.CreateTest(test)

reqBody := `{"user_id":1}`
req := httptest.NewRequest(http.MethodPost, "/tests/1/start", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.StartTestHandler(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected status Created, got %d", w.Code)
}

var attempt models.TestAttempt
json.NewDecoder(w.Body).Decode(&attempt)

if attempt.TestID != test.ID {
t.Errorf("Expected test ID %d, got %d", test.ID, attempt.TestID)
}
if attempt.UserID != user.ID {
t.Errorf("Expected user ID %d, got %d", user.ID, attempt.UserID)
}
}

func TestCannotStartTestTwice(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60, IsPublished: true}
db.CreateTest(test)

// Start test first time
reqBody := `{"user_id":1}`
req := httptest.NewRequest(http.MethodPost, "/tests/1/start", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.StartTestHandler(w, req)

if w.Code != http.StatusCreated {
t.Errorf("First attempt should succeed, got %d", w.Code)
}

// Try to start again without submitting
req = httptest.NewRequest(http.MethodPost, "/tests/1/start", bytes.NewBufferString(reqBody))
w = httptest.NewRecorder()
handler.StartTestHandler(w, req)

if w.Code != http.StatusConflict {
t.Errorf("Expected conflict for second attempt, got %d", w.Code)
}
}

func TestUpdateTestAnswer(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60}
db.CreateTest(test)

attempt := &models.TestAttempt{
UserID:    user.ID,
TestID:    test.ID,
StartTime: time.Now(),
Answers:   make(map[string]int),
}
db.CreateTestAttempt(attempt)

reqBody := `{"question_id":"1","answer_id":42}`
req := httptest.NewRequest(http.MethodPut, "/test-attempts/1/answer", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.UpdateTestAnswerHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

// Verify answer was saved
savedAttempt, _ := db.GetTestAttempt(attempt.ID)
if savedAttempt.Answers["1"] != 42 {
t.Errorf("Expected answer 42, got %d", savedAttempt.Answers["1"])
}
}

func TestSubmitTestManual(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60}
db.CreateTest(test)

// Create questions and answers
q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3}
db.CreateQuestion(q1)
a1 := &models.Answer{QuestionID: q1.ID, AnswerText: "Correct", IsCorrect: true}
a2 := &models.Answer{QuestionID: q1.ID, AnswerText: "Wrong", IsCorrect: false}
db.CreateAnswer(a1)
db.CreateAnswer(a2)

db.AddQuestionToTest(test.ID, q1.ID, 1)

attempt := &models.TestAttempt{
UserID:    user.ID,
TestID:    test.ID,
StartTime: time.Now(),
Answers:   map[string]int{fmt.Sprintf("%d", q1.ID): a1.ID},
}
db.CreateTestAttempt(attempt)

reqBody := `{"is_auto_submit":false}`
req := httptest.NewRequest(http.MethodPost, "/test-attempts/1/submit", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.SubmitTestHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var result models.TestAttempt
json.NewDecoder(w.Body).Decode(&result)

if result.IsAutoSubmit {
t.Error("Expected manual submit, got auto submit")
}
if result.SubmitTime == nil {
t.Error("Expected submit time to be set")
}
if result.Score != 100.0 {
t.Errorf("Expected score 100%%, got %.2f%%", result.Score)
}
}

func TestSubmitTestAutoSubmit(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 1} // 1 minute
db.CreateTest(test)

q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3}
db.CreateQuestion(q1)
a1 := &models.Answer{QuestionID: q1.ID, AnswerText: "Correct", IsCorrect: true}
db.CreateAnswer(a1)
db.AddQuestionToTest(test.ID, q1.ID, 1)

// Start attempt in the past to simulate time expiring
pastTime := time.Now().Add(-2 * time.Minute)
attempt := &models.TestAttempt{
UserID:    user.ID,
TestID:    test.ID,
StartTime: pastTime,
Answers:   make(map[string]int),
}
db.CreateTestAttempt(attempt)

reqBody := `{"is_auto_submit":true}`
req := httptest.NewRequest(http.MethodPost, "/test-attempts/1/submit", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.SubmitTestHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var result models.TestAttempt
json.NewDecoder(w.Body).Decode(&result)

if !result.IsAutoSubmit {
t.Error("Expected auto submit")
}
if result.TimeRemaining != 0 {
t.Errorf("Expected 0 time remaining, got %d", result.TimeRemaining)
}
}

func TestCalculateTestScore(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60}
db.CreateTest(test)

// Create 5 questions
for i := 1; i <= 5; i++ {
q := &models.Question{TopicID: topic.ID, TopicQuestionID: fmt.Sprintf("1.%d", i), QuestionText: fmt.Sprintf("Q%d", i), Difficulty: 3}
db.CreateQuestion(q)
a1 := &models.Answer{QuestionID: q.ID, AnswerText: "Correct", IsCorrect: true}
a2 := &models.Answer{QuestionID: q.ID, AnswerText: "Wrong", IsCorrect: false}
db.CreateAnswer(a1)
db.CreateAnswer(a2)
db.AddQuestionToTest(test.ID, q.ID, i)
}

// Student answers 3 out of 5 correctly
answers := make(map[string]int)
answers["1"] = 1 // correct
answers["2"] = 3 // correct
answers["3"] = 5 // correct
answers["4"] = 8 // wrong
answers["5"] = 10 // wrong

attempt := &models.TestAttempt{
UserID:    user.ID,
TestID:    test.ID,
StartTime: time.Now(),
Answers:   answers,
}
db.CreateTestAttempt(attempt)

reqBody := `{"is_auto_submit":false}`
req := httptest.NewRequest(http.MethodPost, "/test-attempts/1/submit", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.SubmitTestHandler(w, req)

var result models.TestAttempt
json.NewDecoder(w.Body).Decode(&result)

if result.TotalQuestions != 5 {
t.Errorf("Expected 5 total questions, got %d", result.TotalQuestions)
}
if result.CorrectAnswers != 3 {
t.Errorf("Expected 3 correct answers, got %d", result.CorrectAnswers)
}
if result.Score != 60.0 {
t.Errorf("Expected score 60%%, got %.2f%%", result.Score)
}
}

func TestGetTestResults(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60}
db.CreateTest(test)

now := time.Now()
attempt := &models.TestAttempt{
UserID:          user.ID,
TestID:          test.ID,
StartTime:       now,
SubmitTime:      &now,
TotalQuestions:  10,
CorrectAnswers:  7,
Score:           70.0,
IsAutoSubmit:    false,
TimeRemaining:   300,
Answers:         make(map[string]int),
QuestionResults: make(map[string]bool),
}
db.CreateTestAttempt(attempt)

req := httptest.NewRequest(http.MethodGet, "/test-attempts/1/results", nil)
w := httptest.NewRecorder()
handler.GetTestResultsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var result models.TestAttempt
json.NewDecoder(w.Body).Decode(&result)

if result.Score != 70.0 {
t.Errorf("Expected score 70%%, got %.2f%%", result.Score)
}
if result.CorrectAnswers != 7 {
t.Errorf("Expected 7 correct, got %d", result.CorrectAnswers)
}
}

func TestGetTestDetails(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60, Description: "Java Midterm Exam"}
db.CreateTest(test)

// Add 5 questions
for i := 1; i <= 5; i++ {
q := &models.Question{TopicID: topic.ID, TopicQuestionID: fmt.Sprintf("1.%d", i), QuestionText: fmt.Sprintf("Q%d", i), Difficulty: 3}
db.CreateQuestion(q)
db.AddQuestionToTest(test.ID, q.ID, i)
}

req := httptest.NewRequest(http.MethodGet, "/tests/1", nil)
w := httptest.NewRecorder()
handler.GetTestDetailsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var response map[string]interface{}
json.NewDecoder(w.Body).Decode(&response)

testData := response["test"].(map[string]interface{})
if testData["name"].(string) != "Midterm" {
t.Errorf("Expected test name 'Midterm', got %s", testData["name"])
}

questionCount := int(response["question_count"].(float64))
if questionCount != 5 {
t.Errorf("Expected 5 questions, got %d", questionCount)
}
}


// Phase 4: Tag-Based Review Tests

func TestGetTagsWithCounts(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

// Create questions with tags
q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3, Tags: []string{"loops", "arrays"}}
q2 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.2", QuestionText: "Q2", Difficulty: 3, Tags: []string{"loops", "functions"}}
q3 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.3", QuestionText: "Q3", Difficulty: 3, Tags: []string{"arrays"}}
db.CreateQuestion(q1)
db.CreateQuestion(q2)
db.CreateQuestion(q3)

req := httptest.NewRequest(http.MethodGet, "/tags", nil)
w := httptest.NewRecorder()
handler.GetTagsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var tagCounts map[string]int
json.NewDecoder(w.Body).Decode(&tagCounts)

if tagCounts["loops"] != 2 {
t.Errorf("Expected 2 questions with 'loops', got %d", tagCounts["loops"])
}
if tagCounts["arrays"] != 2 {
t.Errorf("Expected 2 questions with 'arrays', got %d", tagCounts["arrays"])
}
if tagCounts["functions"] != 1 {
t.Errorf("Expected 1 question with 'functions', got %d", tagCounts["functions"])
}
}

func TestGetTagStats(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

// Create questions with tags
q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3, Tags: []string{"loops"}}
q2 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.2", QuestionText: "Q2", Difficulty: 3, Tags: []string{"loops"}}
q3 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.3", QuestionText: "Q3", Difficulty: 3, Tags: []string{"arrays"}}
db.CreateQuestion(q1)
db.CreateQuestion(q2)
db.CreateQuestion(q3)

// Create question attempts
qa1 := &models.QuestionAttempt{UserID: user.ID, QuestionID: q1.ID, TopicQuestionID: "1.1", CorrectAttempts: 3, IncorrectAttempts: 1}
qa2 := &models.QuestionAttempt{UserID: user.ID, QuestionID: q2.ID, TopicQuestionID: "1.2", CorrectAttempts: 2, IncorrectAttempts: 2}
qa3 := &models.QuestionAttempt{UserID: user.ID, QuestionID: q3.ID, TopicQuestionID: "1.3", CorrectAttempts: 1, IncorrectAttempts: 4}
db.CreateQuestionAttempt(qa1)
db.CreateQuestionAttempt(qa2)
db.CreateQuestionAttempt(qa3)

req := httptest.NewRequest(http.MethodGet, "/tags/stats?user_id=1", nil)
w := httptest.NewRecorder()
handler.GetTagStatsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var stats []*models.TagStats
json.NewDecoder(w.Body).Decode(&stats)

// Find loops tag
var loopsStats *models.TagStats
for _, s := range stats {
if s.Tag == "loops" {
loopsStats = s
break
}
}

if loopsStats == nil {
t.Fatal("Expected loops tag in stats")
}

// loops: 5 correct, 3 incorrect = 8 total, 62.5% accuracy
if loopsStats.TotalAttempts != 8 {
t.Errorf("Expected 8 total attempts for loops, got %d", loopsStats.TotalAttempts)
}
if loopsStats.CorrectAttempts != 5 {
t.Errorf("Expected 5 correct attempts for loops, got %d", loopsStats.CorrectAttempts)
}
}

func TestGenerateQuizByTagsOR(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

// Create questions with different tags
q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3, Tags: []string{"loops"}}
q2 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.2", QuestionText: "Q2", Difficulty: 3, Tags: []string{"arrays"}}
q3 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.3", QuestionText: "Q3", Difficulty: 3, Tags: []string{"functions"}}
q4 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.4", QuestionText: "Q4", Difficulty: 3, Tags: []string{"loops", "arrays"}}
db.CreateQuestion(q1)
db.CreateQuestion(q2)
db.CreateQuestion(q3)
db.CreateQuestion(q4)

reqBody := `{"tags":["loops","arrays"],"use_and":false,"user_id":1,"count":10}`
req := httptest.NewRequest(http.MethodPost, "/quiz/generate/by-tags", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.GenerateQuizByTagsHandler(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected status Created, got %d", w.Code)
}

var response map[string]interface{}
json.NewDecoder(w.Body).Decode(&response)

questions := response["questions"].([]interface{})
// Should get q1, q2, q4 (OR logic)
if len(questions) != 3 {
t.Errorf("Expected 3 questions with OR logic, got %d", len(questions))
}
}

func TestGenerateQuizByTagsAND(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

// Create questions with different tags
q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3, Tags: []string{"loops"}}
q2 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.2", QuestionText: "Q2", Difficulty: 3, Tags: []string{"arrays"}}
q3 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.3", QuestionText: "Q3", Difficulty: 3, Tags: []string{"functions"}}
q4 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.4", QuestionText: "Q4", Difficulty: 3, Tags: []string{"loops", "arrays"}}
db.CreateQuestion(q1)
db.CreateQuestion(q2)
db.CreateQuestion(q3)
db.CreateQuestion(q4)

reqBody := `{"tags":["loops","arrays"],"use_and":true,"user_id":1,"count":10}`
req := httptest.NewRequest(http.MethodPost, "/quiz/generate/by-tags", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.GenerateQuizByTagsHandler(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected status Created, got %d", w.Code)
}

var response map[string]interface{}
json.NewDecoder(w.Body).Decode(&response)

questions := response["questions"].([]interface{})
// Should get only q4 (AND logic)
if len(questions) != 1 {
t.Errorf("Expected 1 question with AND logic, got %d", len(questions))
}
}

func TestGenerateQuizByTagsNoResults(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Q1", Difficulty: 3, Tags: []string{"loops"}}
db.CreateQuestion(q1)

reqBody := `{"tags":["nonexistent"],"use_and":false,"user_id":1,"count":10}`
req := httptest.NewRequest(http.MethodPost, "/quiz/generate/by-tags", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.GenerateQuizByTagsHandler(w, req)

if w.Code != http.StatusNotFound {
t.Errorf("Expected status Not Found, got %d", w.Code)
}
}

func TestGenerateQuizByTagsLimit(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

user := &models.User{Username: "student1", Email: "student1@test.com"}
topic := &models.Topic{Name: "Java"}
db.CreateUser(user)
db.CreateTopic(topic)

// Create 10 questions with the same tag
for i := 1; i <= 10; i++ {
q := &models.Question{TopicID: topic.ID, TopicQuestionID: fmt.Sprintf("1.%d", i), QuestionText: fmt.Sprintf("Q%d", i), Difficulty: 3, Tags: []string{"loops"}}
db.CreateQuestion(q)
}

reqBody := `{"tags":["loops"],"use_and":false,"user_id":1,"count":5}`
req := httptest.NewRequest(http.MethodPost, "/quiz/generate/by-tags", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.GenerateQuizByTagsHandler(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected status Created, got %d", w.Code)
}

var response map[string]interface{}
json.NewDecoder(w.Body).Decode(&response)

questions := response["questions"].([]interface{})
if len(questions) != 5 {
t.Errorf("Expected quiz limited to 5 questions, got %d", len(questions))
}
}

// Phase 4: Test Management (Instructor) Tests

func TestPublishTest(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60, IsPublished: false}
db.CreateTest(test)

req := httptest.NewRequest(http.MethodPost, "/instructor/tests/1/publish", nil)
w := httptest.NewRecorder()
handler.PublishTestHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

// Verify test is published
publishedTest, _ := db.GetTest(test.ID)
if !publishedTest.IsPublished {
t.Error("Expected test to be published")
}
}

func TestGetTestResultsByTest(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60}
db.CreateTest(test)

// Create users and attempts
for i := 1; i <= 3; i++ {
user := &models.User{Username: fmt.Sprintf("student%d", i), Email: fmt.Sprintf("s%d@test.com", i)}
db.CreateUser(user)

now := time.Now()
attempt := &models.TestAttempt{
UserID:         user.ID,
TestID:         test.ID,
StartTime:      now,
SubmitTime:     &now,
TotalQuestions: 10,
CorrectAnswers: 5 + i,
Score:          float64(50 + i*10),
TimeRemaining:  300,
}
db.CreateTestAttempt(attempt)
}

req := httptest.NewRequest(http.MethodGet, "/instructor/tests/1/results", nil)
w := httptest.NewRecorder()
handler.GetTestResultsInstructorHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var results models.TestResult
json.NewDecoder(w.Body).Decode(&results)

if results.TestName != "Midterm" {
t.Errorf("Expected test name 'Midterm', got %s", results.TestName)
}
if len(results.StudentResults) != 3 {
t.Errorf("Expected 3 student results, got %d", len(results.StudentResults))
}
}

func TestExportTestResults(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

test := &models.Test{Name: "Midterm", TopicID: topic.ID, TimeLimit: 60}
db.CreateTest(test)

user := &models.User{Username: "student1", Email: "s1@test.com"}
db.CreateUser(user)

now := time.Now()
attempt := &models.TestAttempt{
UserID:         user.ID,
TestID:         test.ID,
StartTime:      now,
SubmitTime:     &now,
TotalQuestions: 10,
CorrectAnswers: 7,
Score:          70.0,
TimeRemaining:  300,
}
db.CreateTestAttempt(attempt)

req := httptest.NewRequest(http.MethodGet, "/instructor/tests/1/export", nil)
w := httptest.NewRecorder()
handler.ExportTestResultsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

contentType := w.Header().Get("Content-Type")
if contentType != "text/csv" {
t.Errorf("Expected Content-Type text/csv, got %s", contentType)
}

csvData := w.Body.String()
if !bytes.Contains([]byte(csvData), []byte("student1")) {
t.Error("CSV should contain student username")
}
if !bytes.Contains([]byte(csvData), []byte("70.00")) {
t.Error("CSV should contain score")
}
}

// Phase 4: Question Performance Analytics Tests

func TestGetQuestionAnalytics(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

user := &models.User{Username: "student1", Email: "s1@test.com"}
db.CreateUser(user)

// Create questions
q1 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Easy Q", Difficulty: 2}
q2 := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.2", QuestionText: "Hard Q", Difficulty: 6}
db.CreateQuestion(q1)
db.CreateQuestion(q2)

// Create attempts (q1 very easy, q2 very hard)
qa1 := &models.QuestionAttempt{UserID: user.ID, QuestionID: q1.ID, TopicQuestionID: "1.1", CorrectAttempts: 95, IncorrectAttempts: 5}
qa2 := &models.QuestionAttempt{UserID: user.ID, QuestionID: q2.ID, TopicQuestionID: "1.2", CorrectAttempts: 10, IncorrectAttempts: 90}
db.CreateQuestionAttempt(qa1)
db.CreateQuestionAttempt(qa2)

req := httptest.NewRequest(http.MethodGet, "/instructor/questions/analytics", nil)
w := httptest.NewRecorder()
handler.GetQuestionAnalyticsHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var analytics []*models.QuestionAnalytics
json.NewDecoder(w.Body).Decode(&analytics)

if len(analytics) != 2 {
t.Errorf("Expected 2 question analytics, got %d", len(analytics))
}

// Find q1 analytics
var q1Analytics *models.QuestionAnalytics
for _, a := range analytics {
if a.QuestionID == q1.ID {
q1Analytics = a
break
}
}

if q1Analytics == nil {
t.Fatal("Expected q1 in analytics")
}

if q1Analytics.Accuracy != 95.0 {
t.Errorf("Expected 95%% accuracy for q1, got %.2f%%", q1Analytics.Accuracy)
}
if !q1Analytics.FlagForReview {
t.Error("Expected q1 to be flagged for review (too easy)")
}
}

func TestQuestionAnalyticsTooEasy(t *testing.T) {
db := database.NewInMemoryDB()
_ = NewHandler(db)

topic := &models.Topic{Name: "Java"}
user := &models.User{Username: "student1", Email: "s1@test.com"}
db.CreateTopic(topic)
db.CreateUser(user)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Too Easy", Difficulty: 3}
db.CreateQuestion(q)

qa := &models.QuestionAttempt{UserID: user.ID, QuestionID: q.ID, TopicQuestionID: "1.1", CorrectAttempts: 95, IncorrectAttempts: 5}
db.CreateQuestionAttempt(qa)

analytics, _ := db.GetQuestionDetailedAnalytics(q.ID)

if analytics.Accuracy != 95.0 {
t.Errorf("Expected 95%% accuracy, got %.2f%%", analytics.Accuracy)
}
if !analytics.FlagForReview {
t.Error("Expected question to be flagged for review")
}
if analytics.RecommendedDifficulty != 4 {
t.Errorf("Expected recommended difficulty 4, got %d", analytics.RecommendedDifficulty)
}
}

func TestQuestionAnalyticsTooHard(t *testing.T) {
db := database.NewInMemoryDB()
_ = NewHandler(db)

topic := &models.Topic{Name: "Java"}
user := &models.User{Username: "student1", Email: "s1@test.com"}
db.CreateTopic(topic)
db.CreateUser(user)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Too Hard", Difficulty: 5}
db.CreateQuestion(q)

qa := &models.QuestionAttempt{UserID: user.ID, QuestionID: q.ID, TopicQuestionID: "1.1", CorrectAttempts: 10, IncorrectAttempts: 90}
db.CreateQuestionAttempt(qa)

analytics, _ := db.GetQuestionDetailedAnalytics(q.ID)

if analytics.Accuracy != 10.0 {
t.Errorf("Expected 10%% accuracy, got %.2f%%", analytics.Accuracy)
}
if !analytics.FlagForReview {
t.Error("Expected question to be flagged for review")
}
if analytics.RecommendedDifficulty != 4 {
t.Errorf("Expected recommended difficulty 4, got %d", analytics.RecommendedDifficulty)
}
}

func TestQuestionAnalyticsAppropriate(t *testing.T) {
db := database.NewInMemoryDB()
_ = NewHandler(db)

topic := &models.Topic{Name: "Java"}
user := &models.User{Username: "student1", Email: "s1@test.com"}
db.CreateTopic(topic)
db.CreateUser(user)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Just Right", Difficulty: 4}
db.CreateQuestion(q)

qa := &models.QuestionAttempt{UserID: user.ID, QuestionID: q.ID, TopicQuestionID: "1.1", CorrectAttempts: 75, IncorrectAttempts: 25}
db.CreateQuestionAttempt(qa)

analytics, _ := db.GetQuestionDetailedAnalytics(q.ID)

if analytics.Accuracy != 75.0 {
t.Errorf("Expected 75%% accuracy, got %.2f%%", analytics.Accuracy)
}
if analytics.FlagForReview {
t.Error("Expected question NOT to be flagged (appropriate difficulty)")
}
if analytics.RecommendedDifficulty != 4 {
t.Errorf("Expected recommended difficulty same as current (4), got %d", analytics.RecommendedDifficulty)
}
}

func TestQuestionAnalyticsNoAttempts(t *testing.T) {
db := database.NewInMemoryDB()
_ = NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Never Attempted", Difficulty: 4}
db.CreateQuestion(q)

analytics, _ := db.GetQuestionDetailedAnalytics(q.ID)

if analytics.TotalAttempts != 0 {
t.Errorf("Expected 0 attempts, got %d", analytics.TotalAttempts)
}
if analytics.Accuracy != 0.0 {
t.Errorf("Expected 0%% accuracy, got %.2f%%", analytics.Accuracy)
}
if analytics.FlagForReview {
t.Error("Expected question NOT to be flagged (no data)")
}
}

func TestUpdateQuestionDifficulty(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Test Q", Difficulty: 3}
db.CreateQuestion(q)

reqBody := `{"difficulty":5}`
req := httptest.NewRequest(http.MethodPut, "/instructor/questions/1/difficulty", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.UpdateQuestionDifficultyHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

// Verify difficulty was updated
updated, _ := db.GetQuestion(q.ID)
if updated.Difficulty != 5 {
t.Errorf("Expected difficulty 5, got %d", updated.Difficulty)
}
}

func TestUpdateQuestionDifficultyInvalidRange(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Test Q", Difficulty: 3}
db.CreateQuestion(q)

// Try to set difficulty to 8 (out of range)
reqBody := `{"difficulty":8}`
req := httptest.NewRequest(http.MethodPut, "/instructor/questions/1/difficulty", bytes.NewBufferString(reqBody))
w := httptest.NewRecorder()
handler.UpdateQuestionDifficultyHandler(w, req)

if w.Code == http.StatusOK {
t.Error("Expected error for difficulty out of range")
}
}

func TestGetAnswerDistribution(t *testing.T) {
db := database.NewInMemoryDB()
handler := NewHandler(db)

topic := &models.Topic{Name: "Java"}
db.CreateTopic(topic)

q := &models.Question{TopicID: topic.ID, TopicQuestionID: "1.1", QuestionText: "Test Q", Difficulty: 3}
db.CreateQuestion(q)

req := httptest.NewRequest(http.MethodGet, "/instructor/questions/1/answer-distribution", nil)
w := httptest.NewRecorder()
handler.GetAnswerDistributionHandler(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected status OK, got %d", w.Code)
}

var distribution map[int]int
json.NewDecoder(w.Body).Decode(&distribution)

// Currently returns empty map (not fully implemented)
if distribution == nil {
t.Error("Expected non-nil distribution map")
}
}

