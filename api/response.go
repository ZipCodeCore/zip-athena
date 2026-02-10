package api

import (
	"encoding/json"
	"net/http"
)

// APIResponse represents a standardized API response envelope
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error in the API response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// respondJSON sends a successful JSON response
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: status < 400,
		Data:    data,
	})
}

// respondError sends an error JSON response
func (h *Handler) respondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}
