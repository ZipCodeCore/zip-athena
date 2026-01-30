#!/bin/bash
# Script to apply database migrations and seed learning resources

DB_PATH="${1:-./athena.db}"

echo "Applying learning resources migration to database: $DB_PATH"

# Apply migration
sqlite3 "$DB_PATH" < database/migrations/004_learning_resources.sql

if [ $? -eq 0 ]; then
    echo "✓ Migration applied successfully"
else
    echo "✗ Migration failed"
    exit 1
fi

# Seed learning resources
echo "Seeding learning resources..."
sqlite3 "$DB_PATH" < database/seed_learning_resources.sql

if [ $? -eq 0 ]; then
    echo "✓ Learning resources seeded successfully"
else
    echo "✗ Seeding failed"
    exit 1
fi

echo ""
echo "Database setup complete!"
echo "You can now access the Learning Center at /learning.html"
