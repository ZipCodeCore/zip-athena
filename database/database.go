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
	nextUserID          int
	nextTopicID         int
	nextQuestionID      int
	nextAnswerID        int
	nextQuizAttemptID   int
	nextQAttemptID      int
	nextQuizID          int
	nextUQStateID       int
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
		nextUserID:         1,
		nextTopicID:        1,
		nextQuestionID:     1,
		nextAnswerID:       1,
		nextQuizAttemptID:  1,
		nextQAttemptID:     1,
		nextQuizID:         1,
		nextUQStateID:      1,
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
