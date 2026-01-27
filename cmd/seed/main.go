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

	// Seed Phase 1 data (topics and basic questions)
	if err := seedPhase1Data(db); err != nil {
		log.Fatalf("Failed to seed Phase 1 data: %v", err)
	}

	// Seed Phase 2 data (admin user, quizzes, additional questions)
	if err := seedPhase2Data(db); err != nil {
		log.Fatalf("Failed to seed Phase 2 data: %v", err)
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

	return nil
}
package main

import (
"fmt"
"log"
"time"

"github.com/ZipCodeCore/zip-athena/database"
"github.com/ZipCodeCore/zip-athena/models"
"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
return string(bytes), err
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
"INSERT INTO users (username, email, password, is_admin, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?)",
"admin", "admin@athena.com", adminPassword, true, 0, time.Now(),
)
if err != nil {
return fmt.Errorf("failed to create admin user: %w", err)
}
fmt.Println("Created admin user (username: admin, password: admin)")
} else {
fmt.Println("Admin user already exists")
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
question := models.Question{
TopicID:         javaTopicID,
QuestionText:    q.text,
Difficulty:      q.difficulty,
Tags:            q.tags,
TopicQuestionID: fmt.Sprintf("%d.%d", javaTopicID, startNum+i),
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
{"<class 'list'>", true},
{"<class 'array'>", false},
{"<class 'tuple'>", false},
{"<class 'dict'>", false},
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
question := models.Question{
TopicID:         pythonTopicID,
QuestionText:    q.text,
Difficulty:      q.difficulty,
Tags:            q.tags,
TopicQuestionID: fmt.Sprintf("%d.%d", pythonTopicID, startNum+i),
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

fmt.Println("\nPhase 2 seeding complete!")
return nil
}
