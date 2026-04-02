package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
	"golang.org/x/crypto/bcrypt"
)

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// loadQuestionsFromJSONFile reads questions from a JSON file and creates them in the database for the given topic.
// It skips questions that already exist (matched by question text).
func loadQuestionsFromJSONFile(db *database.SQLiteDB, topicID int, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var qf struct {
		Topic     string `json:"topic"`
		Questions []struct {
			Text       string   `json:"text"`
			Difficulty int      `json:"difficulty"`
			Tags       []string `json:"tags"`
			Answers    []struct {
				Text      string `json:"text"`
				IsCorrect bool   `json:"is_correct"`
			} `json:"answers"`
		} `json:"questions"`
	}

	if err := json.Unmarshal(data, &qf); err != nil {
		return fmt.Errorf("failed to parse JSON in %s: %w", filePath, err)
	}

	// Get existing questions to avoid duplicates
	existingQuestions, _ := db.GetQuestionsByTopic(topicID)
	existingTexts := make(map[string]bool)
	for _, eq := range existingQuestions {
		existingTexts[strings.ToLower(eq.QuestionText)] = true
	}
	startNum := len(existingQuestions) + 1

	imported := 0
	for i, q := range qf.Questions {
		if existingTexts[strings.ToLower(q.Text)] {
			fmt.Printf("Question already exists: %s\n", truncate(q.Text, 50))
			continue
		}

		topicQuestionID := fmt.Sprintf("%d.%d", topicID, startNum+i)
		question := models.Question{
			TopicID:         topicID,
			QuestionText:    q.Text,
			Difficulty:      q.Difficulty,
			Tags:            q.Tags,
			TopicQuestionID: topicQuestionID,
			IsActive:        true,
			CreatedAt:       time.Now(),
		}

		if err := db.CreateQuestion(&question); err != nil {
			log.Printf("Warning: Failed to create question: %v", err)
			continue
		}

		for _, a := range q.Answers {
			answer := models.Answer{
				QuestionID: question.ID,
				AnswerText: a.Text,
				IsCorrect:  a.IsCorrect,
			}
			if err := db.CreateAnswer(&answer); err != nil {
				log.Printf("Warning: Failed to create answer: %v", err)
			}
		}
		imported++
	}

	fmt.Printf("Loaded %d questions from %s\n", imported, filePath)
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func main() {
	// Initialize database
	db, err := database.NewSQLiteDB("./athena.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	fmt.Println("Seeding database with sample data...")

	// Seed Phase 1 data (topics and basic questions)
	if err := seedPhase1Data(db); err != nil {
		log.Fatalf("Failed to seed Phase 1 data: %v", err)
	}

	// Seed Phase 2 data (admin user, quizzes, additional questions)
	if err := seedPhase2Data(db); err != nil {
		log.Fatalf("Failed to seed Phase 2 data: %v", err)
	}

	// Seed quiz attempts with random answers (only save if score > 35%)
	if err := seedQuizAttempts(db); err != nil {
		log.Fatalf("Failed to seed quiz attempts: %v", err)
	}

	fmt.Println("\nDatabase seeding complete!")
	fmt.Println("You can now start the server with: go run main.go")
	fmt.Println("\nAdmin credentials:")
	fmt.Println("  Username: admin")
	fmt.Println("  Password: admin")
}

func seedPhase1Data(db *database.SQLiteDB) error {

	// Create topics
	topics := []models.Topic{
		{Name: "Java Fundamentals", Description: "Core Java programming concepts", Category: "java", CreatedAt: time.Now()},
		{Name: "Python Basics", Description: "Introduction to Python programming", Category: "python", CreatedAt: time.Now()},
		{Name: "SQL Database", Description: "Structured Query Language basics", Category: "sql", CreatedAt: time.Now()},
		{Name: "Java Code Tracing", Description: "Trace Java code snippets and determine the printed output. Each quiz has 5 questions.", Category: "java", CreatedAt: time.Now()},
		{Name: "Python Code Tracing", Description: "Trace Python code snippets and determine the printed output. Each quiz has 5 questions.", Category: "python", CreatedAt: time.Now()},
	}

	// Get existing topics once before the loop
	existingTopics, err := db.ListTopics()
	if err != nil {
		log.Printf("Warning: Failed to list topics: %v", err)
		existingTopics = []*models.Topic{} // Use empty list if query fails
	}

	for i := range topics {
		topicExists := false
		for _, existingTopic := range existingTopics {
			if existingTopic.Name == topics[i].Name {
				topics[i] = *existingTopic
				topicExists = true
				fmt.Printf("Topic already exists: %s (ID: %d)\n", topics[i].Name, topics[i].ID)
				break
			}
		}

		if !topicExists {
			if err := db.CreateTopic(&topics[i]); err != nil {
				log.Printf("Warning: Failed to create topic %s: %v", topics[i].Name, err)
			} else {
				fmt.Printf("Created topic: %s (ID: %d)\n", topics[i].Name, topics[i].ID)
			}
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
		// Get existing questions for this topic
		existingQuestions, _ := db.GetQuestionsByTopic(topicID)
		existingQuestionsMap := make(map[string]bool)
		for _, eq := range existingQuestions {
			existingQuestionsMap[eq.TopicQuestionID] = true
		}

		for i, q := range javaQuestions {
			topicQuestionID := fmt.Sprintf("%d.%d", topicID, i+1)

			// Check if question already exists
			if existingQuestionsMap[topicQuestionID] {
				fmt.Printf("Question already exists: %s (ID: %s)\n", truncate(q.text, 50), topicQuestionID)
				continue
			}

			question := models.Question{
				TopicID:         topicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: topicQuestionID,
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
		// Get existing questions for this topic
		existingQuestions, _ := db.GetQuestionsByTopic(topicID)
		existingQuestionsMap := make(map[string]bool)
		for _, eq := range existingQuestions {
			existingQuestionsMap[eq.TopicQuestionID] = true
		}

		for i, q := range pythonQuestions {
			topicQuestionID := fmt.Sprintf("%d.%d", topicID, i+1)

			// Check if question already exists
			if existingQuestionsMap[topicQuestionID] {
				fmt.Printf("Question already exists: %s (ID: %s)\n", truncate(q.text, 50), topicQuestionID)
				continue
			}

			question := models.Question{
				TopicID:         topicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: topicQuestionID,
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
		// Get existing questions for this topic
		existingQuestions, _ := db.GetQuestionsByTopic(topicID)
		existingQuestionsMap := make(map[string]bool)
		for _, eq := range existingQuestions {
			existingQuestionsMap[eq.TopicQuestionID] = true
		}

		for i, q := range sqlQuestions {
			topicQuestionID := fmt.Sprintf("%d.%d", topicID, i+1)

			// Check if question already exists
			if existingQuestionsMap[topicQuestionID] {
				fmt.Printf("Question already exists: %s (ID: %s)\n", truncate(q.text, 50), topicQuestionID)
				continue
			}

			question := models.Question{
				TopicID:         topicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: topicQuestionID,
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

	if len(topics) > 3 && topics[3].ID > 0 {
		if err := loadQuestionsFromJSONFile(db, topics[3].ID, "./questions/java-code-tracing.json"); err != nil {
			log.Printf("Warning: Failed to load Java Code Tracing questions: %v", err)
		}
	}

	if len(topics) > 4 && topics[4].ID > 0 {
		if err := loadQuestionsFromJSONFile(db, topics[4].ID, "./questions/python-code-tracing.json"); err != nil {
			log.Printf("Warning: Failed to load Python Code Tracing questions: %v", err)
		}
	}

	return nil
}

func seedPhase2Data(db *database.SQLiteDB) error {
	fmt.Println("Seeding Phase 2 data...")

	// Create admin user
	adminPassword, err := hashPassword("admin")
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	// Check if admin already exists
	_, err = db.GetUserByUsername("admin")
	if err == database.ErrNotFound {
		// Create admin user
		_, err = db.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, category, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			"admin", "admin@athena.com", adminPassword, true, "java", 0, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}
		fmt.Println("Created admin user (username: admin, password: admin)")
	} else {
		fmt.Println("Admin user already exists")
	}

	// Create Paul admin user
	paulPassword, err := hashPassword("zipcode0")
	if err != nil {
		return fmt.Errorf("failed to hash Paul password: %w", err)
	}

	_, err = db.GetUserByUsername("paul")
	if err == database.ErrNotFound {
		_, err = db.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, category, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			"paul", "paul@zipcodewilmington.com", paulPassword, true, "java", 0, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create Paul admin user: %w", err)
		}
		fmt.Println("Created Paul admin user (username: paul, password: zipcode0)")
	} else {
		fmt.Println("Paul admin user already exists")
	}

	// Create Kristofer as a data category student (non-admin)
	krisPassword, err := hashPassword("zipcode0")
	if err != nil {
		return fmt.Errorf("failed to hash Kristofer password: %w", err)
	}

	_, err = db.GetUserByUsername("kris")
	if err == database.ErrNotFound {
		_, err = db.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, category, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			"kris", "kris@zipcodewilmington.com", krisPassword, false, "data", 0, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create Kristofer user: %w", err)
		}
		fmt.Println("Created Kristofer as data category student (username: kris, password: zipcode0)")
	} else {
		fmt.Println("Kristofer user already exists")
	}

	// Create regular java student user
	userPassword, err := hashPassword("user")
	if err != nil {
		return fmt.Errorf("failed to hash user password: %w", err)
	}

	// Check if regular user already exists
	_, err = db.GetUserByUsername("user")
	if err == database.ErrNotFound {
		// Create regular user (Java category)
		_, err = db.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, category, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			"user", "user@zipcodewilmington.com", userPassword, false, "java", 0, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create regular user: %w", err)
		}
		fmt.Println("Created regular user (username: user, password: user, category: java)")
	} else {
		fmt.Println("Regular user already exists")
	}

	// Create data student user
	dataPassword, err := hashPassword("data")
	if err != nil {
		return fmt.Errorf("failed to hash data user password: %w", err)
	}

	_, err = db.GetUserByUsername("data")
	if err == database.ErrNotFound {
		// Create data user (Data category)
		_, err = db.DB().Exec(
			"INSERT INTO users (username, email, password, is_admin, category, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
			"data", "data@zipcodewilmington.com", dataPassword, false, "data", 0, time.Now(),
		)
		if err != nil {
			return fmt.Errorf("failed to create data user: %w", err)
		}
		fmt.Println("Created data user (username: data, password: data, category: data)")
	} else {
		fmt.Println("Data user already exists")
	}

	// Create quizzes
	topics, err := db.ListTopics()
	if err != nil {
		return fmt.Errorf("failed to list topics: %w", err)
	}

	for i, topic := range topics {
		quizName := fmt.Sprintf("%s.%d", topic.Name, 1)
		quiz := &models.Quiz{
			Name:        quizName,
			TopicID:     topic.ID,
			Description: fmt.Sprintf("Quiz 1 for %s", topic.Name),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := db.CreateQuiz(quiz); err != nil {
			log.Printf("Warning: Failed to create quiz %s: %v", quizName, err)
		} else {
			fmt.Printf("Created quiz: %s (ID: %d)\n", quizName, quiz.ID)
		}

		// Create a second quiz for each topic
		quiz2Name := fmt.Sprintf("%s.%d", topic.Name, 2)
		quiz2 := &models.Quiz{
			Name:        quiz2Name,
			TopicID:     topic.ID,
			Description: fmt.Sprintf("Quiz 2 for %s - Advanced", topic.Name),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := db.CreateQuiz(quiz2); err != nil {
			log.Printf("Warning: Failed to create quiz %s: %v", quiz2Name, err)
		} else {
			fmt.Printf("Created quiz: %s (ID: %d)\n", quiz2Name, quiz2.ID)
		}

		// For Java and Python, add extra quizzes
		if i < 2 {
			quiz3Name := fmt.Sprintf("%s.%d", topic.Name, 3)
			quiz3 := &models.Quiz{
				Name:        quiz3Name,
				TopicID:     topic.ID,
				Description: fmt.Sprintf("Quiz 3 for %s - Expert", topic.Name),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			if err := db.CreateQuiz(quiz3); err != nil {
				log.Printf("Warning: Failed to create quiz %s: %v", quiz3Name, err)
			} else {
				fmt.Printf("Created quiz: %s (ID: %d)\n", quiz3Name, quiz3.ID)
			}
		}
	}

	// Add more Java questions
	if len(topics) > 0 {
		javaTopicID := topics[0].ID
		existingJavaQuestions, _ := db.GetQuestionsByTopic(javaTopicID)
		existingQuestionsMap := make(map[string]bool)
		for _, eq := range existingJavaQuestions {
			existingQuestionsMap[eq.TopicQuestionID] = true
		}
		startNum := len(existingJavaQuestions) + 1

		additionalJavaQuestions := []struct {
			text       string
			difficulty int
			tags       []string
			answers    []struct {
				text      string
				isCorrect bool
			}
		}{
			{
				text:       "What is the output of: System.out.println(10 + 20 + \"Hello\");",
				difficulty: 2,
				tags:       []string{"operators", "strings"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"30Hello", true},
					{"1020Hello", false},
					{"Hello1020", false},
					{"Compilation error", false},
				},
			},
			{
				text:       "Which of the following is NOT a Java keyword?",
				difficulty: 2,
				tags:       []string{"keywords", "syntax"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"include", true},
					{"extends", false},
					{"implements", false},
					{"interface", false},
				},
			},
			{
				text:       "What is the default value of a boolean variable in Java?",
				difficulty: 1,
				tags:       []string{"variables", "primitives"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"false", true},
					{"true", false},
					{"0", false},
					{"null", false},
				},
			},
			{
				text:       "Which method is the entry point of a Java application?",
				difficulty: 1,
				tags:       []string{"methods", "main-method"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"public static void main(String[] args)", true},
					{"public void start()", false},
					{"public static void run()", false},
					{"public void main()", false},
				},
			},
			{
				text:       "What is inheritance in Java?",
				difficulty: 3,
				tags:       []string{"oop", "inheritance"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"A mechanism where one class acquires properties of another class", true},
					{"A way to hide implementation details", false},
					{"A method to create multiple instances", false},
					{"A technique for error handling", false},
				},
			},
			{
				text:       "Which access modifier makes a member accessible only within the same package?",
				difficulty: 3,
				tags:       []string{"access-modifiers", "packages"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"default (no modifier)", true},
					{"public", false},
					{"private", false},
					{"protected", false},
				},
			},
			{
				text:       "What does the 'final' keyword do when applied to a variable?",
				difficulty: 2,
				tags:       []string{"keywords", "variables"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"Makes the variable a constant", true},
					{"Makes the variable global", false},
					{"Makes the variable private", false},
					{"Makes the variable static", false},
				},
			},
			{
				text:       "Which collection does NOT allow duplicate elements?",
				difficulty: 3,
				tags:       []string{"collections", "data-structures"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"Set", true},
					{"List", false},
					{"ArrayList", false},
					{"LinkedList", false},
				},
			},
			{
				text:       "What is polymorphism in Java?",
				difficulty: 4,
				tags:       []string{"oop", "polymorphism"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"The ability of an object to take many forms", true},
					{"The ability to create multiple classes", false},
					{"The process of hiding data", false},
					{"The process of creating objects", false},
				},
			},
			{
				text:       "What is the purpose of the 'super' keyword?",
				difficulty: 3,
				tags:       []string{"keywords", "inheritance"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"To refer to the parent class", true},
					{"To create a new instance", false},
					{"To define a static method", false},
					{"To import a package", false},
				},
			},
			{
				text:       "Which exception is thrown when dividing by zero with integers?",
				difficulty: 2,
				tags:       []string{"exceptions", "arithmetic"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"ArithmeticException", true},
					{"NullPointerException", false},
					{"NumberFormatException", false},
					{"IllegalArgumentException", false},
				},
			},
			{
				text:       "What is encapsulation?",
				difficulty: 3,
				tags:       []string{"oop", "encapsulation"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"Bundling data and methods that operate on that data", true},
					{"Creating multiple instances of a class", false},
					{"Inheriting from a parent class", false},
					{"Overriding methods", false},
				},
			},
			{
				text:       "Which loop is guaranteed to execute at least once?",
				difficulty: 2,
				tags:       []string{"loops", "control-structures"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"do-while", true},
					{"while", false},
					{"for", false},
					{"enhanced for", false},
				},
			},
			{
				text:       "What is the size of an int in Java?",
				difficulty: 2,
				tags:       []string{"primitives", "data-types"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"32 bits", true},
					{"16 bits", false},
					{"64 bits", false},
					{"8 bits", false},
				},
			},
			{
				text:       "Which keyword is used to prevent method overriding?",
				difficulty: 3,
				tags:       []string{"keywords", "inheritance"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"final", true},
					{"static", false},
					{"private", false},
					{"abstract", false},
				},
			},
		}

		for i, q := range additionalJavaQuestions {
			topicQuestionID := fmt.Sprintf("%d.%d", javaTopicID, startNum+i)

			// Check if question already exists
			if existingQuestionsMap[topicQuestionID] {
				fmt.Printf("Question already exists: %d.%d\n", javaTopicID, startNum+i)
				continue
			}

			question := models.Question{
				TopicID:         javaTopicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: topicQuestionID,
				CreatedAt:       time.Now(),
			}

			if err := db.CreateQuestion(&question); err != nil {
				log.Printf("Warning: Failed to create question: %v", err)
				continue
			}

			fmt.Printf("Created Java question %d.%d\n", javaTopicID, startNum+i)

			for _, a := range q.answers {
				answer := models.Answer{
					QuestionID: question.ID,
					AnswerText: a.text,
					IsCorrect:  a.isCorrect,
				}
				db.CreateAnswer(&answer)
			}
		}
	}

	// Add more Python questions
	if len(topics) > 1 {
		pythonTopicID := topics[1].ID
		existingPythonQuestions, _ := db.GetQuestionsByTopic(pythonTopicID)
		existingQuestionsMap := make(map[string]bool)
		for _, eq := range existingPythonQuestions {
			existingQuestionsMap[eq.TopicQuestionID] = true
		}
		startNum := len(existingPythonQuestions) + 1

		additionalPythonQuestions := []struct {
			text       string
			difficulty int
			tags       []string
			answers    []struct {
				text      string
				isCorrect bool
			}
		}{
			{
				text:       "What is the output of: print(type([]))",
				difficulty: 2,
				tags:       []string{"lists", "types"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"&lt;class 'list'&gt;", true},
					{"&lt;class 'array'&gt;", false},
					{"&lt;class 'tuple'&gt;", false},
					{"&lt;class 'dict'&gt;", false},
				},
			},
			{
				text:       "Which of the following is mutable in Python?",
				difficulty: 2,
				tags:       []string{"data-structures", "mutability"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"list", true},
					{"tuple", false},
					{"string", false},
					{"int", false},
				},
			},
			{
				text:       "What does the 'pass' statement do?",
				difficulty: 1,
				tags:       []string{"keywords", "syntax"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"It's a null operation placeholder", true},
					{"It exits the function", false},
					{"It skips to the next iteration", false},
					{"It raises an exception", false},
				},
			},
			{
				text:       "How do you create a dictionary in Python?",
				difficulty: 1,
				tags:       []string{"dictionaries", "data-structures"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"my_dict = {'key': 'value'}", true},
					{"my_dict = ['key', 'value']", false},
					{"my_dict = ('key', 'value')", false},
					{"my_dict = <'key', 'value'>", false},
				},
			},
			{
				text:       "What is a lambda function?",
				difficulty: 3,
				tags:       []string{"functions", "lambda"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"A small anonymous function", true},
					{"A recursive function", false},
					{"A generator function", false},
					{"A class method", false},
				},
			},
			{
				text:       "Which method adds an element to the end of a list?",
				difficulty: 1,
				tags:       []string{"lists", "methods"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"append()", true},
					{"add()", false},
					{"insert()", false},
					{"extend()", false},
				},
			},
			{
				text:       "What is the output of: print(2 ** 3)",
				difficulty: 2,
				tags:       []string{"operators", "arithmetic"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"8", true},
					{"6", false},
					{"9", false},
					{"23", false},
				},
			},
			{
				text:       "What is list comprehension?",
				difficulty: 3,
				tags:       []string{"lists", "comprehension"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"A concise way to create lists", true},
					{"A method to compress lists", false},
					{"A way to sort lists", false},
					{"A type of loop", false},
				},
			},
			{
				text:       "Which keyword is used for exception handling?",
				difficulty: 2,
				tags:       []string{"exceptions", "keywords"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"try", true},
					{"catch", false},
					{"handle", false},
					{"error", false},
				},
			},
			{
				text:       "What does the 'self' parameter represent?",
				difficulty: 2,
				tags:       []string{"oop", "classes"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"The instance of the class", true},
					{"The class itself", false},
					{"The parent class", false},
					{"A static variable", false},
				},
			},
			{
				text:       "How do you create a tuple with a single element?",
				difficulty: 2,
				tags:       []string{"tuples", "syntax"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"(1,)", true},
					{"(1)", false},
					{"[1]", false},
					{"{1}", false},
				},
			},
			{
				text:       "What is the purpose of __init__?",
				difficulty: 2,
				tags:       []string{"oop", "classes"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"It's a constructor method", true},
					{"It's a destructor method", false},
					{"It's a static method", false},
					{"It's an abstract method", false},
				},
			},
			{
				text:       "Which operator is used for floor division?",
				difficulty: 2,
				tags:       []string{"operators", "arithmetic"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"//", true},
					{"/", false},
					{"%", false},
					{"**", false},
				},
			},
			{
				text:       "What does the 'range(5)' function return?",
				difficulty: 2,
				tags:       []string{"functions", "iteration"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"Numbers from 0 to 4", true},
					{"Numbers from 1 to 5", false},
					{"Numbers from 0 to 5", false},
					{"Numbers from 1 to 4", false},
				},
			},
			{
				text:       "What is inheritance in Python?",
				difficulty: 3,
				tags:       []string{"oop", "inheritance"},
				answers: []struct {
					text      string
					isCorrect bool
				}{
					{"A way to create a class based on another class", true},
					{"A method to copy objects", false},
					{"A type of exception", false},
					{"A way to import modules", false},
				},
			},
		}

		for i, q := range additionalPythonQuestions {
			topicQuestionID := fmt.Sprintf("%d.%d", pythonTopicID, startNum+i)

			// Check if question already exists
			if existingQuestionsMap[topicQuestionID] {
				fmt.Printf("Question already exists: %d.%d\n", pythonTopicID, startNum+i)
				continue
			}

			question := models.Question{
				TopicID:         pythonTopicID,
				QuestionText:    q.text,
				Difficulty:      q.difficulty,
				Tags:            q.tags,
				TopicQuestionID: topicQuestionID,
				CreatedAt:       time.Now(),
			}

			if err := db.CreateQuestion(&question); err != nil {
				log.Printf("Warning: Failed to create question: %v", err)
				continue
			}

			fmt.Printf("Created Python question %d.%d\n", pythonTopicID, startNum+i)

			for _, a := range q.answers {
				answer := models.Answer{
					QuestionID: question.ID,
					AnswerText: a.text,
					IsCorrect:  a.isCorrect,
				}
				db.CreateAnswer(&answer)
			}
		}
	}

	// Assign questions to quizzes
	fmt.Println("\nAssigning questions to quizzes...")
	quizzes, err := db.ListQuizzes()
	if err != nil {
		return fmt.Errorf("failed to list quizzes: %w", err)
	}

	for _, quiz := range quizzes {
		// Get questions for this quiz's topic
		questions, err := db.GetQuestionsByTopic(quiz.TopicID)
		if err != nil {
			log.Printf("Warning: Failed to get questions for topic %d: %v", quiz.TopicID, err)
			continue
		}

		if len(questions) == 0 {
			log.Printf("Warning: No questions found for topic %d", quiz.TopicID)
			continue
		}

		// Assign 5-10 random questions to each quiz
		numQuestions := 5 + rand.Intn(6) // 5-10 questions
		if numQuestions > len(questions) {
			numQuestions = len(questions)
		}

		// Shuffle and pick questions
		questionIndices := rand.Perm(len(questions))
		assignedCount := 0
		for i := 0; i < numQuestions && i < len(questionIndices); i++ {
			question := questions[questionIndices[i]]
			err := db.AddQuestionToQuiz(quiz.ID, question.ID, i+1)
			if err != nil {
				log.Printf("Warning: Failed to add question %d to quiz %d: %v", question.ID, quiz.ID, err)
			} else {
				assignedCount++
			}
		}
		fmt.Printf("  Assigned %d questions to quiz '%s'\n", assignedCount, quiz.Name)
	}

	fmt.Println("\nPhase 2 seeding complete!")
	return nil
}

func seedQuizAttempts(db *database.SQLiteDB) error {
	fmt.Println("\n=== Seeding Quiz Attempts ===")

	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Get all users
	users, err := db.ListUsers()
	if err != nil {
		return fmt.Errorf("failed to list users: %v", err)
	}

	// Get all quizzes
	quizzes, err := db.ListQuizzes()
	if err != nil {
		return fmt.Errorf("failed to list quizzes: %v", err)
	}

	if len(quizzes) == 0 {
		return fmt.Errorf("no quizzes found in database")
	}

	// For each user, create 1-2 successful quiz attempts
	for _, user := range users {
		fmt.Printf("\nCreating quiz attempts for user: %s (ID: %d)\n", user.Username, user.ID)

		attemptsNeeded := 1 + rand.Intn(2) // 1 or 2 attempts
		successfulAttempts := 0
		totalTries := 0
		maxTries := 100 // Safety limit

		for successfulAttempts < attemptsNeeded && totalTries < maxTries {
			totalTries++

			// Pick a random quiz
			quiz := quizzes[rand.Intn(len(quizzes))]

			// Get questions for this quiz
			questions, err := db.GetQuizQuestions(quiz.ID)
			if err != nil || len(questions) == 0 {
				fmt.Printf("  Warning: Quiz %d has no questions, skipping\n", quiz.ID)
				continue
			}

			// Create quiz attempt with random answers
			attempt := &models.QuizAttempt{
				UserID:          user.ID,
				StartTime:       time.Now().Add(-time.Duration(rand.Intn(168)) * time.Hour), // Random time in last week
				TotalQuestions:  len(questions),
				QuestionIDs:     make([]string, 0),
				QuestionResults: make(map[string]bool),
			}

			correctCount := 0

			// Answer each question with a slight bias toward correct answers
			// This simulates students who are somewhat competent (not pure random guessing)
			for _, question := range questions {
				// Get answers for this question
				answers, err := db.GetAnswersByQuestion(question.ID)
				if err != nil || len(answers) == 0 {
					continue
				}

				// 60% chance to pick the correct answer, 40% chance to pick randomly
				// This gives us expected ~60% accuracy which is above our 35% threshold
				var selectedAnswer *models.Answer
				if rand.Float64() < 0.60 {
					// Find and pick the correct answer
					for _, ans := range answers {
						if ans.IsCorrect {
							selectedAnswer = ans
							break
						}
					}
					if selectedAnswer == nil {
						// Fallback to random if no correct answer found
						selectedAnswer = answers[rand.Intn(len(answers))]
					}
				} else {
					// Pick a random answer
					selectedAnswer = answers[rand.Intn(len(answers))]
				}

				// Record the question
				attempt.QuestionIDs = append(attempt.QuestionIDs, question.TopicQuestionID)

				// Check if correct
				isCorrect := selectedAnswer.IsCorrect
				attempt.QuestionResults[question.TopicQuestionID] = isCorrect

				if isCorrect {
					correctCount++
				}
			}

			// Calculate percentage
			percentage := float64(correctCount) / float64(len(questions)) * 100

			// Only save if percentage > 35% (lower threshold to get some data)
			if percentage > 30.0 {
				attempt.CorrectAnswers = correctCount
				completedTime := attempt.StartTime.Add(time.Duration(5+rand.Intn(20)) * time.Minute)
				attempt.CompletedTime = &completedTime

				// Create the quiz attempt
				if err := db.CreateQuizAttempt(attempt); err != nil {
					fmt.Printf("  Error creating quiz attempt: %v\n", err)
					continue
				}

				fmt.Printf("  ✓ Created attempt %d/%d: Quiz '%s', Score: %d/%d (%.1f%%)\n",
					successfulAttempts, attemptsNeeded, quiz.Name, correctCount, len(questions), percentage)

				// Also create question_attempts records for each question
				for _, question := range questions {
					isCorrect := attempt.QuestionResults[question.TopicQuestionID]

					// Create question attempt record
					qa := &models.QuestionAttempt{
						UserID:             user.ID,
						QuestionID:         question.ID,
						TopicQuestionID:    question.TopicQuestionID,
						CorrectAttempts:    0,
						IncorrectAttempts:  0,
						LastAttemptTime:    completedTime,
						LastAttemptCorrect: isCorrect,
					}

					if isCorrect {
						qa.CorrectAttempts = 1
					} else {
						qa.IncorrectAttempts = 1
					}

					db.CreateQuestionAttempt(qa)
				}
			} else {
				fmt.Printf("  ✗ Attempt failed (%.1f%% < 35%%), retrying...\n", percentage)
			}
		}

		if successfulAttempts < attemptsNeeded {
			fmt.Printf("  Warning: Only created %d/%d attempts after %d tries\n", successfulAttempts, attemptsNeeded, maxTries)
		}
	}

	fmt.Println("\nQuiz attempts seeding complete!")
	return nil
}
