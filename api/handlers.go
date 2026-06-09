package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/logger"
	"github.com/ZipCodeCore/zip-athena/models"
)

// Handler contains the API handlers for the application
type Handler struct {
	db         database.Database
	logger     *logger.Logger
	jwtManager *JWTManager
}

// SetJWTManager sets the JWT manager for token generation and validation
func (h *Handler) SetJWTManager(jm *JWTManager) {
	h.jwtManager = jm
}

// NewHandler creates a new API handler
func NewHandler(db database.Database) *Handler {
	return &Handler{
		db:     db,
		logger: logger.GetDefault(),
	}
}

// NewHandlerWithLogger creates a new API handler with a custom logger
func NewHandlerWithLogger(db database.Database, log *logger.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: log,
	}
}

// CreateUserHandler handles POST /users
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Username    string `json:"username"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		IsAdmin     bool   `json:"is_admin"`
		IsActive    *bool  `json:"is_active"`
		Category    string `json:"category"`
		TotalPoints int    `json:"total_points"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.User{
		Username:    body.Username,
		Email:       body.Email,
		Password:    body.Password,
		IsAdmin:     body.IsAdmin,
		IsActive:    true,
		IsActiveSet: body.IsActive != nil,
		Category:    body.Category,
		TotalPoints: body.TotalPoints,
	}
	if body.IsActive != nil {
		user.IsActive = *body.IsActive
	}

	// Hash password if provided
	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
				"username": user.Username,
			}).WithError(err).Error("Failed to hash password")
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
	}

	user.CreatedAt = time.Now()
	if err := h.db.CreateUser(&user); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Don't send password back
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUserHandler handles GET /users/{id}
func (h *Handler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUser(id)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateTopicHandler handles POST /topics
func (h *Handler) CreateTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var topic models.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topic.CreatedAt = time.Now()
	// Default to active for new topics
	topic.IsActive = true

	if err := h.db.CreateTopic(&topic); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create topic")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(topic)
}

// ListTopicsHandler handles GET /topics
func (h *Handler) ListTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Require user_id parameter for topic filtering
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	// Parse and validate user_id
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id parameter", http.StatusBadRequest)
		return
	}

	// Get user to determine category and admin status
	user, err := h.db.GetUser(userID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get user")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get all topics from database
	topics, err := h.db.ListTopics()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list topics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter topics based on user category
	filteredTopics := []*models.Topic{}
	for _, topic := range topics {
		if userCanAccessTopic(user, topic) {
			filteredTopics = append(filteredTopics, topic)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredTopics)
}

// GetTopicHandler handles GET /topics/{id}
func (h *Handler) GetTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	topic, err := h.db.GetTopic(id)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get topic")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

// CreateQuestionHandler handles POST /questions
func (h *Handler) CreateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	question.CreatedAt = time.Now()
	if err := h.db.CreateQuestion(&question); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create question")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(question)
}

// GetQuestionsByTopicHandler handles GET /questions?topic_id={id}
func (h *Handler) GetQuestionsByTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topicIDStr := r.URL.Query().Get("topic_id")
	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	questions, err := h.db.GetQuestionsByTopic(topicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// CreateQuizAttemptHandler handles POST /quiz-attempts
func (h *Handler) CreateQuizAttemptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var attempt models.QuizAttempt
	if err := json.NewDecoder(r.Body).Decode(&attempt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	attempt.StartTime = time.Now()
	if err := h.db.CreateQuizAttempt(&attempt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attempt)
}

// GetQuizAttemptsByUserHandler handles GET /quiz-attempts?user_id={id}
func (h *Handler) GetQuizAttemptsByUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	attempts, err := h.db.GetQuizAttemptsByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempts)
}
