#!/bin/bash

# Script to fix Kris user in the database
# This updates Kris to be a data category student (not admin)

DB_PATH="${1:-./athena.db}"

if [ ! -f "$DB_PATH" ]; then
    echo "Error: Database file not found: $DB_PATH"
    echo "Usage: $0 [path/to/athena.db]"
    exit 1
fi

echo "Fixing Kris user in database: $DB_PATH"

# Apply the migration
sqlite3 "$DB_PATH" < database/migrations/008_fix_kris_user_category.sql

if [ $? -eq 0 ]; then
    echo "✓ Successfully updated Kris to be a data category student"
    echo ""
    echo "Verification:"
    sqlite3 "$DB_PATH" "SELECT id, username, category, is_admin FROM users WHERE username='kris';"
else
    echo "✗ Failed to update database"
    exit 1
fi
