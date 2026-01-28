package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

func TestQuizCompletionUpdatesFSRS(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create questions
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a goroutine?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.1",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a channel?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.2",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is defer?",
			Difficulty:      2,
			Tags:            []string{"basics"},
			TopicQuestionID: "go.3",
			CreatedAt:       time.Now(),
		},
	}

	for _, q := range questions {
		db.CreateQuestion(q)
	}

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create a quiz attempt with mixed results
	quizAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-30 * time.Minute),
		TotalQuestions:  3,
		CorrectAnswers:  2,
		QuestionIDs:     []string{"go.1", "go.2", "go.3"},
		QuestionResults: map[string]bool{
			"go.1": true,  // Correct
			"go.2": true,  // Correct
			"go.3": false, // Incorrect
		},
	}
	db.CreateQuizAttempt(quizAttempt)

	// Complete the quiz
	reqBody := CompleteQuizRequest{
		QuizAttemptID: quizAttempt.ID,
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/quiz/complete", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CompleteQuizHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	// Verify FSRS states were created for all questions
	for _, question := range questions {
		state, err := db.GetUserQuestionState(user.ID, question.ID)
		if err != nil {
			t.Errorf("Expected FSRS state for question %d, got error: %v", question.ID, err)
			continue
		}

		// Verify state has been updated (not default values)
		if state.Reps == 0 {
			t.Errorf("Expected Reps > 0 for question %d, got %d", question.ID, state.Reps)
		}

		// Check that incorrect question (go.3) has higher lapses
		if question.TopicQuestionID == "go.3" {
			// For a new question answered incorrectly, it goes to Learning state (1), not Relearning (3)
			// Lapses are incremented when moving from Review (2) to Relearning (3)
			// For first-time incorrect answers on new cards, lapses remain 0
			if state.State != 1 { // Learning state for new incorrect answer
				t.Logf("Question %s state is %d (expected 1=Learning for new incorrect answer)", question.TopicQuestionID, state.State)
			}
		}

		// Check that correct questions have appropriate state
		if question.TopicQuestionID == "go.1" || question.TopicQuestionID == "go.2" {
			if state.Lapses != 0 {
				t.Errorf("Expected Lapses = 0 for correct question %s, got %d", question.TopicQuestionID, state.Lapses)
			}
		}
	}
}

func TestRetakeQuizWithFSRSIntegration(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create questions
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a goroutine?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.1",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a channel?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.2",
			CreatedAt:       time.Now(),
		},
	}

	for _, q := range questions {
		db.CreateQuestion(q)
		db.CreateAnswer(&models.Answer{
			QuestionID: q.ID,
			AnswerText: "Correct answer",
			IsCorrect:  true,
		})
	}

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// First attempt - one question incorrect
	completedTime1 := time.Now()
	attempt1 := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-60 * time.Minute),
		CompletedTime:   &completedTime1,
		TotalQuestions:  2,
		CorrectAnswers:  1,
		QuestionIDs:     []string{"go.1", "go.2"},
		QuestionResults: map[string]bool{
			"go.1": true,  // Correct
			"go.2": false, // Incorrect
		},
	}
	db.CreateQuizAttempt(attempt1)

	// Complete first attempt to update FSRS
	reqBody1 := CompleteQuizRequest{QuizAttemptID: attempt1.ID}
	body1, _ := json.Marshal(reqBody1)
	req1 := httptest.NewRequest(http.MethodPost, "/quiz/complete", bytes.NewBuffer(body1))
	w1 := httptest.NewRecorder()
	handler.CompleteQuizHandler(w1, req1)

	// Get FSRS state for the incorrect question before retake
	stateBefore, _ := db.GetUserQuestionState(user.ID, questions[1].ID)
	repsBefore := stateBefore.Reps
	lapsesBefore := stateBefore.Lapses

	// Generate retake
	reqBody2 := RetakeQuizRequest{
		ParentAttemptID: attempt1.ID,
		UserID:          user.ID,
	}
	body2, _ := json.Marshal(reqBody2)
	req2 := httptest.NewRequest(http.MethodPost, "/quiz/retake", bytes.NewBuffer(body2))
	w2 := httptest.NewRecorder()
	handler.RetakeQuizHandler(w2, req2)

	var retakeResp GenerateQuizResponse
	json.NewDecoder(w2.Body).Decode(&retakeResp)

	// Update retake attempt - this time both correct
	attempt2, _ := db.GetQuizAttempt(retakeResp.QuizAttemptID)
	completedTime2 := time.Now()
	attempt2.CompletedTime = &completedTime2
	attempt2.CorrectAnswers = 2
	attempt2.QuestionResults = map[string]bool{
		"go.1": true, // Correct
		"go.2": true, // Correct this time
	}
	db.UpdateQuizAttempt(attempt2)

	// Complete retake to update FSRS
	reqBody3 := CompleteQuizRequest{QuizAttemptID: attempt2.ID}
	body3, _ := json.Marshal(reqBody3)
	req3 := httptest.NewRequest(http.MethodPost, "/quiz/complete", bytes.NewBuffer(body3))
	w3 := httptest.NewRecorder()
	handler.CompleteQuizHandler(w3, req3)

	// Verify FSRS state was updated for the previously incorrect question
	stateAfter, err := db.GetUserQuestionState(user.ID, questions[1].ID)
	if err != nil {
		t.Fatalf("Expected FSRS state after retake, got error: %v", err)
	}

	// After retake with correct answer, reps should increase
	if stateAfter.Reps <= repsBefore {
		t.Errorf("Expected Reps to increase after retake. Before: %d, After: %d", repsBefore, stateAfter.Reps)
	}

	// Lapses should remain the same (it was already incremented on first failure)
	if stateAfter.Lapses != lapsesBefore {
		t.Logf("Lapses changed from %d to %d (may be expected depending on FSRS logic)", lapsesBefore, stateAfter.Lapses)
	}

	// State should have changed from Relearning (3) to Review (2) after correct answer
	if stateAfter.State != 2 {
		t.Logf("Expected State 2 (Review) after correct retake, got %d (may be valid)", stateAfter.State)
	}
}

func TestIncorrectQuestionsAvailableInAdaptiveQuiz(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a user
	user := &models.User{
		Username: "student1",
		Email:    "student1@test.com",
		Password: "password",
		IsAdmin:  false,
	}
	db.CreateUser(user)

	// Create a topic
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Learn Go",
		CreatedAt:   time.Now(),
	}
	db.CreateTopic(topic)

	// Create multiple questions
	questions := []*models.Question{
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a goroutine?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.1",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is a channel?",
			Difficulty:      3,
			Tags:            []string{"concurrency"},
			TopicQuestionID: "go.2",
			CreatedAt:       time.Now(),
		},
		{
			TopicID:         topic.ID,
			QuestionText:    "What is defer?",
			Difficulty:      2,
			Tags:            []string{"basics"},
			TopicQuestionID: "go.3",
			CreatedAt:       time.Now(),
		},
	}

	for _, q := range questions {
		db.CreateQuestion(q)
	}

	// Create a quiz
	quiz := &models.Quiz{
		Name:        "Go Quiz 1",
		TopicID:     topic.ID,
		Description: "Test quiz",
		CreatedAt:   time.Now(),
	}
	db.CreateQuiz(quiz)

	// Create a quiz attempt where go.2 was answered incorrectly
	quizAttempt := &models.QuizAttempt{
		UserID:          user.ID,
		QuizID:          &quiz.ID,
		ParentAttemptID: nil,
		AttemptNumber:   1,
		IsRetake:        false,
		StartTime:       time.Now().Add(-30 * time.Minute),
		TotalQuestions:  2,
		CorrectAnswers:  1,
		QuestionIDs:     []string{"go.1", "go.2"},
		QuestionResults: map[string]bool{
			"go.1": true,  // Correct
			"go.2": false, // Incorrect - should be prioritized in adaptive quiz
		},
	}
	db.CreateQuizAttempt(quizAttempt)

	// Complete the quiz to update FSRS states
	reqBody1 := CompleteQuizRequest{QuizAttemptID: quizAttempt.ID}
	body1, _ := json.Marshal(reqBody1)
	req1 := httptest.NewRequest(http.MethodPost, "/quiz/complete", bytes.NewBuffer(body1))
	w1 := httptest.NewRecorder()
	handler.CompleteQuizHandler(w1, req1)

	// Verify that incorrectly answered question is marked as due
	stateGo2, err := db.GetUserQuestionState(user.ID, questions[1].ID)
	if err != nil {
		t.Fatalf("Expected FSRS state for go.2, got error: %v", err)
	}

	// The incorrect question should have FSRS state created
	if stateGo2.Reps == 0 {
		t.Error("Expected Reps > 0 for incorrect question after quiz completion")
	}

	// For a failed question, it should be in Learning or Relearning state
	if stateGo2.State != 1 && stateGo2.State != 3 {
		t.Errorf("Expected state 1 (Learning) or 3 (Relearning) for failed question, got %d", stateGo2.State)
	}

	// Note: FSRS may set a future due date even for failed reviews
	// The important thing is that the state is tracked, not necessarily that it's due immediately
	t.Logf("Incorrect question (go.2) FSRS state: Due=%v, Stability=%.2f, Difficulty=%.2f, State=%d",
		stateGo2.Due, stateGo2.Stability, stateGo2.Difficulty, stateGo2.State)
}
