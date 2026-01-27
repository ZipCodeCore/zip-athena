package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ZipCodeCore/zip-athena/api"
	"github.com/ZipCodeCore/zip-athena/database"
)

func main() {
	// Initialize database
	db := database.NewInMemoryDB()

	// Create API handler
	handler := api.NewHandler(db)

	// Setup routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateUserHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetUserHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateTopicHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.ListTopicsHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/questions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handler.CreateQuestionHandler(w, r)
		} else if r.Method == http.MethodGet {
			handler.GetQuestionsByTopicHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

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
	log.Fatal(http.ListenAndServe(port, nil))
}
