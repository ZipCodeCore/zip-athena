package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd := flag.String("cmd", "up", "Migration command: up, down, version, force")
	steps := flag.Int("steps", 1, "Number of migrations (for down)")
	dbPath := flag.String("db", "athena.db", "Path to database file")
	migrationsPath := flag.String("migrations", "database/migrations", "Path to migrations directory")
	version := flag.Int("version", 0, "Migration version (for force)")
	flag.Parse()

	// Open database connection directly
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create migration instance
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+*migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrator: %v", err)
	}

	switch *cmd {
	case "up":
		fmt.Println("Running migrations...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		fmt.Printf("Rolling back %d migration(s)...\n", *steps)
		if err := m.Steps(-*steps); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to roll back migrations: %v", err)
		}
		fmt.Printf("Rolled back %d migration(s) successfully\n", *steps)

	case "version":
		v, dirty, err := m.Version()
		if err == migrate.ErrNilVersion {
			fmt.Println("Current migration version: none (no migrations applied)")
		} else if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		} else {
			if dirty {
				fmt.Printf("Current migration version: %d (dirty - needs fixing)\n", v)
			} else {
				fmt.Printf("Current migration version: %d\n", v)
			}
		}

	case "force":
		if *version < 0 {
			log.Fatalf("Version must be >= 0")
		}
		fmt.Printf("Forcing migration version to %d...\n", *version)
		if err := m.Force(*version); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		fmt.Printf("Migration version forced to %d\n", *version)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", *cmd)
		fmt.Fprintf(os.Stderr, "Usage: migrate -cmd <up|down|version|force> [options]\n")
		os.Exit(1)
	}
}
