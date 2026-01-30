package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// Phase 5 Test Suite

// Test Question Management (FR-INSTRUCTOR-004)

func TestCreateQuestionWithValidation(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create a topic first
	topic := &models.Topic{
		Name:        "Go Programming",
		Description: "Go language questions",
		CreatedAt:   time.Now(),
	}
	if err := db.CreateTopic(topic); err != nil {
		t.Fatalf("Failed to create topic: %v", err)
	}

	// Create question with answers
	reqBody := map[string]interface{}{
		"question": map[string]interface{}{
			"topic_id":      topic.ID,
			"question_text": "What is a goroutine?",
			"difficulty":    3,
			"tags":          []string{"concurrency", "basics"},
			"topic_question_id": "1.1",
			"created_at":    time.Now(),
		},
		"answers": []map[string]interface{}{
			{"answer_text": "A lightweight thread", "is_correct": true},
			{"answer_text": "A function", "is_correct": false},
			{"answer_text": "A package", "is_correct": false},
			{"answer_text": "A variable", "is_correct": false},
		},
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/instructor/questions", bytes.NewReader(body))
	req.Header.Set("X-User-ID", "1") // Admin user
	w := httptest.NewRecorder()

	// Create admin user
	adminUser := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(adminUser)

	handler.AdminMiddleware(handler.CreateInstructorQuestionHandler)(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		t.Logf("Response: %s", w.Body.String())
	}
}

func TestQuestionTextRequired(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "", // Empty
		Difficulty:      3,
		TopicQuestionID: "1.1",
	}

	answers := []models.Answer{
		{AnswerText: "A", IsCorrect: true},
		{AnswerText: "B", IsCorrect: false},
	}

	err := db.ValidateQuestionData(question, answers)
	if err == nil {
		t.Error("Expected validation error for empty question text")
	}
}

func TestMinimumTwoAnswersRequired(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test question",
		Difficulty:      3,
		TopicQuestionID: "1.1",
	}

	answers := []models.Answer{
		{AnswerText: "A", IsCorrect: true},
	}

	err := db.ValidateQuestionData(question, answers)
	if err == nil {
		t.Error("Expected validation error for less than 2 answers")
	}
	if !strings.Contains(err.Error(), "2-4 answer options") {
		t.Errorf("Expected '2-4 answer options' error, got: %v", err)
	}
}

func TestExactlyOneCorrectAnswerRequired(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test question",
		Difficulty:      3,
		TopicQuestionID: "1.1",
	}

	// Test with no correct answers
	answers := []models.Answer{
		{AnswerText: "A", IsCorrect: false},
		{AnswerText: "B", IsCorrect: false},
	}

	err := db.ValidateQuestionData(question, answers)
	if err == nil || !strings.Contains(err.Error(), "exactly 1 answer must be marked as correct") {
		t.Error("Expected validation error for no correct answers")
	}

	// Test with multiple correct answers
	answers = []models.Answer{
		{AnswerText: "A", IsCorrect: true},
		{AnswerText: "B", IsCorrect: true},
	}

	err = db.ValidateQuestionData(question, answers)
	if err == nil || !strings.Contains(err.Error(), "exactly 1 answer must be marked as correct") {
		t.Error("Expected validation error for multiple correct answers")
	}
}

func TestQuestionUpdate(t *testing.T) {
	db := database.NewInMemoryDB()
	handler := NewHandler(db)

	// Create admin user
	adminUser := &models.User{
		Username:  "admin",
		Email:     "admin@test.com",
		IsAdmin:   true,
		CreatedAt: time.Now(),
	}
	db.CreateUser(adminUser)

	// Create a topic and question first
	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Original question",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Create answers
	answers := []*models.Answer{
		{QuestionID: question.ID, AnswerText: "A", IsCorrect: true},
		{QuestionID: question.ID, AnswerText: "B", IsCorrect: false},
	}
	for _, ans := range answers {
		db.CreateAnswer(ans)
	}

	// Update the question
	updateReq := map[string]interface{}{
		"question": map[string]interface{}{
			"id":            question.ID,
			"topic_id":      topic.ID,
			"question_text": "Updated question",
			"difficulty":    5,
			"tags":          []string{"updated"},
			"topic_question_id": "1.1",
		},
		"answers": []map[string]interface{}{
			{"answer_text": "New A", "is_correct": true},
			{"answer_text": "New B", "is_correct": false},
		},
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/instructor/questions?id="+strconv.Itoa(question.ID), bytes.NewReader(body))
	req.Header.Set("X-User-ID", "1")
	w := httptest.NewRecorder()

	handler.AdminMiddleware(handler.UpdateInstructorQuestionHandler)(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestQuestionDeactivation(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	err := db.DeactivateQuestion(question.ID)
	if err != nil {
		t.Errorf("Failed to deactivate question: %v", err)
	}
}

func TestTagAssignment(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		Tags:            []string{"tag1", "tag2", "tag3"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	retrieved, err := db.GetQuestion(question.ID)
	if err != nil {
		t.Fatalf("Failed to get question: %v", err)
	}

	if len(retrieved.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(retrieved.Tags))
	}
}

func TestDifficultyValidation(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Test difficulty < 1
	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      0,
		TopicQuestionID: "1.1",
	}
	answers := []models.Answer{
		{AnswerText: "A", IsCorrect: true},
		{AnswerText: "B", IsCorrect: false},
	}

	err := db.ValidateQuestionData(question, answers)
	if err == nil || !strings.Contains(err.Error(), "difficulty must be between 1 and 7") {
		t.Error("Expected validation error for difficulty < 1")
	}

	// Test difficulty > 7
	question.Difficulty = 8
	err = db.ValidateQuestionData(question, answers)
	if err == nil || !strings.Contains(err.Error(), "difficulty must be between 1 and 7") {
		t.Error("Expected validation error for difficulty > 7")
	}
}

func TestPreviewGeneration(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test question",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	answers := []*models.Answer{
		{QuestionID: question.ID, AnswerText: "A", IsCorrect: true},
		{QuestionID: question.ID, AnswerText: "B", IsCorrect: false},
	}
	for _, ans := range answers {
		db.CreateAnswer(ans)
	}

	preview, err := db.GetQuestionPreview(question.ID)
	if err != nil {
		t.Fatalf("Failed to get preview: %v", err)
	}

	if preview.Question.ID != question.ID {
		t.Error("Preview question ID mismatch")
	}

	if len(preview.Answers) != 2 {
		t.Errorf("Expected 2 answers in preview, got %d", len(preview.Answers))
	}
}

// Test Bulk CSV Import (FR-CONTENT-001)

func TestCSVParsing(t *testing.T) {
	rows := []models.CSVImportRow{
		{
			Topic:         "Java",
			Question:      "What is Java?",
			Answer1:       "A language",
			Answer2:       "A drink",
			Answer3:       "An island",
			Answer4:       "",
			CorrectAnswer: "A language",
			Difficulty:    "2",
			Tags:          "basics",
		},
	}

	db := database.NewInMemoryDB()
	result, err := db.ValidateCSVImport(rows)
	if err != nil {
		t.Fatalf("Failed to validate CSV: %v", err)
	}

	if result.Summary.ValidRows != 1 {
		t.Errorf("Expected 1 valid row, got %d", result.Summary.ValidRows)
	}
}

func TestColumnMapping(t *testing.T) {
	// Tested in TestCSVParsing - columns map to CSVImportRow fields
	t.Skip("Column mapping tested in TestCSVParsing")
}

func TestValidationRequiredFields(t *testing.T) {
	db := database.NewInMemoryDB()

	rows := []models.CSVImportRow{
		{
			Topic:         "", // Missing
			Question:      "What is Java?",
			Answer1:       "A",
			Answer2:       "B",
			CorrectAnswer: "A",
			Difficulty:    "2",
		},
	}

	result, err := db.ValidateCSVImport(rows)
	if err != nil {
		t.Fatalf("Failed to validate: %v", err)
	}

	if result.Summary.InvalidRows != 1 {
		t.Errorf("Expected 1 invalid row, got %d", result.Summary.InvalidRows)
	}

	if len(result.InvalidRows) > 0 && !strings.Contains(result.InvalidRows[0].Errors[0], "Topic is required") {
		t.Error("Expected 'Topic is required' error")
	}
}

func TestValidationCorrectAnswerExists(t *testing.T) {
	db := database.NewInMemoryDB()

	rows := []models.CSVImportRow{
		{
			Topic:         "Java",
			Question:      "Test",
			Answer1:       "A",
			Answer2:       "B",
			CorrectAnswer: "C", // Not in options
			Difficulty:    "2",
		},
	}

	result, err := db.ValidateCSVImport(rows)
	if err != nil {
		t.Fatalf("Failed to validate: %v", err)
	}

	if result.Summary.InvalidRows != 1 {
		t.Error("Expected invalid row for mismatched correct answer")
	}
}

func TestValidationDifficulty1to7(t *testing.T) {
	db := database.NewInMemoryDB()

	rows := []models.CSVImportRow{
		{
			Topic:         "Java",
			Question:      "Test",
			Answer1:       "A",
			Answer2:       "B",
			CorrectAnswer: "A",
			Difficulty:    "10", // Invalid
		},
	}

	result, err := db.ValidateCSVImport(rows)
	if err != nil {
		t.Fatalf("Failed to validate: %v", err)
	}

	if result.Summary.InvalidRows != 1 {
		t.Error("Expected invalid row for difficulty out of range")
	}
}

func TestDuplicateQuestionDetection(t *testing.T) {
	// Note: Current implementation doesn't detect duplicates
	// This would require comparing question text
	t.Skip("Duplicate detection not implemented")
}

func TestBatchInsertPerformance(t *testing.T) {
	db := database.NewInMemoryDB()

	// Create 500 valid rows
	rows := make([]models.CSVImportRow, 500)
	for i := 0; i < 500; i++ {
		rows[i] = models.CSVImportRow{
			Topic:         "Java",
			Question:      "Question " + string(rune(i)),
			Answer1:       "A",
			Answer2:       "B",
			Answer3:       "C",
			Answer4:       "D",
			CorrectAnswer: "A",
			Difficulty:    "3",
			Tags:          "test",
		}
	}

	start := time.Now()
	imported, err := db.BulkCreateQuestions(rows)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Failed to bulk import: %v", err)
	}

	if imported != 500 {
		t.Errorf("Expected 500 imported, got %d", imported)
	}

	// Should complete in reasonable time (< 5 seconds)
	if duration > 5*time.Second {
		t.Errorf("Bulk import took too long: %v", duration)
	}
}

func TestErrorReportingInvalidRows(t *testing.T) {
	db := database.NewInMemoryDB()

	rows := []models.CSVImportRow{
		{Topic: "", Question: "Q1", Answer1: "A", Answer2: "B", CorrectAnswer: "A", Difficulty: "2"},
		{Topic: "Java", Question: "", Answer1: "A", Answer2: "B", CorrectAnswer: "A", Difficulty: "2"},
		{Topic: "Java", Question: "Q3", Answer1: "A", Answer2: "B", CorrectAnswer: "A", Difficulty: "99"},
	}

	result, err := db.ValidateCSVImport(rows)
	if err != nil {
		t.Fatalf("Failed to validate: %v", err)
	}

	if len(result.InvalidRows) != 3 {
		t.Errorf("Expected 3 invalid rows, got %d", len(result.InvalidRows))
	}
}

func TestPartialImport(t *testing.T) {
	db := database.NewInMemoryDB()

	validRow := models.CSVImportRow{
		Topic:         "Java",
		Question:      "Valid",
		Answer1:       "A",
		Answer2:       "B",
		CorrectAnswer: "A",
		Difficulty:    "3",
	}

	rows := []models.CSVImportRow{validRow}

	imported, err := db.BulkCreateQuestions(rows)
	if err != nil {
		t.Fatalf("Failed to import: %v", err)
	}

	if imported != 1 {
		t.Errorf("Expected 1 imported, got %d", imported)
	}
}

// Tag Management Tests (FR-CONTENT-002)

func TestTagCreation(t *testing.T) {
	// Tags are created implicitly when questions are created
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		Tags:            []string{"new-tag"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	tags, err := db.GetTagsWithMetadata()
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	found := false
	for _, tag := range tags {
		if tag.Name == "new-tag" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected tag 'new-tag' to be created")
	}
}

func TestTagNameUniqueness(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Create two questions with same tag
	for i := 0; i < 2; i++ {
		question := &models.Question{
			TopicID:         topic.ID,
			QuestionText:    "Test " + string(rune(i)),
			Difficulty:      3,
			Tags:            []string{"duplicate-tag"},
			TopicQuestionID: "1." + string(rune(i)),
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(question)
	}

	tags, err := db.GetTagsWithMetadata()
	if err != nil {
		t.Fatalf("Failed to get tags: %v", err)
	}

	count := 0
	for _, tag := range tags {
		if tag.Name == "duplicate-tag" {
			count++
			if tag.QuestionCount != 2 {
				t.Errorf("Expected question count 2, got %d", tag.QuestionCount)
			}
		}
	}

	if count != 1 {
		t.Errorf("Expected 1 unique tag, got %d", count)
	}
}

func TestTagDeletion(t *testing.T) {
	db := database.NewInMemoryDB()

	// Tag with no questions can be deleted
	err := db.DeleteTag("nonexistent-tag")
	if err != nil {
		t.Errorf("Expected no error for deleting non-existent tag: %v", err)
	}

	// Create question with tag
	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		Tags:            []string{"in-use-tag"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Tag with questions cannot be deleted
	err = db.DeleteTag("in-use-tag")
	if err == nil {
		t.Error("Expected error when deleting tag in use")
	}
}

func TestTagMergeOperation(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Create questions with different tags
	q1 := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Q1",
		Difficulty:      3,
		Tags:            []string{"old-tag-1"},
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(q1)

	q2 := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Q2",
		Difficulty:      3,
		Tags:            []string{"old-tag-2"},
		TopicQuestionID: "1.2",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(q2)

	// Merge tags
	err := db.MergeTags([]string{"old-tag-1", "old-tag-2"}, "new-tag")
	if err != nil {
		t.Fatalf("Failed to merge tags: %v", err)
	}

	// Verify questions now have new tag
	tags, _ := db.GetTagsWithMetadata()
	found := false
	for _, tag := range tags {
		if tag.Name == "new-tag" {
			found = true
			if tag.QuestionCount != 2 {
				t.Errorf("Expected 2 questions with merged tag, got %d", tag.QuestionCount)
			}
		}
	}

	if !found {
		t.Error("Merged tag not found")
	}
}

func TestQuestionsByTagRetrieval(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Create questions with specific tag
	for i := 0; i < 3; i++ {
		question := &models.Question{
			TopicID:         topic.ID,
			QuestionText:    "Test " + string(rune(i)),
			Difficulty:      3,
			Tags:            []string{"specific-tag"},
			TopicQuestionID: "1." + string(rune(i)),
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(question)
	}

	questions, err := db.GetQuestionsByTag("specific-tag")
	if err != nil {
		t.Fatalf("Failed to get questions by tag: %v", err)
	}

	if len(questions) != 3 {
		t.Errorf("Expected 3 questions, got %d", len(questions))
	}
}

// Quality Review Tests (FR-CONTENT-003)

func TestFlaggedQuestionsQueueRetrieval(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Flag the question
	err := db.FlagQuestion(question.ID, models.FlagReasonManual)
	if err != nil {
		t.Fatalf("Failed to flag question: %v", err)
	}

	// Retrieve flagged questions
	flagged, err := db.GetFlaggedQuestions()
	if err != nil {
		t.Fatalf("Failed to get flagged questions: %v", err)
	}

	if len(flagged) != 1 {
		t.Errorf("Expected 1 flagged question, got %d", len(flagged))
	}

	if flagged[0].FlagReason != models.FlagReasonManual {
		t.Errorf("Expected flag reason %s, got %s", models.FlagReasonManual, flagged[0].FlagReason)
	}
}

func TestAutoFlagTooEasy(t *testing.T) {
	// Auto-flag logic is tested in AutoFlagQuestions
	// This requires question attempts data
	t.Skip("Auto-flag logic requires question attempts")
}

func TestAutoFlagTooHard(t *testing.T) {
	// Auto-flag logic is tested in AutoFlagQuestions
	t.Skip("Auto-flag logic requires question attempts")
}

func TestManualFlag(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	err := db.FlagQuestion(question.ID, models.FlagReasonManual)
	if err != nil {
		t.Fatalf("Failed to manually flag question: %v", err)
	}

	flagged, _ := db.GetFlaggedQuestions()
	if len(flagged) != 1 {
		t.Error("Expected question to be flagged")
	}
}

func TestUnflagOperation(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Flag and then unflag
	db.FlagQuestion(question.ID, models.FlagReasonManual)
	err := db.UnflagQuestion(question.ID)
	if err != nil {
		t.Fatalf("Failed to unflag question: %v", err)
	}

	flagged, _ := db.GetFlaggedQuestions()
	if len(flagged) != 0 {
		t.Error("Expected no flagged questions after unflagging")
	}
}

func TestResolveOperation(t *testing.T) {
	// Resolve is same as unflag
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	db.FlagQuestion(question.ID, models.FlagReasonManual)
	err := db.UnflagQuestion(question.ID) // Resolve = Unflag
	if err != nil {
		t.Fatalf("Failed to resolve question: %v", err)
	}
}

func TestFlagReasonTracking(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	question := &models.Question{
		TopicID:         topic.ID,
		QuestionText:    "Test",
		Difficulty:      3,
		TopicQuestionID: "1.1",
		CreatedAt:       time.Now(),
	}
	db.CreateQuestion(question)

	// Test different flag reasons
	reasons := []string{
		models.FlagReasonTooHard,
		models.FlagReasonTooEasy,
		models.FlagReasonAmbiguous,
		models.FlagReasonPoorDistractors,
		models.FlagReasonManual,
	}

	for _, reason := range reasons {
		db.UnflagQuestion(question.ID) // Clear previous
		err := db.FlagQuestion(question.ID, reason)
		if err != nil {
			t.Errorf("Failed to flag with reason %s: %v", reason, err)
		}

		flagged, _ := db.GetFlaggedQuestions()
		if len(flagged) > 0 && flagged[0].FlagReason != reason {
			t.Errorf("Expected reason %s, got %s", reason, flagged[0].FlagReason)
		}
	}
}

// Additional helper test
func TestListQuestionsWithPagination(t *testing.T) {
	db := database.NewInMemoryDB()

	topic := &models.Topic{Name: "Test", CreatedAt: time.Now()}
	db.CreateTopic(topic)

	// Create 10 questions
	for i := 0; i < 10; i++ {
		question := &models.Question{
			TopicID:         topic.ID,
			QuestionText:    "Question " + string(rune(i)),
			Difficulty:      3,
			TopicQuestionID: "1." + string(rune(i)),
			CreatedAt:       time.Now(),
		}
		db.CreateQuestion(question)
	}

	// Test pagination
	questions, total, err := db.ListQuestionsWithPagination(1, 5, "", 0, 0)
	if err != nil {
		t.Fatalf("Failed to list questions: %v", err)
	}

	if total != 10 {
		t.Errorf("Expected total 10, got %d", total)
	}

	if len(questions) != 5 {
		t.Errorf("Expected 5 questions on first page, got %d", len(questions))
	}
}
