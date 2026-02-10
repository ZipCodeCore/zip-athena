package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZipCodeCore/zip-athena/api"
	"github.com/ZipCodeCore/zip-athena/config"
	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/logger"
	"github.com/ZipCodeCore/zip-athena/router"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(&cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	// Set as default logger
	if err := logger.InitDefault(&cfg.Logging); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to set default logger: %v\n", err)
		os.Exit(1)
	}

	log.WithFields(map[string]interface{}{
		"port":      cfg.Server.Port,
		"db_type":   cfg.Database.Type,
		"db_path":   cfg.Database.Path,
		"log_level": cfg.Logging.Level,
	}).Info("Starting Athena API server")

	// Initialize database
	var db database.Database
	if cfg.Database.Type != "sqlite" {
		log.Fatal("Only SQLite is supported. Set DB_TYPE=sqlite")
	} else {
		db, err = database.NewSQLiteDB(cfg.Database.Path)
		if err != nil {
			log.WithError(err).Fatal("Failed to initialize database")
		}
	}
	defer db.Close()

	log.Info("Database initialized successfully")

	// Create API handler with logger
	handler := api.NewHandlerWithLogger(db, log)

	// Initialize JWT manager
	jwtManager := api.NewJWTManager(&cfg.Security)
	handler.SetJWTManager(jwtManager)

	// Create router with all routes
	r := router.New(handler)

	// Initialize rate limiter (60 requests/minute, burst of 20)
	rateLimiter := api.NewRateLimiter(60, 20)

	// Apply global middleware chain:
	//   SecurityHeaders -> CORS -> RateLimit -> RequestSize -> ContentType -> JWTOptional -> route handler
	globalHandler := api.Chain(r.ServeHTTP,
		api.SecurityHeadersMiddleware(),
		api.CORSMiddleware(&cfg.CORS),
		api.RateLimitMiddleware(rateLimiter),
		api.RequestSizeMiddleware(10*1024*1024), // 10MB max request body
		api.ContentTypeMiddleware([]string{"application/json", "multipart/form-data", "text/csv"}),
		api.JWTOptionalMiddleware(jwtManager),
	)

	// Create HTTP server with timeouts
	server := &http.Server{
		Addr:         cfg.Server.GetAddress(),
		Handler:      globalHandler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Graceful shutdown handling
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.WithError(err).Error("Server shutdown error")
		}

		log.Info("Server stopped")
	}()

	log.Infof("Athena API server listening on %s", cfg.Server.GetAddress())
	log.Infof("Database: %s", cfg.Database.Path)

	// Start server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("Server failed to start")
	}
}
