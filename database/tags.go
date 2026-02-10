package database

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/ZipCodeCore/zip-athena/models"
)

func (db *SQLiteDB) GetAllTags() ([]string, error) {
	rows, err := db.db.Query("SELECT DISTINCT tags FROM questions WHERE tags IS NOT NULL AND tags != '[]'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagSet := make(map[string]bool)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}

func (db *SQLiteDB) GetTagsWithCounts() (map[string]int, error) {
	rows, err := db.db.Query("SELECT tags FROM questions WHERE tags IS NOT NULL AND tags != '[]'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			tagCounts[tag]++
		}
	}

	return tagCounts, rows.Err()
}

func (db *SQLiteDB) GetTagStats(userID int) ([]*models.TagStats, error) {
	query := `
		SELECT q.tags, qa.correct_attempts, qa.incorrect_attempts
		FROM question_attempts qa
		JOIN questions q ON qa.question_id = q.id
		WHERE qa.user_id = ? AND q.tags IS NOT NULL AND q.tags != '[]'
	`

	rows, err := db.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagStats := make(map[string]*models.TagStats)
	for rows.Next() {
		var tagsJSON string
		var correctAttempts, incorrectAttempts int

		if err := rows.Scan(&tagsJSON, &correctAttempts, &incorrectAttempts); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			if _, exists := tagStats[tag]; !exists {
				tagStats[tag] = &models.TagStats{Tag: tag}
			}

			stats := tagStats[tag]
			stats.CorrectAttempts += correctAttempts
			stats.TotalAttempts += correctAttempts + incorrectAttempts
		}
	}

	// Get question counts
	tagCounts, _ := db.GetTagsWithCounts()

	// Calculate accuracy
	result := make([]*models.TagStats, 0, len(tagStats))
	for tag, stats := range tagStats {
		stats.QuestionCount = tagCounts[tag]
		if stats.TotalAttempts > 0 {
			stats.Accuracy = float64(stats.CorrectAttempts) / float64(stats.TotalAttempts) * 100
		}
		result = append(result, stats)
	}

	return result, rows.Err()
}

func (db *SQLiteDB) GetQuestionsByTags(tags []string, useAnd bool) ([]*models.Question, error) {
	rows, err := db.db.Query("SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions WHERE tags IS NOT NULL AND tags != '[]'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		q := &models.Question{}
		var tagsJSON string

		err := rows.Scan(&q.ID, &q.TopicID, &q.QuestionText, &q.Difficulty, &tagsJSON, &q.TopicQuestionID, &q.CreatedAt)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(tagsJSON), &q.Tags)

		// Check if question matches tags
		if useAnd {
			hasAllTags := true
			for _, tag := range tags {
				found := false
				for _, qTag := range q.Tags {
					if qTag == tag {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if hasAllTags {
				questions = append(questions, q)
			}
		} else {
			hasAnyTag := false
			for _, tag := range tags {
				for _, qTag := range q.Tags {
					if qTag == tag {
						hasAnyTag = true
						break
					}
				}
				if hasAnyTag {
					break
				}
			}
			if hasAnyTag {
				questions = append(questions, q)
			}
		}
	}

	return questions, rows.Err()
}

func (db *SQLiteDB) GetTagsWithMetadata() ([]*models.Tag, error) {
	// Get all tags from questions
	rows, err := db.db.Query("SELECT tags FROM questions WHERE tags IS NOT NULL AND tags != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tagCounts := make(map[string]int)
	for rows.Next() {
		var tagsJSON string
		if err := rows.Scan(&tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			tagCounts[tag]++
		}
	}

	var result []*models.Tag
	for name, count := range tagCounts {
		result = append(result, &models.Tag{
			Name:          name,
			QuestionCount: count,
		})
	}

	return result, nil
}

func (db *SQLiteDB) MergeTags(sourceTags []string, targetTag string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get all questions with any of the source tags
	rows, err := tx.Query("SELECT id, tags FROM questions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var tagsJSON string
		if err := rows.Scan(&id, &tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		// Check if any source tag is present
		hasSourceTag := false
		newTags := []string{}
		for _, tag := range tags {
			isSourceTag := false
			for _, sourceTag := range sourceTags {
				if strings.EqualFold(tag, sourceTag) {
					isSourceTag = true
					hasSourceTag = true
					break
				}
			}
			if !isSourceTag {
				newTags = append(newTags, tag)
			}
		}

		if hasSourceTag {
			// Add target tag if not already present
			hasTarget := false
			for _, tag := range newTags {
				if strings.EqualFold(tag, targetTag) {
					hasTarget = true
					break
				}
			}
			if !hasTarget {
				newTags = append(newTags, targetTag)
			}

			// Update question with new tags
			newTagsJSON, err := json.Marshal(newTags)
			if err != nil {
				continue
			}

			_, err = tx.Exec("UPDATE questions SET tags = ? WHERE id = ?", string(newTagsJSON), id)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (db *SQLiteDB) DeleteTag(tagName string) error {
	// Check if any questions use this tag
	rows, err := db.db.Query("SELECT id, tags FROM questions")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var tagsJSON string
		if err := rows.Scan(&id, &tagsJSON); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		for _, tag := range tags {
			if strings.EqualFold(tag, tagName) {
				return errors.New("cannot delete tag: tag is in use by questions")
			}
		}
	}

	return nil
}

func (db *SQLiteDB) GetQuestionsByTag(tagName string) ([]*models.Question, error) {
	rows, err := db.db.Query("SELECT id, topic_id, question_text, difficulty, tags, topic_question_id, created_at FROM questions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question := &models.Question{}
		var tagsJSON string
		if err := rows.Scan(&question.ID, &question.TopicID, &question.QuestionText, &question.Difficulty, &tagsJSON, &question.TopicQuestionID, &question.CreatedAt); err != nil {
			continue
		}

		var tags []string
		if err := json.Unmarshal([]byte(tagsJSON), &tags); err != nil {
			continue
		}

		question.Tags = tags

		// Check if tag matches
		for _, tag := range tags {
			if strings.EqualFold(tag, tagName) {
				question.IsActive = true
				questions = append(questions, question)
				break
			}
		}
	}

	return questions, nil
}

