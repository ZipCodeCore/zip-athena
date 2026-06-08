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

func TestSubmitAnswerHandlerCountsQuestionOnlyOnce(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	user := &models.User{
		Username: "student",
		Email:    "student@example.com",
		Password: "password",
		Category: "java",
	}
	if err := db.CreateUser(user); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	topic := &models.Topic{
		Name:        "Java Basics",
		Description: "test topic",
		Category:    "java",
		CreatedAt:   time.Now(),
	}
	if err := db.CreateTopic(topic); err != nil {
		t.Fatalf("failed to create topic: %v", err)
	}

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "What is Java?",
		Difficulty:      1,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	if err := db.CreateQuestion(question); err != nil {
		t.Fatalf("failed to create question: %v", err)
	}

	correctAnswer := &models.Answer{
		QuestionID: question.ID,
		AnswerText: "A programming language",
		IsCorrect:  true,
	}
	if err := db.CreateAnswer(correctAnswer); err != nil {
		t.Fatalf("failed to create correct answer: %v", err)
	}

	wrongAnswer := &models.Answer{
		QuestionID: question.ID,
		AnswerText: "A database",
		IsCorrect:  false,
	}
	if err := db.CreateAnswer(wrongAnswer); err != nil {
		t.Fatalf("failed to create wrong answer: %v", err)
	}

	attempt := &models.QuizAttempt{
		UserID:          user.ID,
		StartTime:       time.Now(),
		TotalQuestions:  1,
		QuestionIDs:     []string{question.TopicQuestionID},
		QuestionResults: make(map[string]bool),
		AttemptNumber:   1,
	}
	if err := db.CreateQuizAttempt(attempt); err != nil {
		t.Fatalf("failed to create quiz attempt: %v", err)
	}

	submit := func(answerID int) SubmitAnswerResponse {
		reqBody := SubmitAnswerRequest{
			QuizAttemptID: attempt.ID,
			QuestionID:    question.ID,
			AnswerID:      answerID,
			UserID:        user.ID,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/quiz/submit-answer", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.SubmitAnswerHandler(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d. body: %s", w.Code, w.Body.String())
		}

		var resp SubmitAnswerResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		return resp
	}

	first := submit(correctAnswer.ID)
	if !first.Correct {
		t.Fatal("expected first submission to be correct")
	}

	second := submit(correctAnswer.ID)
	if !second.Correct {
		t.Fatal("expected second submission to be correct")
	}

	updatedAttempt, err := db.GetQuizAttempt(attempt.ID)
	if err != nil {
		t.Fatalf("failed to fetch updated attempt: %v", err)
	}

	if updatedAttempt.CorrectAnswers != 1 {
		t.Fatalf("expected correct_answers to remain 1 after duplicate correct submissions, got %d", updatedAttempt.CorrectAnswers)
	}
}
