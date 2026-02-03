package logger

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ZipCodeCore/zip-athena/config"
	"github.com/sirupsen/logrus"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

const (
	// RequestIDKey is the context key for request IDs
	RequestIDKey contextKey = "request_id"
)

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// New creates a new configured logger instance
func New(cfg *config.LoggingConfig) (*Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	// Set formatter
	if cfg.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	}

	// Set output
	if cfg.Output == "stdout" {
		logger.SetOutput(os.Stdout)
	} else if cfg.Output == "stderr" {
		logger.SetOutput(os.Stderr)
	} else {
		// Assume it's a file path
		file, err := os.OpenFile(cfg.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)
	}

	return &Logger{Logger: logger}, nil
}

// WithRequestID returns a logger with the request ID field set
func (l *Logger) WithRequestID(requestID string) *logrus.Entry {
	return l.WithField("request_id", requestID)
}

// WithContext returns a logger with fields from context (e.g., request ID)
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
	entry := l.WithFields(logrus.Fields{})

	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	return entry
}

// WithUser returns a logger with user ID field set
func (l *Logger) WithUser(userID int) *logrus.Entry {
	return l.WithField("user_id", userID)
}

// WithError returns a logger with error field set
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// WithFields returns a logger with custom fields
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// HTTP Logging Helpers

// LogRequest logs an HTTP request with standard fields
func (l *Logger) LogRequest(r *http.Request, statusCode int, duration time.Duration) {
	requestID := GetRequestID(r.Context())

	l.WithFields(logrus.Fields{
		"request_id": requestID,
		"method":     r.Method,
		"path":       r.URL.Path,
		"query":      r.URL.RawQuery,
		"status":     statusCode,
		"duration":   duration.Milliseconds(),
		"ip":         getClientIP(r),
		"user_agent": r.UserAgent(),
	}).Info("HTTP request")
}

// LogRequestStart logs the start of an HTTP request
func (l *Logger) LogRequestStart(r *http.Request) {
	requestID := GetRequestID(r.Context())

	l.WithFields(logrus.Fields{
		"request_id": requestID,
		"method":     r.Method,
		"path":       r.URL.Path,
		"ip":         getClientIP(r),
	}).Debug("HTTP request started")
}

// LogRequestError logs an HTTP request error
func (l *Logger) LogRequestError(r *http.Request, err error, statusCode int) {
	requestID := GetRequestID(r.Context())

	l.WithFields(logrus.Fields{
		"request_id": requestID,
		"method":     r.Method,
		"path":       r.URL.Path,
		"status":     statusCode,
		"error":      err.Error(),
	}).Error("HTTP request error")
}

// Database Logging Helpers

// LogDBQuery logs a database query
func (l *Logger) LogDBQuery(ctx context.Context, query string, duration time.Duration, err error) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"query":    query,
		"duration": duration.Milliseconds(),
	})

	if err != nil {
		entry.WithError(err).Error("Database query failed")
	} else {
		entry.Debug("Database query executed")
	}
}

// LogDBOperation logs a database operation (create, update, delete)
func (l *Logger) LogDBOperation(ctx context.Context, operation, entity string, entityID int, err error) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"operation": operation,
		"entity":    entity,
		"entity_id": entityID,
	})

	if err != nil {
		entry.WithError(err).Error("Database operation failed")
	} else {
		entry.Info("Database operation successful")
	}
}

// Business Logic Logging Helpers

// LogQuizGeneration logs quiz generation
func (l *Logger) LogQuizGeneration(ctx context.Context, userID, quizID int, questionCount int) {
	l.WithContext(ctx).WithFields(logrus.Fields{
		"user_id":        userID,
		"quiz_id":        quizID,
		"question_count": questionCount,
	}).Info("Quiz generated")
}

// LogQuizSubmission logs quiz submission
func (l *Logger) LogQuizSubmission(ctx context.Context, userID, quizID int, score float64) {
	l.WithContext(ctx).WithFields(logrus.Fields{
		"user_id": userID,
		"quiz_id": quizID,
		"score":   score,
	}).Info("Quiz submitted")
}

// LogAuthentication logs authentication attempts
func (l *Logger) LogAuthentication(ctx context.Context, username string, success bool, reason string) {
	entry := l.WithContext(ctx).WithFields(logrus.Fields{
		"username": username,
		"success":  success,
	})

	if !success {
		entry = entry.WithField("reason", reason)
		entry.Warn("Authentication failed")
	} else {
		entry.Info("Authentication successful")
	}
}

// Request ID Generation and Retrieval

// GenerateRequestID generates a new random request ID
func GenerateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp if random generation fails
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// Middleware

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.HandlerFunc, log *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if request ID already exists in header
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = GenerateRequestID()
		}

		// Add request ID to context
		ctx := WithRequestID(r.Context(), requestID)
		r = r.WithContext(ctx)

		// Add request ID to response header
		w.Header().Set("X-Request-ID", requestID)

		next(w, r)
	}
}

// LoggingMiddleware logs all HTTP requests with timing information
func LoggingMiddleware(next http.HandlerFunc, log *Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request start
		log.LogRequestStart(r)

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Call the next handler
		next(wrapped, r)

		// Log request completion
		duration := time.Since(start)
		log.LogRequest(r, wrapped.statusCode, duration)
	}
}

// responseWriterWrapper wraps http.ResponseWriter to capture status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// Utility functions

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

// Default logger instance (for backward compatibility)
var defaultLogger *Logger

// InitDefault initializes the default logger
func InitDefault(cfg *config.LoggingConfig) error {
	logger, err := New(cfg)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

// GetDefault returns the default logger instance
func GetDefault() *Logger {
	if defaultLogger == nil {
		// Create a basic logger if none exists
		logger := logrus.New()
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		logger.SetOutput(os.Stdout)
		logger.SetLevel(logrus.InfoLevel)
		defaultLogger = &Logger{Logger: logger}
	}
	return defaultLogger
}

// Convenience functions using the default logger

// Debug logs a debug message
func Debug(args ...interface{}) {
	GetDefault().Debug(args...)
}

// Info logs an info message
func Info(args ...interface{}) {
	GetDefault().Info(args...)
}

// Warn logs a warning message
func Warn(args ...interface{}) {
	GetDefault().Warn(args...)
}

// Error logs an error message
func Error(args ...interface{}) {
	GetDefault().Error(args...)
}

// Fatal logs a fatal message and exits
func Fatal(args ...interface{}) {
	GetDefault().Fatal(args...)
}

// Debugf logs a formatted debug message
func Debugf(format string, args ...interface{}) {
	GetDefault().Debugf(format, args...)
}

// Infof logs a formatted info message
func Infof(format string, args ...interface{}) {
	GetDefault().Infof(format, args...)
}

// Warnf logs a formatted warning message
func Warnf(format string, args ...interface{}) {
	GetDefault().Warnf(format, args...)
}

// Errorf logs a formatted error message
func Errorf(format string, args ...interface{}) {
	GetDefault().Errorf(format, args...)
}

// Fatalf logs a formatted fatal message and exits
func Fatalf(format string, args ...interface{}) {
	GetDefault().Fatalf(format, args...)
}

// WithField returns a logger with a single field
func WithField(key string, value interface{}) *logrus.Entry {
	return GetDefault().WithField(key, value)
}

// WithFields returns a logger with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return GetDefault().WithFields(fields)
}

// WithError returns a logger with an error field
func WithError(err error) *logrus.Entry {
	return GetDefault().WithError(err)
}

// CombineMiddleware combines request ID and logging middleware
func CombineMiddleware(next http.HandlerFunc, log *Logger) http.HandlerFunc {
	return RequestIDMiddleware(LoggingMiddleware(next, log), log)
}

// CloseLogger closes any open file handles
func (l *Logger) Close() error {
	if closer, ok := l.Out.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
