package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending migrations from the specified directory
func (db *SQLiteDB) RunMigrations(migrationsPath string) error {
	driver, err := sqlite3.WithInstance(db.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("running migrations: %w", err)
	}

	return nil
}

// MigrateUp applies all pending migrations
func (db *SQLiteDB) MigrateUp(migrationsPath string) error {
	return db.RunMigrations(migrationsPath)
}

// MigrateDown rolls back the last n migrations
func (db *SQLiteDB) MigrateDown(migrationsPath string, steps int) error {
	driver, err := sqlite3.WithInstance(db.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	if steps <= 0 {
		return fmt.Errorf("steps must be a positive number")
	}

	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rolling back migrations: %w", err)
	}

	return nil
}

// MigrationVersion returns the current migration version
func (db *SQLiteDB) MigrationVersion(migrationsPath string) (uint, bool, error) {
	driver, err := sqlite3.WithInstance(db.db, &sqlite3.Config{})
	if err != nil {
		return 0, false, fmt.Errorf("creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return 0, false, fmt.Errorf("creating migrator: %w", err)
	}

	version, dirty, err := m.Version()
	if err == migrate.ErrNilVersion {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, fmt.Errorf("getting migration version: %w", err)
	}

	return version, dirty, nil
}

// ForceMigrationVersion sets the migration version without running migrations
// This is useful for fixing dirty migration states
func (db *SQLiteDB) ForceMigrationVersion(migrationsPath string, version int) error {
	driver, err := sqlite3.WithInstance(db.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		return fmt.Errorf("creating migrator: %w", err)
	}

	if err := m.Force(version); err != nil {
		return fmt.Errorf("forcing migration version: %w", err)
	}

	return nil
}
