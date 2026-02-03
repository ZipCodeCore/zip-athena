package api

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ZipCodeCore/zip-athena/config"
)

// Middleware type for chaining
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Chain applies middlewares to a handler
func Chain(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// CORSMiddleware adds CORS headers based on configuration
func CORSMiddleware(cfg *config.CORSConfig) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if !cfg.Enabled {
				next(w, r)
				return
			}

			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == "*" {
					allowed = true
					w.Header().Set("Access-Control-Allow-Origin", "*")
					break
				}
				if allowedOrigin == origin {
					allowed = true
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
				// Support wildcard subdomains (e.g., http://localhost:*)
				if strings.Contains(allowedOrigin, "*") {
					prefix := strings.Split(allowedOrigin, "*")[0]
					if strings.HasPrefix(origin, prefix) {
						allowed = true
						w.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			if allowed {
				// Set allowed methods
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))

				// Set allowed headers
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))

				// Set credentials
				if cfg.AllowedCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				// Set max age for preflight cache
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next(w, r)
		}
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Prevent MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Prevent clickjacking
			w.Header().Set("X-Frame-Options", "DENY")

			// Enable XSS protection
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Content Security Policy
			csp := "default-src 'self'; " +
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; " +
				"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; " +
				"connect-src 'self'; " +
				"frame-ancestors 'none';"
			w.Header().Set("Content-Security-Policy", csp)

			// Referrer Policy
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Permissions Policy (formerly Feature Policy)
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

			// Strict Transport Security (only for HTTPS)
			if r.TLS != nil {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}

			next(w, r)
		}
	}
}

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.RWMutex
	rate     int           // requests per minute
	burst    int           // maximum burst size
	cleanup  time.Duration // cleanup interval
}

// Visitor represents a rate-limited client
type Visitor struct {
	tokens   int
	lastSeen time.Time
	mu       sync.Mutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
		cleanup:  5 * time.Minute,
	}

	// Start cleanup goroutine
	go rl.cleanupVisitors()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.RLock()
	visitor, exists := rl.visitors[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		visitor = &Visitor{
			tokens:   rl.burst,
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = visitor
		rl.mu.Unlock()
	}

	visitor.mu.Lock()
	defer visitor.mu.Unlock()

	// Refill tokens based on time elapsed
	now := time.Now()
	elapsed := now.Sub(visitor.lastSeen)
	tokensToAdd := int(elapsed.Minutes() * float64(rl.rate))

	if tokensToAdd > 0 {
		visitor.tokens = min(visitor.tokens+tokensToAdd, rl.burst)
		visitor.lastSeen = now
	}

	// Check if we have tokens available
	if visitor.tokens > 0 {
		visitor.tokens--
		return true
	}

	return false
}

// cleanupVisitors removes old visitors
func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, visitor := range rl.visitors {
			visitor.mu.Lock()
			if time.Since(visitor.lastSeen) > rl.cleanup {
				delete(rl.visitors, ip)
			}
			visitor.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware applies rate limiting to requests
func RateLimitMiddleware(limiter *RateLimiter) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.Allow(ip) {
				http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
				return
			}

			next(w, r)
		}
	}
}

// getClientIP extracts the real client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Remove port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// RequestSizeMiddleware limits the size of request bodies
func RequestSizeMiddleware(maxBytes int64) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next(w, r)
		}
	}
}

// ContentTypeMiddleware validates the Content-Type header for POST/PUT requests
func ContentTypeMiddleware(allowedTypes []string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Only check for POST, PUT, PATCH requests
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				contentType := r.Header.Get("Content-Type")

				// Remove charset if present
				if idx := strings.Index(contentType, ";"); idx != -1 {
					contentType = contentType[:idx]
				}
				contentType = strings.TrimSpace(contentType)

				allowed := false
				for _, allowedType := range allowedTypes {
					if contentType == allowedType {
						allowed = true
						break
					}
				}

				if !allowed && len(allowedTypes) > 0 {
					http.Error(w, "Invalid Content-Type. Expected: "+strings.Join(allowedTypes, " or "), http.StatusUnsupportedMediaType)
					return
				}
			}

			next(w, r)
		}
	}
}

// Helper function for min (Go 1.21+)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
