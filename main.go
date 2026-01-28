package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ZipCodeCore/zip-athena/api"
	"github.com/ZipCodeCore/zip-athena/database"
)

func main() {
	// Initialize database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./athena.db"
	}

	db, err := database.NewSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create API handler
	handler := api.NewHandler(db)

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Authentication routes
	http.HandleFunc("/auth/register", handler.RegisterHandler)
	http.HandleFunc("/auth/login", handler.LoginHandler)

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

	port := ":9206"
	fmt.Printf("Athena API server starting on port %s\n", port)
	fmt.Printf("Database: %s\n", dbPath)
	log.Fatal(http.ListenAndServe(port, nil))
}
