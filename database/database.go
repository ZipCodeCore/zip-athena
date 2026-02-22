package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

var (
	ErrNotFound      = errors.New("record not found")
	ErrAlreadyExists = errors.New("record already exists")
)

// UserRepository handles user-related database operations
type UserRepository interface {
	CreateUser(user *models.User) error
	GetUser(id int) (*models.User, error)
	UpdateUserPoints(id int, points int) error
	ListUsers() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id int) error
}

// TopicRepository handles topic-related database operations
type TopicRepository interface {
	CreateTopic(topic *models.Topic) error
	GetTopic(id int) (*models.Topic, error)
	ListTopics() ([]*models.Topic, error)
	UpdateTopic(topic *models.Topic) error
	DeleteTopic(id int) error
	ReorderTopics(topicIDs []int) error
}

// QuestionRepository handles question-related database operations
type QuestionRepository interface {
	CreateQuestion(question *models.Question) error
	GetQuestion(id int) (*models.Question, error)
	GetQuestionsByTopic(topicID int) ([]*models.Question, error)
	ListQuestions() ([]*models.Question, error)
	UpdateQuestion(question *models.Question) error
	DeleteQuestion(id int) error
	ListQuestionsWithPagination(page int, perPage int, search string, topicID int, difficulty int) ([]*models.Question, int, error)
	DeactivateQuestion(questionID int) error
	GetQuestionPreview(questionID int) (*models.QuestionPreview, error)
	ValidateQuestionData(question *models.Question, answers []models.Answer) error
}

// AnswerRepository handles answer-related database operations
type AnswerRepository interface {
	CreateAnswer(answer *models.Answer) error
	GetAnswer(id int) (*models.Answer, error)
	GetAnswersByQuestion(questionID int) ([]*models.Answer, error)
	UpdateAnswer(answer *models.Answer) error
	DeleteAnswer(id int) error
}

// QuizRepository handles quiz and quiz attempt operations
type QuizRepository interface {
	// Quiz CRUD
	CreateQuiz(quiz *models.Quiz) error
	GetQuiz(id int) (*models.Quiz, error)
	ListQuizzes() ([]*models.Quiz, error)
	UpdateQuiz(quiz *models.Quiz) error
	DeleteQuiz(id int) error

	// Quiz-Question relationships
	AddQuestionToQuiz(quizID, questionID, order int) error
	RemoveQuestionFromQuiz(quizID, questionID int) error
	GetQuizQuestions(quizID int) ([]*models.Question, error)

	// Quiz Attempts
	CreateQuizAttempt(attempt *models.QuizAttempt) error
	GetQuizAttempt(id int) (*models.QuizAttempt, error)
	GetQuizAttemptsByUser(userID int) ([]*models.QuizAttempt, error)
	GetAllQuizAttempts() ([]*models.QuizAttempt, error)
	UpdateQuizAttempt(attempt *models.QuizAttempt) error

	// Quiz History
	GetQuizHistoryByUser(userID int, page int, perPage int) ([]*models.QuizAttempt, error)
	GetQuizHistoryByTopic(userID int, topicID int, page int, perPage int) ([]*models.QuizAttempt, error)
	GetQuizAttemptWithDetails(attemptID int) (*models.QuizAttempt, error)
}

// QuestionAttemptRepository handles question attempt operations
type QuestionAttemptRepository interface {
	CreateQuestionAttempt(attempt *models.QuestionAttempt) error
	GetQuestionAttempt(id int) (*models.QuestionAttempt, error)
	GetQuestionAttemptsByUser(userID int) ([]*models.QuestionAttempt, error)
	UpdateQuestionAttempt(attempt *models.QuestionAttempt) error
}

// TestRepository handles test and test attempt operations
type TestRepository interface {
	// Test CRUD
	CreateTest(test *models.Test) error
	GetTest(id int) (*models.Test, error)
	ListTests() ([]*models.Test, error)
	UpdateTest(test *models.Test) error
	DeleteTest(id int) error

	// Test-Question relationships
	AddQuestionToTest(testID, questionID, order int) error
	RemoveQuestionFromTest(testID, questionID int) error
	GetTestQuestions(testID int) ([]*models.Question, error)

	// Test Attempts
	CreateTestAttempt(attempt *models.TestAttempt) error
	GetTestAttempt(id int) (*models.TestAttempt, error)
	UpdateTestAttemptAnswer(attemptID int, questionID string, answerID int) error
	SubmitTestAttempt(attemptID int, isAutoSubmit bool) error
	GetTestResults(attemptID int) (*models.TestAttempt, error)
	GetAvailableTests(userID int) ([]*models.Test, error)
	GetTestAttemptsByUser(userID int, testID int) ([]*models.TestAttempt, error)

	// Test Management
	PublishTest(testID int) error
	GetTestResultsByTest(testID int) (*models.TestResult, error)
	ExportTestResults(testID int) ([]byte, error)
}

// FSRSRepository handles FSRS (spaced repetition) state operations
type FSRSRepository interface {
	CreateUserQuestionState(state *models.UserQuestionState) error
	GetUserQuestionState(userID, questionID int) (*models.UserQuestionState, error)
	GetUserQuestionStatesByUser(userID int) ([]*models.UserQuestionState, error)
	GetDueUserQuestionStates(userID int) ([]*models.UserQuestionState, error)
	UpdateUserQuestionState(state *models.UserQuestionState) error
}

// AnalyticsRepository handles analytics and reporting operations
type AnalyticsRepository interface {
	GetStudentProgress(userID int) (*models.StudentProgress, error)
	GetCohortAnalytics() (*models.CohortAnalytics, error)
	ListStudentsWithMetrics(searchQuery string, sortBy string, filterCohort string) ([]*models.StudentProgress, error)
	GetStudentQuizHistory(userID int) ([]*models.QuizAttempt, error)
	GetStudentWeakAreas(userID int) (*models.WeakAreas, error)
	GetQuestionAnalytics() ([]*models.QuestionAnalytics, error)
	GetQuestionDetailedAnalytics(questionID int) (*models.QuestionAnalytics, error)
	GetAnswerDistribution(questionID int) (map[int]int, error)
	UpdateQuestionDifficultyFromAnalytics(questionID int, newDifficulty int) error
}

// TagRepository handles tag-related operations
type TagRepository interface {
	GetAllTags() ([]string, error)
	GetTagsWithCounts() (map[string]int, error)
	GetTagStats(userID int) ([]*models.TagStats, error)
	GetQuestionsByTags(tags []string, useAnd bool) ([]*models.Question, error)
	GetTagsWithMetadata() ([]*models.Tag, error)
	MergeTags(sourceTags []string, targetTag string) error
	DeleteTag(tagName string) error
	GetQuestionsByTag(tagName string) ([]*models.Question, error)
}

// QualityRepository handles question quality and review operations
type QualityRepository interface {
	GetFlaggedQuestions() ([]*models.FlaggedQuestion, error)
	FlagQuestion(questionID int, reason string) error
	UnflagQuestion(questionID int) error
	AutoFlagQuestions() error
}

// AssignedQuizRepository handles quiz assignment operations
type AssignedQuizRepository interface {
	AssignQuizToUser(quizID, userID, assignedBy int) error
	UnassignQuizFromUser(quizID, userID int) error
	GetAssignedQuizzesForUser(userID int) ([]*models.AssignedQuizInfo, error)
	GetAssignedQuizUserIDs(quizID int) ([]int, error)
	GetAssignedQuizStatus(quizID int) (*models.AssignedQuizStatus, error)
	ListAssignedQuizzes() ([]*models.AssignedQuiz, error)
}

// BulkImportRepository handles bulk import operations
type BulkImportRepository interface {
	ValidateCSVImport(rows []models.CSVImportRow) (*models.CSVImportResult, error)
	BulkCreateQuestions(rows []models.CSVImportRow) (int, error)
}

// Database interface composes all repository interfaces
type Database interface {
	// Core database access
	DB() *sql.DB
	Close() error

	// Composed repositories
	UserRepository
	TopicRepository
	QuestionRepository
	AnswerRepository
	QuizRepository
	QuestionAttemptRepository
	TestRepository
	FSRSRepository
	AnalyticsRepository
	TagRepository
	QualityRepository
	BulkImportRepository
	AssignedQuizRepository
}

// InMemoryDB implements the Database interface using in-memory storage
type InMemoryDB struct {
	users              map[int]*models.User
	topics             map[int]*models.Topic
	questions          map[int]*models.Question
	answers            map[int]*models.Answer
	quizAttempts       map[int]*models.QuizAttempt
	questionAttempts   map[int]*models.QuestionAttempt
	quizzes            map[int]*models.Quiz
	userQuestionStates map[string]*models.UserQuestionState // key: "userID-questionID"
	tests              map[int]*models.Test                 // Phase 3
	quizQuestions      map[string]*models.QuizQuestion      // key: "quizID-questionID" Phase 3
	testQuestions      map[string]*models.TestQuestion      // key: "testID-questionID" Phase 3
	testAttempts       map[int]*models.TestAttempt          // Phase 4
	flaggedQuestions   map[int]*models.FlaggedQuestion      // Phase 5
	assignedQuizzes    map[int]*models.AssignedQuiz         // key: id
	nextUserID         int
	nextTopicID        int
	nextQuestionID     int
	nextAnswerID       int
	nextQuizAttemptID  int
	nextQAttemptID     int
	nextQuizID         int
	nextUQStateID      int
	nextTestID         int // Phase 3
	nextQuizQuestionID int // Phase 3
	nextTestQuestionID int // Phase 3
	nextTestAttemptID  int // Phase 4
	nextAssignedQuizID int
	mu                 sync.RWMutex
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
		testAttempts:       make(map[int]*models.TestAttempt),
		flaggedQuestions:   make(map[int]*models.FlaggedQuestion),
		assignedQuizzes:    make(map[int]*models.AssignedQuiz),
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
		nextAssignedQuizID: 1,
	}
}

// DB returns nil for in-memory database (no underlying SQL connection)
func (db *InMemoryDB) DB() *sql.DB {
	return nil
}

// Close closes the database connection (no-op for in-memory)
func (db *InMemoryDB) Close() error {
	return nil
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

	// Get max display_order and increment
	if topic.DisplayOrder == 0 {
		maxOrder := 0
		for _, t := range db.topics {
			if t.DisplayOrder > maxOrder {
				maxOrder = t.DisplayOrder
			}
		}
		topic.DisplayOrder = maxOrder + 1
	}

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
		// Only include active topics
		if topic.IsActive {
			topics = append(topics, topic)
		}
	}

	// Sort by display_order
	for i := 0; i < len(topics); i++ {
		for j := i + 1; j < len(topics); j++ {
			if topics[i].DisplayOrder > topics[j].DisplayOrder {
				topics[i], topics[j] = topics[j], topics[i]
			}
		}
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

func (db *InMemoryDB) GetAllQuizAttempts() ([]*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempts := make([]*models.QuizAttempt, 0, len(db.quizAttempts))
	for _, attempt := range db.quizAttempts {
		attempts = append(attempts, attempt)
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

	topic, exists := db.topics[id]
	if !exists {
		return ErrNotFound
	}
	// Soft delete: set is_active to false
	topic.IsActive = false
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

// getTestQuestionsUnlocked is an internal helper that assumes the caller already holds a lock
func (db *InMemoryDB) getTestQuestionsUnlocked(testID int) ([]*models.Question, error) {
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

func (db *InMemoryDB) GetTestQuestions(testID int) ([]*models.Question, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	return db.getTestQuestionsUnlocked(testID)
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
		UserID:         userID,
		Username:       user.Username,
		TotalPoints:    user.TotalPoints,
		TopicBreakdown: make(map[string]float64),
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
	topicAttempts := make(map[int]int) // topicID -> total attempts
	topicCorrect := make(map[int]int)  // topicID -> correct attempts

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
		TotalCards:  totalCards,
		DueCards:    dueCards,
		NewCards:    newCards,
		ReviewCards: reviewCards,
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

// ListStudentsWithMetrics returns a list of all students with their performance metrics
func (db *InMemoryDB) ListStudentsWithMetrics(searchQuery string, sortBy string, filterCohort string) ([]*models.StudentProgress, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	students := make([]*models.StudentProgress, 0)

	for _, user := range db.users {
		// Skip admin users
		if user.IsAdmin {
			continue
		}

		// Apply search filter if provided
		if searchQuery != "" {
			// Simple case-insensitive search in username and email
			queryLower := strings.ToLower(searchQuery)
			usernameLower := strings.ToLower(user.Username)
			emailLower := strings.ToLower(user.Email)

			if !strings.Contains(usernameLower, queryLower) && !strings.Contains(emailLower, queryLower) {
				continue
			}
		}

		// Get student progress
		progress, err := db.GetStudentProgress(user.ID)
		if err != nil {
			continue
		}

		students = append(students, progress)
	}

	// Note: Sorting by sortBy parameter would be implemented here
	// For simplicity, returning unsorted list

	return students, nil
}

// GetStudentQuizHistory returns all quiz attempts for a student
func (db *InMemoryDB) GetStudentQuizHistory(userID int) ([]*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Verify user exists
	if _, exists := db.users[userID]; !exists {
		return nil, ErrNotFound
	}

	attempts := make([]*models.QuizAttempt, 0)
	for _, attempt := range db.quizAttempts {
		if attempt.UserID == userID {
			attempts = append(attempts, attempt)
		}
	}

	// Sort by start time descending (most recent first)
	// Simple bubble sort for small datasets
	for i := 0; i < len(attempts); i++ {
		for j := i + 1; j < len(attempts); j++ {
			if attempts[i].StartTime.Before(attempts[j].StartTime) {
				attempts[i], attempts[j] = attempts[j], attempts[i]
			}
		}
	}

	return attempts, nil
}

// GetStudentWeakAreas identifies topics and tags where student is struggling (< 60% accuracy)
func (db *InMemoryDB) GetStudentWeakAreas(userID int) (*models.WeakAreas, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	user, exists := db.users[userID]
	if !exists {
		return nil, ErrNotFound
	}

	weakAreas := &models.WeakAreas{
		UserID:     userID,
		Username:   user.Username,
		WeakTopics: make([]models.TopicWeakness, 0),
		WeakTags:   make([]models.TagWeakness, 0),
	}

	// Calculate topic accuracy
	topicStats := make(map[int]struct {
		correct int
		total   int
	})

	for _, qa := range db.questionAttempts {
		if qa.UserID == userID {
			question, exists := db.questions[qa.QuestionID]
			if !exists {
				continue
			}

			stats := topicStats[question.TopicID]
			stats.correct += qa.CorrectAttempts
			stats.total += qa.CorrectAttempts + qa.IncorrectAttempts
			topicStats[question.TopicID] = stats
		}
	}

	// Identify weak topics (< 60% accuracy)
	for topicID, stats := range topicStats {
		if stats.total > 0 {
			accuracy := float64(stats.correct) / float64(stats.total) * 100
			if accuracy < 60 {
				topic, exists := db.topics[topicID]
				if exists {
					weakAreas.WeakTopics = append(weakAreas.WeakTopics, models.TopicWeakness{
						TopicID:       topicID,
						TopicName:     topic.Name,
						Accuracy:      accuracy,
						TotalAttempts: stats.total,
					})
				}
			}
		}
	}

	// Calculate tag accuracy
	tagStats := make(map[string]struct {
		correct int
		total   int
	})

	for _, qa := range db.questionAttempts {
		if qa.UserID == userID {
			question, exists := db.questions[qa.QuestionID]
			if !exists {
				continue
			}

			for _, tag := range question.Tags {
				stats := tagStats[tag]
				stats.correct += qa.CorrectAttempts
				stats.total += qa.CorrectAttempts + qa.IncorrectAttempts
				tagStats[tag] = stats
			}
		}
	}

	// Identify weak tags (< 60% accuracy)
	for tag, stats := range tagStats {
		if stats.total > 0 {
			accuracy := float64(stats.correct) / float64(stats.total) * 100
			if accuracy < 60 {
				weakAreas.WeakTags = append(weakAreas.WeakTags, models.TagWeakness{
					Tag:           tag,
					Accuracy:      accuracy,
					TotalAttempts: stats.total,
				})
			}
		}
	}

	return weakAreas, nil
}

// ReorderTopics updates the display_order for topics
func (db *InMemoryDB) ReorderTopics(topicIDs []int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// First validate all topics exist
	for _, topicID := range topicIDs {
		if _, exists := db.topics[topicID]; !exists {
			return ErrNotFound
		}
	}

	// Then update display orders atomically
	for i, topicID := range topicIDs {
		db.topics[topicID].DisplayOrder = i
	}
	return nil
}

// Phase 4: Quiz History Methods

// GetQuizHistoryByUser retrieves quiz attempts for a user with pagination
func (db *InMemoryDB) GetQuizHistoryByUser(userID int, page int, perPage int) ([]*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var attempts []*models.QuizAttempt
	for _, attempt := range db.quizAttempts {
		if attempt.UserID == userID && attempt.CompletedTime != nil {
			attempts = append(attempts, attempt)
		}
	}

	// Sort by completed time descending (newest first)
	for i := 0; i < len(attempts); i++ {
		for j := i + 1; j < len(attempts); j++ {
			if attempts[i].CompletedTime.Before(*attempts[j].CompletedTime) {
				attempts[i], attempts[j] = attempts[j], attempts[i]
			}
		}
	}

	// Apply pagination
	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(attempts) {
		return []*models.QuizAttempt{}, nil
	}
	if end > len(attempts) {
		end = len(attempts)
	}

	return attempts[start:end], nil
}

// GetQuizHistoryByTopic retrieves quiz attempts for a user filtered by topic
func (db *InMemoryDB) GetQuizHistoryByTopic(userID int, topicID int, page int, perPage int) ([]*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var attempts []*models.QuizAttempt
	for _, attempt := range db.quizAttempts {
		if attempt.UserID == userID && attempt.CompletedTime != nil {
			// Check if any question in the attempt belongs to the topic
			hasTopicQuestion := false
			for qid := range attempt.QuestionResults {
				// Find question by topic_question_id
				for _, q := range db.questions {
					if q.TopicQuestionID == qid && q.TopicID == topicID {
						hasTopicQuestion = true
						break
					}
				}
				if hasTopicQuestion {
					break
				}
			}
			if hasTopicQuestion {
				attempts = append(attempts, attempt)
			}
		}
	}

	// Sort by completed time descending
	for i := 0; i < len(attempts); i++ {
		for j := i + 1; j < len(attempts); j++ {
			if attempts[i].CompletedTime.Before(*attempts[j].CompletedTime) {
				attempts[i], attempts[j] = attempts[j], attempts[i]
			}
		}
	}

	// Apply pagination
	start := (page - 1) * perPage
	end := start + perPage
	if start >= len(attempts) {
		return []*models.QuizAttempt{}, nil
	}
	if end > len(attempts) {
		end = len(attempts)
	}

	return attempts[start:end], nil
}

// GetQuizAttemptWithDetails retrieves a quiz attempt with full details
func (db *InMemoryDB) GetQuizAttemptWithDetails(attemptID int) (*models.QuizAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempt, exists := db.quizAttempts[attemptID]
	if !exists {
		return nil, ErrNotFound
	}

	return attempt, nil
}

// Phase 4: Timed Test Taking Methods

// CreateTestAttempt creates a new test attempt
func (db *InMemoryDB) CreateTestAttempt(attempt *models.TestAttempt) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.nextTestAttemptID++
	attempt.ID = db.nextTestAttemptID
	db.testAttempts[attempt.ID] = attempt
	return nil
}

// GetTestAttempt retrieves a test attempt by ID
func (db *InMemoryDB) GetTestAttempt(id int) (*models.TestAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	attempt, exists := db.testAttempts[id]
	if !exists {
		return nil, ErrNotFound
	}
	return attempt, nil
}

// UpdateTestAttemptAnswer updates an answer in a test attempt
func (db *InMemoryDB) UpdateTestAttemptAnswer(attemptID int, questionID string, answerID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	attempt, exists := db.testAttempts[attemptID]
	if !exists {
		return ErrNotFound
	}

	if attempt.Answers == nil {
		attempt.Answers = make(map[string]int)
	}
	attempt.Answers[questionID] = answerID
	return nil
}

// SubmitTestAttempt submits a test attempt and calculates the score
func (db *InMemoryDB) SubmitTestAttempt(attemptID int, isAutoSubmit bool) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	attempt, exists := db.testAttempts[attemptID]
	if !exists {
		return ErrNotFound
	}

	now := time.Now()
	attempt.SubmitTime = &now
	attempt.IsAutoSubmit = isAutoSubmit

	// Calculate time remaining
	test, exists := db.tests[attempt.TestID]
	if !exists {
		return ErrNotFound
	}

	timeLimitSeconds := test.TimeLimit * 60
	elapsed := int(now.Sub(attempt.StartTime).Seconds())
	attempt.TimeRemaining = timeLimitSeconds - elapsed
	if attempt.TimeRemaining < 0 {
		attempt.TimeRemaining = 0
	}

	// Get test questions to calculate score (use unlocked version since we already hold the lock)
	testQuestions, err := db.getTestQuestionsUnlocked(attempt.TestID)
	if err != nil {
		return err
	}

	attempt.TotalQuestions = len(testQuestions)
	attempt.CorrectAnswers = 0
	attempt.QuestionResults = make(map[string]bool)

	// Check each answer
	for _, q := range testQuestions {
		questionKey := fmt.Sprintf("%d", q.ID)
		answerID, hasAnswer := attempt.Answers[questionKey]

		if hasAnswer {
			// Check if answer is correct
			answer, exists := db.answers[answerID]
			if exists && answer.QuestionID == q.ID && answer.IsCorrect {
				attempt.CorrectAnswers++
				attempt.QuestionResults[questionKey] = true
			} else {
				attempt.QuestionResults[questionKey] = false
			}
		} else {
			attempt.QuestionResults[questionKey] = false
		}
	}

	// Calculate score
	if attempt.TotalQuestions > 0 {
		attempt.Score = float64(attempt.CorrectAnswers) / float64(attempt.TotalQuestions) * 100
	}

	return nil
}

// GetTestResults retrieves the results of a test attempt
func (db *InMemoryDB) GetTestResults(attemptID int) (*models.TestAttempt, error) {
	return db.GetTestAttempt(attemptID)
}

// GetAvailableTests retrieves available tests for a student
func (db *InMemoryDB) GetAvailableTests(userID int) ([]*models.Test, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	now := time.Now()
	var available []*models.Test

	for _, test := range db.tests {
		// Only show published tests
		if !test.IsPublished {
			continue
		}

		// Check if within availability window
		if test.StartDate != nil && now.Before(*test.StartDate) {
			continue
		}
		if test.EndDate != nil && now.After(*test.EndDate) {
			continue
		}

		available = append(available, test)
	}

	return available, nil
}

// GetTestAttemptsByUser retrieves all test attempts for a user for a specific test
func (db *InMemoryDB) GetTestAttemptsByUser(userID int, testID int) ([]*models.TestAttempt, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var attempts []*models.TestAttempt
	for _, attempt := range db.testAttempts {
		if attempt.UserID == userID && attempt.TestID == testID {
			attempts = append(attempts, attempt)
		}
	}

	return attempts, nil
}

// Phase 4: Tag-Based Review Methods

// GetAllTags retrieves all unique tags from all questions
func (db *InMemoryDB) GetAllTags() ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tagSet := make(map[string]bool)
	for _, q := range db.questions {
		for _, tag := range q.Tags {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetTagsWithCounts retrieves all tags with their question counts
func (db *InMemoryDB) GetTagsWithCounts() (map[string]int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tagCounts := make(map[string]int)
	for _, q := range db.questions {
		for _, tag := range q.Tags {
			tagCounts[tag]++
		}
	}

	return tagCounts, nil
}

// GetTagStats retrieves accuracy statistics for each tag for a user
func (db *InMemoryDB) GetTagStats(userID int) ([]*models.TagStats, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tagStats := make(map[string]*models.TagStats)

	// Calculate stats from question attempts
	for _, qa := range db.questionAttempts {
		if qa.UserID == userID {
			question, exists := db.questions[qa.QuestionID]
			if !exists {
				continue
			}

			for _, tag := range question.Tags {
				if _, exists := tagStats[tag]; !exists {
					tagStats[tag] = &models.TagStats{
						Tag: tag,
					}
				}

				stats := tagStats[tag]
				stats.CorrectAttempts += qa.CorrectAttempts
				stats.TotalAttempts += qa.CorrectAttempts + qa.IncorrectAttempts
			}
		}
	}

	// Calculate accuracy and question counts
	tagCounts, _ := db.GetTagsWithCounts()
	result := make([]*models.TagStats, 0, len(tagStats))
	for tag, stats := range tagStats {
		stats.QuestionCount = tagCounts[tag]
		if stats.TotalAttempts > 0 {
			stats.Accuracy = float64(stats.CorrectAttempts) / float64(stats.TotalAttempts) * 100
		}
		result = append(result, stats)
	}

	return result, nil
}

// GetQuestionsByTags retrieves questions that match the given tags
func (db *InMemoryDB) GetQuestionsByTags(tags []string, useAnd bool) ([]*models.Question, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var result []*models.Question

	for _, q := range db.questions {
		if useAnd {
			// AND logic: question must have all tags
			hasAllTags := true
			for _, tag := range tags {
				found := false
				for _, qTag := range q.Tags {
					if qTag == tag {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if hasAllTags {
				result = append(result, q)
			}
		} else {
			// OR logic: question must have at least one tag
			hasAnyTag := false
			for _, tag := range tags {
				for _, qTag := range q.Tags {
					if qTag == tag {
						hasAnyTag = true
						break
					}
				}
				if hasAnyTag {
					break
				}
			}
			if hasAnyTag {
				result = append(result, q)
			}
		}
	}

	return result, nil
}

// Phase 4: Test Management (Instructor) Methods

// PublishTest sets a test as published
func (db *InMemoryDB) PublishTest(testID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	test, exists := db.tests[testID]
	if !exists {
		return ErrNotFound
	}

	test.IsPublished = true
	test.UpdatedAt = time.Now()
	return nil
}

// GetTestResultsByTest retrieves all test results for a specific test
func (db *InMemoryDB) GetTestResultsByTest(testID int) (*models.TestResult, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	test, exists := db.tests[testID]
	if !exists {
		return nil, ErrNotFound
	}

	result := &models.TestResult{
		TestID:         testID,
		TestName:       test.Name,
		StudentResults: make([]models.StudentTestResult, 0),
	}

	for _, attempt := range db.testAttempts {
		if attempt.TestID == testID && attempt.SubmitTime != nil {
			user, exists := db.users[attempt.UserID]
			if !exists {
				continue
			}

			timeTaken := 0
			if test.TimeLimit > 0 {
				timeTaken = (test.TimeLimit * 60) - attempt.TimeRemaining
			}

			result.StudentResults = append(result.StudentResults, models.StudentTestResult{
				UserID:         attempt.UserID,
				Username:       user.Username,
				Score:          attempt.Score,
				CorrectAnswers: attempt.CorrectAnswers,
				TotalQuestions: attempt.TotalQuestions,
				SubmitTime:     *attempt.SubmitTime,
				TimeTaken:      timeTaken,
			})
		}
	}

	return result, nil
}

// ExportTestResults exports test results as CSV
func (db *InMemoryDB) ExportTestResults(testID int) ([]byte, error) {
	results, err := db.GetTestResultsByTest(testID)
	if err != nil {
		return nil, err
	}

	var csv strings.Builder
	csv.WriteString("Username,Score,Correct Answers,Total Questions,Submit Time,Time Taken (seconds)\n")

	for _, sr := range results.StudentResults {
		csv.WriteString(fmt.Sprintf("%s,%.2f,%d,%d,%s,%d\n",
			sr.Username,
			sr.Score,
			sr.CorrectAnswers,
			sr.TotalQuestions,
			sr.SubmitTime.Format(time.RFC3339),
			sr.TimeTaken,
		))
	}

	return []byte(csv.String()), nil
}

// Phase 4: Question Performance Analytics Methods

// GetQuestionAnalytics retrieves analytics for all questions
func (db *InMemoryDB) GetQuestionAnalytics() ([]*models.QuestionAnalytics, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	analytics := make([]*models.QuestionAnalytics, 0)

	for _, q := range db.questions {
		qa, err := db.getQuestionAnalyticsInternal(q.ID)
		if err != nil {
			continue
		}
		analytics = append(analytics, qa)
	}

	return analytics, nil
}

// GetQuestionDetailedAnalytics retrieves detailed analytics for a single question
func (db *InMemoryDB) GetQuestionDetailedAnalytics(questionID int) (*models.QuestionAnalytics, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getQuestionAnalyticsInternal(questionID)
}

// getQuestionAnalyticsInternal is an internal helper (must be called with lock held)
func (db *InMemoryDB) getQuestionAnalyticsInternal(questionID int) (*models.QuestionAnalytics, error) {
	question, exists := db.questions[questionID]
	if !exists {
		return nil, ErrNotFound
	}

	topic, _ := db.topics[question.TopicID]
	topicName := ""
	if topic != nil {
		topicName = topic.Name
	}

	analytics := &models.QuestionAnalytics{
		QuestionID:         questionID,
		TopicQuestionID:    question.TopicQuestionID,
		QuestionText:       question.QuestionText,
		TopicID:            question.TopicID,
		TopicName:          topicName,
		Tags:               question.Tags,
		CurrentDifficulty:  question.Difficulty,
		AnswerDistribution: make(map[int]int),
	}

	// Calculate attempts and accuracy
	totalAttempts := 0
	correctAttempts := 0

	for _, qa := range db.questionAttempts {
		if qa.QuestionID == questionID {
			totalAttempts += qa.CorrectAttempts + qa.IncorrectAttempts
			correctAttempts += qa.CorrectAttempts
		}
	}

	analytics.TotalAttempts = totalAttempts
	analytics.CorrectAttempts = correctAttempts

	if totalAttempts > 0 {
		analytics.Accuracy = float64(correctAttempts) / float64(totalAttempts) * 100
	}

	// Calculate recommended difficulty based on accuracy
	// Only flag for review if there are enough attempts (at least 10)
	if totalAttempts >= 10 {
		if analytics.Accuracy > 90 {
			analytics.RecommendedDifficulty = question.Difficulty + 1
			if analytics.RecommendedDifficulty > 7 {
				analytics.RecommendedDifficulty = 7
			}
			analytics.FlagForReview = true
		} else if analytics.Accuracy < 30 {
			analytics.RecommendedDifficulty = question.Difficulty - 1
			if analytics.RecommendedDifficulty < 1 {
				analytics.RecommendedDifficulty = 1
			}
			analytics.FlagForReview = true
		} else {
			analytics.RecommendedDifficulty = question.Difficulty
		}
	} else {
		// Not enough data to make a recommendation
		analytics.RecommendedDifficulty = question.Difficulty
	}

	// Get answer distribution from quiz/test attempts
	answerDistribution, _ := db.getAnswerDistributionInternal(questionID)
	analytics.AnswerDistribution = answerDistribution

	return analytics, nil
}

// GetAnswerDistribution retrieves distribution of answers chosen for a question
func (db *InMemoryDB) GetAnswerDistribution(questionID int) (map[int]int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.getAnswerDistributionInternal(questionID)
}

// getAnswerDistributionInternal is an internal helper (must be called with lock held)
func (db *InMemoryDB) getAnswerDistributionInternal(questionID int) (map[int]int, error) {
	distribution := make(map[int]int)

	// This would require tracking individual answer selections
	// For now, we'll return empty distribution as we don't track this in quiz attempts
	// In a real implementation, you'd need to track answer_id in QuizAttempt

	return distribution, nil
}

// UpdateQuestionDifficultyFromAnalytics updates a question's difficulty
func (db *InMemoryDB) UpdateQuestionDifficultyFromAnalytics(questionID int, newDifficulty int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	question, exists := db.questions[questionID]
	if !exists {
		return ErrNotFound
	}

	if newDifficulty < 1 || newDifficulty > 7 {
		return fmt.Errorf("difficulty must be between 1 and 7")
	}

	question.Difficulty = newDifficulty
	return nil
}

// Phase 5: Question Management - InMemoryDB implementations

func (db *InMemoryDB) ListQuestionsWithPagination(page int, perPage int, search string, topicID int, difficulty int) ([]*models.Question, int, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var filtered []*models.Question
	for _, q := range db.questions {
		// Apply filters
		if search != "" && !strings.Contains(strings.ToLower(q.QuestionText), strings.ToLower(search)) {
			continue
		}
		if topicID > 0 && q.TopicID != topicID {
			continue
		}
		if difficulty > 0 && q.Difficulty != difficulty {
			continue
		}
		filtered = append(filtered, q)
	}

	total := len(filtered)

	// Apply pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 30
	}
	start := (page - 1) * perPage
	end := start + perPage

	if start >= len(filtered) {
		return []*models.Question{}, total, nil
	}
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[start:end], total, nil
}

func (db *InMemoryDB) DeactivateQuestion(questionID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.questions[questionID]; !exists {
		return ErrNotFound
	}
	// For in-memory, we just mark it (IsActive field would be added to model)
	return nil
}

func (db *InMemoryDB) GetQuestionPreview(questionID int) (*models.QuestionPreview, error) {
	question, err := db.GetQuestion(questionID)
	if err != nil {
		return nil, err
	}

	answers, err := db.GetAnswersByQuestion(questionID)
	if err != nil {
		return nil, err
	}

	answerList := make([]models.Answer, len(answers))
	for i, ans := range answers {
		answerList[i] = *ans
	}

	return &models.QuestionPreview{
		Question: *question,
		Answers:  answerList,
	}, nil
}

func (db *InMemoryDB) ValidateQuestionData(question *models.Question, answers []models.Answer) error {
	if question.QuestionText == "" {
		return errors.New("question text is required")
	}
	if len(question.QuestionText) > 1000 {
		return errors.New("question text must be 1000 characters or less")
	}

	_, err := db.GetTopic(question.TopicID)
	if err != nil {
		return errors.New("invalid topic ID")
	}

	if question.Difficulty < 1 || question.Difficulty > 7 {
		return errors.New("difficulty must be between 1 and 7")
	}

	if len(answers) < 2 || len(answers) > 4 {
		return errors.New("must have 2-4 answer options")
	}

	correctCount := 0
	for _, answer := range answers {
		if answer.IsCorrect {
			correctCount++
		}
	}

	if correctCount != 1 {
		return errors.New("exactly 1 answer must be marked as correct")
	}

	return nil
}

// Phase 5: Bulk CSV Import - InMemoryDB implementations

func (db *InMemoryDB) ValidateCSVImport(rows []models.CSVImportRow) (*models.CSVImportResult, error) {
	result := &models.CSVImportResult{
		ValidRows:   []models.CSVImportRow{},
		InvalidRows: []models.CSVImportError{},
		Summary: models.CSVImportSummary{
			TotalRows: len(rows),
		},
	}

	for i, row := range rows {
		errors := db.validateCSVRow(row)
		if len(errors) > 0 {
			result.InvalidRows = append(result.InvalidRows, models.CSVImportError{
				Row:    i + 2,
				Data:   row,
				Errors: errors,
			})
			result.Summary.InvalidRows++
		} else {
			result.ValidRows = append(result.ValidRows, row)
			result.Summary.ValidRows++
		}
	}

	return result, nil
}

func (db *InMemoryDB) validateCSVRow(row models.CSVImportRow) []string {
	var errors []string

	if row.Topic == "" {
		errors = append(errors, "Topic is required")
	}
	if row.Question == "" {
		errors = append(errors, "Question is required")
	}
	if row.Answer1 == "" || row.Answer2 == "" {
		errors = append(errors, "At least 2 answers are required")
	}

	if row.CorrectAnswer == "" {
		errors = append(errors, "Correct answer is required")
	} else {
		validAnswer := false
		if row.CorrectAnswer == row.Answer1 || row.CorrectAnswer == row.Answer2 ||
			row.CorrectAnswer == row.Answer3 || row.CorrectAnswer == row.Answer4 {
			validAnswer = true
		}
		if !validAnswer {
			errors = append(errors, "Correct answer must match one of the answer options")
		}
	}

	if row.Difficulty == "" {
		errors = append(errors, "Difficulty is required")
	} else {
		diff, err := strconv.Atoi(row.Difficulty)
		if err != nil || diff < 1 || diff > 7 {
			errors = append(errors, "Difficulty must be a number between 1 and 7")
		}
	}

	return errors
}

func (db *InMemoryDB) BulkCreateQuestions(rows []models.CSVImportRow) (int, error) {
	imported := 0
	for _, row := range rows {
		// Get or create topic
		topicID := db.getOrCreateTopicInMemory(row.Topic)

		// Parse difficulty
		difficulty, _ := strconv.Atoi(row.Difficulty)

		// Parse tags
		var tags []string
		if row.Tags != "" {
			tagList := strings.Split(row.Tags, ",")
			for _, tag := range tagList {
				trimmed := strings.TrimSpace(tag)
				if trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
		}

		// Create question
		question := &models.Question{
			TopicID:         topicID,
			QuestionText:    row.Question,
			Difficulty:      difficulty,
			Tags:            tags,
			TopicQuestionID: fmt.Sprintf("%d.%d", topicID, time.Now().UnixNano()%10000),
			CreatedAt:       time.Now(),
		}

		if err := db.CreateQuestion(question); err != nil {
			return imported, err
		}

		// Create answers
		answers := []string{row.Answer1, row.Answer2, row.Answer3, row.Answer4}
		for _, answer := range answers {
			if answer == "" {
				continue
			}
			isCorrect := answer == row.CorrectAnswer
			ans := &models.Answer{
				QuestionID: question.ID,
				AnswerText: answer,
				IsCorrect:  isCorrect,
			}
			if err := db.CreateAnswer(ans); err != nil {
				return imported, err
			}
		}

		imported++
	}

	return imported, nil
}

func (db *InMemoryDB) getOrCreateTopicInMemory(topicName string) int {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Find existing topic
	for _, topic := range db.topics {
		if topic.Name == topicName {
			return topic.ID
		}
	}

	// Create new topic
	topic := &models.Topic{
		Name:      topicName,
		IsActive:  true,
		CreatedAt: time.Now(),
	}
	topic.ID = db.nextTopicID
	db.nextTopicID++
	db.topics[topic.ID] = topic

	return topic.ID
}

// Phase 5: Tag Management - InMemoryDB implementations

func (db *InMemoryDB) GetTagsWithMetadata() ([]*models.Tag, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	tagCounts := make(map[string]int)
	for _, q := range db.questions {
		for _, tag := range q.Tags {
			tagCounts[tag]++
		}
	}

	var result []*models.Tag
	for name, count := range tagCounts {
		result = append(result, &models.Tag{
			Name:          name,
			QuestionCount: count,
		})
	}

	return result, nil
}

func (db *InMemoryDB) MergeTags(sourceTags []string, targetTag string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	for _, q := range db.questions {
		hasSourceTag := false
		newTags := []string{}

		for _, tag := range q.Tags {
			isSourceTag := false
			for _, sourceTag := range sourceTags {
				if strings.EqualFold(tag, sourceTag) {
					isSourceTag = true
					hasSourceTag = true
					break
				}
			}
			if !isSourceTag {
				newTags = append(newTags, tag)
			}
		}

		if hasSourceTag {
			// Add target tag if not already present
			hasTarget := false
			for _, tag := range newTags {
				if strings.EqualFold(tag, targetTag) {
					hasTarget = true
					break
				}
			}
			if !hasTarget {
				newTags = append(newTags, targetTag)
			}

			q.Tags = newTags
		}
	}

	return nil
}

func (db *InMemoryDB) DeleteTag(tagName string) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Check if any questions use this tag
	for _, q := range db.questions {
		for _, tag := range q.Tags {
			if strings.EqualFold(tag, tagName) {
				return errors.New("cannot delete tag: tag is in use by questions")
			}
		}
	}

	return nil
}

func (db *InMemoryDB) GetQuestionsByTag(tagName string) ([]*models.Question, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var questions []*models.Question
	for _, q := range db.questions {
		for _, tag := range q.Tags {
			if strings.EqualFold(tag, tagName) {
				questions = append(questions, q)
				break
			}
		}
	}

	return questions, nil
}

// Phase 5: Quality Review - InMemoryDB implementations

func (db *InMemoryDB) GetFlaggedQuestions() ([]*models.FlaggedQuestion, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	flagged := make([]*models.FlaggedQuestion, 0, len(db.flaggedQuestions))
	for _, fq := range db.flaggedQuestions {
		flagged = append(flagged, fq)
	}
	return flagged, nil
}

func (db *InMemoryDB) FlagQuestion(questionID int, reason string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.questions[questionID]; !exists {
		return ErrNotFound
	}

	// Store flagged question
	db.flaggedQuestions[questionID] = &models.FlaggedQuestion{
		QuestionID: questionID,
		FlagReason: reason,
		FlaggedAt:  time.Now(),
	}
	return nil
}

func (db *InMemoryDB) UnflagQuestion(questionID int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.flaggedQuestions, questionID)
	return nil
}

func (db *InMemoryDB) AutoFlagQuestions() error {
	// First phase: collect statistics with read lock
	db.mu.RLock()

	// Count attempts and correct answers per question
	questionStats := make(map[int]struct {
		totalAttempts int
		correctCount  int
	})

	for _, qa := range db.questionAttempts {
		stats := questionStats[qa.QuestionID]
		stats.totalAttempts++
		if qa.LastAttemptCorrect {
			stats.correctCount++
		}
		questionStats[qa.QuestionID] = stats
	}

	db.mu.RUnlock()

	// Second phase: flag questions (requires write lock in FlagQuestion)
	for questionID, stats := range questionStats {
		if stats.totalAttempts < 10 {
			continue
		}

		accuracy := (float64(stats.correctCount) / float64(stats.totalAttempts)) * 100

		var reason string
		if accuracy < 30 {
			reason = models.FlagReasonTooHard
		} else if accuracy > 90 {
			reason = models.FlagReasonTooEasy
		} else {
			continue
		}

		// FlagQuestion will acquire its own write lock
		db.FlagQuestion(questionID, reason)
	}

	return nil
}

// AssignedQuiz operations (InMemoryDB)

func (db *InMemoryDB) AssignQuizToUser(quizID, userID, assignedBy int) error {
db.mu.Lock()
defer db.mu.Unlock()

// Check if already assigned
for _, aq := range db.assignedQuizzes {
if aq.QuizID == quizID && aq.UserID == userID {
return nil // already assigned, idempotent
}
}

aq := &models.AssignedQuiz{
ID:         db.nextAssignedQuizID,
QuizID:     quizID,
UserID:     userID,
AssignedAt: time.Now(),
AssignedBy: assignedBy,
}
db.nextAssignedQuizID++
db.assignedQuizzes[aq.ID] = aq
return nil
}

func (db *InMemoryDB) UnassignQuizFromUser(quizID, userID int) error {
db.mu.Lock()
defer db.mu.Unlock()

for id, aq := range db.assignedQuizzes {
if aq.QuizID == quizID && aq.UserID == userID {
delete(db.assignedQuizzes, id)
return nil
}
}
return ErrNotFound
}

func (db *InMemoryDB) GetAssignedQuizzesForUser(userID int) ([]*models.AssignedQuizInfo, error) {
db.mu.RLock()
defer db.mu.RUnlock()

var result []*models.AssignedQuizInfo
for _, aq := range db.assignedQuizzes {
if aq.UserID != userID {
continue
}

quiz, exists := db.quizzes[aq.QuizID]
if !exists {
continue
}

info := &models.AssignedQuizInfo{
AssignedQuiz:    *aq,
QuizName:        quiz.Name,
QuizDescription: quiz.Description,
QuizCategory:    quiz.Category,
}

// Find the latest completed quiz attempt for this quiz+user
for _, attempt := range db.quizAttempts {
if attempt.QuizID != nil && *attempt.QuizID == aq.QuizID && attempt.UserID == userID && attempt.CompletedTime != nil {
if info.CompletedAt == nil || attempt.CompletedTime.After(*info.CompletedAt) {
info.CompletedAt = attempt.CompletedTime
ca := attempt.CorrectAnswers
tq := attempt.TotalQuestions
info.CorrectAnswers = &ca
info.TotalQuestions = &tq
if tq > 0 {
score := float64(ca) / float64(tq) * 100.0
info.Score = &score
}
}
}
}

result = append(result, info)
}
return result, nil
}

func (db *InMemoryDB) GetAssignedQuizUserIDs(quizID int) ([]int, error) {
db.mu.RLock()
defer db.mu.RUnlock()

var userIDs []int
for _, aq := range db.assignedQuizzes {
if aq.QuizID == quizID {
userIDs = append(userIDs, aq.UserID)
}
}
return userIDs, nil
}

func (db *InMemoryDB) GetAssignedQuizStatus(quizID int) (*models.AssignedQuizStatus, error) {
db.mu.RLock()
defer db.mu.RUnlock()

quiz, exists := db.quizzes[quizID]
if !exists {
return nil, ErrNotFound
}

status := &models.AssignedQuizStatus{
QuizID:          quiz.ID,
QuizName:        quiz.Name,
QuizCategory:    quiz.Category,
QuizDescription: quiz.Description,
StudentResults:  []models.AssignedQuizStudentResult{},
}

for _, aq := range db.assignedQuizzes {
if aq.QuizID != quizID {
continue
}

user, exists := db.users[aq.UserID]
if !exists {
continue
}

r := models.AssignedQuizStudentResult{
UserID:     user.ID,
Username:   user.Username,
AssignedAt: aq.AssignedAt,
}

// Find latest completed attempt
for _, attempt := range db.quizAttempts {
if attempt.QuizID != nil && *attempt.QuizID == quizID && attempt.UserID == aq.UserID && attempt.CompletedTime != nil {
if r.CompletedAt == nil || attempt.CompletedTime.After(*r.CompletedAt) {
r.CompletedAt = attempt.CompletedTime
ca := attempt.CorrectAnswers
tq := attempt.TotalQuestions
r.CorrectAnswers = &ca
r.TotalQuestions = &tq
if tq > 0 {
score := float64(ca) / float64(tq) * 100.0
r.Score = &score
}
}
}
}

status.StudentResults = append(status.StudentResults, r)
}

return status, nil
}

func (db *InMemoryDB) ListAssignedQuizzes() ([]*models.AssignedQuiz, error) {
db.mu.RLock()
defer db.mu.RUnlock()

result := make([]*models.AssignedQuiz, 0, len(db.assignedQuizzes))
for _, aq := range db.assignedQuizzes {
result = append(result, aq)
}
return result, nil
}
