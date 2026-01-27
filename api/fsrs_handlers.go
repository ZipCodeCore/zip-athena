package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/fsrs"
	"github.com/ZipCodeCore/zip-athena/models"
)

// ReviewQuestionRequest represents a request to review a question with FSRS
type ReviewQuestionRequest struct {
	UserID     int `json:"user_id"`
	QuestionID int `json:"question_id"`
	Rating     int `json:"rating"` // 1=Again, 2=Hard, 3=Good, 4=Easy
}

// ReviewQuestionResponse contains the updated FSRS state after review
type ReviewQuestionResponse struct {
	State        *models.UserQuestionState `json:"state"`
	NextReview   time.Time                 `json:"next_review"`
	IntervalDays int                       `json:"interval_days"`
}

// GetDueQuestionsRequest for getting due questions
type GetDueQuestionsRequest struct {
	UserID  int `json:"user_id"`
	TopicID int `json:"topic_id,omitempty"`
	Limit   int `json:"limit,omitempty"`
}

// GetDueQuestionsResponse contains questions due for review
type GetDueQuestionsResponse struct {
	Questions []QuizQuestion `json:"questions"`
	Count     int            `json:"count"`
}

// ReviewQuestionHandler handles POST /fsrs/review
func (h *Handler) ReviewQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReviewQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate rating
	if req.Rating < fsrs.RatingAgain || req.Rating > fsrs.RatingEasy {
		http.Error(w, "Invalid rating. Must be 1-4", http.StatusBadRequest)
		return
	}

	// Get or create user question state
	state, err := h.db.GetUserQuestionState(req.UserID, req.QuestionID)
	if err != nil {
		if err == database.ErrNotFound {
			// Create new state
			state = &models.UserQuestionState{
				UserID:     req.UserID,
				QuestionID: req.QuestionID,
				Due:        time.Now(),
				State:      0, // New
			}
			if err := h.db.CreateUserQuestionState(state); err != nil {
				http.Error(w, "Failed to create question state", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
	switch req.Rating {
	case fsrs.RatingAgain:
		schedulingInfo = schedulingCards.Again
	case fsrs.RatingHard:
		schedulingInfo = schedulingCards.Hard
	case fsrs.RatingGood:
		schedulingInfo = schedulingCards.Good
	case fsrs.RatingEasy:
		schedulingInfo = schedulingCards.Easy
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
		http.Error(w, "Failed to update question state", http.StatusInternalServerError)
		return
	}

	response := ReviewQuestionResponse{
		State:        state,
		NextReview:   state.Due,
		IntervalDays: state.ScheduledDays,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDueQuestionsHandler handles GET /fsrs/due?user_id={id}
func (h *Handler) GetDueQuestionsHandler(w http.ResponseWriter, r *http.Request) {
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

	topicIDStr := r.URL.Query().Get("topic_id")
	var topicID int
	if topicIDStr != "" {
		topicID, err = strconv.Atoi(topicIDStr)
		if err != nil {
			http.Error(w, "Invalid topic ID", http.StatusBadRequest)
			return
		}
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusBadRequest)
			return
		}
	}

	// Get due states
	states, err := h.db.GetDueUserQuestionStates(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build quiz questions from due states
	quizQuestions := make([]QuizQuestion, 0)
	for i, state := range states {
		if limit > 0 && i >= limit {
			break
		}

		question, err := h.db.GetQuestion(state.QuestionID)
		if err != nil {
			continue
		}

		// Filter by topic if specified
		if topicID > 0 && question.TopicID != topicID {
			continue
		}

		answers, err := h.db.GetAnswersByQuestion(question.ID)
		if err != nil {
			continue
		}

		quizQuestions = append(quizQuestions, QuizQuestion{
			Question: question,
			Answers:  answers,
		})
	}

	response := GetDueQuestionsResponse{
		Questions: quizQuestions,
		Count:     len(quizQuestions),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GenerateAdaptiveQuizHandler handles POST /fsrs/quiz/generate
// Generates a quiz using FSRS to prioritize due questions
func (h *Handler) GenerateAdaptiveQuizHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GenerateQuizRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.TopicID == 0 {
		http.Error(w, "user_id and topic_id are required", http.StatusBadRequest)
		return
	}

	if req.NumQuestions == 0 {
		req.NumQuestions = 10 // Default to 10 questions
	}

	// Get all questions for the topic
	allQuestions, err := h.db.GetQuestionsByTopic(req.TopicID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(allQuestions) == 0 {
		http.Error(w, "No questions found for this topic", http.StatusNotFound)
		return
	}

	// Get user's question states
	userStates, err := h.db.GetUserQuestionStatesByUser(req.UserID)
	if err != nil {
		userStates = []*models.UserQuestionState{}
	}

	// Build a map of question ID to state
	stateMap := make(map[int]*models.UserQuestionState)
	for _, state := range userStates {
		stateMap[state.QuestionID] = state
	}

	// Prioritize questions:
	// 1. Due for review (have state and due <= now)
	// 2. New questions (no state)
	// 3. Future questions (have state and due > now)
	now := time.Now()
	dueQuestions := make([]*models.Question, 0)
	newQuestions := make([]*models.Question, 0)
	futureQuestions := make([]*models.Question, 0)

	for _, q := range allQuestions {
		if state, exists := stateMap[q.ID]; exists {
			if !state.Due.After(now) {
				dueQuestions = append(dueQuestions, q)
			} else {
				futureQuestions = append(futureQuestions, q)
			}
		} else {
			newQuestions = append(newQuestions, q)
		}
	}

	// Select questions: prioritize due, then new, then future
	selectedQuestions := make([]*models.Question, 0, req.NumQuestions)
	selectedQuestions = append(selectedQuestions, dueQuestions...)
	if len(selectedQuestions) < req.NumQuestions {
		selectedQuestions = append(selectedQuestions, newQuestions...)
	}
	if len(selectedQuestions) < req.NumQuestions {
		selectedQuestions = append(selectedQuestions, futureQuestions...)
	}

	// Trim to requested number
	if len(selectedQuestions) > req.NumQuestions {
		selectedQuestions = selectedQuestions[:req.NumQuestions]
	}

	// Create quiz attempt
	questionIDs := make([]string, len(selectedQuestions))
	for i, q := range selectedQuestions {
		questionIDs[i] = q.TopicQuestionID
	}

	quizAttempt := &models.QuizAttempt{
		UserID:          req.UserID,
		StartTime:       time.Now(),
		TotalQuestions:  len(selectedQuestions),
		CorrectAnswers:  0,
		QuestionIDs:     questionIDs,
		QuestionResults: make(map[string]bool),
	}

	if err := h.db.CreateQuizAttempt(quizAttempt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build response with questions and answers
	quizQuestions := make([]QuizQuestion, len(selectedQuestions))
	for i, q := range selectedQuestions {
		answers, err := h.db.GetAnswersByQuestion(q.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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
