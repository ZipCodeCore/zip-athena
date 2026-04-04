package database

import (
	"database/sql"

	"github.com/ZipCodeCore/zip-athena/models"
)

// CreateTopic creates a new topic in the database
func (db *SQLiteDB) CreateTopic(topic *models.Topic) error {
	// Default to active for new topics
	if topic.ID == 0 {
		topic.IsActive = true
	}

	// Get max display_order and increment
	if topic.DisplayOrder == 0 {
		var maxOrder int
		err := db.db.QueryRow("SELECT COALESCE(MAX(display_order), 0) FROM topics").Scan(&maxOrder)
		if err == nil {
			topic.DisplayOrder = maxOrder + 1
		}
	}

	result, err := db.db.Exec(
		"INSERT INTO topics (name, description, category, is_active, display_order, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		topic.Name, topic.Description, topic.Category, topic.IsActive, topic.DisplayOrder, topic.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	topic.ID = int(id)
	return nil
}

// GetTopic retrieves a topic by ID
func (db *SQLiteDB) GetTopic(id int) (*models.Topic, error) {
	topic := &models.Topic{}
	err := db.db.QueryRow(
		"SELECT id, name, description, category, is_active, display_order, created_at FROM topics WHERE id = ?",
		id,
	).Scan(&topic.ID, &topic.Name, &topic.Description, &topic.Category, &topic.IsActive, &topic.DisplayOrder, &topic.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return topic, nil
}

// ListAllTopics retrieves all topics (active and inactive) ordered by display order
func (db *SQLiteDB) ListAllTopics() ([]*models.Topic, error) {
	rows, err := db.db.Query("SELECT id, name, description, category, is_active, display_order, created_at FROM topics ORDER BY display_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*models.Topic
	for rows.Next() {
		topic := &models.Topic{}
		if err := rows.Scan(&topic.ID, &topic.Name, &topic.Description, &topic.Category, &topic.IsActive, &topic.DisplayOrder, &topic.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, rows.Err()
}

// ListTopics retrieves all active topics ordered by display order
func (db *SQLiteDB) ListTopics() ([]*models.Topic, error) {
	rows, err := db.db.Query("SELECT id, name, description, category, is_active, display_order, created_at FROM topics WHERE is_active = 1 ORDER BY display_order ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topics []*models.Topic
	for rows.Next() {
		topic := &models.Topic{}
		if err := rows.Scan(&topic.ID, &topic.Name, &topic.Description, &topic.Category, &topic.IsActive, &topic.DisplayOrder, &topic.CreatedAt); err != nil {
			return nil, err
		}
		topics = append(topics, topic)
	}
	return topics, rows.Err()
}

// UpdateTopic updates a topic's information
func (db *SQLiteDB) UpdateTopic(topic *models.Topic) error {
	result, err := db.db.Exec(
		"UPDATE topics SET name = ?, description = ?, category = ?, is_active = ?, display_order = ? WHERE id = ?",
		topic.Name, topic.Description, topic.Category, topic.IsActive, topic.DisplayOrder, topic.ID,
	)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// DeleteTopic soft deletes a topic by setting is_active to false
func (db *SQLiteDB) DeleteTopic(id int) error {
	result, err := db.db.Exec("UPDATE topics SET is_active = 0 WHERE id = ?", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

// ReorderTopics updates the display order of topics
func (db *SQLiteDB) ReorderTopics(topicIDs []int) error {
	// First validate all topics exist
	for _, topicID := range topicIDs {
		var exists int
		err := db.db.QueryRow("SELECT COUNT(*) FROM topics WHERE id = ?", topicID).Scan(&exists)
		if err != nil {
			return err
		}
		if exists == 0 {
			return ErrNotFound
		}
	}

	// Update display_order for each topic based on its position in the array
	for i, topicID := range topicIDs {
		_, err := db.db.Exec("UPDATE topics SET display_order = ? WHERE id = ?", i, topicID)
		if err != nil {
			return err
		}
	}
	return nil
}
