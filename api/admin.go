package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// Admin middleware to check if user is admin
func (h *Handler) AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// In a real app, you'd check session/JWT token here
		// For now, we'll check a simple header
		userIDStr := r.Header.Get("X-User-ID")
		if userIDStr == "" {
			http.Error(w, "Unauthorized: User ID required", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		user, err := h.db.GetUser(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if !user.IsAdmin {
			http.Error(w, "Forbidden: Admin access required", http.StatusForbidden)
			return
		}

		next(w, r)
	}
}

// ListUsersHandler handles GET /admin/users
func (h *Handler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.db.ListUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove passwords from response
	for _, user := range users {
		user.Password = ""
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUserHandler handles PUT /admin/users/{id}
func (h *Handler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id
	if err := h.db.UpdateUser(&user); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUserHandler handles DELETE /admin/users/{id}
func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteUser(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListQuestionsHandler handles GET /admin/questions
func (h *Handler) ListQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	questions, err := h.db.ListQuestions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// UpdateQuestionHandler handles PUT /admin/questions/{id}
func (h *Handler) UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var question models.Question
	if err := json.NewDecoder(r.Body).Decode(&question); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	question.ID = id
	if err := h.db.UpdateQuestion(&question); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

// DeleteQuestionHandler handles DELETE /admin/questions/{id}
func (h *Handler) DeleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteQuestion(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Quiz CRUD handlers
// ListQuizzesHandler handles GET /admin/quizzes
func (h *Handler) ListQuizzesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quizzes, err := h.db.ListQuizzes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizzes)
}

// CreateQuizHandler handles POST /admin/quizzes
func (h *Handler) CreateQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quiz.CreatedAt = time.Now()
	quiz.UpdatedAt = time.Now()
	if err := h.db.CreateQuiz(&quiz); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(quiz)
}

// UpdateQuizHandler handles PUT /admin/quizzes/{id}
func (h *Handler) UpdateQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quiz.ID = id
	quiz.UpdatedAt = time.Now()
	if err := h.db.UpdateQuiz(&quiz); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
}

// DeleteQuizHandler handles DELETE /admin/quizzes/{id}
func (h *Handler) DeleteQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteQuiz(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Topic CRUD handlers
// UpdateTopicHandler handles PUT /admin/topics/{id}
func (h *Handler) UpdateTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	var topic models.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	topic.ID = id
	if err := h.db.UpdateTopic(&topic); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

// DeleteTopicHandler handles DELETE /admin/topics/{id}
func (h *Handler) DeleteTopicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTopic(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Answer CRUD handlers
// UpdateAnswerHandler handles PUT /admin/answers/{id}
func (h *Handler) UpdateAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	var answer models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	answer.ID = id
	if err := h.db.UpdateAnswer(&answer); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Answer not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answer)
}

// DeleteAnswerHandler handles DELETE /admin/answers/{id}
func (h *Handler) DeleteAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteAnswer(id); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Answer not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
