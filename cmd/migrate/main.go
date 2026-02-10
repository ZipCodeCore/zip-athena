package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ZipCodeCore/zip-athena/database"
)

func main() {
	cmd := flag.String("cmd", "up", "Migration command: up, down, version, force")
	steps := flag.Int("steps", 1, "Number of migrations (for down)")
	dbPath := flag.String("db", "athena.db", "Path to database file")
	migrationsPath := flag.String("migrations", "database/migrations", "Path to migrations directory")
	version := flag.Int("version", 0, "Migration version (for force)")
	flag.Parse()

	db, err := database.NewSQLiteDB(*dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	switch *cmd {
	case "up":
		fmt.Println("Running migrations...")
		if err := db.MigrateUp(*migrationsPath); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		fmt.Printf("Rolling back %d migration(s)...\n", *steps)
		if err := db.MigrateDown(*migrationsPath, *steps); err != nil {
			log.Fatalf("Failed to roll back migrations: %v", err)
		}
		fmt.Printf("Rolled back %d migration(s) successfully\n", *steps)

	case "version":
		v, dirty, err := db.MigrationVersion(*migrationsPath)
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		if dirty {
			fmt.Printf("Current migration version: %d (dirty - needs fixing)\n", v)
		} else {
			fmt.Printf("Current migration version: %d\n", v)
		}

	case "force":
		if *version < 0 {
			log.Fatalf("Version must be >= 0")
		}
		fmt.Printf("Forcing migration version to %d...\n", *version)
		if err := db.ForceMigrationVersion(*migrationsPath, *version); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		fmt.Printf("Migration version forced to %d\n", *version)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", *cmd)
		fmt.Fprintf(os.Stderr, "Usage: migrate -cmd <up|down|version|force> [options]\n")
		os.Exit(1)
	}
}
