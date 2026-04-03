package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ZipCodeCore/zip-athena/database"
	"github.com/ZipCodeCore/zip-athena/models"
)

// QuestionFile represents a JSON file containing questions for a topic
type QuestionFile struct {
	Topic     string         `json:"topic"` // Topic name (must exist in database)
	Questions []QuestionData `json:"questions"`
}

// QuestionData represents a single question with its answers
type QuestionData struct {
	Text       string       `json:"text"`
	Difficulty int          `json:"difficulty"` // 1-7
	Tags       []string     `json:"tags"`
	Answers    []AnswerData `json:"answers"`
}

// AnswerData represents an answer option
type AnswerData struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

func main() {
	// Parse command line flags
	dbPath := flag.String("db", "./athena.db", "Path to the SQLite database")
	questionsDir := flag.String("dir", "./questions", "Directory containing question JSON files")
	specificFile := flag.String("file", "", "Load a specific JSON file instead of a directory")
	dryRun := flag.Bool("dry-run", false, "Validate files without importing")
	flag.Parse()

	// Initialize database
	db, err := database.NewSQLiteDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Get existing topics for validation
	topics, err := db.ListTopics()
	if err != nil {
		log.Fatalf("Failed to list topics: %v", err)
	}
	topicMap := make(map[string]*models.Topic)
	for _, t := range topics {
		topicMap[strings.ToLower(t.Name)] = t
	}

	fmt.Println("=== Question Bulk Loader ===")
	fmt.Printf("Database: %s\n", *dbPath)

	if *dryRun {
		fmt.Println("Mode: DRY RUN (no changes will be made)")
	}

	var files []string

	if *specificFile != "" {
		// Load specific file
		files = []string{*specificFile}
	} else {
		// Find all JSON files in directory
		files, err = filepath.Glob(filepath.Join(*questionsDir, "*.json"))
		if err != nil {
			log.Fatalf("Failed to find JSON files: %v", err)
		}
	}

	if len(files) == 0 {
		fmt.Println("No JSON files found to process")
		os.Exit(0)
	}

	fmt.Printf("Found %d file(s) to process\n\n", len(files))

	totalQuestions := 0
	totalImported := 0
	totalSkipped := 0
	totalErrors := 0

	for _, file := range files {
		fmt.Printf("Processing: %s\n", filepath.Base(file))

		qf, err := loadQuestionFile(file)
		if err != nil {
			fmt.Printf("  ERROR: Failed to load file: %v\n", err)
			totalErrors++
			continue
		}

		// Find the topic
		topic, ok := topicMap[strings.ToLower(qf.Topic)]
		if !ok {
			if *dryRun {
				// In dry-run mode, simulate topic creation so validation/import preview can continue.
				topic = &models.Topic{ID: -1, Name: qf.Topic}
				fmt.Printf("  Topic: %s (would be auto-created)\n", qf.Topic)
			} else {
				topic = &models.Topic{
					Name:        qf.Topic,
					Description: fmt.Sprintf("Auto-created by loadquestions from %s", filepath.Base(file)),
					Category:    inferTopicCategory(qf.Topic),
					IsActive:    true,
					CreatedAt:   time.Now(),
				}

				if err := db.CreateTopic(topic); err != nil {
					fmt.Printf("  ERROR: Failed to auto-create topic '%s': %v\n", qf.Topic, err)
					fmt.Printf("  Available topics: %s\n", getTopicNames(topics))
					totalErrors++
					continue
				}

				topics = append(topics, topic)
				topicMap[strings.ToLower(topic.Name)] = topic
				fmt.Printf("  Topic: %s (auto-created, ID: %d)\n", topic.Name, topic.ID)
			}
		} else {
			fmt.Printf("  Topic: %s (ID: %d)\n", topic.Name, topic.ID)
		}
		fmt.Printf("  Questions in file: %d\n", len(qf.Questions))

		// Get existing questions to avoid duplicates
		existingQuestions, _ := db.GetQuestionsByTopic(topic.ID)
		existingTexts := make(map[string]bool)
		for _, eq := range existingQuestions {
			existingTexts[strings.ToLower(eq.QuestionText)] = true
		}

		imported := 0
		skipped := 0
		startNum := len(existingQuestions) + 1

		for i, q := range qf.Questions {
			totalQuestions++

			// Validate question
			if err := validateQuestion(q); err != nil {
				fmt.Printf("    Question %d: INVALID - %v\n", i+1, err)
				totalErrors++
				continue
			}

			// Check for duplicates
			if existingTexts[strings.ToLower(q.Text)] {
				skipped++
				totalSkipped++
				continue
			}

			if *dryRun {
				imported++
				totalImported++
				continue
			}

			// Create the question
			topicQuestionID := fmt.Sprintf("%d.%d", topic.ID, startNum+imported)
			question := models.Question{
				TopicID:         topic.ID,
				QuestionText:    q.Text,
				Difficulty:      q.Difficulty,
				Tags:            q.Tags,
				TopicQuestionID: topicQuestionID,
				CreatedAt:       time.Now(),
			}

			if err := db.CreateQuestion(&question); err != nil {
				fmt.Printf("    Question %d: ERROR - %v\n", i+1, err)
				totalErrors++
				continue
			}

			// Create answers
			for _, a := range q.Answers {
				answer := models.Answer{
					QuestionID: question.ID,
					AnswerText: a.Text,
					IsCorrect:  a.IsCorrect,
				}
				if err := db.CreateAnswer(&answer); err != nil {
					fmt.Printf("    Answer error for question %d: %v\n", question.ID, err)
				}
			}

			imported++
			totalImported++
		}

		fmt.Printf("  Imported: %d, Skipped (duplicates): %d\n\n", imported, skipped)
	}

	fmt.Println("=== Summary ===")
	fmt.Printf("Total questions processed: %d\n", totalQuestions)
	fmt.Printf("Successfully imported: %d\n", totalImported)
	fmt.Printf("Skipped (duplicates): %d\n", totalSkipped)
	fmt.Printf("Errors: %d\n", totalErrors)

	if *dryRun {
		fmt.Println("\nDry run complete. Run without --dry-run to import questions.")
	}
}

func loadQuestionFile(path string) (*QuestionFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var qf QuestionFile
	if err := json.Unmarshal(data, &qf); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &qf, nil
}

func validateQuestion(q QuestionData) error {
	if strings.TrimSpace(q.Text) == "" {
		return fmt.Errorf("question text is empty")
	}

	if q.Difficulty < 1 || q.Difficulty > 7 {
		return fmt.Errorf("difficulty must be between 1 and 7 (got %d)", q.Difficulty)
	}

	if len(q.Answers) < 2 {
		return fmt.Errorf("at least 2 answers required (got %d)", len(q.Answers))
	}

	if len(q.Answers) > 6 {
		return fmt.Errorf("maximum 6 answers allowed (got %d)", len(q.Answers))
	}

	correctCount := 0
	for _, a := range q.Answers {
		if strings.TrimSpace(a.Text) == "" {
			return fmt.Errorf("answer text is empty")
		}
		if a.IsCorrect {
			correctCount++
		}
	}

	if correctCount == 0 {
		return fmt.Errorf("at least one answer must be marked as correct")
	}

	if correctCount > 1 {
		return fmt.Errorf("only one answer should be correct (got %d)", correctCount)
	}

	return nil
}

func getTopicNames(topics []*models.Topic) string {
	names := make([]string, len(topics))
	for i, t := range topics {
		names[i] = t.Name
	}
	return strings.Join(names, ", ")
}

func inferTopicCategory(topicName string) string {
	name := strings.ToLower(topicName)
	switch {
	case strings.Contains(name, "python"):
		return "python"
	case strings.Contains(name, "sql"):
		return "sql"
	default:
		return "java"
	}
}
