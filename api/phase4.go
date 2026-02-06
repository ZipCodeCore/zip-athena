package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

// Phase 4: Quiz Attempt History Handlers

// GetQuizHistoryHandler handles GET /students/quiz-history
func (h *Handler) GetQuizHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user ID from query parameter or session
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Get optional topic filter
	topicIDStr := r.URL.Query().Get("topic_id")

	// Get pagination parameters
	page := 1
	perPage := 10
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if perPageStr := r.URL.Query().Get("per_page"); perPageStr != "" {
		if pp, err := strconv.Atoi(perPageStr); err == nil && pp > 0 && pp <= 100 {
			perPage = pp
		}
	}

	var attempts []*models.QuizAttempt
	if topicIDStr != "" {
		topicID, err := strconv.Atoi(topicIDStr)
		if err != nil {
			http.Error(w, "Invalid topic_id", http.StatusBadRequest)
			return
		}
		attempts, err = h.db.GetQuizHistoryByTopic(userID, topicID, page, perPage)
		if err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get quiz history by topic")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		attempts, err = h.db.GetQuizHistoryByUser(userID, page, perPage)
		if err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get quiz history")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempts)
}

// GetQuizAttemptDetailsHandler handles GET /students/quiz-history/{id}
func (h *Handler) GetQuizAttemptDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get attempt ID from URL path
	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	attemptIDStr := parts[len(parts)-1]
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		http.Error(w, "Invalid attempt ID", http.StatusBadRequest)
		return
	}

	attempt, err := h.db.GetQuizAttemptWithDetails(attemptID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get quiz attempt details")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Load questions if quiz is in progress (for resume functionality)
	var quizQuestions []QuizQuestion
	if attempt.CompletedTime == nil && len(attempt.QuestionIDs) > 0 {
		// Get all questions
		allQuestions, err := h.db.ListQuestions()
		if err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to list questions")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Build a map of TopicQuestionID -> Question
		questionMap := make(map[string]*models.Question)
		for _, q := range allQuestions {
			questionMap[q.TopicQuestionID] = q
		}

		// Reconstruct questions from QuestionIDs
		for _, topicQuestionID := range attempt.QuestionIDs {
			if q, exists := questionMap[topicQuestionID]; exists {
				// Get answers for this question
				answers, err := h.db.GetAnswersByQuestion(q.ID)
				if err != nil {
					h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get answers for question")
					continue
				}

				quizQuestions = append(quizQuestions, QuizQuestion{
					Question: q,
					Answers:  answers,
				})
			}
		}
	}

	// Create response with both attempt and questions
	response := struct {
		*models.QuizAttempt
		Questions []QuizQuestion `json:"questions,omitempty"`
	}{
		QuizAttempt: attempt,
		Questions:   quizQuestions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Phase 4: Timed Test Taking Handlers

// GetAvailableTestsHandler handles GET /tests/available
func (h *Handler) GetAvailableTestsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	tests, err := h.db.GetAvailableTests(userID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get available tests")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tests)
}

// GetTestDetailsHandler handles GET /tests/{id}
func (h *Handler) GetTestDetailsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	testIDStr := parts[len(parts)-1]
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		http.Error(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	test, err := h.db.GetTest(testID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get test")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get test questions
	questions, err := h.db.GetTestQuestions(testID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get test questions")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"test":           test,
		"question_count": len(questions),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StartTestHandler handles POST /tests/{id}/start
func (h *Handler) StartTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	testIDStr := parts[2]
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		http.Error(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already has an attempt for this test
	existingAttempts, err := h.db.GetTestAttemptsByUser(req.UserID, testID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to check existing attempts")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check for incomplete attempts
	for _, attempt := range existingAttempts {
		if attempt.SubmitTime == nil {
			http.Error(w, "You already have an incomplete attempt for this test", http.StatusConflict)
			return
		}
	}

	// Create new test attempt
	attempt := &models.TestAttempt{
		UserID:          req.UserID,
		TestID:          testID,
		StartTime:       time.Now(),
		Answers:         make(map[string]int),
		QuestionResults: make(map[string]bool),
	}

	if err := h.db.CreateTestAttempt(attempt); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create test attempt")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(attempt)
}

// UpdateTestAnswerHandler handles PUT /test-attempts/{id}/answer
func (h *Handler) UpdateTestAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	attemptIDStr := parts[2]
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		http.Error(w, "Invalid attempt ID", http.StatusBadRequest)
		return
	}

	var req struct {
		QuestionID string `json:"question_id"`
		AnswerID   int    `json:"answer_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateTestAttemptAnswer(attemptID, req.QuestionID, req.AnswerID); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update test answer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

// SubmitTestHandler handles POST /test-attempts/{id}/submit
func (h *Handler) SubmitTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	attemptIDStr := parts[2]
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		http.Error(w, "Invalid attempt ID", http.StatusBadRequest)
		return
	}

	var req struct {
		IsAutoSubmit bool `json:"is_auto_submit"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.db.SubmitTestAttempt(attemptID, req.IsAutoSubmit); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to submit test")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get results
	results, err := h.db.GetTestResults(attemptID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get test results")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// GetTestResultsHandler handles GET /test-attempts/{id}/results
func (h *Handler) GetTestResultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	attemptIDStr := parts[2]
	attemptID, err := strconv.Atoi(attemptIDStr)
	if err != nil {
		http.Error(w, "Invalid attempt ID", http.StatusBadRequest)
		return
	}

	results, err := h.db.GetTestResults(attemptID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get test results")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// Phase 4: Tag-Based Review Handlers

// GetTagsHandler handles GET /tags
func (h *Handler) GetTagsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tagCounts, err := h.db.GetTagsWithCounts()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get tags")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tagCounts)
}

// GetTagStatsHandler handles GET /tags/stats
func (h *Handler) GetTagStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	stats, err := h.db.GetTagStats(userID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get tag stats")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GenerateQuizByTagsHandler handles POST /quiz/generate/by-tags
func (h *Handler) GenerateQuizByTagsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Tags   []string `json:"tags"`
		UseAnd bool     `json:"use_and"` // true for AND logic, false for OR
		UserID int      `json:"user_id"`
		Count  int      `json:"count"` // number of questions to include
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Tags) == 0 {
		http.Error(w, "At least one tag is required", http.StatusBadRequest)
		return
	}

	if req.Count <= 0 {
		req.Count = 10 // default
	}

	questions, err := h.db.GetQuestionsByTags(req.Tags, req.UseAnd)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get questions by tags")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(questions) == 0 {
		http.Error(w, "No questions found for the given tags", http.StatusNotFound)
		return
	}

	// Remove duplicate questions based on TopicQuestionID
	uniqueQuestions := make(map[string]*models.Question)
	for _, q := range questions {
		if _, exists := uniqueQuestions[q.TopicQuestionID]; !exists {
			uniqueQuestions[q.TopicQuestionID] = q
		}
	}

	// Convert map back to slice
	questions = make([]*models.Question, 0, len(uniqueQuestions))
	for _, q := range uniqueQuestions {
		questions = append(questions, q)
	}

	// Limit to requested count
	if len(questions) > req.Count {
		questions = questions[:req.Count]
	}

	// Create quiz attempt
	attempt := &models.QuizAttempt{
		UserID:          req.UserID,
		StartTime:       time.Now(),
		TotalQuestions:  len(questions),
		QuestionIDs:     make([]string, len(questions)),
		QuestionResults: make(map[string]bool),
		AttemptNumber:   1,
	}

	for i, q := range questions {
		attempt.QuestionIDs[i] = q.TopicQuestionID
	}

	if err := h.db.CreateQuizAttempt(attempt); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create quiz attempt")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"attempt_id": attempt.ID,
		"questions":  questions,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Phase 4: Instructor Test Management Handlers

// PublishTestHandler handles POST /instructor/tests/{id}/publish
func (h *Handler) PublishTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	testIDStr := parts[3]
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		http.Error(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	if err := h.db.PublishTest(testID); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to publish test")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Test published successfully"}`))
}

// GetTestResultsInstructorHandler handles GET /instructor/tests/{id}/results
func (h *Handler) GetTestResultsInstructorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	testIDStr := parts[3]
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		http.Error(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	results, err := h.db.GetTestResultsByTest(testID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get test results")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// ExportTestResultsHandler handles GET /instructor/tests/{id}/export
func (h *Handler) ExportTestResultsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	testIDStr := parts[3]
	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		http.Error(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	csvData, err := h.db.ExportTestResults(testID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to export test results")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=test_%d_results.csv", testID))
	w.Write(csvData)
}

// Phase 4: Question Performance Analytics Handlers

// GetQuestionAnalyticsHandler handles GET /instructor/questions/analytics
func (h *Handler) GetQuestionAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	analytics, err := h.db.GetQuestionAnalytics()
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get question analytics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// GetQuestionDetailedAnalyticsHandler handles GET /instructor/questions/{id}/analytics
func (h *Handler) GetQuestionDetailedAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	questionIDStr := parts[3]
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	analytics, err := h.db.GetQuestionDetailedAnalytics(questionID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get question analytics")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// GetAnswerDistributionHandler handles GET /instructor/questions/{id}/answer-distribution
func (h *Handler) GetAnswerDistributionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	questionIDStr := parts[3]
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	distribution, err := h.db.GetAnswerDistribution(questionID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get answer distribution")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(distribution)
}

// UpdateQuestionDifficultyHandler handles PUT /instructor/questions/{id}/difficulty
func (h *Handler) UpdateQuestionDifficultyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	parts := strings.Split(path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	questionIDStr := parts[3]
	questionID, err := strconv.Atoi(questionIDStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Difficulty int `json:"difficulty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.db.UpdateQuestionDifficultyFromAnalytics(questionID, req.Difficulty); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to update question difficulty")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success","message":"Difficulty updated successfully"}`))
}
