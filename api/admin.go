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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list users")
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

	// Hash password if provided
	if user.Password != "" {
		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to hash password for user")
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword
	}

	user.ID = id
	if err := h.db.UpdateUser(&user); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update user")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to delete user")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list questions")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update question")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to delete question")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list quizzes")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizzes)
}

// GetQuizHandler handles GET /quizzes/{id}
func (h *Handler) GetQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid quiz ID", http.StatusBadRequest)
		return
	}

	quiz, err := h.db.GetQuiz(id)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get quiz")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quiz)
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create quiz")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update quiz")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to delete quiz")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Topic CRUD handlers

// CreateAdminTopicHandler handles POST /admin/topics
// Creates a new topic with validation, for admin use.
func (h *Handler) CreateAdminTopicHandler(w http.ResponseWriter, r *http.Request) {
	var topic models.Topic
	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if topic.Name == "" {
		http.Error(w, "Topic name is required", http.StatusBadRequest)
		return
	}
	if len(topic.Name) > 200 {
		http.Error(w, "Topic name too long (max 200 characters)", http.StatusBadRequest)
		return
	}

	topic.CreatedAt = time.Now()
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

// ListAdminTopicsHandler handles GET /admin/topics
// Returns all active topics without user-category filtering, for admin use.
func (h *Handler) ListAdminTopicsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topics, err := h.db.ListAllTopics()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list topics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Ensure we return an empty array instead of null when there are no topics
	if topics == nil {
		topics = []*models.Topic{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topics)
}

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

	// Preserve fields that the request body does not supply so that a
	// partial update cannot accidentally reset the display order.
	// is_active is intentionally taken from the request body so that
	// admins can activate/deactivate topics via the edit form.
	existing, err := h.db.GetTopic(id)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get topic for update")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	topic.DisplayOrder = existing.DisplayOrder
	topic.CreatedAt = existing.CreatedAt

	if err := h.db.UpdateTopic(&topic); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update topic")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to delete topic")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Answer CRUD handlers
// ListAnswersHandler handles GET /admin/answers?question_id=X
func (h *Handler) ListAnswersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	questionIDStr := r.URL.Query().Get("question_id")
	if questionIDStr == "" {
		http.Error(w, "question_id is required", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	answers, err := h.db.GetAnswersByQuestion(questionID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
			"question_id": questionID,
		}).WithError(err).Error("Failed to get answers for question")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(answers)
}

// CreateAnswerHandler handles POST /admin/answers
func (h *Handler) CreateAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var answer models.Answer
	if err := json.NewDecoder(r.Body).Decode(&answer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if answer.QuestionID <= 0 {
		http.Error(w, "question_id must be a positive integer", http.StatusBadRequest)
		return
	}
	if answer.AnswerText == "" {
		http.Error(w, "answer_text is required", http.StatusBadRequest)
		return
	}

	if err := h.db.CreateAnswer(&answer); err != nil {
		h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
			"question_id": answer.QuestionID,
		}).WithError(err).Error("Failed to create answer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(answer)
}

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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update answer")
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
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to delete answer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ImportQuestionsJSONHandler handles POST /admin/questions/import-json
// This handler accepts a JSON file in the format used by the CLI tool
// along with a topic_id parameter to assign all questions to that topic
func (h *Handler) ImportQuestionsJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get topic_id from form
	topicIDStr := r.FormValue("topic_id")
	if topicIDStr == "" {
		http.Error(w, "topic_id is required", http.StatusBadRequest)
		return
	}

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		http.Error(w, "Invalid topic ID", http.StatusBadRequest)
		return
	}

	// Verify topic exists
	topic, err := h.db.GetTopic(topicID)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Topic not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get topic")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the file from the form
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse JSON file
	var questionData struct {
		Topic     string `json:"topic"`
		Questions []struct {
			Text       string   `json:"text"`
			Difficulty int      `json:"difficulty"`
			Tags       []string `json:"tags"`
			Answers    []struct {
				Text      string `json:"text"`
				IsCorrect bool   `json:"is_correct"`
			} `json:"answers"`
		} `json:"questions"`
	}

	if err := json.NewDecoder(file).Decode(&questionData); err != nil {
		http.Error(w, "Failed to parse JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(questionData.Questions) == 0 {
		http.Error(w, "No questions found in JSON file", http.StatusBadRequest)
		return
	}

	// Get existing questions to check for duplicates
	existingQuestions, _ := h.db.GetQuestionsByTopic(topicID)
	existingTexts := make(map[string]bool)
	for _, eq := range existingQuestions {
		existingTexts[eq.QuestionText] = true
	}

	importedCount := 0
	skippedCount := 0
	errorCount := 0
	var errors []string

	// Import questions
	startNum := len(existingQuestions) + 1
	for i, q := range questionData.Questions {
		// Validate question
		if q.Text == "" {
			errorCount++
			errors = append(errors, "Question "+(strconv.Itoa(i+1))+": question text is empty")
			continue
		}

		if q.Difficulty < 1 || q.Difficulty > 7 {
			errorCount++
			errors = append(errors, "Question "+(strconv.Itoa(i+1))+": difficulty must be between 1 and 7")
			continue
		}

		if len(q.Answers) < 2 || len(q.Answers) > 6 {
			errorCount++
			errors = append(errors, "Question "+(strconv.Itoa(i+1))+": must have between 2 and 6 answers")
			continue
		}

		// Validate answers
		correctCount := 0
		hasEmptyAnswer := false
		for _, a := range q.Answers {
			if a.Text == "" {
				hasEmptyAnswer = true
				break
			}
			if a.IsCorrect {
				correctCount++
			}
		}

		if hasEmptyAnswer {
			errorCount++
			errors = append(errors, "Question "+(strconv.Itoa(i+1))+": answer text is empty")
			continue
		}

		if correctCount != 1 {
			errorCount++
			errors = append(errors, "Question "+(strconv.Itoa(i+1))+": exactly one answer must be correct")
			continue
		}

		// Check for duplicates
		if existingTexts[q.Text] {
			skippedCount++
			continue
		}

		// Create question with sequential numbering based on actual imported count
		topicQuestionID := strconv.Itoa(topicID) + "." + strconv.Itoa(startNum+importedCount)
		question := models.Question{
			TopicID:         topicID,
			QuestionText:    q.Text,
			Difficulty:      q.Difficulty,
			Tags:            q.Tags,
			TopicQuestionID: topicQuestionID,
			CreatedAt:       time.Now(),
		}

		if err := h.db.CreateQuestion(&question); err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create question during JSON import")
			errorCount++
			errors = append(errors, "Question "+(strconv.Itoa(i+1))+": failed to create - "+err.Error())
			continue
		}

		// Create answers
		for _, a := range q.Answers {
			answer := models.Answer{
				QuestionID: question.ID,
				AnswerText: a.Text,
				IsCorrect:  a.IsCorrect,
			}
			if err := h.db.CreateAnswer(&answer); err != nil {
				h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create answer during JSON import")
				errorCount++
				errors = append(errors, "Question "+(strconv.Itoa(i+1))+": failed to create answer - "+err.Error())
			}
		}

		importedCount++
		existingTexts[q.Text] = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"imported": importedCount,
		"skipped":  skippedCount,
		"errors":   errorCount,
		"messages": errors,
		"topic":    topic.Name,
	})
}
