package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ZipCodeCore/zip-athena/api"
	"github.com/ZipCodeCore/zip-athena/config"
	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/logger"
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
		"port":     cfg.Server.Port,
		"db_type":  cfg.Database.Type,
		"db_path":  cfg.Database.Path,
		"log_level": cfg.Logging.Level,
	}).Info("Starting Athena API server")

	// Initialize database
	var db database.Database
	if cfg.Database.Type == "postgres" {
		log.Fatal("PostgreSQL support not yet implemented")
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

	// Health and metrics endpoints (no auth required)
	http.HandleFunc("/health", handler.HealthHandler)
	http.HandleFunc("/ready", handler.ReadyHandler)
	http.HandleFunc("/metrics", handler.MetricsHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Authentication routes
	http.HandleFunc("/auth/register", handler.RegisterHandler)
	http.HandleFunc("/auth/login", handler.LoginHandler)
	http.HandleFunc("/auth/logout", handler.LogoutHandler)
	http.HandleFunc("/auth/refresh", handler.RefreshTokenHandler)

	// User routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateUserHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Topic routes
	http.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTopicHandler(w, r)
		} else if r.Method == http.MethodGet {
			// Check if ID parameter is provided
			if r.URL.Query().Get("id") != "" {
				handler.GetTopicHandler(w, r)
			} else {
				handler.ListTopicsHandler(w, r)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Question routes
	http.HandleFunc("/questions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateQuestionHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetQuestionsByTopicHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Quiz routes
	http.HandleFunc("/quiz/generate", handler.GenerateQuizHandler)
	http.HandleFunc("/quiz/submit-answer", handler.SubmitAnswerHandler)
	http.HandleFunc("/quiz/complete", handler.CompleteQuizHandler)
	http.HandleFunc("/quiz/retake", handler.RetakeQuizHandler)
	http.HandleFunc("/progress", handler.GetProgressHandler)

	// Quizzes route (for getting quiz details by ID)
	http.HandleFunc("/quizzes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handler.GetQuizHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Quiz attempt routes
	http.HandleFunc("/quiz-attempts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateQuizAttemptHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetQuizAttemptsByUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Phase 2: FSRS routes
	http.HandleFunc("/fsrs/review", handler.ReviewQuestionHandler)
	http.HandleFunc("/fsrs/due", handler.GetDueQuestionsHandler)
	http.HandleFunc("/fsrs/quiz/generate", handler.GenerateAdaptiveQuizHandler)

	// Phase 2: Admin routes
	http.HandleFunc("/admin/users", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListUsersHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.CreateUserHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateUserHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteUserHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/admin/questions", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListQuestionsHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.CreateQuestionHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateQuestionHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteQuestionHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/admin/quizzes", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListQuizzesHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.CreateQuizHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateQuizHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteQuizHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/admin/topics", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				handler.UpdateTopicHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteTopicHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/admin/answers", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListAnswersHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.CreateAnswerHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateAnswerHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteAnswerHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 3: Instructor Analytics routes
	http.HandleFunc("/instructor/students", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.ListStudentsHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/students/analytics", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.GetStudentAnalyticsHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/students/quiz-history", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.GetStudentQuizHistoryHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/students/weak-areas", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.GetStudentWeakAreasHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/student-progress", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.GetStudentProgressHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/cohort-analytics", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.GetCohortAnalyticsHandler(w, r)
		})(w, r)
	})

	// Phase 3: Import/Export routes
	http.HandleFunc("/instructor/quizzes/export", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.ExportQuizHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/quizzes/import", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.ImportQuizHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/instructor/questions/export", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			handler.ExportQuestionsHandler(w, r)
		})(w, r)
	})

	http.HandleFunc("/instructor/questions/import", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.ImportQuestionsHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 3: Topic Management routes (instructor)
	http.HandleFunc("/instructor/topics", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.CreateInstructorTopicHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateInstructorTopicHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteInstructorTopicHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/instructor/topics/reorder", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				handler.ReorderTopicsHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 3: Test management routes
	http.HandleFunc("/admin/tests", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListTestsHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.CreateTestHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateTestHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteTestHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 3: Quiz-Question relationship routes
	http.HandleFunc("/admin/quiz-questions", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.GetQuizQuestionsHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.AddQuestionToQuizHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.RemoveQuestionFromQuizHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 3: Test-Question relationship routes
	http.HandleFunc("/admin/test-questions", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.GetTestQuestionsHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.AddQuestionToTestHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.RemoveQuestionFromTestHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 4: Quiz Attempt History routes
	http.HandleFunc("/students/quiz-history", handler.GetQuizHistoryHandler)
	http.HandleFunc("/students/quiz-history/", handler.GetQuizAttemptDetailsHandler)

	// Phase 4: Timed Test Taking routes
	http.HandleFunc("/tests/available", handler.GetAvailableTestsHandler)
	http.HandleFunc("/tests/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/start") {
			handler.StartTestHandler(w, r)
		} else {
			handler.GetTestDetailsHandler(w, r)
		}
	})
	http.HandleFunc("/test-attempts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/answer") {
			handler.UpdateTestAnswerHandler(w, r)
		} else if strings.Contains(r.URL.Path, "/submit") {
			handler.SubmitTestHandler(w, r)
		} else if strings.Contains(r.URL.Path, "/results") {
			handler.GetTestResultsHandler(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	// Phase 4: Tag-Based Review routes
	http.HandleFunc("/tags", handler.GetTagsHandler)
	http.HandleFunc("/tags/stats", handler.GetTagStatsHandler)
	http.HandleFunc("/quiz/generate/by-tags", handler.GenerateQuizByTagsHandler)

	// Learning Resources routes
	http.HandleFunc("/learning/outline", handler.GetCurriculumOutlineHandler)
	http.HandleFunc("/learning/phase/", handler.GetPhaseResourcesHandler)
	http.HandleFunc("/learning/resource/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/content") {
			handler.GetLearningResourceContentHandler(w, r)
		} else {
			handler.GetLearningResourceHandler(w, r)
		}
	})
	http.HandleFunc("/learning/search", handler.SearchLearningResourcesHandler)
	http.HandleFunc("/learning/track-view", handler.TrackLearningViewHandler)
	http.HandleFunc("/quiz/question/", func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/learning") {
			handler.GetQuestionLearningLinksHandler(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	})

	// Phase 4: Instructor Test Management routes
	http.HandleFunc("/instructor/tests/", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/publish") {
				handler.PublishTestHandler(w, r)
			} else if strings.Contains(r.URL.Path, "/results") {
				handler.GetTestResultsInstructorHandler(w, r)
			} else if strings.Contains(r.URL.Path, "/export") {
				handler.ExportTestResultsHandler(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		})(w, r)
	})

	// Phase 4: Question Performance Analytics routes
	http.HandleFunc("/instructor/questions/analytics", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(handler.GetQuestionAnalyticsHandler)(w, r)
	})

	// Phase 5: Instructor Question Management routes (must come before the path-based matcher)
	http.HandleFunc("/instructor/questions", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListInstructorQuestionsHandler(w, r)
			} else if r.Method == http.MethodPost {
				handler.CreateInstructorQuestionHandler(w, r)
			} else if r.Method == http.MethodPut {
				handler.UpdateInstructorQuestionHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeactivateInstructorQuestionHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	http.HandleFunc("/instructor/questions/", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.URL.Path, "/analytics") {
				handler.GetQuestionDetailedAnalyticsHandler(w, r)
			} else if strings.Contains(r.URL.Path, "/answer-distribution") {
				handler.GetAnswerDistributionHandler(w, r)
			} else if strings.Contains(r.URL.Path, "/difficulty") {
				handler.UpdateQuestionDifficultyHandler(w, r)
			} else if strings.Contains(r.URL.Path, "/preview") {
				// Phase 5: Question preview
				handler.PreviewInstructorQuestionHandler(w, r)
			} else {
				http.Error(w, "Not found", http.StatusNotFound)
			}
		})(w, r)
	})

	// Phase 5: Admin CSV Import routes
	http.HandleFunc("/admin/questions/import", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.ImportQuestionsCSVHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/questions/import/validate", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.ValidateQuestionsCSVHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/questions/import/confirm", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.ConfirmImportQuestionsCSVHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 5: Admin Tag Management routes
	http.HandleFunc("/admin/tags", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.ListTagsHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.DeleteTagHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/tags/merge", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.MergeTagsHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/tags/questions", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.GetQuestionsByTagHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Phase 5: Admin Quality Review routes
	http.HandleFunc("/admin/questions/flagged", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				handler.GetFlaggedQuestionsHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/questions/auto-flag", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.AutoFlagQuestionsHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/questions/flag", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost {
				handler.FlagQuestionHandler(w, r)
			} else if r.Method == http.MethodDelete {
				handler.UnflagQuestionHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})
	http.HandleFunc("/admin/questions/resolve", func(w http.ResponseWriter, r *http.Request) {
		handler.AdminMiddleware(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPut {
				handler.ResolveQuestionHandler(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		})(w, r)
	})

	// Initialize rate limiter (60 requests/minute, burst of 20)
	rateLimiter := api.NewRateLimiter(60, 20)

	// Apply global middleware chain:
	//   SecurityHeaders -> CORS -> RateLimit -> RequestSize -> ContentType -> JWTOptional -> route handler
	globalHandler := api.Chain(http.DefaultServeMux.ServeHTTP,
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
