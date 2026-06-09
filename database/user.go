package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ZipCodeCore/zip-athena/models"
)

// CreateUser creates a new user in the database
func (db *SQLiteDB) CreateUser(user *models.User) error {
	// Normalize category to lowercase for storage
	user.Category = strings.ToLower(user.Category)
	if !user.IsActiveSet {
		user.IsActive = true
	}

	result, err := db.db.Exec(
		"INSERT INTO users (username, email, password, is_admin, is_active, category, total_points, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		user.Username, user.Email, user.Password, user.IsAdmin, user.IsActive, user.Category, user.TotalPoints, user.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = int(id)
	return nil
}

// GetUser retrieves a user by ID
func (db *SQLiteDB) GetUser(id int) (*models.User, error) {
	user := &models.User{}
	err := db.db.QueryRow(
		"SELECT id, username, email, is_admin, is_active, category, total_points, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.IsActive, &user.Category, &user.TotalPoints, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username (helper method)
func (db *SQLiteDB) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := db.db.QueryRow(
		"SELECT id, username, email, password, is_admin, is_active, category, total_points, created_at FROM users WHERE username = ?",
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.IsAdmin, &user.IsActive, &user.Category, &user.TotalPoints, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUserPoints updates a user's point total
func (db *SQLiteDB) UpdateUserPoints(id int, points int) error {
	result, err := db.db.Exec("UPDATE users SET total_points = ? WHERE id = ?", points, id)
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

// ListUsers retrieves all users
func (db *SQLiteDB) ListUsers() ([]*models.User, error) {
	rows, err := db.db.Query("SELECT id, username, email, is_admin, is_active, category, total_points, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsAdmin, &user.IsActive, &user.Category, &user.TotalPoints, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

// UpdateUser updates a user's information
func (db *SQLiteDB) UpdateUser(user *models.User) error {
	// Normalize category to lowercase for comparison and storage
	user.Category = strings.ToLower(user.Category)

	// Validate category (after normalization)
	if user.Category != "" && user.Category != "java" && user.Category != "data" {
		return fmt.Errorf("invalid category: must be 'java' or 'data' (case-insensitive)")
	}

	// If password is provided, update it along with other fields
	if user.Password != "" {
		result, err := db.db.Exec(
			"UPDATE users SET username = ?, email = ?, password = ?, is_admin = ?, is_active = ?, category = ?, total_points = ? WHERE id = ?",
			user.Username, user.Email, user.Password, user.IsAdmin, user.IsActive, user.Category, user.TotalPoints, user.ID,
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
	} else {
		// Don't update password if not provided
		result, err := db.db.Exec(
			"UPDATE users SET username = ?, email = ?, is_admin = ?, is_active = ?, category = ?, total_points = ? WHERE id = ?",
			user.Username, user.Email, user.IsAdmin, user.IsActive, user.Category, user.TotalPoints, user.ID,
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
	}
	return nil
}

// DeleteUser deletes a user by ID
func (db *SQLiteDB) DeleteUser(id int) error {
	result, err := db.db.Exec("DELETE FROM users WHERE id = ?", id)
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
