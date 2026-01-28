package database

import (
	"errors"
	"fmt"
	"time"
	"github.com/ZipCodeCore/zip-athena/models"
	"sync"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

// Database interface defines the contract for data operations
type Database interface {
	// User operations
	CreateUser(user *models.User) error
	GetUser(id int) (*models.User, error)
	UpdateUserPoints(id int, points int) error

	// Topic operations
	CreateTopic(topic *models.Topic) error
	GetTopic(id int) (*models.Topic, error)
	ListTopics() ([]*models.Topic, error)

	// Question operations
	CreateQuestion(question *models.Question) error
	GetQuestion(id int) (*models.Question, error)
	GetQuestionsByTopic(topicID int) ([]*models.Question, error)

	// Answer operations
	CreateAnswer(answer *models.Answer) error
	GetAnswersByQuestion(questionID int) ([]*models.Answer, error)

	// QuizAttempt operations
	CreateQuizAttempt(attempt *models.QuizAttempt) error
	GetQuizAttempt(id int) (*models.QuizAttempt, error)
	GetQuizAttemptsByUser(userID int) ([]*models.QuizAttempt, error)
	UpdateQuizAttempt(attempt *models.QuizAttempt) error

	// QuestionAttempt operations
	CreateQuestionAttempt(attempt *models.QuestionAttempt) error
	GetQuestionAttempt(id int) (*models.QuestionAttempt, error)
	GetQuestionAttemptsByUser(userID int) ([]*models.QuestionAttempt, error)
	UpdateQuestionAttempt(attempt *models.QuestionAttempt) error

	// Quiz operations (Phase 2)
	CreateQuiz(quiz *models.Quiz) error
	GetQuiz(id int) (*models.Quiz, error)
	ListQuizzes() ([]*models.Quiz, error)
	UpdateQuiz(quiz *models.Quiz) error
	DeleteQuiz(id int) error

	// UserQuestionState operations (Phase 2 - FSRS)
	CreateUserQuestionState(state *models.UserQuestionState) error
	GetUserQuestionState(userID, questionID int) (*models.UserQuestionState, error)
	GetUserQuestionStatesByUser(userID int) ([]*models.UserQuestionState, error)
	GetDueUserQuestionStates(userID int) ([]*models.UserQuestionState, error)
	UpdateUserQuestionState(state *models.UserQuestionState) error
	
	// Admin operations (Phase 2)
	ListUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
	ListQuestions() ([]*models.Question, error)
	UpdateQuestion(question *models.Question) error
	DeleteQuestion(id int) error
	UpdateAnswer(answer *models.Answer) error
	DeleteAnswer(id int) error
	UpdateTopic(topic *models.Topic) error
	DeleteTopic(id int) error
	
	// Phase 3: Test operations
	CreateTest(test *models.Test) error
	GetTest(id int) (*models.Test, error)
	ListTests() ([]*models.Test, error)
	UpdateTest(test *models.Test) error
	DeleteTest(id int) error
	
	// Phase 3: Many-to-many Quiz-Question relationships
	AddQuestionToQuiz(quizID, questionID, order int) error
	RemoveQuestionFromQuiz(quizID, questionID int) error
	GetQuizQuestions(quizID int) ([]*models.Question, error)
	
	// Phase 3: Many-to-many Test-Question relationships
	AddQuestionToTest(testID, questionID, order int) error
	RemoveQuestionFromTest(testID, questionID int) error
	GetTestQuestions(testID int) ([]*models.Question, error)
	
	// Phase 3: Analytics operations
	GetStudentProgress(userID int) (*models.StudentProgress, error)
	GetCohortAnalytics() (*models.CohortAnalytics, error)
	GetAnswer(id int) (*models.Answer, error)
}

// InMemoryDB implements the Database interface using in-memory storage
type InMemoryDB struct {
	users               map[int]*models.User
	topics              map[int]*models.Topic
	questions           map[int]*models.Question
	answers             map[int]*models.Answer
	quizAttempts        map[int]*models.QuizAttempt
	questionAttempts    map[int]*models.QuestionAttempt
	quizzes             map[int]*models.Quiz
	userQuestionStates  map[string]*models.UserQuestionState // key: "userID-questionID"
	tests               map[int]*models.Test                 // Phase 3
	quizQuestions       map[string]*models.QuizQuestion      // key: "quizID-questionID" Phase 3
	testQuestions       map[string]*models.TestQuestion      // key: "testID-questionID" Phase 3
	nextUserID          int
	nextTopicID         int
	nextQuestionID      int
	nextAnswerID        int
	nextQuizAttemptID   int
	nextQAttemptID      int
	nextQuizID          int
	nextUQStateID       int
	nextTestID          int // Phase 3
	nextQuizQuestionID  int // Phase 3
	nextTestQuestionID  int // Phase 3
	mu                  sync.RWMutex
}

// NewInMemoryDB creates a new in-memory database
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		users:              make(map[int]*models.User),
		topics:             make(map[int]*models.Topic),
		questions:          make(map[int]*models.Question),
		answers:            make(map[int]*models.Answer),
		quizAttempts:       make(map[int]*models.QuizAttempt),
		questionAttempts:   make(map[int]*models.QuestionAttempt),
		quizzes:            make(map[int]*models.Quiz),
		userQuestionStates: make(map[string]*models.UserQuestionState),
		tests:              make(map[int]*models.Test),
		quizQuestions:      make(map[string]*models.QuizQuestion),
		testQuestions:      make(map[string]*models.TestQuestion),
		nextUserID:         1,
		nextTopicID:        1,
		nextQuestionID:     1,
		nextAnswerID:       1,
		nextQuizAttemptID:  1,
		nextQAttemptID:     1,
		nextQuizID:         1,
		nextUQStateID:      1,
		nextTestID:         1,
		nextQuizQuestionID: 1,
		nextTestQuestionID: 1,
	}
}

// User operations
func (db *InMemoryDB) CreateUser(user *models.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	user.ID = db.nextUserID
	db.nextUserID++
	db.users[user.ID] = user
	return nil
}

func (db *InMemoryDB) GetUser(id int) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, exists := db.users[id]
	if !exists {
		return nil, ErrNotFound
	}
	return user, nil
}

func (db *InMemoryDB) UpdateUserPoints(id int, points int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	user, exists := db.users[id]
	if !exists {
		return ErrNotFound
	}
	user.TotalPoints = points
	return nil
}

// Topic operations
func (db *InMemoryDB) CreateTopic(topic *models.Topic) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	topic.ID = db.nextTopicID
	db.nextTopicID++
	db.topics[topic.ID] = topic
	return nil
}

func (db *InMemoryDB) GetTopic(id int) (*models.Topic, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	topic, exists := db.topics[id]
	if !exists {
		return nil, ErrNotFound
	}
	return topic, nil
}

func (db *InMemoryDB) ListTopics() ([]*models.Topic, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	topics := make([]*models.Topic, 0, len(db.topics))
	for _, topic := range db.topics {
		topics = append(topics, topic)
	}
	return topics, nil
}

// Question operations
func (db *InMemoryDB) CreateQuestion(question *models.Question) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	question.ID = db.nextQuestionID
	db.nextQuestionID++
	db.questions[question.ID] = question
	return nil
}

func (db *InMemoryDB) GetQuestion(id int) (*models.Question, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	question, exists := db.questions[id]
	if !exists {
		return nil, ErrNotFound
	}
	return question, nil
}

func (db *InMemoryDB) GetQuestionsByTopic(topicID int) ([]*models.Question, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	questions := make([]*models.Question, 0)
	for _, question := range db.questions {
		if question.TopicID == topicID {
			questions = append(questions, question)
		}
	}
	return questions, nil
}

// Answer operations
func (db *InMemoryDB) CreateAnswer(answer *models.Answer) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	answer.ID = db.nextAnswerID
	db.nextAnswerID++
	db.answers[answer.ID] = answer
	return nil
}

func (db *InMemoryDB) GetAnswersByQuestion(questionID int) ([]*models.Answer, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	answers := make([]*models.Answer, 0)
	for _, answer := range db.answers {
		if answer.QuestionID == questionID {
			answers = append(answers, answer)
		}
	}
	return answers, nil
}

// QuizAttempt operations
func (db *InMemoryDB) CreateQuizAttempt(attempt *models.QuizAttempt) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	attempt.ID = db.nextQuizID
	db.nextQuizID++
	db.quizAttempts[attempt.ID] = attempt
	return nil
}

func (db *InMemoryDB) GetQuizAttempt(id int) (*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempt, exists := db.quizAttempts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return attempt, nil
}

func (db *InMemoryDB) GetQuizAttemptsByUser(userID int) ([]*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempts := make([]*models.QuizAttempt, 0)
	for _, attempt := range db.quizAttempts {
		if attempt.UserID == userID {
			attempts = append(attempts, attempt)
		}
	}
	return attempts, nil
}

func (db *InMemoryDB) UpdateQuizAttempt(attempt *models.QuizAttempt) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.quizAttempts[attempt.ID]; !exists {
		return ErrNotFound
	}
	db.quizAttempts[attempt.ID] = attempt
	return nil
}

// QuestionAttempt operations
func (db *InMemoryDB) CreateQuestionAttempt(attempt *models.QuestionAttempt) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	attempt.ID = db.nextQAttemptID
	db.nextQAttemptID++
	db.questionAttempts[attempt.ID] = attempt
	return nil
}

func (db *InMemoryDB) GetQuestionAttempt(id int) (*models.QuestionAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempt, exists := db.questionAttempts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return attempt, nil
}

func (db *InMemoryDB) GetQuestionAttemptsByUser(userID int) ([]*models.QuestionAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempts := make([]*models.QuestionAttempt, 0)
	for _, attempt := range db.questionAttempts {
		if attempt.UserID == userID {
			attempts = append(attempts, attempt)
		}
	}
	return attempts, nil
}

func (db *InMemoryDB) UpdateQuestionAttempt(attempt *models.QuestionAttempt) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.questionAttempts[attempt.ID]; !exists {
		return ErrNotFound
	}
	db.questionAttempts[attempt.ID] = attempt
	return nil
}

// Quiz operations
func (db *InMemoryDB) CreateQuiz(quiz *models.Quiz) error {
db.mu.Lock()
defer db.mu.Unlock()

quiz.ID = db.nextQuizID
db.nextQuizID++
db.quizzes[quiz.ID] = quiz
return nil
}

func (db *InMemoryDB) GetQuiz(id int) (*models.Quiz, error) {
db.mu.RLock()
defer db.mu.RUnlock()

quiz, exists := db.quizzes[id]
if !exists {
return nil, ErrNotFound
}
return quiz, nil
}

func (db *InMemoryDB) ListQuizzes() ([]*models.Quiz, error) {
db.mu.RLock()
defer db.mu.RUnlock()

quizzes := make([]*models.Quiz, 0, len(db.quizzes))
for _, quiz := range db.quizzes {
quizzes = append(quizzes, quiz)
}
return quizzes, nil
}

func (db *InMemoryDB) UpdateQuiz(quiz *models.Quiz) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.quizzes[quiz.ID]; !exists {
return ErrNotFound
}
db.quizzes[quiz.ID] = quiz
return nil
}

func (db *InMemoryDB) DeleteQuiz(id int) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.quizzes[id]; !exists {
return ErrNotFound
}
delete(db.quizzes, id)
return nil
}

// UserQuestionState operations
func (db *InMemoryDB) CreateUserQuestionState(state *models.UserQuestionState) error {
db.mu.Lock()
defer db.mu.Unlock()

state.ID = db.nextUQStateID
db.nextUQStateID++
key := fmt.Sprintf("%d-%d", state.UserID, state.QuestionID)
db.userQuestionStates[key] = state
return nil
}

func (db *InMemoryDB) GetUserQuestionState(userID, questionID int) (*models.UserQuestionState, error) {
db.mu.RLock()
defer db.mu.RUnlock()

key := fmt.Sprintf("%d-%d", userID, questionID)
state, exists := db.userQuestionStates[key]
if !exists {
return nil, ErrNotFound
}
return state, nil
}

func (db *InMemoryDB) GetUserQuestionStatesByUser(userID int) ([]*models.UserQuestionState, error) {
db.mu.RLock()
defer db.mu.RUnlock()

states := make([]*models.UserQuestionState, 0)
for _, state := range db.userQuestionStates {
if state.UserID == userID {
states = append(states, state)
}
}
return states, nil
}

func (db *InMemoryDB) GetDueUserQuestionStates(userID int) ([]*models.UserQuestionState, error) {
db.mu.RLock()
defer db.mu.RUnlock()

now := time.Now()
states := make([]*models.UserQuestionState, 0)
for _, state := range db.userQuestionStates {
if state.UserID == userID && !state.Due.After(now) {
states = append(states, state)
}
}
return states, nil
}

func (db *InMemoryDB) UpdateUserQuestionState(state *models.UserQuestionState) error {
db.mu.Lock()
defer db.mu.Unlock()

key := fmt.Sprintf("%d-%d", state.UserID, state.QuestionID)
if _, exists := db.userQuestionStates[key]; !exists {
return ErrNotFound
}
db.userQuestionStates[key] = state
return nil
}

// Admin operations - Users
func (db *InMemoryDB) ListUsers() ([]*models.User, error) {
db.mu.RLock()
defer db.mu.RUnlock()

users := make([]*models.User, 0, len(db.users))
for _, user := range db.users {
users = append(users, user)
}
return users, nil
}

func (db *InMemoryDB) UpdateUser(user *models.User) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.users[user.ID]; !exists {
return ErrNotFound
}
db.users[user.ID] = user
return nil
}

func (db *InMemoryDB) DeleteUser(id int) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.users[id]; !exists {
return ErrNotFound
}
delete(db.users, id)
return nil
}

// Admin operations - Questions
func (db *InMemoryDB) ListQuestions() ([]*models.Question, error) {
db.mu.RLock()
defer db.mu.RUnlock()

questions := make([]*models.Question, 0, len(db.questions))
for _, question := range db.questions {
questions = append(questions, question)
}
return questions, nil
}

func (db *InMemoryDB) UpdateQuestion(question *models.Question) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.questions[question.ID]; !exists {
return ErrNotFound
}
db.questions[question.ID] = question
return nil
}

func (db *InMemoryDB) DeleteQuestion(id int) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.questions[id]; !exists {
return ErrNotFound
}
delete(db.questions, id)
return nil
}

// Admin operations - Answers
func (db *InMemoryDB) UpdateAnswer(answer *models.Answer) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.answers[answer.ID]; !exists {
return ErrNotFound
}
db.answers[answer.ID] = answer
return nil
}

func (db *InMemoryDB) DeleteAnswer(id int) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.answers[id]; !exists {
return ErrNotFound
}
delete(db.answers, id)
return nil
}

// Admin operations - Topics
func (db *InMemoryDB) UpdateTopic(topic *models.Topic) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.topics[topic.ID]; !exists {
return ErrNotFound
}
db.topics[topic.ID] = topic
return nil
}

func (db *InMemoryDB) DeleteTopic(id int) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.topics[id]; !exists {
return ErrNotFound
}
delete(db.topics, id)
return nil
}

func (db *InMemoryDB) GetAnswer(id int) (*models.Answer, error) {
db.mu.RLock()
defer db.mu.RUnlock()

answer, exists := db.answers[id]
if !exists {
return nil, ErrNotFound
}
return answer, nil
}

// Phase 3: Test operations
func (db *InMemoryDB) CreateTest(test *models.Test) error {
db.mu.Lock()
defer db.mu.Unlock()

test.ID = db.nextTestID
db.nextTestID++
db.tests[test.ID] = test
return nil
}

func (db *InMemoryDB) GetTest(id int) (*models.Test, error) {
db.mu.RLock()
defer db.mu.RUnlock()

test, exists := db.tests[id]
if !exists {
return nil, ErrNotFound
}
return test, nil
}

func (db *InMemoryDB) ListTests() ([]*models.Test, error) {
db.mu.RLock()
defer db.mu.RUnlock()

tests := make([]*models.Test, 0, len(db.tests))
for _, test := range db.tests {
tests = append(tests, test)
}
return tests, nil
}

func (db *InMemoryDB) UpdateTest(test *models.Test) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.tests[test.ID]; !exists {
return ErrNotFound
}
db.tests[test.ID] = test
return nil
}

func (db *InMemoryDB) DeleteTest(id int) error {
db.mu.Lock()
defer db.mu.Unlock()

if _, exists := db.tests[id]; !exists {
return ErrNotFound
}
delete(db.tests, id)
return nil
}

// Phase 3: Many-to-many Quiz-Question relationships
func (db *InMemoryDB) AddQuestionToQuiz(quizID, questionID, order int) error {
db.mu.Lock()
defer db.mu.Unlock()

// Verify quiz and question exist
if _, exists := db.quizzes[quizID]; !exists {
return ErrNotFound
}
if _, exists := db.questions[questionID]; !exists {
return ErrNotFound
}

key := fmt.Sprintf("%d-%d", quizID, questionID)
qq := &models.QuizQuestion{
ID:         db.nextQuizQuestionID,
QuizID:     quizID,
QuestionID: questionID,
Order:      order,
}
db.nextQuizQuestionID++
db.quizQuestions[key] = qq
return nil
}

func (db *InMemoryDB) RemoveQuestionFromQuiz(quizID, questionID int) error {
db.mu.Lock()
defer db.mu.Unlock()

key := fmt.Sprintf("%d-%d", quizID, questionID)
if _, exists := db.quizQuestions[key]; !exists {
return ErrNotFound
}
delete(db.quizQuestions, key)
return nil
}

func (db *InMemoryDB) GetQuizQuestions(quizID int) ([]*models.Question, error) {
db.mu.RLock()
defer db.mu.RUnlock()

// Collect quiz questions with their order
type orderedQuestion struct {
question *models.Question
order    int
}
orderedQuestions := make([]orderedQuestion, 0)

for _, qq := range db.quizQuestions {
if qq.QuizID == quizID {
if question, exists := db.questions[qq.QuestionID]; exists {
orderedQuestions = append(orderedQuestions, orderedQuestion{question, qq.Order})
}
}
}

// Sort by order
for i := 0; i < len(orderedQuestions); i++ {
for j := i + 1; j < len(orderedQuestions); j++ {
if orderedQuestions[i].order > orderedQuestions[j].order {
orderedQuestions[i], orderedQuestions[j] = orderedQuestions[j], orderedQuestions[i]
}
}
}

questions := make([]*models.Question, len(orderedQuestions))
for i, oq := range orderedQuestions {
questions[i] = oq.question
}

return questions, nil
}

// Phase 3: Many-to-many Test-Question relationships
func (db *InMemoryDB) AddQuestionToTest(testID, questionID, order int) error {
db.mu.Lock()
defer db.mu.Unlock()

// Verify test and question exist
if _, exists := db.tests[testID]; !exists {
return ErrNotFound
}
if _, exists := db.questions[questionID]; !exists {
return ErrNotFound
}

key := fmt.Sprintf("%d-%d", testID, questionID)
tq := &models.TestQuestion{
ID:         db.nextTestQuestionID,
TestID:     testID,
QuestionID: questionID,
Order:      order,
}
db.nextTestQuestionID++
db.testQuestions[key] = tq
return nil
}

func (db *InMemoryDB) RemoveQuestionFromTest(testID, questionID int) error {
db.mu.Lock()
defer db.mu.Unlock()

key := fmt.Sprintf("%d-%d", testID, questionID)
if _, exists := db.testQuestions[key]; !exists {
return ErrNotFound
}
delete(db.testQuestions, key)
return nil
}

func (db *InMemoryDB) GetTestQuestions(testID int) ([]*models.Question, error) {
db.mu.RLock()
defer db.mu.RUnlock()

// Collect test questions with their order
type orderedQuestion struct {
question *models.Question
order    int
}
orderedQuestions := make([]orderedQuestion, 0)

for _, tq := range db.testQuestions {
if tq.TestID == testID {
if question, exists := db.questions[tq.QuestionID]; exists {
orderedQuestions = append(orderedQuestions, orderedQuestion{question, tq.Order})
}
}
}

// Sort by order
for i := 0; i < len(orderedQuestions); i++ {
for j := i + 1; j < len(orderedQuestions); j++ {
if orderedQuestions[i].order > orderedQuestions[j].order {
orderedQuestions[i], orderedQuestions[j] = orderedQuestions[j], orderedQuestions[i]
}
}
}

questions := make([]*models.Question, len(orderedQuestions))
for i, oq := range orderedQuestions {
questions[i] = oq.question
}

return questions, nil
}

// Phase 3: Analytics operations
func (db *InMemoryDB) GetStudentProgress(userID int) (*models.StudentProgress, error) {
db.mu.RLock()
defer db.mu.RUnlock()

user, exists := db.users[userID]
if !exists {
return nil, ErrNotFound
}

progress := &models.StudentProgress{
UserID:           userID,
Username:         user.Username,
TotalPoints:      user.TotalPoints,
TopicBreakdown:   make(map[string]float64),
}

// Count quiz attempts
var completedQuizzes int
var totalScore float64
recentActivity := make([]models.QuizAttempt, 0)

for _, attempt := range db.quizAttempts {
if attempt.UserID == userID {
progress.TotalQuizzes++
if attempt.CompletedTime != nil {
completedQuizzes++
if attempt.TotalQuestions > 0 {
score := float64(attempt.CorrectAnswers) / float64(attempt.TotalQuestions) * 100
totalScore += score
}
recentActivity = append(recentActivity, *attempt)
}
}
}

progress.CompletedQuizzes = completedQuizzes
if completedQuizzes > 0 {
progress.AverageScore = totalScore / float64(completedQuizzes)
}

// Sort and limit recent activity to 5
// Simple bubble sort by start time (most recent first)
for i := 0; i < len(recentActivity); i++ {
for j := i + 1; j < len(recentActivity); j++ {
if recentActivity[i].StartTime.Before(recentActivity[j].StartTime) {
recentActivity[i], recentActivity[j] = recentActivity[j], recentActivity[i]
}
}
}
if len(recentActivity) > 5 {
recentActivity = recentActivity[:5]
}
progress.RecentActivity = recentActivity

// Calculate question accuracy
var totalAttempts, correctAttempts int
for _, qa := range db.questionAttempts {
if qa.UserID == userID {
totalAttempts += qa.CorrectAttempts + qa.IncorrectAttempts
correctAttempts += qa.CorrectAttempts
}
}
if totalAttempts > 0 {
progress.QuestionAccuracy = float64(correctAttempts) / float64(totalAttempts) * 100
}

// Calculate topic breakdown
topicAttempts := make(map[int]int)      // topicID -> total attempts
topicCorrect := make(map[int]int)       // topicID -> correct attempts

for _, qa := range db.questionAttempts {
if qa.UserID == userID {
if question, exists := db.questions[qa.QuestionID]; exists {
topicAttempts[question.TopicID] += qa.CorrectAttempts + qa.IncorrectAttempts
topicCorrect[question.TopicID] += qa.CorrectAttempts
}
}
}

for topicID, total := range topicAttempts {
if topic, exists := db.topics[topicID]; exists && total > 0 {
accuracy := float64(topicCorrect[topicID]) / float64(total) * 100
progress.TopicBreakdown[topic.Name] = accuracy
}
}

// Calculate FSRS metrics
var totalCards, dueCards, newCards, reviewCards int
var totalStability, totalDifficulty float64
var cardCount int

now := time.Now()
for _, state := range db.userQuestionStates {
if state.UserID == userID {
totalCards++
if state.State == 0 {
newCards++
} else if state.State == 2 {
reviewCards++
}
if !state.Due.After(now) {
dueCards++
}
totalStability += state.Stability
totalDifficulty += state.Difficulty
cardCount++
}
}

progress.FSRSMetrics = models.FSRSMetrics{
TotalCards:   totalCards,
DueCards:     dueCards,
NewCards:     newCards,
ReviewCards:  reviewCards,
}
if cardCount > 0 {
progress.FSRSMetrics.AvgStability = totalStability / float64(cardCount)
progress.FSRSMetrics.AvgDifficulty = totalDifficulty / float64(cardCount)
}

return progress, nil
}

func (db *InMemoryDB) GetCohortAnalytics() (*models.CohortAnalytics, error) {
db.mu.RLock()
defer db.mu.RUnlock()

analytics := &models.CohortAnalytics{
TopicDifficulty: make(map[string]float64),
QuestionStats:   make([]models.QuestionStatistic, 0),
}

// Count students (non-admin users)
studentIDs := make([]int, 0)
for _, user := range db.users {
if !user.IsAdmin {
analytics.TotalStudents++
studentIDs = append(studentIDs, user.ID)
}
}

// Count active students (activity in last 7 days)
sevenDaysAgo := time.Now().AddDate(0, 0, -7)
activeStudents := make(map[int]bool)
for _, attempt := range db.quizAttempts {
if attempt.StartTime.After(sevenDaysAgo) {
activeStudents[attempt.UserID] = true
}
}
analytics.ActiveStudents = len(activeStudents)

// Calculate average cohort score and total quizzes
var totalScore float64
var scoredQuizzes int

for _, attempt := range db.quizAttempts {
if attempt.CompletedTime != nil {
analytics.TotalQuizzesTaken++
if attempt.TotalQuestions > 0 {
score := float64(attempt.CorrectAnswers) / float64(attempt.TotalQuestions) * 100
totalScore += score
scoredQuizzes++
}
}
}

if scoredQuizzes > 0 {
analytics.AverageCohortScore = totalScore / float64(scoredQuizzes)
}

// Get top and bottom performers
studentScores := make(map[int]float64)
studentQuizCounts := make(map[int]int)

for _, attempt := range db.quizAttempts {
if attempt.CompletedTime != nil && attempt.TotalQuestions > 0 {
score := float64(attempt.CorrectAnswers) / float64(attempt.TotalQuestions) * 100
studentScores[attempt.UserID] += score
studentQuizCounts[attempt.UserID]++
}
}

// Calculate average scores for each student
type studentScore struct {
userID   int
avgScore float64
}
scores := make([]studentScore, 0)
for userID, totalScore := range studentScores {
if count := studentQuizCounts[userID]; count > 0 {
scores = append(scores, studentScore{userID, totalScore / float64(count)})
}
}

// Sort by score (descending)
for i := 0; i < len(scores); i++ {
for j := i + 1; j < len(scores); j++ {
if scores[i].avgScore < scores[j].avgScore {
scores[i], scores[j] = scores[j], scores[i]
}
}
}

// Get top 5 performers
topCount := 5
if len(scores) < topCount {
topCount = len(scores)
}
topPerformers := make([]models.StudentProgress, 0, topCount)
for i := 0; i < topCount; i++ {
if progress, err := db.GetStudentProgress(scores[i].userID); err == nil {
topPerformers = append(topPerformers, *progress)
}
}
analytics.TopPerformers = topPerformers

// Get bottom 5 performers
bottomCount := 5
if len(scores) < bottomCount {
bottomCount = len(scores)
}
bottomPerformers := make([]models.StudentProgress, 0, bottomCount)
for i := len(scores) - 1; i >= len(scores)-bottomCount && i >= 0; i-- {
if progress, err := db.GetStudentProgress(scores[i].userID); err == nil {
bottomPerformers = append(bottomPerformers, *progress)
}
}
analytics.BottomPerformers = bottomPerformers

// Calculate topic difficulty (avg FSRS difficulty per topic)
topicDifficulties := make(map[int][]float64) // topicID -> []difficulty values

for _, state := range db.userQuestionStates {
if question, exists := db.questions[state.QuestionID]; exists {
topicDifficulties[question.TopicID] = append(topicDifficulties[question.TopicID], state.Difficulty)
}
}

for topicID, difficulties := range topicDifficulties {
if topic, exists := db.topics[topicID]; exists && len(difficulties) > 0 {
var sum float64
for _, d := range difficulties {
sum += d
}
analytics.TopicDifficulty[topic.Name] = sum / float64(len(difficulties))
}
}

// Calculate question statistics
questionAttemptCounts := make(map[int]int)
questionCorrectCounts := make(map[int]int)

for _, qa := range db.questionAttempts {
totalAttempts := qa.CorrectAttempts + qa.IncorrectAttempts
questionAttemptCounts[qa.QuestionID] += totalAttempts
questionCorrectCounts[qa.QuestionID] += qa.CorrectAttempts
}

for questionID, totalAttempts := range questionAttemptCounts {
if question, exists := db.questions[questionID]; exists {
var accuracy float64
if totalAttempts > 0 {
accuracy = float64(questionCorrectCounts[questionID]) / float64(totalAttempts) * 100
}

// Calculate average FSRS difficulty for this question
var totalDifficulty float64
var difficultyCount int
for _, state := range db.userQuestionStates {
if state.QuestionID == questionID {
totalDifficulty += state.Difficulty
difficultyCount++
}
}
var avgDifficulty float64
if difficultyCount > 0 {
avgDifficulty = totalDifficulty / float64(difficultyCount)
}

stat := models.QuestionStatistic{
QuestionID:      questionID,
TopicQuestionID: question.TopicQuestionID,
QuestionText:    question.QuestionText,
TotalAttempts:   totalAttempts,
CorrectAttempts: questionCorrectCounts[questionID],
Accuracy:        accuracy,
AvgDifficulty:   avgDifficulty,
}
analytics.QuestionStats = append(analytics.QuestionStats, stat)
}
}

return analytics, nil
}
