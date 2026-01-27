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
			handler.ListTopicsHandler(w, r)
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
	http.HandleFunc("/progress", handler.GetProgressHandler)

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

	port := ":8080"
	fmt.Printf("Athena API server starting on port %s\n", port)
	fmt.Printf("Database: %s\n", dbPath)
	log.Fatal(http.ListenAndServe(port, nil))
}

