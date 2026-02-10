# Database Migrations

This document describes the database migration system for Athena, which uses [golang-migrate](https://github.com/golang-migrate/migrate) for version control and schema management.

## Overview

The migration system provides:
- **Version tracking**: Tracks which migrations have been applied
- **Rollback capability**: Can roll back migrations if needed
- **Portable migrations**: SQL files work across environments
- **Clean separation**: No inline schema mixed with migration logic

## Migration Files

Migration files are stored in `database/migrations/` and follow the naming convention:
```
{version}_{description}.up.sql    - Forward migration
{version}_{description}.down.sql  - Rollback migration
```

Current migrations:
- **001_initial_schema**: Creates all base tables (users, topics, questions, etc.)
- **002_add_indexes**: Adds performance indexes for learning resources
- **003_fix_learning_resources**: Activates phases 2-3 and fixes file paths

## Using the Migration CLI

The `migrate` command-line tool provides migration management:

### Build the migrate tool
```bash
go build -o migrate ./cmd/migrate
```

### Check migration version
```bash
./migrate -db athena.db -cmd version
```

### Run all pending migrations
```bash
./migrate -db athena.db -cmd up
```

### Roll back migrations
```bash
# Roll back 1 migration
./migrate -db athena.db -cmd down -steps 1

# Roll back 2 migrations
./migrate -db athena.db -cmd down -steps 2
```

### Force migration version (for fixing dirty state)
```bash
./migrate -db athena.db -cmd force -version 3
```

### Options
- `-db`: Path to database file (default: `athena.db`)
- `-migrations`: Path to migrations directory (default: `database/migrations`)
- `-cmd`: Command to run (`up`, `down`, `version`, `force`)
- `-steps`: Number of migrations to roll back (for `down`)
- `-version`: Target version (for `force`)

## Creating New Migrations

To create a new migration:

1. **Create up migration file**:
   ```sql
   -- database/migrations/004_add_new_feature.up.sql
   CREATE TABLE IF NOT EXISTS new_table (
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       name TEXT NOT NULL
   );
   ```

2. **Create down migration file**:
   ```sql
   -- database/migrations/004_add_new_feature.down.sql
   DROP TABLE IF EXISTS new_table;
   ```

3. **Test the migration**:
   ```bash
   # Test on a copy of your database
   cp athena.db athena.db.backup
   ./migrate -db athena.db -cmd up
   
   # Verify the change
   sqlite3 athena.db ".schema new_table"
   
   # Test rollback
   ./migrate -db athena.db -cmd down -steps 1
   ```

## Automatic Migrations

The application automatically runs pending migrations on startup when creating a database connection:

```go
db, err := database.NewSQLiteDB("athena.db")
// Migrations are applied automatically
```

This ensures that:
- Fresh databases get the complete schema
- Existing databases are updated to the latest version
- No manual migration step needed for development

## Migration State

Migrations are tracked in the `schema_migrations` table:
```sql
CREATE TABLE schema_migrations (
    version bigint NOT NULL PRIMARY KEY,
    dirty boolean NOT NULL
);
```

- `version`: Current migration version number
- `dirty`: Whether the migration failed mid-execution

If a migration fails (dirty state):
1. Fix the issue manually or with SQL
2. Use `force` command to set the correct version
3. Re-run migrations with `up` command

## Best Practices

1. **Always create both up and down migrations** for rollback capability
2. **Use idempotent operations** where possible (e.g., `CREATE TABLE IF NOT EXISTS`)
3. **Test migrations on a copy** of production data before applying
4. **Keep migrations small and focused** - one feature per migration
5. **Never modify existing migration files** - create new ones instead
6. **Use transactions** for complex migrations to ensure atomicity
7. **Document breaking changes** in migration comments

## Migrating from Old System

The previous system used:
- Inline schema in `initSchema()`
- Inline migrations in `runMigrations()`
- Separate SQL files (004-010) without version tracking

These have been consolidated into migrations 001-003:
- **001**: All base tables from inline schema
- **002**: Performance indexes
- **003**: Data fixes (learning resources paths, activations)

For existing databases, migrations 001-003 are idempotent and safe to run.

## Troubleshooting

### "No change" error
This is normal - it means all migrations are already applied.

### Dirty migration state
If a migration fails partway through:
```bash
# Check current state
./migrate -db athena.db -cmd version
# Output: Current migration version: 2 (dirty - needs fixing)

# Fix the issue manually, then force the version
./migrate -db athena.db -cmd force -version 2

# Continue with migrations
./migrate -db athena.db -cmd up
```

### Migration not found
Ensure you're running the command from the project root directory where `database/migrations/` exists, or specify the migrations path:
```bash
./migrate -db athena.db -migrations /path/to/database/migrations -cmd up
```

## References

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [SQLite migration guide](https://github.com/golang-migrate/migrate/tree/master/database/sqlite3)
