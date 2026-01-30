package api

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ZipCodeCore/zip-athena/models"
)

// Phase 5: Instructor Question Management Handlers

// ListInstructorQuestionsHandler handles GET /instructor/questions
func (h *Handler) ListInstructorQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")
	search := r.URL.Query().Get("search")
	topicIDStr := r.URL.Query().Get("topic_id")
	difficultyStr := r.URL.Query().Get("difficulty")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 30
	}

	topicID, _ := strconv.Atoi(topicIDStr)
	difficulty, _ := strconv.Atoi(difficultyStr)

	questions, total, err := h.db.ListQuestionsWithPagination(page, perPage, search, topicID, difficulty)
	if err != nil {
		log.Printf("ERROR: Failed to list questions: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"questions":  questions,
		"total":      total,
		"page":       page,
		"per_page":   perPage,
		"total_pages": (total + perPage - 1) / perPage,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateInstructorQuestionHandler handles POST /instructor/questions
func (h *Handler) CreateInstructorQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Question models.Question `json:"question"`
		Answers  []models.Answer `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate question data
	if err := h.db.ValidateQuestionData(&req.Question, req.Answers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create question
	if err := h.db.CreateQuestion(&req.Question); err != nil {
		log.Printf("ERROR: Failed to create question: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create answers
	for i := range req.Answers {
		req.Answers[i].QuestionID = req.Question.ID
		if err := h.db.CreateAnswer(&req.Answers[i]); err != nil {
			log.Printf("ERROR: Failed to create answer: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"question": req.Question,
		"answers":  req.Answers,
	})
}

// UpdateInstructorQuestionHandler handles PUT /instructor/questions/{id}
func (h *Handler) UpdateInstructorQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Question models.Question `json:"question"`
		Answers  []models.Answer `json:"answers"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set ID from URL
	req.Question.ID = id

	// Validate question data
	if err := h.db.ValidateQuestionData(&req.Question, req.Answers); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update question
	if err := h.db.UpdateQuestion(&req.Question); err != nil {
		log.Printf("ERROR: Failed to update question: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get existing answers
	existingAnswers, err := h.db.GetAnswersByQuestion(id)
	if err != nil {
		log.Printf("ERROR: Failed to get existing answers: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete old answers
	for _, answer := range existingAnswers {
		if err := h.db.DeleteAnswer(answer.ID); err != nil {
			log.Printf("ERROR: Failed to delete answer: %v\n", err)
		}
	}

	// Create new answers
	for i := range req.Answers {
		req.Answers[i].QuestionID = id
		if err := h.db.CreateAnswer(&req.Answers[i]); err != nil {
			log.Printf("ERROR: Failed to create answer: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"question": req.Question,
		"answers":  req.Answers,
	})
}

// DeactivateInstructorQuestionHandler handles DELETE /instructor/questions/{id}
func (h *Handler) DeactivateInstructorQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Deactivate question (soft delete)
	if err := h.db.DeactivateQuestion(id); err != nil {
		log.Printf("ERROR: Failed to deactivate question: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// PreviewInstructorQuestionHandler handles GET /instructor/questions/{id}/preview
func (h *Handler) PreviewInstructorQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	preview, err := h.db.GetQuestionPreview(id)
	if err != nil {
		log.Printf("ERROR: Failed to get question preview: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(preview)
}

// Phase 5: Admin CSV Import Handlers

// ImportQuestionsCSVHandler handles POST /admin/questions/import
func (h *Handler) ImportQuestionsCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to parse CSV: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(records) < 2 {
		http.Error(w, "CSV must have at least a header row and one data row", http.StatusBadRequest)
		return
	}

	// Parse rows
	var rows []models.CSVImportRow
	for i, record := range records[1:] { // Skip header
		if len(record) < 9 {
			continue // Skip invalid rows
		}

		rows = append(rows, models.CSVImportRow{
			Topic:         record[0],
			Question:      record[1],
			Answer1:       record[2],
			Answer2:       record[3],
			Answer3:       record[4],
			Answer4:       record[5],
			CorrectAnswer: record[6],
			Difficulty:    record[7],
			Tags:          record[8],
		})

		// Limit to prevent abuse
		if i >= 500 {
			break
		}
	}

	// Validate
	result, err := h.db.ValidateCSVImport(rows)
	if err != nil {
		log.Printf("ERROR: Failed to validate CSV import: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Import valid rows
	imported, err := h.db.BulkCreateQuestions(result.ValidRows)
	if err != nil {
		log.Printf("ERROR: Failed to bulk create questions: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result.Summary.ImportedRows = imported

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ValidateQuestionsCSVHandler handles POST /admin/questions/import/validate
func (h *Handler) ValidateQuestionsCSVHandler(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to parse CSV: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(records) < 2 {
		http.Error(w, "CSV must have at least a header row and one data row", http.StatusBadRequest)
		return
	}

	// Parse rows
	var rows []models.CSVImportRow
	for i, record := range records[1:] { // Skip header
		if len(record) < 9 {
			continue // Skip invalid rows
		}

		rows = append(rows, models.CSVImportRow{
			Topic:         record[0],
			Question:      record[1],
			Answer1:       record[2],
			Answer2:       record[3],
			Answer3:       record[4],
			Answer4:       record[5],
			CorrectAnswer: record[6],
			Difficulty:    record[7],
			Tags:          record[8],
		})

		// Limit to prevent abuse
		if i >= 500 {
			break
		}
	}

	// Validate only (don't import)
	result, err := h.db.ValidateCSVImport(rows)
	if err != nil {
		log.Printf("ERROR: Failed to validate CSV import: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// ConfirmImportQuestionsCSVHandler handles POST /admin/questions/import/confirm
func (h *Handler) ConfirmImportQuestionsCSVHandler(w http.ResponseWriter, r *http.Request) {
	var rows []models.CSVImportRow
	if err := json.NewDecoder(r.Body).Decode(&rows); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Import confirmed rows
	imported, err := h.db.BulkCreateQuestions(rows)
	if err != nil {
		log.Printf("ERROR: Failed to bulk create questions: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"imported": imported,
	})
}

// Phase 5: Admin Tag Management Handlers

// ListTagsHandler handles GET /admin/tags
func (h *Handler) ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := h.db.GetTagsWithMetadata()
	if err != nil {
		log.Printf("ERROR: Failed to list tags: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// MergeTagsHandler handles POST /admin/tags/merge
func (h *Handler) MergeTagsHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SourceTags []string `json:"source_tags"`
		TargetTag  string   `json:"target_tag"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.SourceTags) == 0 {
		http.Error(w, "source_tags is required", http.StatusBadRequest)
		return
	}

	if req.TargetTag == "" {
		http.Error(w, "target_tag is required", http.StatusBadRequest)
		return
	}

	if err := h.db.MergeTags(req.SourceTags, req.TargetTag); err != nil {
		log.Printf("ERROR: Failed to merge tags: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Tags merged successfully",
	})
}

// DeleteTagHandler handles DELETE /admin/tags/{name}
func (h *Handler) DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tagName := r.URL.Query().Get("name")
	if tagName == "" {
		http.Error(w, "tag name is required", http.StatusBadRequest)
		return
	}

	if err := h.db.DeleteTag(tagName); err != nil {
		log.Printf("ERROR: Failed to delete tag: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetQuestionsByTagHandler handles GET /admin/tags/{name}/questions
func (h *Handler) GetQuestionsByTagHandler(w http.ResponseWriter, r *http.Request) {
	tagName := r.URL.Query().Get("name")
	if tagName == "" {
		http.Error(w, "tag name is required", http.StatusBadRequest)
		return
	}

	questions, err := h.db.GetQuestionsByTag(tagName)
	if err != nil {
		log.Printf("ERROR: Failed to get questions by tag: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

// Phase 5: Admin Quality Review Handlers

// GetFlaggedQuestionsHandler handles GET /admin/questions/flagged
func (h *Handler) GetFlaggedQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	flagged, err := h.db.GetFlaggedQuestions()
	if err != nil {
		log.Printf("ERROR: Failed to get flagged questions: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(flagged)
}

// FlagQuestionHandler handles POST /admin/questions/{id}/flag
func (h *Handler) FlagQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Reason == "" {
		req.Reason = models.FlagReasonManual
	}

	if err := h.db.FlagQuestion(id, req.Reason); err != nil {
		log.Printf("ERROR: Failed to flag question: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Question flagged successfully",
	})
}

// UnflagQuestionHandler handles DELETE /admin/questions/{id}/flag
func (h *Handler) UnflagQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := h.db.UnflagQuestion(id); err != nil {
		log.Printf("ERROR: Failed to unflag question: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ResolveQuestionHandler handles PUT /admin/questions/{id}/resolve
func (h *Handler) ResolveQuestionHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Resolve is same as unflag
	if err := h.db.UnflagQuestion(id); err != nil {
		log.Printf("ERROR: Failed to resolve question: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AutoFlagQuestionsHandler handles POST /admin/questions/auto-flag
func (h *Handler) AutoFlagQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	if err := h.db.AutoFlagQuestions(); err != nil {
		log.Printf("ERROR: Failed to auto-flag questions: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Questions auto-flagged successfully",
	})
}
