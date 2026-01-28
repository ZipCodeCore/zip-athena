package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/fsrs"
	"github.com/ZipCodeCore/zip-athena/models"
)

// QuizQuestion represents a question with its answers for a quiz
type QuizQuestion struct {
	Question *models.Question `json:"question"`
	Answers  []*models.Answer `json:"answers"`
}

// GenerateQuizRequest represents the request to generate a quiz
type GenerateQuizRequest struct {
	UserID      int `json:"user_id"`
	TopicID     int `json:"topic_id"`
	NumQuestions int `json:"num_questions"`
}

// GenerateQuizResponse represents the response containing quiz questions
type GenerateQuizResponse struct {
	QuizAttemptID int            `json:"quiz_attempt_id"`
	Questions     []QuizQuestion `json:"questions"`
}

// SubmitAnswerRequest represents a single answer submission
type SubmitAnswerRequest struct {
	QuizAttemptID int `json:"quiz_attempt_id"`
	QuestionID    int `json:"question_id"`
	AnswerID      int `json:"answer_id"`
	UserID        int `json:"user_id"`
}

// SubmitAnswerResponse represents the response for an answer submission
type SubmitAnswerResponse struct {
	Correct       bool   `json:"correct"`
	CorrectAnswer string `json:"correct_answer,omitempty"`
}

// CompleteQuizRequest represents completing a quiz
type CompleteQuizRequest struct {
	QuizAttemptID int `json:"quiz_attempt_id"`
}

// RetakeQuizRequest represents requesting a quiz retake
type RetakeQuizRequest struct {
	ParentAttemptID int `json:"parent_attempt_id"`
	UserID          int `json:"user_id"`
}

// GenerateQuizHandler handles POST /quiz/generate
func (h *Handler) GenerateQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate request
	if req.UserID == 0 || req.TopicID == 0 {
		http.Error(w, "user_id and topic_id are required", http.StatusBadRequest)
		return
	}

	if req.NumQuestions == 0 {
		req.NumQuestions = 10 // Default to 10 questions
	}

	// Get topic to use in quiz name
	topic, err := h.db.GetTopic(req.TopicID)
	if err != nil {
		http.Error(w, "Topic not found", http.StatusNotFound)
		return
	}

	// Get questions for the topic
	questions, err := h.db.GetQuestionsByTopic(req.TopicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(questions) == 0 {
		http.Error(w, "No questions found for this topic", http.StatusNotFound)
		return
	}

	// Shuffle and limit questions
	rand.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})

	if len(questions) > req.NumQuestions {
		questions = questions[:req.NumQuestions]
	}

	// Create or use existing quiz for this topic
	quiz := &models.Quiz{
		Name:        fmt.Sprintf("%s Quiz", topic.Name),
		TopicID:     req.TopicID,
		Description: fmt.Sprintf("Auto-generated quiz for %s", topic.Name),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	if err := h.db.CreateQuiz(quiz); err != nil {
		// If quiz creation fails, it might already exist, but continue anyway
		// We'll just not set the quiz_id
	}

	// Create quiz attempt
	questionIDs := make([]string, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.TopicQuestionID
	}

	quizAttempt := &models.QuizAttempt{
		UserID:          req.UserID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now(),
		TotalQuestions:  len(questions),
		CorrectAnswers:  0,
		QuestionIDs:     questionIDs,
		QuestionResults: make(map[string]bool),
	}

	if err := h.db.CreateQuizAttempt(quizAttempt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with questions and answers
	quizQuestions := make([]QuizQuestion, len(questions))
	for i, q := range questions {
		answers, err := h.db.GetAnswersByQuestion(q.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Shuffle answers
		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
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

// SubmitAnswerHandler handles POST /quiz/submit-answer
func (h *Handler) SubmitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the quiz attempt
	quizAttempt, err := h.db.GetQuizAttempt(req.QuizAttemptID)
	if err != nil {
		http.Error(w, "Quiz attempt not found", http.StatusNotFound)
		return
	}

	// Verify user
	if quizAttempt.UserID != req.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the question
	question, err := h.db.GetQuestion(req.QuestionID)
	if err != nil {
		http.Error(w, "Question not found", http.StatusNotFound)
		return
	}

	// Get all answers for the question
	answers, err := h.db.GetAnswersByQuestion(req.QuestionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the submitted answer is correct
	var correct bool
	var correctAnswerText string
	for _, ans := range answers {
		if ans.IsCorrect {
			correctAnswerText = ans.AnswerText
		}
		if ans.ID == req.AnswerID && ans.IsCorrect {
			correct = true
		}
	}

	// Update quiz attempt results
	if quizAttempt.QuestionResults == nil {
		quizAttempt.QuestionResults = make(map[string]bool)
	}
	quizAttempt.QuestionResults[question.TopicQuestionID] = correct

	if correct {
		quizAttempt.CorrectAnswers++
	}

	// Save quiz attempt updates to database
	if err := h.db.UpdateQuizAttempt(quizAttempt); err != nil {
		http.Error(w, "Failed to update quiz attempt", http.StatusInternalServerError)
		return
	}

	// Update or create question attempt
	questionAttempts, err := h.db.GetQuestionAttemptsByUser(req.UserID)
	if err != nil {
		questionAttempts = []*models.QuestionAttempt{}
	}

	var existingAttempt *models.QuestionAttempt
	for _, qa := range questionAttempts {
		if qa.QuestionID == req.QuestionID {
			existingAttempt = qa
			break
		}
	}

	if existingAttempt != nil {
		if correct {
			existingAttempt.CorrectAttempts++
		} else {
			existingAttempt.IncorrectAttempts++
		}
		existingAttempt.LastAttemptTime = time.Now()
		existingAttempt.LastAttemptCorrect = correct
		
		// Update existing attempt in database
		if err := h.db.UpdateQuestionAttempt(existingAttempt); err != nil {
			// Log error but don't fail the request
			fmt.Printf("Warning: Failed to update question attempt: %v\n", err)
		}
	} else {
		newAttempt := &models.QuestionAttempt{
			UserID:             req.UserID,
			QuestionID:         req.QuestionID,
			TopicQuestionID:    question.TopicQuestionID,
			CorrectAttempts:    0,
			IncorrectAttempts:  0,
			LastAttemptTime:    time.Now(),
			LastAttemptCorrect: correct,
		}
		if correct {
			newAttempt.CorrectAttempts = 1
		} else {
			newAttempt.IncorrectAttempts = 1
		}
		h.db.CreateQuestionAttempt(newAttempt)
	}

	response := SubmitAnswerResponse{
		Correct: correct,
	}
	if !correct {
		response.CorrectAnswer = correctAnswerText
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CompleteQuizHandler handles POST /quiz/complete
func (h *Handler) CompleteQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CompleteQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the quiz attempt
	quizAttempt, err := h.db.GetQuizAttempt(req.QuizAttemptID)
	if err != nil {
		http.Error(w, "Quiz attempt not found", http.StatusNotFound)
		return
	}

	// Mark as completed
	now := time.Now()
	quizAttempt.CompletedTime = &now
	
	// Save the updated quiz attempt
	if err := h.db.UpdateQuizAttempt(quizAttempt); err != nil {
		http.Error(w, "Failed to update quiz attempt", http.StatusInternalServerError)
		return
	}

	// Update user's total points
	user, err := h.db.GetUser(quizAttempt.UserID)
	if err == nil {
		user.TotalPoints += quizAttempt.CorrectAnswers
		h.db.UpdateUserPoints(user.ID, user.TotalPoints)
	}

	// Update FSRS state for each question in the quiz
	h.updateFSRSStatesForQuiz(quizAttempt)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizAttempt)
}

// updateFSRSStatesForQuiz updates FSRS state for all questions in a quiz attempt
func (h *Handler) updateFSRSStatesForQuiz(quizAttempt *models.QuizAttempt) {
	// Get all questions to map TopicQuestionID to QuestionID
	allQuestions, err := h.db.ListQuestions()
	if err != nil {
		// Log error but don't block quiz completion
		fmt.Printf("Warning: Failed to load questions for FSRS update: %v\n", err)
		return
	}

	questionMap := make(map[string]*models.Question)
	for _, q := range allQuestions {
		questionMap[q.TopicQuestionID] = q
	}

	// Track failures for logging
	var failedUpdates []string

	// Process each question in the quiz
	for topicQuestionID, wasCorrect := range quizAttempt.QuestionResults {
		question, exists := questionMap[topicQuestionID]
		if !exists {
			failedUpdates = append(failedUpdates, fmt.Sprintf("%s (question not found)", topicQuestionID))
			continue
		}

		// Determine FSRS rating based on correctness
		// Correct = Good (3), Incorrect = Again (1)
		rating := 1 // Again (incorrect)
		if wasCorrect {
			rating = 3 // Good (correct)
		}

		// Get or create user question state
		state, err := h.db.GetUserQuestionState(quizAttempt.UserID, question.ID)
		if err != nil {
			if err == database.ErrNotFound {
				// Create new state for this question
				state = &models.UserQuestionState{
					UserID:     quizAttempt.UserID,
					QuestionID: question.ID,
					Due:        time.Now(),
					State:      0, // New
				}
				if err := h.db.CreateUserQuestionState(state); err != nil {
					failedUpdates = append(failedUpdates, fmt.Sprintf("%s (create state failed: %v)", topicQuestionID, err))
					continue
				}
			} else {
				failedUpdates = append(failedUpdates, fmt.Sprintf("%s (get state failed: %v)", topicQuestionID, err))
				continue
			}
		}

		// Convert to FSRS card
		card := fsrs.Card{
			Due:           state.Due,
			Stability:     state.Stability,
			Difficulty:    state.Difficulty,
			ElapsedDays:   state.ElapsedDays,
			ScheduledDays: state.ScheduledDays,
			Reps:          state.Reps,
			Lapses:        state.Lapses,
			State:         state.State,
			LastReview:    state.LastReview,
		}

		// Apply FSRS algorithm
		params := fsrs.NewDefaultParameters()
		schedulingCards := params.Repeat(card, time.Now())

		// Get the scheduling info based on rating
		var schedulingInfo fsrs.SchedulingInfo
		if rating == 1 {
			schedulingInfo = schedulingCards.Again
		} else {
			schedulingInfo = schedulingCards.Good
		}

		// Update state with new values
		state.Due = schedulingInfo.Card.Due
		state.Stability = schedulingInfo.Card.Stability
		state.Difficulty = schedulingInfo.Card.Difficulty
		state.ElapsedDays = schedulingInfo.Card.ElapsedDays
		state.ScheduledDays = schedulingInfo.Card.ScheduledDays
		state.Reps = schedulingInfo.Card.Reps
		state.Lapses = schedulingInfo.Card.Lapses
		state.State = schedulingInfo.Card.State
		state.LastReview = schedulingInfo.Card.LastReview

		// Save updated state
		if err := h.db.UpdateUserQuestionState(state); err != nil {
			failedUpdates = append(failedUpdates, fmt.Sprintf("%s (update failed: %v)", topicQuestionID, err))
		}
	}

	// Log any failures
	if len(failedUpdates) > 0 {
		fmt.Printf("Warning: Failed to update FSRS for %d question(s) in quiz attempt %d:\n", len(failedUpdates), quizAttempt.ID)
		for _, failure := range failedUpdates {
			fmt.Printf("  - %s\n", failure)
		}
	}
}

// RetakeQuizHandler handles POST /quiz/retake
// Generates a new quiz attempt with the same questions as a previous attempt
func (h *Handler) RetakeQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RetakeQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the parent quiz attempt
	parentAttempt, err := h.db.GetQuizAttempt(req.ParentAttemptID)
	if err != nil {
		http.Error(w, "Parent quiz attempt not found", http.StatusNotFound)
		return
	}

	// Verify user
	if parentAttempt.UserID != req.UserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Check if parent quiz is completed
	if parentAttempt.CompletedTime == nil {
		http.Error(w, "Cannot retake an incomplete quiz", http.StatusBadRequest)
		return
	}

	// Check if student scored 100% (no retake needed)
	if parentAttempt.CorrectAnswers == parentAttempt.TotalQuestions {
		http.Error(w, "Cannot retake a quiz with 100% score", http.StatusBadRequest)
		return
	}

	// Get the questions from the parent attempt using their TopicQuestionIDs
	questions := make([]*models.Question, 0, len(parentAttempt.QuestionIDs))
	allQuestions, err := h.db.ListQuestions()
	if err != nil {
		http.Error(w, "Failed to retrieve questions", http.StatusInternalServerError)
		return
	}

	// Build a map of TopicQuestionID -> Question for efficient lookup
	questionMap := make(map[string]*models.Question)
	for _, q := range allQuestions {
		questionMap[q.TopicQuestionID] = q
	}

	// Reconstruct the question list from parent attempt
	for _, topicQuestionID := range parentAttempt.QuestionIDs {
		if q, exists := questionMap[topicQuestionID]; exists {
			questions = append(questions, q)
		}
	}

	if len(questions) == 0 {
		http.Error(w, "No questions found for retake", http.StatusNotFound)
		return
	}

	// Calculate attempt number
	attemptNumber := parentAttempt.AttemptNumber + 1

	// Create new quiz attempt as a retake
	retakeAttempt := &models.QuizAttempt{
		UserID:          req.UserID,
		QuizID:          parentAttempt.QuizID,
		ParentAttemptID: &req.ParentAttemptID,
		AttemptNumber:   attemptNumber,
		IsRetake:        true,
		StartTime:       time.Now(),
		TotalQuestions:  len(questions),
		CorrectAnswers:  0,
		QuestionIDs:     parentAttempt.QuestionIDs, // Use same questions as parent
		QuestionResults: make(map[string]bool),
	}

	if err := h.db.CreateQuizAttempt(retakeAttempt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with questions and answers
	quizQuestions := make([]QuizQuestion, len(questions))
	for i, q := range questions {
		answers, err := h.db.GetAnswersByQuestion(q.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Shuffle answers
		rand.Shuffle(len(answers), func(i, j int) {
			answers[i], answers[j] = answers[j], answers[i]
		})

		quizQuestions[i] = QuizQuestion{
			Question: q,
			Answers:  answers,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(GenerateQuizResponse{
		QuizAttemptID: retakeAttempt.ID,
		Questions:     quizQuestions,
	})
}

// GetProgressHandler handles GET /progress?user_id={id}
func (h *Handler) GetProgressHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get user
	user, err := h.db.GetUser(userID)
	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get quiz attempts
	quizAttempts, err := h.db.GetQuizAttemptsByUser(userID)
	if err != nil {
		quizAttempts = []*models.QuizAttempt{}
	}

	// Enrich quiz attempts with quiz information
	type QuizAttemptWithQuiz struct {
		*models.QuizAttempt
		QuizName        string `json:"quiz_name,omitempty"`
		QuizDescription string `json:"quiz_description,omitempty"`
	}

	enrichedAttempts := make([]QuizAttemptWithQuiz, 0, len(quizAttempts))
	for _, qa := range quizAttempts {
		enriched := QuizAttemptWithQuiz{
			QuizAttempt: qa,
		}
		
		// If the attempt has a quiz_id, fetch the quiz details
		if qa.QuizID != nil {
			quiz, err := h.db.GetQuiz(*qa.QuizID)
			if err == nil {
				enriched.QuizName = quiz.Name
				enriched.QuizDescription = quiz.Description
			}
		}
		
		enrichedAttempts = append(enrichedAttempts, enriched)
	}

	// Get question attempts
	questionAttempts, err := h.db.GetQuestionAttemptsByUser(userID)
	if err != nil {
		questionAttempts = []*models.QuestionAttempt{}
	}

	// Calculate statistics
	totalQuestions := 0
	totalCorrect := 0
	for _, qa := range quizAttempts {
		// Count questions answered from all quizzes (completed and in-progress)
		// For completed quizzes, all questions were answered
		// For in-progress quizzes, count only the questions that have been answered
		if qa.CompletedTime != nil {
			// Completed quiz - all questions were answered
			totalQuestions += qa.TotalQuestions
			totalCorrect += qa.CorrectAnswers
		} else {
			// In-progress quiz - count only answered questions (those in QuestionResults)
			answeredCount := len(qa.QuestionResults)
			totalQuestions += answeredCount
			totalCorrect += qa.CorrectAnswers
		}
	}

	progress := map[string]interface{}{
		"user":              user,
		"total_points":      user.TotalPoints,
		"total_quizzes":     len(quizAttempts),
		"total_questions":   totalQuestions,
		"total_correct":     totalCorrect,
		"quiz_attempts":     enrichedAttempts,
		"question_attempts": questionAttempts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
