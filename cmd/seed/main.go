package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func main() {
	// Initialize database
	db, err := database.NewSQLiteDB("./athena.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Seeding database with sample data...")

	// Create topics
	topics := []models.Topic{
		{Name: "Java Fundamentals", Description: "Core Java programming concepts", CreatedAt: time.Now()},
		{Name: "Python Basics", Description: "Introduction to Python programming", CreatedAt: time.Now()},
		{Name: "SQL Database", Description: "Structured Query Language basics", CreatedAt: time.Now()},
	}

	for i := range topics {
		if err := db.CreateTopic(&topics[i]); err != nil {
			log.Printf("Warning: Failed to create topic %s: %v", topics[i].Name, err)
		} else {
			fmt.Printf("Created topic: %s (ID: %d)\n", topics[i].Name, topics[i].ID)
		}
	}

	// Create Java questions
	javaQuestions := []struct {
		text       string
		difficulty int
		tags       []string
		answers    []struct {
			text      string
			isCorrect bool
		}
	}{
		{
			text:       "What is the correct way to declare a variable in Java?",
			difficulty: 1,
			tags:       []string{"variables", "syntax"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"int x = 5;", true},
				{"variable x = 5;", false},
				{"x = 5;", false},
				{"declare int x = 5;", false},
			},
		},
		{
			text:       "Which keyword is used to create a class in Java?",
			difficulty: 1,
			tags:       []string{"classes", "syntax"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"class", true},
				{"Class", false},
				{"new", false},
				{"object", false},
			},
		},
		{
			text:       "What does 'public static void main(String[] args)' mean?",
			difficulty: 2,
			tags:       []string{"methods", "main-method"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"Entry point of a Java program", true},
				{"A constructor", false},
				{"A variable declaration", false},
				{"A class definition", false},
			},
		},
		{
			text:       "What is the result of 5 / 2 in Java (integer division)?",
			difficulty: 2,
			tags:       []string{"operators", "arithmetic"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"2", true},
				{"2.5", false},
				{"3", false},
				{"Error", false},
			},
		},
		{
			text:       "Which of the following is a valid Java for loop?",
			difficulty: 3,
			tags:       []string{"loops", "control-structures"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"for (int i = 0; i < 10; i++) { }", true},
				{"for i in range(10): { }", false},
				{"for (i = 0; i < 10) { }", false},
				{"loop (int i = 0; i < 10; i++) { }", false},
			},
		},
	}

	if len(topics) > 0 {
		topicID := topics[0].ID
		for i, q := range javaQuestions {
			question := models.Question{
				TopicID:         topicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: fmt.Sprintf("%d.%d", topicID, i+1),
				CreatedAt:       time.Now(),
			}

			if err := db.CreateQuestion(&question); err != nil {
				log.Printf("Warning: Failed to create question: %v", err)
				continue
			}

			fmt.Printf("Created question: %s (ID: %d)\n", truncate(question.QuestionText, 50), question.ID)

			// Create answers for the question
			for _, a := range q.answers {
				answer := models.Answer{
					QuestionID: question.ID,
					AnswerText: a.text,
					IsCorrect:  a.isCorrect,
				}
				if err := db.CreateAnswer(&answer); err != nil {
					log.Printf("Warning: Failed to create answer: %v", err)
				}
			}
		}
	}

	// Create Python questions
	pythonQuestions := []struct {
		text       string
		difficulty int
		tags       []string
		answers    []struct {
			text      string
			isCorrect bool
		}
	}{
		{
			text:       "How do you print 'Hello World' in Python?",
			difficulty: 1,
			tags:       []string{"print", "syntax"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"print('Hello World')", true},
				{"echo 'Hello World'", false},
				{"System.out.println('Hello World')", false},
				{"console.log('Hello World')", false},
			},
		},
		{
			text:       "What is the correct way to create a list in Python?",
			difficulty: 1,
			tags:       []string{"lists", "data-structures"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"my_list = [1, 2, 3]", true},
				{"my_list = (1, 2, 3)", false},
				{"my_list = {1, 2, 3}", false},
				{"my_list = <1, 2, 3>", false},
			},
		},
		{
			text:       "Which keyword is used to define a function in Python?",
			difficulty: 1,
			tags:       []string{"functions", "syntax"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"def", true},
				{"function", false},
				{"func", false},
				{"define", false},
			},
		},
	}

	if len(topics) > 1 {
		topicID := topics[1].ID
		for i, q := range pythonQuestions {
			question := models.Question{
				TopicID:         topicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: fmt.Sprintf("%d.%d", topicID, i+1),
				CreatedAt:       time.Now(),
			}

			if err := db.CreateQuestion(&question); err != nil {
				log.Printf("Warning: Failed to create question: %v", err)
				continue
			}

			fmt.Printf("Created question: %s (ID: %d)\n", truncate(question.QuestionText, 50), question.ID)

			// Create answers for the question
			for _, a := range q.answers {
				answer := models.Answer{
					QuestionID: question.ID,
					AnswerText: a.text,
					IsCorrect:  a.isCorrect,
				}
				if err := db.CreateAnswer(&answer); err != nil {
					log.Printf("Warning: Failed to create answer: %v", err)
				}
			}
		}
	}

	// Create SQL questions
	sqlQuestions := []struct {
		text       string
		difficulty int
		tags       []string
		answers    []struct {
			text      string
			isCorrect bool
		}
	}{
		{
			text:       "Which SQL statement is used to retrieve data from a database?",
			difficulty: 1,
			tags:       []string{"select", "queries"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"SELECT", true},
				{"GET", false},
				{"RETRIEVE", false},
				{"FETCH", false},
			},
		},
		{
			text:       "Which clause is used to filter results in SQL?",
			difficulty: 2,
			tags:       []string{"where", "queries"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"WHERE", true},
				{"FILTER", false},
				{"HAVING", false},
				{"IF", false},
			},
		},
		{
			text:       "What does JOIN do in SQL?",
			difficulty: 3,
			tags:       []string{"joins", "queries"},
			answers: []struct {
				text      string
				isCorrect bool
			}{
				{"Combines rows from two or more tables", true},
				{"Adds new columns to a table", false},
				{"Deletes rows from a table", false},
				{"Creates a new table", false},
			},
		},
	}

	if len(topics) > 2 {
		topicID := topics[2].ID
		for i, q := range sqlQuestions {
			question := models.Question{
				TopicID:         topicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: fmt.Sprintf("%d.%d", topicID, i+1),
				CreatedAt:       time.Now(),
			}

			if err := db.CreateQuestion(&question); err != nil {
				log.Printf("Warning: Failed to create question: %v", err)
				continue
			}

			fmt.Printf("Created question: %s (ID: %d)\n", truncate(question.QuestionText, 50), question.ID)

			// Create answers for the question
			for _, a := range q.answers {
				answer := models.Answer{
					QuestionID: question.ID,
					AnswerText: a.text,
					IsCorrect:  a.isCorrect,
				}
				if err := db.CreateAnswer(&answer); err != nil {
					log.Printf("Warning: Failed to create answer: %v", err)
				}
			}
		}
	}

	fmt.Println("\nDatabase seeding complete!")
	fmt.Println("You can now start the server with: go run main.go")
}
