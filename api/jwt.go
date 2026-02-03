package api

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ZipCodeCore/zip-athena/config"
	"github.com/golang-jwt/jwt/v5"
)

// TokenClaims represents the custom claims in a JWT token
type TokenClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	TokenID  string `json:"jti"` // Token ID for blacklisting

	jwt.RegisteredClaims
}

// TokenPair holds both access and refresh tokens
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`  // Access token expiry in seconds
	TokenType    string `json:"token_type"`  // Always "Bearer"
}

// TokenBlacklist keeps track of invalidated tokens
type TokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

// JWTManager handles all JWT operations
type JWTManager struct {
	secretKey       []byte
	accessExpiry    time.Duration
	refreshExpiry   time.Duration
	blacklist       *TokenBlacklist
}

var (
	ErrTokenExpired    = errors.New("token has expired")
	ErrTokenInvalid    = errors.New("token is invalid")
	ErrTokenBlacklisted = errors.New("token has been revoked")
	ErrNoToken         = errors.New("no token provided")
)

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.SecurityConfig) *JWTManager {
	return &JWTManager{
		secretKey:     []byte(cfg.JWTSecret),
		accessExpiry:  cfg.JWTExpiration,
		refreshExpiry: cfg.JWTExpiration * 7, // Refresh token lasts 7x longer
		blacklist: &TokenBlacklist{
			tokens: make(map[string]time.Time),
		},
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (jm *JWTManager) GenerateTokenPair(userID int, username string, isAdmin bool) (*TokenPair, error) {
	accessToken, err := jm.generateToken(userID, username, isAdmin, jm.accessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jm.generateToken(userID, username, isAdmin, jm.refreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(jm.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateToken creates a signed JWT token
func (jm *JWTManager) generateToken(userID int, username string, isAdmin bool, expiry time.Duration) (string, error) {
	claims := TokenClaims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		TokenID:  generateUUID(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "athena",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jm.secretKey)
}

// ValidateToken validates a token string and returns the claims
func (jm *JWTManager) ValidateToken(tokenString string) (*TokenClaims, error) {
	if tokenString == "" {
		return nil, ErrNoToken
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return jm.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	// Check if token is blacklisted
	if jm.blacklist.IsBlacklisted(claims.TokenID) {
		return nil, ErrTokenBlacklisted
	}

	return claims, nil
}

// RefreshToken generates a new access token from a valid refresh token
func (jm *JWTManager) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	claims, err := jm.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, err
	}

	// Blacklist the old refresh token
	jm.blacklist.Add(claims.TokenID, jm.refreshExpiry)

	// Generate new token pair
	return jm.GenerateTokenPair(claims.UserID, claims.Username, claims.IsAdmin)
}

// RevokeToken adds a token to the blacklist
func (jm *JWTManager) RevokeToken(tokenString string) error {
	claims, err := jm.ValidateToken(tokenString)
	if err != nil {
		// If token is already invalid/expired, that's fine
		return nil
	}

	// Calculate remaining time until expiry
	remainingTime := time.Until(claims.ExpiresAt.Time)
	if remainingTime < 0 {
		remainingTime = 0
	}

	jm.blacklist.Add(claims.TokenID, remainingTime)
	return nil
}

// ExtractToken extracts the token from the Authorization header
func ExtractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// Also check cookie as fallback
		cookie, err := r.Cookie("athena_token")
		if err == nil && cookie.Value != "" {
			return cookie.Value
		}
		return ""
	}

	// Expect "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// JWTAuthMiddleware validates JWT tokens on protected routes
func JWTAuthMiddleware(jm *JWTManager) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := ExtractToken(r)

			claims, err := jm.ValidateToken(tokenString)
			if err != nil {
				if err == ErrNoToken {
					http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
				} else if err == ErrTokenExpired {
					http.Error(w, `{"error":"token expired","detail":"use refresh_token to obtain a new access token"}`, http.StatusUnauthorized)
				} else if err == ErrTokenBlacklisted {
					http.Error(w, `{"error":"token revoked"}`, http.StatusUnauthorized)
				} else {
					http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
				}
				return
			}

			// Set user info in request header for downstream handlers
			r.Header.Set("X-User-ID", strings.Replace(r.Header.Get("X-User-ID"), "", "", 0))
			r.Header.Set("X-User-ID", itoa(claims.UserID))
			r.Header.Set("X-Username", claims.Username)
			if claims.IsAdmin {
				r.Header.Set("X-Is-Admin", "true")
			}

			next(w, r)
		}
	}
}

// JWTOptionalMiddleware attempts JWT validation but allows unauthenticated access
func JWTOptionalMiddleware(jm *JWTManager) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString := ExtractToken(r)

			if tokenString != "" {
				claims, err := jm.ValidateToken(tokenString)
				if err == nil {
					r.Header.Set("X-User-ID", itoa(claims.UserID))
					r.Header.Set("X-Username", claims.Username)
					if claims.IsAdmin {
						r.Header.Set("X-Is-Admin", "true")
					}
				}
			}

			next(w, r)
		}
	}
}

// TokenBlacklist methods

// Add adds a token to the blacklist
func (tb *TokenBlacklist) Add(tokenID string, ttl time.Duration) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens[tokenID] = time.Now().Add(ttl)
}

// IsBlacklisted checks if a token is in the blacklist
func (tb *TokenBlacklist) IsBlacklisted(tokenID string) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()

	expiresAt, exists := tb.tokens[tokenID]
	if !exists {
		return false
	}

	// If the blacklist entry itself has expired, it's not really blacklisted anymore
	if time.Now().After(expiresAt) {
		return false
	}

	return true
}

// Cleanup removes expired entries from the blacklist
func (tb *TokenBlacklist) Cleanup() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	for tokenID, expiresAt := range tb.tokens {
		if time.Now().After(expiresAt) {
			delete(tb.tokens, tokenID)
		}
	}
}

// Helper: generateUUID generates a simple unique ID
func generateUUID() string {
	return formatTime(time.Now()) + "-" + formatInt(randInt())
}

func formatTime(t time.Time) string {
	return t.Format("20060102150405")
}

func formatInt(n int) string {
	return itoa(n)
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	negative := i < 0
	if negative {
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}

func randInt() int {
	return int(time.Now().UnixNano() % 1000000)
}
