#!/bin/bash
# Script to fix topic visibility issue
# This ensures all topics are marked as active so students can see them

DB_PATH="${1:-./athena.db}"

echo "========================================="
echo "Fixing Topic Visibility"
echo "========================================="
echo "Database: $DB_PATH"
echo ""

if [ ! -f "$DB_PATH" ]; then
    echo "Error: Database file not found: $DB_PATH"
    echo "Usage: $0 [path/to/database.db]"
    exit 1
fi

echo "Current topic status:"
echo "---------------------"
sqlite3 "$DB_PATH" "SELECT id, name, category, is_active FROM topics ORDER BY category, name;"
echo ""

echo "Applying migration 007_ensure_topics_active.sql..."
sqlite3 "$DB_PATH" < database/migrations/007_ensure_topics_active.sql

if [ $? -eq 0 ]; then
    echo "✓ Migration applied successfully"
else
    echo "✗ Migration failed"
    exit 1
fi

echo ""
echo "Updated topic status:"
echo "---------------------"
sqlite3 "$DB_PATH" "SELECT id, name, category, is_active FROM topics ORDER BY category, name;"
echo ""

echo "========================================="
echo "Topic visibility fix complete!"
echo "========================================="
echo ""
echo "Students should now be able to see topics based on their category:"
echo "  - Java students: Java and SQL topics"
echo "  - Data students: Python and SQL topics"
echo ""
