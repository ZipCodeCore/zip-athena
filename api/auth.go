package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
	"golang.org/x/crypto/bcrypt"
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

// hashPassword creates a bcrypt hash of the password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPasswordHash compares a password with a hash
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		log.Printf("ERROR: Failed to hash password for user %s: %v\n", req.Username, err)
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user with hashed password
	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
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
			log.Printf("ERROR: Failed to create user %s: %v\n", user.Username, err)
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		id, _ := result.LastInsertId()
		user.ID = int(id)
	} else {
		// Fallback for in-memory DB
		if err := h.db.CreateUser(user); err != nil {
			log.Printf("ERROR: Failed to create user %s in in-memory DB: %v\n", user.Username, err)
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
			// Don't log "not found" as it could be username enumeration attempts
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		log.Printf("ERROR: Failed to get user %s during login: %v\n", req.Username, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify password
	if !checkPasswordHash(req.Password, user.Password) {
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
