package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// Phase 3: Instructor Analytics Handlers

// ListStudentsHandler returns a list of all students with their metrics
func (h *Handler) ListStudentsHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	searchQuery := r.URL.Query().Get("q")
	sortBy := r.URL.Query().Get("sort")
	filterCohort := r.URL.Query().Get("cohort")
	
	students, err := h.db.ListStudentsWithMetrics(searchQuery, sortBy, filterCohort)
	if err != nil {
		log.Printf("ERROR: Failed to list students: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetStudentAnalyticsHandler returns comprehensive analytics for a specific student
func (h *Handler) GetStudentAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
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
	
	// Get student progress (same as GetStudentProgressHandler)
	progress, err := h.db.GetStudentProgress(userID)
	if err == database.ErrNotFound {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Failed to get student analytics (user_id: %d): %v\n", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GetStudentQuizHistoryHandler returns quiz attempt history for a student
func (h *Handler) GetStudentQuizHistoryHandler(w http.ResponseWriter, r *http.Request) {
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
	
	history, err := h.db.GetStudentQuizHistory(userID)
	if err == database.ErrNotFound {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Failed to get student quiz history (user_id: %d): %v\n", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// GetStudentWeakAreasHandler returns topics and tags where student is struggling
func (h *Handler) GetStudentWeakAreasHandler(w http.ResponseWriter, r *http.Request) {
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
	
	weakAreas, err := h.db.GetStudentWeakAreas(userID)
	if err == database.ErrNotFound {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Failed to get student weak areas (user_id: %d): %v\n", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weakAreas)
}

// GetStudentProgressHandler returns progress analytics for a specific student
func (h *Handler) GetStudentProgressHandler(w http.ResponseWriter, r *http.Request) {
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

	progress, err := h.db.GetStudentProgress(userID)
	if err == database.ErrNotFound {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Failed to get student progress (user_id: %d): %v\n", userID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(progress)
}

// GetCohortAnalyticsHandler returns analytics for all students
func (h *Handler) GetCohortAnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	analytics, err := h.db.GetCohortAnalytics()
	if err != nil {
		log.Printf("ERROR: Failed to get cohort analytics: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}

// Phase 3: Import/Export Handlers

// ExportQuizHandler exports a quiz with its questions and answers as JSON
func (h *Handler) ExportQuizHandler(w http.ResponseWriter, r *http.Request) {
	quizIDStr := r.URL.Query().Get("id")
	if quizIDStr == "" {
		http.Error(w, "quiz id is required", http.StatusBadRequest)
		return
	}

	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		http.Error(w, "Invalid quiz id", http.StatusBadRequest)
		return
	}

	// Get quiz
	quiz, err := h.db.GetQuiz(quizID)
	if err == database.ErrNotFound {
		http.Error(w, "Quiz not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Failed to get quiz (ID: %d): %v\n", quizID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get questions for this quiz
	questions, err := h.db.GetQuizQuestions(quizID)
	if err != nil {
		log.Printf("ERROR: Failed to get quiz questions (quiz_id: %d): %v\n", quizID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build export structure with answers
	questionsWithAnswers := make([]models.QuestionWithAnswers, 0, len(questions))
	for _, q := range questions {
		answers, err := h.db.GetAnswersByQuestion(q.ID)
		if err != nil {
			log.Printf("ERROR: Failed to get answers for quiz export (quiz_id: %d, question_id: %d): %v\n", quizID, q.ID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert answer pointers to values
		answerValues := make([]models.Answer, len(answers))
		for i, a := range answers {
			answerValues[i] = *a
		}

		questionsWithAnswers = append(questionsWithAnswers, models.QuestionWithAnswers{
			Question: *q,
			Answers:  answerValues,
		})
	}

	export := models.QuizExport{
		Quiz:      *quiz,
		Questions: questionsWithAnswers,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=\"quiz-"+quiz.Name+".json\"")
	json.NewEncoder(w).Encode(export)
}

// ImportQuizHandler imports a quiz with questions and answers from JSON
func (h *Handler) ImportQuizHandler(w http.ResponseWriter, r *http.Request) {
	var importData models.QuizImport
	if err := json.NewDecoder(r.Body).Decode(&importData); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if importData.Name == "" {
		http.Error(w, "Quiz name is required", http.StatusBadRequest)
		return
	}
	if importData.TopicID == 0 {
		http.Error(w, "Topic ID is required", http.StatusBadRequest)
		return
	}

	// Create the quiz
	quiz := &models.Quiz{
		Name:        importData.Name,
		TopicID:     importData.TopicID,
		Description: importData.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.db.CreateQuiz(quiz); err != nil {
		log.Printf("ERROR: Failed to create quiz during import: %v\n", err)
		http.Error(w, "Failed to create quiz: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Import questions and link them to the quiz
	for order, qwa := range importData.Questions {
		// Create question
		question := &models.Question{
			TopicID:         importData.TopicID,
			QuestionText:    qwa.Question.QuestionText,
			Difficulty:      qwa.Question.Difficulty,
			Tags:            qwa.Question.Tags,
			TopicQuestionID: qwa.Question.TopicQuestionID,
			CreatedAt:       time.Now(),
		}

		if err := h.db.CreateQuestion(question); err != nil {
			log.Printf("ERROR: Failed to create question during quiz import (quiz_id: %d): %v\n", quiz.ID, err)
			http.Error(w, "Failed to create question: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Add question to quiz
		if err := h.db.AddQuestionToQuiz(quiz.ID, question.ID, order); err != nil {
			log.Printf("ERROR: Failed to link question to quiz (quiz_id: %d, question_id: %d): %v\n", quiz.ID, question.ID, err)
			http.Error(w, "Failed to link question to quiz: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create answers
		for _, answer := range qwa.Answers {
			a := &models.Answer{
				QuestionID: question.ID,
				AnswerText: answer.AnswerText,
				IsCorrect:  answer.IsCorrect,
			}
			if err := h.db.CreateAnswer(a); err != nil {
				log.Printf("ERROR: Failed to create answer during quiz import (question_id: %d): %v\n", question.ID, err)
				http.Error(w, "Failed to create answer: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "Quiz imported successfully",
		"quiz_id":      quiz.ID,
		"quiz_name":    quiz.Name,
		"question_count": len(importData.Questions),
	})
}

// ExportQuestionsHandler exports all questions with answers as JSON
func (h *Handler) ExportQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all questions
	questions, err := h.db.ListQuestions()
	if err != nil {
		log.Printf("ERROR: Failed to list questions for export: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build export structure with answers
	questionsWithAnswers := make([]models.QuestionWithAnswers, 0, len(questions))
	for _, q := range questions {
		answers, err := h.db.GetAnswersByQuestion(q.ID)
		if err != nil {
			log.Printf("ERROR: Failed to get answers (question_id: %d) for questions export: %v\n", q.ID, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert answer pointers to values
		answerValues := make([]models.Answer, len(answers))
		for i, a := range answers {
			answerValues[i] = *a
		}

		questionsWithAnswers = append(questionsWithAnswers, models.QuestionWithAnswers{
			Question: *q,
			Answers:  answerValues,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=\"questions.json\"")
	json.NewEncoder(w).Encode(questionsWithAnswers)
}

// ImportQuestionsHandler imports questions with answers from JSON
func (h *Handler) ImportQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	var questionsWithAnswers []models.QuestionWithAnswers
	if err := json.NewDecoder(r.Body).Decode(&questionsWithAnswers); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	importedCount := 0
	for _, qwa := range questionsWithAnswers {
		// Create question
		question := &models.Question{
			TopicID:         qwa.Question.TopicID,
			QuestionText:    qwa.Question.QuestionText,
			Difficulty:      qwa.Question.Difficulty,
			Tags:            qwa.Question.Tags,
			TopicQuestionID: qwa.Question.TopicQuestionID,
			CreatedAt:       time.Now(),
		}

		if err := h.db.CreateQuestion(question); err != nil {
			log.Printf("ERROR: Failed to create question during questions import: %v\n", err)
			http.Error(w, "Failed to create question: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Create answers
		for _, answer := range qwa.Answers {
			a := &models.Answer{
				QuestionID: question.ID,
				AnswerText: answer.AnswerText,
				IsCorrect:  answer.IsCorrect,
			}
			if err := h.db.CreateAnswer(a); err != nil {
				log.Printf("ERROR: Failed to create answer during questions import (question_id: %d): %v\n", question.ID, err)
				http.Error(w, "Failed to create answer: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		importedCount++
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Questions imported successfully",
		"count":   importedCount,
	})
}

// Phase 3: Test Management Handlers

// CreateTestHandler creates a new test
func (h *Handler) CreateTestHandler(w http.ResponseWriter, r *http.Request) {
	var test models.Test
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()

	if err := h.db.CreateTest(&test); err != nil {
		log.Printf("ERROR: Failed to create test: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(test)
}

// GetTestHandler retrieves a test by ID
func (h *Handler) GetTestHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "test id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	test, err := h.db.GetTest(id)
	if err == database.ErrNotFound {
		http.Error(w, "Test not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("ERROR: Failed to get test (ID: %d): %v\n", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}

// ListTestsHandler lists all tests
func (h *Handler) ListTestsHandler(w http.ResponseWriter, r *http.Request) {
	tests, err := h.db.ListTests()
	if err != nil {
		log.Printf("ERROR: Failed to list tests: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tests)
}

// UpdateTestHandler updates a test
func (h *Handler) UpdateTestHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "test id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	var test models.Test
	if err := json.NewDecoder(r.Body).Decode(&test); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	test.ID = id
	test.UpdatedAt = time.Now()

	if err := h.db.UpdateTest(&test); err == database.ErrNotFound {
		http.Error(w, "Test not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("ERROR: Failed to update test (ID: %d): %v\n", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(test)
}

// DeleteTestHandler deletes a test
func (h *Handler) DeleteTestHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "test id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid test id", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTest(id); err == database.ErrNotFound {
		http.Error(w, "Test not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("ERROR: Failed to delete test (ID: %d): %v\n", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddQuestionToQuizHandler adds a question to a quiz
func (h *Handler) AddQuestionToQuizHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		QuizID     int `json:"quiz_id"`
		QuestionID int `json:"question_id"`
		Order      int `json:"order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate IDs are positive
	if req.QuizID <= 0 || req.QuestionID <= 0 {
		http.Error(w, "quiz_id and question_id must be positive integers", http.StatusBadRequest)
		return
	}

	if err := h.db.AddQuestionToQuiz(req.QuizID, req.QuestionID, req.Order); err != nil {
		log.Printf("ERROR: Failed to add question to quiz (quiz_id: %d, question_id: %d): %v\n", req.QuizID, req.QuestionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Question added to quiz successfully"})
}

// RemoveQuestionFromQuizHandler removes a question from a quiz
func (h *Handler) RemoveQuestionFromQuizHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		QuizID     int `json:"quiz_id"`
		QuestionID int `json:"question_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate IDs are positive
	if req.QuizID <= 0 || req.QuestionID <= 0 {
		http.Error(w, "quiz_id and question_id must be positive integers", http.StatusBadRequest)
		return
	}

	if err := h.db.RemoveQuestionFromQuiz(req.QuizID, req.QuestionID); err != nil {
		log.Printf("ERROR: Failed to remove question from quiz (quiz_id: %d, question_id: %d): %v\n", req.QuizID, req.QuestionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetQuizQuestionsHandler gets all questions for a quiz
func (h *Handler) GetQuizQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	quizIDStr := r.URL.Query().Get("quiz_id")
	if quizIDStr == "" {
		http.Error(w, "quiz_id is required", http.StatusBadRequest)
		return
	}

	quizID, err := strconv.Atoi(quizIDStr)
	if err != nil {
		http.Error(w, "Invalid quiz_id", http.StatusBadRequest)
		return
	}

	questions, err := h.db.GetQuizQuestions(quizID)
	if err != nil {
		log.Printf("ERROR: Failed to get questions for quiz (quiz_id: %d): %v\n", quizID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// AddQuestionToTestHandler adds a question to a test
func (h *Handler) AddQuestionToTestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TestID     int `json:"test_id"`
		QuestionID int `json:"question_id"`
		Order      int `json:"order"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate IDs are positive
	if req.TestID <= 0 || req.QuestionID <= 0 {
		http.Error(w, "test_id and question_id must be positive integers", http.StatusBadRequest)
		return
	}

	if err := h.db.AddQuestionToTest(req.TestID, req.QuestionID, req.Order); err != nil {
		log.Printf("ERROR: Failed to add question to test (test_id: %d, question_id: %d): %v\n", req.TestID, req.QuestionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Question added to test successfully"})
}

// RemoveQuestionFromTestHandler removes a question from a test
func (h *Handler) RemoveQuestionFromTestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TestID     int `json:"test_id"`
		QuestionID int `json:"question_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate IDs are positive
	if req.TestID <= 0 || req.QuestionID <= 0 {
		http.Error(w, "test_id and question_id must be positive integers", http.StatusBadRequest)
		return
	}

	if err := h.db.RemoveQuestionFromTest(req.TestID, req.QuestionID); err != nil {
		log.Printf("ERROR: Failed to remove question from test (test_id: %d, question_id: %d): %v\n", req.TestID, req.QuestionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTestQuestionsHandler gets all questions for a test
func (h *Handler) GetTestQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	testIDStr := r.URL.Query().Get("test_id")
	if testIDStr == "" {
		http.Error(w, "test_id is required", http.StatusBadRequest)
		return
	}

	testID, err := strconv.Atoi(testIDStr)
	if err != nil {
		http.Error(w, "Invalid test_id", http.StatusBadRequest)
		return
	}

	questions, err := h.db.GetTestQuestions(testID)
	if err != nil {
		log.Printf("ERROR: Failed to get questions for test (test_id: %d): %v\n", testID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}
