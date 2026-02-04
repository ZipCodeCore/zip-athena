package api

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/models"
)

// GetCurriculumOutlineHandler returns the complete curriculum organized by phases
func (h *Handler) GetCurriculumOutlineHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get language from query parameter (default to both)
	language := r.URL.Query().Get("language")
	if language != "" && language != "java" && language != "python" {
		http.Error(w, "Invalid language. Must be 'java' or 'python'", http.StatusBadRequest)
		return
	}

	// Get user_id from query parameter for filtering (optional)
	userIDStr := r.URL.Query().Get("user_id")
	var user *models.User
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err == nil {
			user, _ = h.db.GetUser(userID)
		}
	}

	// Filter language based on user category for non-admins
	if user != nil && !user.IsAdmin {
		if user.Category == "java" {
			language = "java"
		} else if user.Category == "data" {
			language = "python"
		}
	}

	// Define phase metadata
	phaseInfo := map[int]struct {
		Title       string
		Description string
	}{
		1: {"Core Fundamentals", "Build a strong foundation with essential programming concepts"},
		2: {"Object-Oriented Programming", "Master classes, objects, inheritance, and polymorphism"},
		3: {"Advanced Data Handling", "Work with files, exceptions, and complex data structures"},
		4: {"Advanced Topics", "Explore generics, concurrency, functional programming, and more"},
		5: {"Specialized Topics", "Learn best practices, testing, and professional development skills"},
	}

	phases := make([]models.PhaseOutline, 0)

	for phase := 1; phase <= 5; phase++ {
		info := phaseInfo[phase]
		outline := models.PhaseOutline{
			Phase:       phase,
			Title:       info.Title,
			Description: info.Description,
			Resources:   make([]models.LearningResource, 0),
		}

		// Build query
		query := `
			SELECT id, language, phase, section_number, section_title, 
			       file_path, display_order, is_active, created_at
			FROM learning_resources
			WHERE phase = ? AND is_active = 1
		`
		args := []interface{}{phase}

		if language != "" {
			query += " AND language = ?"
			args = append(args, language)
		}

		query += " ORDER BY language, display_order"

		rows, err := h.db.DB().Query(query, args...)
		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var resource models.LearningResource
			err := rows.Scan(
				&resource.ID,
				&resource.Language,
				&resource.Phase,
				&resource.SectionNumber,
				&resource.SectionTitle,
				&resource.FilePath,
				&resource.DisplayOrder,
				&resource.IsActive,
				&resource.CreatedAt,
			)
			if err != nil {
				continue
			}
			outline.Resources = append(outline.Resources, resource)
		}

		// Only include phases that have resources
		if len(outline.Resources) > 0 {
			phases = append(phases, outline)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"phases": phases,
	})
}

// GetPhaseResourcesHandler returns all resources for a specific phase
func (h *Handler) GetPhaseResourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract phase from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	phase, err := strconv.Atoi(pathParts[3])
	if err != nil || phase < 1 || phase > 5 {
		http.Error(w, "Invalid phase number", http.StatusBadRequest)
		return
	}

	language := r.URL.Query().Get("language")
	userID := getUserIDFromContext(r)

	query := `
		SELECT id, language, phase, section_number, section_title, 
		       file_path, display_order, is_active, created_at
		FROM learning_resources
		WHERE phase = ? AND is_active = 1
	`
	args := []interface{}{phase}

	if language != "" {
		query += " AND language = ?"
		args = append(args, language)
	}

	query += " ORDER BY language, display_order"

	rows, err := h.db.DB().Query(query, args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	resources := make([]models.LearningResourceWithProgress, 0)

	for rows.Next() {
		var resource models.LearningResource
		err := rows.Scan(
			&resource.ID,
			&resource.Language,
			&resource.Phase,
			&resource.SectionNumber,
			&resource.SectionTitle,
			&resource.FilePath,
			&resource.DisplayOrder,
			&resource.IsActive,
			&resource.CreatedAt,
		)
		if err != nil {
			continue
		}

		// Get user progress if logged in
		progress := models.LearningResourceWithProgress{
			LearningResource:     resource,
			ViewCount:            0,
			CompletionPercentage: 0,
		}

		if userID > 0 {
			var viewCount int
			var lastViewed time.Time
			var completion float64

			err := h.db.DB().QueryRow(`
				SELECT view_count, last_viewed_at, completion_percentage
				FROM user_learning_progress
				WHERE user_id = ? AND learning_resource_id = ?
			`, userID, resource.ID).Scan(&viewCount, &lastViewed, &completion)

			if err == nil {
				progress.ViewCount = viewCount
				progress.LastViewedAt = &lastViewed
				progress.CompletionPercentage = completion
			}
		}

		resources = append(resources, progress)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"phase":     phase,
		"resources": resources,
	})
}

// GetLearningResourceHandler returns metadata for a specific learning resource
func (h *Handler) GetLearningResourceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract resource ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: URL=%s, pathParts=%v", r.URL.Path, pathParts)
	if len(pathParts) < 4 {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: Path too short, len=%d", len(pathParts))
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	resourceID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: Failed to parse ID from '%s': %v", pathParts[3], err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: Fetching resource ID=%d", resourceID)

	userID := getUserIDFromContext(r)

	var resource models.LearningResource
	err = h.db.DB().QueryRow(`
		SELECT id, language, phase, section_number, section_title, 
		       file_path, display_order, is_active, created_at
		FROM learning_resources
		WHERE id = ? AND is_active = 1
	`, resourceID).Scan(
		&resource.ID,
		&resource.Language,
		&resource.Phase,
		&resource.SectionNumber,
		&resource.SectionTitle,
		&resource.FilePath,
		&resource.DisplayOrder,
		&resource.IsActive,
		&resource.CreatedAt,
	)

	if err == sql.ErrNoRows {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: Resource ID=%d not found", resourceID)
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceHandler: Found resource: %+v", resource)

	// Build response with progress
	response := models.LearningResourceWithProgress{
		LearningResource:     resource,
		ViewCount:            0,
		CompletionPercentage: 0,
	}

	if userID > 0 {
		var viewCount int
		var lastViewed time.Time
		var completion float64

		err := h.db.DB().QueryRow(`
			SELECT view_count, last_viewed_at, completion_percentage
			FROM user_learning_progress
			WHERE user_id = ? AND learning_resource_id = ?
		`, userID, resource.ID).Scan(&viewCount, &lastViewed, &completion)

		if err == nil {
			response.ViewCount = viewCount
			response.LastViewedAt = &lastViewed
			response.CompletionPercentage = completion
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetLearningResourceContentHandler returns the markdown content of a learning resource
func (h *Handler) GetLearningResourceContentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract resource ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: URL=%s, pathParts=%v", r.URL.Path, pathParts)
	if len(pathParts) < 4 {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Path too short, len=%d", len(pathParts))
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	resourceID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Failed to parse ID from '%s': %v", pathParts[3], err)
		http.Error(w, "Invalid resource ID", http.StatusBadRequest)
		return
	}
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Fetching content for resource ID=%d", resourceID)

	var filePath string
	err = h.db.DB().QueryRow(`
		SELECT file_path
		FROM learning_resources
		WHERE id = ? AND is_active = 1
	`, resourceID).Scan(&filePath)

	if err == sql.ErrNoRows {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Resource ID=%d not found", resourceID)
		http.Error(w, "Resource not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Database error: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Found file_path: %s", filePath)

	// Read the markdown file
	fullPath := filepath.Join("learning", filePath)
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Reading file: %s", fullPath)
	content, err := ioutil.ReadFile(fullPath)
	if err != nil {
		h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: File read error: %v", err)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	h.logger.WithContext(r.Context()).Debugf("GetLearningResourceContentHandler: Successfully read %d bytes", len(content))

	w.Header().Set("Content-Type", "text/markdown")
	w.Write(content)
}

// TrackLearningViewHandler records that a user viewed a learning resource
func (h *Handler) TrackLearningViewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		// Allow anonymous users - just return success without tracking
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "tracking skipped for anonymous user"})
		return
	}

	var req struct {
		ResourceID           int     `json:"resource_id"`
		CompletionPercentage float64 `json:"completion_percentage"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if progress record exists
	var exists bool
	err := h.db.DB().QueryRow(`
		SELECT EXISTS(SELECT 1 FROM user_learning_progress 
		              WHERE user_id = ? AND learning_resource_id = ?)
	`, userID, req.ResourceID).Scan(&exists)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if exists {
		// Update existing record
		_, err = h.db.DB().Exec(`
			UPDATE user_learning_progress
			SET last_viewed_at = CURRENT_TIMESTAMP,
			    view_count = view_count + 1,
			    completion_percentage = ?
			WHERE user_id = ? AND learning_resource_id = ?
		`, req.CompletionPercentage, userID, req.ResourceID)
	} else {
		// Insert new record
		_, err = h.db.DB().Exec(`
			INSERT INTO user_learning_progress 
			(user_id, learning_resource_id, completion_percentage)
			VALUES (?, ?, ?)
		`, userID, req.ResourceID, req.CompletionPercentage)
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "View tracked successfully",
	})
}

// GetQuestionLearningLinksHandler returns learning resources linked to a question
func (h *Handler) GetQuestionLearningLinksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract question ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	rows, err := h.db.DB().Query(`
		SELECT lr.id, lr.language, lr.phase, lr.section_number, 
		       lr.section_title, lr.file_path, lr.display_order, 
		       lr.is_active, lr.created_at, qll.relevance_score
		FROM learning_resources lr
		JOIN question_learning_links qll ON lr.id = qll.learning_resource_id
		WHERE qll.question_id = ? AND lr.is_active = 1
		ORDER BY qll.relevance_score DESC
	`, questionID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ResourceWithRelevance struct {
		models.LearningResource
		RelevanceScore float64 `json:"relevance_score"`
	}

	resources := make([]ResourceWithRelevance, 0)

	for rows.Next() {
		var resource models.LearningResource
		var relevance float64

		err := rows.Scan(
			&resource.ID,
			&resource.Language,
			&resource.Phase,
			&resource.SectionNumber,
			&resource.SectionTitle,
			&resource.FilePath,
			&resource.DisplayOrder,
			&resource.IsActive,
			&resource.CreatedAt,
			&relevance,
		)

		if err != nil {
			continue
		}

		resources = append(resources, ResourceWithRelevance{
			LearningResource: resource,
			RelevanceScore:   relevance,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"question_id": questionID,
		"resources":   resources,
	})
}

// SearchLearningResourcesHandler searches learning resources by query
func (h *Handler) SearchLearningResourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	language := r.URL.Query().Get("language")

	// Get user_id from query parameter for filtering (optional)
	userIDStr := r.URL.Query().Get("user_id")
	var user *models.User
	if userIDStr != "" {
		userID, err := strconv.Atoi(userIDStr)
		if err == nil {
			user, _ = h.db.GetUser(userID)
		}
	}

	// Filter language based on user category for non-admins
	if user != nil && !user.IsAdmin {
		if user.Category == "java" {
			language = "java"
		} else if user.Category == "data" {
			language = "python"
		}
	}

	sqlQuery := `
		SELECT id, language, phase, section_number, section_title, 
		       file_path, display_order, is_active, created_at
		FROM learning_resources
		WHERE is_active = 1 AND (
			section_title LIKE ? OR 
			section_number LIKE ?
		)
	`
	args := []interface{}{
		"%" + query + "%",
		"%" + query + "%",
	}

	if language != "" {
		sqlQuery += " AND language = ?"
		args = append(args, language)
	}

	sqlQuery += " ORDER BY phase, display_order LIMIT 20"

	rows, err := h.db.DB().Query(sqlQuery, args...)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	resources := make([]models.LearningResource, 0)

	for rows.Next() {
		var resource models.LearningResource
		err := rows.Scan(
			&resource.ID,
			&resource.Language,
			&resource.Phase,
			&resource.SectionNumber,
			&resource.SectionTitle,
			&resource.FilePath,
			&resource.DisplayOrder,
			&resource.IsActive,
			&resource.CreatedAt,
		)
		if err != nil {
			continue
		}
		resources = append(resources, resource)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":   query,
		"results": resources,
		"count":   len(resources),
	})
}

// Helper function to get user ID from request context
// This should be implemented based on your authentication system
func getUserIDFromContext(r *http.Request) int {
	// TODO: Implement based on your auth system
	// For now, try to get from session or token
	_, err := r.Cookie("session_token")
	if err != nil {
		return 0
	}

	// This is a simplified version - you should implement proper session handling
	// For now, return 0 (not authenticated)
	return 0
}
