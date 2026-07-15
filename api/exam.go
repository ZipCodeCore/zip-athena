package api

import (
	cryptorand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mathrand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// generateToken creates a cryptographically random hex token.
func generateToken() (string, error) {
	b := make([]byte, 16)
	if _, err := cryptorand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateExamHandler handles POST /instructor/exams/generate
// Creates a new exam by randomly selecting questions that match the given
// categories/topics and difficulty levels, then saves the exam.
func (h *Handler) GenerateExamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Name             string   `json:"name"`
		Categories       []string `json:"categories"`
		TopicIDs         []int    `json:"topic_ids"`
		DifficultyLevels []int    `json:"difficulty_levels"`
		QuestionCount    int      `json:"question_count"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if body.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if len(body.Categories) == 0 && len(body.TopicIDs) == 0 {
		http.Error(w, "at least one category or topic is required", http.StatusBadRequest)
		return
	}
	if len(body.DifficultyLevels) == 0 {
		http.Error(w, "at least one difficulty_level is required", http.StatusBadRequest)
		return
	}
	count := body.QuestionCount
	if count <= 0 {
		count = 30
	}

	var questions []*models.Question
	var err error
	if len(body.TopicIDs) > 0 {
		questions, err = h.db.GetQuestionsByTopicIDs(body.TopicIDs, body.DifficultyLevels, count)
	} else {
		questions, err = h.db.GetQuestionsByCategories(body.Categories, body.DifficultyLevels, count)
	}
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to fetch questions for exam generation")
		http.Error(w, "Failed to fetch questions: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(questions) == 0 {
		http.Error(w, "No questions found matching the selected criteria", http.StatusUnprocessableEntity)
		return
	}

	questionIDs := make([]int, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.ID
	}

	token, err := generateToken()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to generate exam token")
		http.Error(w, "Failed to generate exam token", http.StatusInternalServerError)
		return
	}

	assignerIDStr := r.Header.Get("X-User-ID")
	assignerID, _ := strconv.Atoi(assignerIDStr)

	// Record categories: use provided categories, or derive from topic IDs
	examCategories := body.Categories
	if len(body.TopicIDs) > 0 {
		examCategories = deriveTopicNames(questions)
	}

	exam := &models.Exam{
		Name:             body.Name,
		Categories:       examCategories,
		DifficultyLevels: body.DifficultyLevels,
		QuestionIDs:      questionIDs,
		QuestionCount:    len(questionIDs),
		AccessToken:      token,
		CreatedBy:        assignerID,
		CreatedAt:        time.Now(),
	}

	if err := h.db.CreateExam(exam); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create exam")
		http.Error(w, "Failed to create exam: "+err.Error(), http.StatusInternalServerError)
		return
	}

	examURL := fmt.Sprintf("/exam.html?token=%s", exam.AccessToken)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":        true,
		"exam_id":        exam.ID,
		"exam_name":      exam.Name,
		"question_count": exam.QuestionCount,
		"access_token":   exam.AccessToken,
		"exam_url":       examURL,
	})
}

// deriveTopicNames builds a deduplicated list of category labels from questions.
func deriveTopicNames(questions []*models.Question) []string {
	seen := map[string]bool{}
	var result []string
	for _, q := range questions {
		if len(q.Tags) > 0 && !seen[q.Tags[0]] {
			seen[q.Tags[0]] = true
			result = append(result, q.Tags[0])
		}
	}
	return result
}

// ListExamsHandler handles GET /instructor/exams
func (h *Handler) ListExamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	exams, err := h.db.ListExams()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list exams")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exams)
}

// AssignExamToStudentsHandler handles POST /instructor/exams/assign?exam_id={id}
// Assigns an exam to one or more students.
func (h *Handler) AssignExamToStudentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	examIDStr := r.URL.Query().Get("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil || examID <= 0 {
		http.Error(w, "Invalid exam_id", http.StatusBadRequest)
		return
	}

	if _, err := h.db.GetExam(examID); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Exam not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	assignerIDStr := r.Header.Get("X-User-ID")
	assignerID, _ := strconv.Atoi(assignerIDStr)

	var body struct {
		UserIDs []int `json:"user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if len(body.UserIDs) == 0 {
		http.Error(w, "user_ids must not be empty", http.StatusBadRequest)
		return
	}

	assigned := 0
	var assignErrors []string
	for _, uid := range body.UserIDs {
		if err := h.db.AssignExamToUser(examID, uid, assignerID); err != nil {
			assignErrors = append(assignErrors, fmt.Sprintf("user %d: %s", uid, err.Error()))
		} else {
			assigned++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  true,
		"assigned": assigned,
		"errors":   assignErrors,
	})
}

// GetExamStatusHandler handles GET /instructor/exams/status?exam_id={id}
// Returns results for all students assigned to the given exam.
func (h *Handler) GetExamStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	examIDStr := r.URL.Query().Get("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil || examID <= 0 {
		http.Error(w, "Invalid exam_id", http.StatusBadRequest)
		return
	}

	status, err := h.db.GetExamStatus(examID)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Exam not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get exam status")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetAssignedExamsHandler handles GET /assigned-exams?user_id={id}
// Returns exams assigned to a student.
func (h *Handler) GetAssignedExamsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	exams, err := h.db.GetAssignedExamsForUser(userID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get assigned exams")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exams)
}

// StartExamHandler handles POST /exam/start?token={token}&user_id={id}
// Verifies the student is assigned to the exam, creates an attempt, and returns
// the questions (answers randomized).
func (h *Handler) StartExamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token is required", http.StatusBadRequest)
		return
	}
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	exam, err := h.db.GetExamByToken(token)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Exam not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.db.GetUser(userID)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !user.IsActive {
		http.Error(w, "User account is inactive", http.StatusForbidden)
		return
	}

	if !user.IsAdmin {
		assignedExams, err := h.db.GetAssignedExamsForUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		assigned := false
		for _, ae := range assignedExams {
			if ae.ExamID == exam.ID {
				assigned = true
				break
			}
		}
		if !assigned {
			http.Error(w, "Exam not assigned to this user", http.StatusForbidden)
			return
		}
	}

	// Build question list with answers
	type examAnswer struct {
		ID         int    `json:"id"`
		AnswerText string `json:"answer_text"`
	}
	type examQuestion struct {
		ID           int          `json:"id"`
		QuestionText string       `json:"question_text"`
		Category     string       `json:"category"`
		Difficulty   int          `json:"difficulty"`
		Answers      []examAnswer `json:"answers"`
	}

	questions := make([]examQuestion, 0, len(exam.QuestionIDs))
	for _, qid := range exam.QuestionIDs {
		q, err := h.db.GetQuestion(qid)
		if err != nil {
			continue
		}
		answers, err := h.db.GetAnswersByQuestion(qid)
		if err != nil {
			continue
		}
		// Randomize answer order
		mathrand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

		eas := make([]examAnswer, len(answers))
		for i, a := range answers {
			eas[i] = examAnswer{ID: a.ID, AnswerText: a.AnswerText}
		}
		categoryLabel := ""
		if topic, err := h.db.GetTopic(q.TopicID); err == nil {
			categoryLabel = strings.TrimSpace(topic.Name)
			if categoryLabel == "" {
				categoryLabel = strings.TrimSpace(topic.Category)
			}
		} else {
			h.logger.WithContext(r.Context()).WithError(err).WithField("question_id", q.ID).WithField("topic_id", q.TopicID).Warn("Failed to load topic for exam question")
		}
		if categoryLabel == "" && len(q.Tags) > 0 {
			categoryLabel = strings.TrimSpace(q.Tags[0])
		}
		questions = append(questions, examQuestion{
			ID:           q.ID,
			QuestionText: q.QuestionText,
			Category:     categoryLabel,
			Difficulty:   q.Difficulty,
			Answers:      eas,
		})
	}

	attempt := &models.ExamAttempt{
		ExamID:          exam.ID,
		UserID:          userID,
		StartedAt:       time.Now(),
		TotalQuestions:  len(questions),
		Answers:         map[string]int{},
		QuestionResults: map[string]bool{},
	}
	if err := h.db.CreateExamAttempt(attempt); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create exam attempt")
		http.Error(w, "Failed to start exam: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"attempt_id":     attempt.ID,
		"exam_id":        exam.ID,
		"exam_name":      exam.Name,
		"question_count": len(questions),
		"questions":      questions,
	})
}

// SubmitExamHandler handles POST /exam/submit
// Body: { "attempt_id": int, "answers": { "questionID": answerID, ... } }
func (h *Handler) SubmitExamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		AttemptID int            `json:"attempt_id"`
		Answers   map[string]int `json:"answers"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	if body.AttemptID <= 0 {
		http.Error(w, "attempt_id is required", http.StatusBadRequest)
		return
	}
	if body.Answers == nil {
		body.Answers = map[string]int{}
	}

	if err := h.db.SubmitExamAttempt(body.AttemptID, body.Answers); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Attempt not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to submit exam attempt")
		http.Error(w, "Failed to submit exam: "+err.Error(), http.StatusInternalServerError)
		return
	}

	attempt, err := h.db.GetExamAttempt(body.AttemptID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"attempt_id":       attempt.ID,
		"exam_id":          attempt.ExamID,
		"score":            attempt.Score,
		"correct_answers":  attempt.CorrectAnswers,
		"total_questions":  attempt.TotalQuestions,
		"completed_at":     attempt.CompletedAt,
		"question_results": attempt.QuestionResults,
	})
}

// GetAvailableExamCategoriesHandler handles GET /instructor/exams/categories
// Returns all distinct topic categories available in the database.
func (h *Handler) GetAvailableExamCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	topics, err := h.db.ListAllTopics()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list topics for categories")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	seen := map[string]bool{}
	var categories []string
	for _, t := range topics {
		if t.Category != "" && !seen[t.Category] {
			seen[t.Category] = true
			categories = append(categories, t.Category)
		}
	}
	if categories == nil {
		categories = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetAvailableExamTopicsHandler handles GET /instructor/exams/topics
// Returns all active topics (id, name, category) available for exam generation.
func (h *Handler) GetAvailableExamTopicsHandler(w http.ResponseWriter, r *http.Request) {
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

	type topicOption struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	result := make([]topicOption, 0, len(topics))
	for _, t := range topics {
		result = append(result, topicOption{
			ID:       t.ID,
			Name:     t.Name,
			Category: t.Category,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// DeleteExamHandler handles DELETE /instructor/exams?exam_id={id}
// Deletes an exam and its assignments/attempts.
func (h *Handler) DeleteExamHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	examIDStr := r.URL.Query().Get("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil || examID <= 0 {
		http.Error(w, "Invalid exam_id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteExam(examID); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Exam not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to delete exam")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}

// UnassignExamFromUserHandler handles DELETE /instructor/exams/unassign?exam_id={id}&user_id={id}
// Removes an exam assignment from a specific user.
func (h *Handler) UnassignExamFromUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	examIDStr := r.URL.Query().Get("exam_id")
	examID, err := strconv.Atoi(examIDStr)
	if err != nil || examID <= 0 {
		http.Error(w, "Invalid exam_id", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.db.UnassignExamFromUser(examID, userID); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Assignment not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to unassign exam")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
}
