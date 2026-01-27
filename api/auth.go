package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// RegisterRequest represents the request body for registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the request body for login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the response for authentication
type AuthResponse struct {
	User    *models.User `json:"user"`
	Message string       `json:"message"`
}

// hashPassword creates a simple hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// RegisterHandler handles POST /auth/register
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Create user with hashed password
	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashPassword(req.Password),
		IsAdmin:     false,
		TotalPoints: 0,
		CreatedAt:   time.Now(),
	}

	// Try to create the user in SQLite database
	if sqlDB, ok := h.db.(*database.SQLiteDB); ok {
		result, err := sqlDB.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			user.Username, user.Email, user.Password, user.IsAdmin, user.TotalPoints, user.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		id, _ := result.LastInsertId()
		user.ID = int(id)
	} else {
		// Fallback for in-memory DB
		if err := h.db.CreateUser(user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Don't send password back
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(AuthResponse{
		User:    user,
		Message: "Registration successful",
	})
}

// LoginHandler handles POST /auth/login
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get user by username
	var user *models.User
	var err error

	if sqlDB, ok := h.db.(*database.SQLiteDB); ok {
		user, err = sqlDB.GetUserByUsername(req.Username)
	} else {
		http.Error(w, "Login not supported with in-memory database", http.StatusInternalServerError)
		return
	}

	if err != nil {
		if err == database.ErrNotFound {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify password
	hashedPassword := hashPassword(req.Password)
	if user.Password != hashedPassword {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Don't send password back
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		User:    user,
		Message: "Login successful",
	})
}
