package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
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

	// Create quiz attempt
	questionIDs := make([]string, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.TopicQuestionID
	}

	quizAttempt := &models.QuizAttempt{
		UserID:          req.UserID,
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

	// Update user's total points
	user, err := h.db.GetUser(quizAttempt.UserID)
	if err == nil {
		user.TotalPoints += quizAttempt.CorrectAnswers
		h.db.UpdateUserPoints(user.ID, user.TotalPoints)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quizAttempt)
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
		"quiz_attempts":     quizAttempts,
		"question_attempts": questionAttempts,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}
