package api

import (
	"encoding/json"
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
	User         *models.User `json:"user"`
	Message      string       `json:"message"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	ExpiresIn    int64        `json:"expires_in,omitempty"`
	TokenType    string       `json:"token_type,omitempty"`
}

// RefreshRequest represents the request body for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
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
		h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
			"username": req.Username,
			"error":    err.Error(),
		}).Error("Failed to hash password during registration")
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Create user with hashed password
	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    hashedPassword,
		IsAdmin:     false,
		IsActive:    true,
		IsActiveSet: true,
		TotalPoints: 0,
		CreatedAt:   time.Now(),
	}

	// Try to create the user in SQLite database
	if sqlDB, ok := h.db.(*database.SQLiteDB); ok {
		result, err := sqlDB.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, is_active, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			user.Username, user.Email, user.Password, user.IsAdmin, user.IsActive, user.TotalPoints, user.CreatedAt,
		)
		if err != nil {
			h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
				"username": user.Username,
				"email":    user.Email,
				"error":    err.Error(),
			}).Error("Failed to create user in database")
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		id, _ := result.LastInsertId()
		user.ID = int(id)
	} else {
		// Fallback for in-memory DB
		if err := h.db.CreateUser(user); err != nil {
			h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
				"username": user.Username,
				"error":    err.Error(),
			}).Error("Failed to create user in in-memory database")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Don't send password back
	user.Password = ""

	response := AuthResponse{
		User:    user,
		Message: "Registration successful",
	}

	if h.jwtManager != nil {
		if tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, user.IsAdmin); err == nil {
			response.AccessToken = tokenPair.AccessToken
			response.RefreshToken = tokenPair.RefreshToken
			response.ExpiresIn = tokenPair.ExpiresIn
			response.TokenType = tokenPair.TokenType
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
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
			// Log failed authentication attempt but don't reveal if username exists
			h.logger.LogAuthentication(r.Context(), req.Username, false, "user not found")
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		h.logger.WithContext(r.Context()).WithFields(map[string]interface{}{
			"username": req.Username,
			"error":    err.Error(),
		}).Error("Database error during login")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify password
	if !checkPasswordHash(req.Password, user.Password) {
		h.logger.LogAuthentication(r.Context(), req.Username, false, "invalid password")
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	if !user.IsActive {
		h.logger.LogAuthentication(r.Context(), req.Username, false, "inactive user")
		http.Error(w, "Account is inactive", http.StatusUnauthorized)
		return
	}

	// Log successful authentication
	h.logger.LogAuthentication(r.Context(), req.Username, true, "")

	// Don't send password back
	user.Password = ""

	response := AuthResponse{
		User:    user,
		Message: "Login successful",
	}

	if h.jwtManager != nil {
		if tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, user.IsAdmin); err == nil {
			response.AccessToken = tokenPair.AccessToken
			response.RefreshToken = tokenPair.RefreshToken
			response.ExpiresIn = tokenPair.ExpiresIn
			response.TokenType = tokenPair.TokenType
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LogoutHandler handles POST /auth/logout
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.jwtManager != nil {
		tokenString := ExtractToken(r)
		if tokenString != "" {
			h.jwtManager.RevokeToken(tokenString)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}

// RefreshTokenHandler handles POST /auth/refresh
func (h *Handler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.jwtManager == nil {
		http.Error(w, "Authentication not configured", http.StatusInternalServerError)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "Refresh token is required", http.StatusBadRequest)
		return
	}

	tokenPair, err := h.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}
