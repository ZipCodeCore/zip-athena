package database

import (
	"errors"
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
}

// InMemoryDB implements the Database interface using in-memory storage
type InMemoryDB struct {
	users            map[int]*models.User
	topics           map[int]*models.Topic
	questions        map[int]*models.Question
	answers          map[int]*models.Answer
	quizAttempts     map[int]*models.QuizAttempt
	questionAttempts map[int]*models.QuestionAttempt
	nextUserID       int
	nextTopicID      int
	nextQuestionID   int
	nextAnswerID     int
	nextQuizID       int
	nextQAttemptID   int
	mu               sync.RWMutex
}

// NewInMemoryDB creates a new in-memory database
func NewInMemoryDB() *InMemoryDB {
	return &InMemoryDB{
		users:            make(map[int]*models.User),
		topics:           make(map[int]*models.Topic),
		questions:        make(map[int]*models.Question),
		answers:          make(map[int]*models.Answer),
		quizAttempts:     make(map[int]*models.QuizAttempt),
		questionAttempts: make(map[int]*models.QuestionAttempt),
		nextUserID:       1,
		nextTopicID:      1,
		nextQuestionID:   1,
		nextAnswerID:     1,
		nextQuizID:       1,
		nextQAttemptID:   1,
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
