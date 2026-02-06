#!/bin/bash

# Run all SQL migrations against a SQLite database
# Usage: ./run_migrations.sh <database_file>

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 <database_file>"
    exit 1
fi

DB_FILE="$1"
MIGRATIONS_DIR="$(dirname "$0")/database/migrations"

if [ ! -f "$DB_FILE" ]; then
    echo "Error: Database file '$DB_FILE' not found"
    exit 1
fi

if [ ! -d "$MIGRATIONS_DIR" ]; then
    echo "Error: Migrations directory '$MIGRATIONS_DIR' not found"
    exit 1
fi

echo "Running migrations on: $DB_FILE"

# Run each migration file in sorted order
for migration in $(ls "$MIGRATIONS_DIR"/*.sql 2>/dev/null | sort); do
    echo "Applying: $(basename "$migration")"
    sqlite3 "$DB_FILE" < "$migration"
done

echo "All migrations applied successfully"
