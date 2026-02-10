package database

import (
	"testing"
)

func TestQueryBuilder_Basic(t *testing.T) {
	query, args := NewQuery("users").
		Where("id = ?", 1).
		Build()

	expected := "SELECT * FROM users WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 || args[0] != 1 {
		t.Errorf("Expected args [1], got %v", args)
	}
}

func TestQueryBuilder_MultipleConditions(t *testing.T) {
	query, args := NewQuery("questions").
		Where("topic_id = ?", 5).
		Where("difficulty = ?", 3).
		Where("is_active = ?", true).
		Build()

	expected := "SELECT * FROM questions WHERE topic_id = ? AND difficulty = ? AND is_active = ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}
}

func TestQueryBuilder_OrderBy(t *testing.T) {
	query, args := NewQuery("users").
		Where("is_admin = ?", false).
		OrderBy("created_at DESC").
		Build()

	expected := "SELECT * FROM users WHERE is_admin = ? ORDER BY created_at DESC"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestQueryBuilder_Limit(t *testing.T) {
	query, args := NewQuery("topics").
		Where("is_active = ?", true).
		OrderBy("display_order ASC").
		Limit(10).
		Build()

	expected := "SELECT * FROM topics WHERE is_active = ? ORDER BY display_order ASC LIMIT 10"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestQueryBuilder_LimitOffset(t *testing.T) {
	query, args := NewQuery("questions").
		Where("topic_id = ?", 3).
		OrderBy("created_at DESC").
		Limit(20).
		Offset(40).
		Build()

	expected := "SELECT * FROM questions WHERE topic_id = ? ORDER BY created_at DESC LIMIT 20 OFFSET 40"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestQueryBuilder_SelectColumns(t *testing.T) {
	query, args := NewQuery("users").
		Select("id", "username", "email").
		Where("is_admin = ?", true).
		Build()

	expected := "SELECT id, username, email FROM users WHERE is_admin = ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestQueryBuilder_Count(t *testing.T) {
	query, args := NewQuery("quiz_attempts").
		Where("user_id = ?", 5).
		Where("completed_time IS NOT NULL").
		BuildCount()

	expected := "SELECT COUNT(*) FROM quiz_attempts WHERE user_id = ? AND completed_time IS NOT NULL"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestQueryBuilder_NoConditions(t *testing.T) {
	query, args := NewQuery("topics").
		OrderBy("name ASC").
		Build()

	expected := "SELECT * FROM topics ORDER BY name ASC"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(args))
	}
}

func TestUpdateBuilder_Basic(t *testing.T) {
	query, args := NewUpdate("users").
		Set("total_points", 100).
		Where("id = ?", 5).
		Build()

	expected := "UPDATE users SET total_points = ? WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 2 || args[0] != 100 || args[1] != 5 {
		t.Errorf("Expected args [100, 5], got %v", args)
	}
}

func TestUpdateBuilder_MultipleSets(t *testing.T) {
	query, args := NewUpdate("questions").
		Set("question_text", "Updated question").
		Set("difficulty", 5).
		Set("is_active", true).
		Where("id = ?", 10).
		Build()

	expected := "UPDATE questions SET question_text = ?, difficulty = ?, is_active = ? WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 4 {
		t.Errorf("Expected 4 args, got %d", len(args))
	}
}

func TestUpdateBuilder_MultipleConditions(t *testing.T) {
	query, args := NewUpdate("topics").
		Set("is_active", false).
		Where("category = ?", "java").
		Where("display_order > ?", 100).
		Build()

	expected := "UPDATE topics SET is_active = ? WHERE category = ? AND display_order > ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 3 {
		t.Errorf("Expected 3 args, got %d", len(args))
	}
}

func TestDeleteBuilder_Basic(t *testing.T) {
	query, args := NewDelete("answers").
		Where("id = ?", 15).
		Build()

	expected := "DELETE FROM answers WHERE id = ?"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 || args[0] != 15 {
		t.Errorf("Expected args [15], got %v", args)
	}
}

func TestDeleteBuilder_MultipleConditions(t *testing.T) {
	query, args := NewDelete("quiz_attempts").
		Where("user_id = ?", 5).
		Where("completed_time IS NULL").
		Build()

	expected := "DELETE FROM quiz_attempts WHERE user_id = ? AND completed_time IS NULL"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 1 {
		t.Errorf("Expected 1 arg, got %d", len(args))
	}
}

func TestDeleteBuilder_NoConditions(t *testing.T) {
	query, args := NewDelete("temp_data").
		Build()

	expected := "DELETE FROM temp_data"
	if query != expected {
		t.Errorf("Expected query %q, got %q", expected, query)
	}

	if len(args) != 0 {
		t.Errorf("Expected 0 args, got %d", len(args))
	}
}
