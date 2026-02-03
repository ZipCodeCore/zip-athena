#!/bin/bash
# Fix Phase 2 and Phase 3 Learning Resources
# This script applies the migration to fix file paths and activate Phase 2 and 3 resources

set -e

DB_PATH="${1:-athena.db}"

echo "Applying migration to fix Phase 2 and 3 learning resources..."
echo "Database: $DB_PATH"

# Check if database exists
if [ ! -f "$DB_PATH" ]; then
    echo "Error: Database file not found: $DB_PATH"
    exit 1
fi

# Show current state
echo ""
echo "Current state of Phase 2 and 3 resources:"
sqlite3 "$DB_PATH" "SELECT id, language, phase, section_number, section_title, file_path, is_active FROM learning_resources WHERE phase IN (2, 3) ORDER BY phase, language, section_number;"

# Apply migration
echo ""
echo "Applying migration..."
sqlite3 "$DB_PATH" < database/migrations/005_fix_phase2_phase3_learning_resources.sql

# Show updated state
echo ""
echo "Updated state of Phase 2 and 3 resources:"
sqlite3 "$DB_PATH" "SELECT id, language, phase, section_number, section_title, file_path, is_active FROM learning_resources WHERE phase IN (2, 3) ORDER BY phase, language, section_number;"

echo ""
echo "Migration completed successfully!"
echo ""
echo "Verify files exist:"
ls -1 learning/java/8-oop-classes-and-objects.md \
      learning/java/9-oop-inheritance.md \
      learning/java/10-oop-polymorphism.md \
      learning/java/11-interfaces.md \
      learning/java/12-packages.md \
      learning/java/13-exception-handling.md \
      learning/java/14-file-io.md \
      learning/java/15-collections-framework.md \
      learning/python/14-oop-classes-and-objects.md \
      learning/python/15-oop-inheritance.md \
      learning/python/16-oop-polymorphism.md \
      learning/python/11-modules-and-packages.md \
      learning/python/12-file-io.md \
      learning/python/13-exception-handling.md 2>&1 | head -1
