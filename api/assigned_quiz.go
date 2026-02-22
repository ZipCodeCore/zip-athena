package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// ImportAndCreateQuizFromJSONHandler handles POST /instructor/assigned-quizzes/import
// Imports questions from a JSON file, creates a topic (or reuses existing), and creates a named quiz.
func (h *Handler) ImportAndCreateQuizFromJSONHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (10 MB max)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	quizName := r.FormValue("quiz_name")
	if quizName == "" {
		http.Error(w, "quiz_name is required", http.StatusBadRequest)
		return
	}

	category := r.FormValue("category")
	if category == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}

	description := r.FormValue("description")

	// Get the JSON file
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

	// Validate all questions first
	var validationErrors []string
	for i, q := range questionData.Questions {
		if q.Text == "" {
			validationErrors = append(validationErrors, "Question "+strconv.Itoa(i+1)+": question text is empty")
			continue
		}
		if q.Difficulty < 1 || q.Difficulty > 7 {
			validationErrors = append(validationErrors, "Question "+strconv.Itoa(i+1)+": difficulty must be between 1 and 7")
			continue
		}
		if len(q.Answers) < 2 || len(q.Answers) > 6 {
			validationErrors = append(validationErrors, "Question "+strconv.Itoa(i+1)+": must have between 2 and 6 answers")
			continue
		}
		correctCount := 0
		for _, a := range q.Answers {
			if a.IsCorrect {
				correctCount++
			}
		}
		if correctCount != 1 {
			validationErrors = append(validationErrors, "Question "+strconv.Itoa(i+1)+": exactly one answer must be marked correct")
		}
	}

	if len(validationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":  false,
			"errors":   validationErrors,
		})
		return
	}

	// Create a topic for this quiz
	topicName := quizName
	if questionData.Topic != "" {
		topicName = questionData.Topic
	}

	topic := &models.Topic{
		Name:        topicName,
		Description: description,
		Category:    category,
		IsActive:    true,
		CreatedAt:   time.Now(),
	}
	if err := h.db.CreateTopic(topic); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create topic for quiz import")
		http.Error(w, "Failed to create topic: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the quiz
	quiz := &models.Quiz{
		Name:        quizName,
		TopicID:     topic.ID,
		Category:    category,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := h.db.CreateQuiz(quiz); err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create quiz during import")
		http.Error(w, "Failed to create quiz: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Import questions and link to quiz
	importedCount := 0
	for i, q := range questionData.Questions {
		topicQuestionID := strconv.Itoa(topic.ID) + "." + strconv.Itoa(i+1)
		question := models.Question{
			TopicID:         topic.ID,
			QuestionText:    q.Text,
			Difficulty:      q.Difficulty,
			Tags:            q.Tags,
			TopicQuestionID: topicQuestionID,
			IsActive:        true,
			CreatedAt:       time.Now(),
		}

		if err := h.db.CreateQuestion(&question); err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create question during quiz import")
			continue
		}

		// Create answers (first answer in JSON is correct per spec, but we mark via is_correct field)
		for _, a := range q.Answers {
			answer := models.Answer{
				QuestionID: question.ID,
				AnswerText: a.Text,
				IsCorrect:  a.IsCorrect,
			}
			if err := h.db.CreateAnswer(&answer); err != nil {
				h.logger.WithContext(r.Context()).WithError(err).Error("Failed to create answer during quiz import")
			}
		}

		// Link question to quiz
		if err := h.db.AddQuestionToQuiz(quiz.ID, question.ID, i+1); err != nil {
			h.logger.WithContext(r.Context()).WithError(err).Error("Failed to link question to quiz")
		}

		importedCount++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"quiz_id":   quiz.ID,
		"quiz_name": quiz.Name,
		"topic_id":  topic.ID,
		"category":  category,
		"imported":  importedCount,
	})
}

// AssignQuizToStudentsHandler handles POST /instructor/assigned-quizzes/{id}/assign
// Assigns a quiz to one or more students.
func (h *Handler) AssignQuizToStudentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get quiz ID from query param
	quizIDStr := r.URL.Query().Get("quiz_id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil || quizID <= 0 {
		http.Error(w, "Invalid quiz_id", http.StatusBadRequest)
		return
	}

	// Verify quiz exists
	if _, err := h.db.GetQuiz(quizID); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get assigning user from header
	assignerIDStr := r.Header.Get("X-User-ID")
	assignerID, _ := strconv.Atoi(assignerIDStr)

	// Parse body: list of user IDs
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
	for _, userID := range body.UserIDs {
		if err := h.db.AssignQuizToUser(quizID, userID, assignerID); err != nil {
			assignErrors = append(assignErrors, "user "+strconv.Itoa(userID)+": "+err.Error())
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

// UnassignQuizFromStudentHandler handles DELETE /instructor/assigned-quizzes/unassign
func (h *Handler) UnassignQuizFromStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quizIDStr := r.URL.Query().Get("quiz_id")
	userIDStr := r.URL.Query().Get("user_id")

	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil || quizID <= 0 {
		http.Error(w, "Invalid quiz_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	if err := h.db.UnassignQuizFromUser(quizID, userID); err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Assignment not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetAssignedQuizzesHandler handles GET /assigned-quizzes?user_id={id}
// Returns the list of quizzes assigned to a student.
func (h *Handler) GetAssignedQuizzesHandler(w http.ResponseWriter, r *http.Request) {
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

	quizzes, err := h.db.GetAssignedQuizzesForUser(userID)
	if err != nil {
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get assigned quizzes for user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if quizzes == nil {
		quizzes = []*models.AssignedQuizInfo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizzes)
}

// TakeAssignedQuizHandler handles POST /assigned-quizzes/start?quiz_id={id}&user_id={id}
// Generates a quiz attempt for an assigned quiz with randomized answer order.
func (h *Handler) TakeAssignedQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quizIDStr := r.URL.Query().Get("quiz_id")
	userIDStr := r.URL.Query().Get("user_id")

	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil || quizID <= 0 {
		http.Error(w, "Invalid quiz_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Verify the quiz is assigned to this user
	assignedQuizzes, err := h.db.GetAssignedQuizzesForUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var assignedQuiz *models.AssignedQuizInfo
	for _, aq := range assignedQuizzes {
		if aq.QuizID == quizID {
			assignedQuiz = aq
			break
		}
	}

	if assignedQuiz == nil {
		http.Error(w, "Quiz not assigned to this user", http.StatusForbidden)
		return
	}

	// Get questions for this quiz
	questions, err := h.db.GetQuizQuestions(quizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(questions) == 0 {
		http.Error(w, "No questions found for this quiz", http.StatusNotFound)
		return
	}

	// Build question IDs list
	questionIDs := make([]string, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.TopicQuestionID
	}

	// Create quiz attempt
	quizAttempt := &models.QuizAttempt{
		UserID:          userID,
		QuizID:          &quizID,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now(),
		TotalQuestions:  len(questions),
		CorrectAnswers:  0,
		Description:     assignedQuiz.QuizName,
		QuestionIDs:     questionIDs,
		QuestionResults: make(map[string]bool),
	}

	if err := h.db.CreateQuizAttempt(quizAttempt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with questions and randomized answers
	quizQuestions := make([]QuizQuestion, len(questions))
	for i, q := range questions {
		answers, err := h.db.GetAnswersByQuestion(q.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Randomize answer order
		rand.Shuffle(len(answers), func(a, b int) {
			answers[a], answers[b] = answers[b], answers[a]
		})

		quizQuestions[i] = QuizQuestion{
			Question: q,
			Answers:  answers,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenerateQuizResponse{
		QuizAttemptID: quizAttempt.ID,
		Questions:     quizQuestions,
	})
}

// GetAssignedQuizStatusHandler handles GET /instructor/assigned-quizzes/status?quiz_id={id}
// Returns results for all students assigned to a specific quiz.
func (h *Handler) GetAssignedQuizStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quizIDStr := r.URL.Query().Get("quiz_id")
	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil || quizID <= 0 {
		http.Error(w, "Invalid quiz_id", http.StatusBadRequest)
		return
	}

	status, err := h.db.GetAssignedQuizStatus(quizID)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Quiz not found", http.StatusNotFound)
			return
		}
		h.logger.WithContext(r.Context()).WithError(err).Error("Failed to get assigned quiz status")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// ListAssignedQuizzesStatusHandler handles GET /instructor/assigned-quizzes/status
// Returns a summary of all assigned quizzes.
func (h *Handler) ListAssignedQuizzesStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	quizzes, err := h.db.ListQuizzes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For each quiz, get the number of assigned students and completions
	type QuizSummary struct {
		QuizID          int    `json:"quiz_id"`
		QuizName        string `json:"quiz_name"`
		QuizCategory    string `json:"quiz_category"`
		QuizDescription string `json:"quiz_description"`
		AssignedCount   int    `json:"assigned_count"`
		CompletedCount  int    `json:"completed_count"`
	}

	var summaries []QuizSummary
	for _, quiz := range quizzes {
		userIDs, err := h.db.GetAssignedQuizUserIDs(quiz.ID)
		if err != nil {
			continue
		}
		if len(userIDs) == 0 {
			continue
		}

		status, err := h.db.GetAssignedQuizStatus(quiz.ID)
		if err != nil {
			continue
		}

		completedCount := 0
		for _, r := range status.StudentResults {
			if r.CompletedAt != nil {
				completedCount++
			}
		}

		summaries = append(summaries, QuizSummary{
			QuizID:          quiz.ID,
			QuizName:        quiz.Name,
			QuizCategory:    quiz.Category,
			QuizDescription: quiz.Description,
			AssignedCount:   len(userIDs),
			CompletedCount:  completedCount,
		})
	}

	if summaries == nil {
		summaries = []QuizSummary{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}
