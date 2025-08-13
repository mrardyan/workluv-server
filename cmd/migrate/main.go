package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"workluv/pkg/config"
	"workluv/pkg/logger"

	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migration", "directory with migration files")
)

func main() {
	flags.Usage = usage
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Failed to parse flags: %v", err)
	}

	args := flags.Args()
	if len(args) == 0 {
		flags.Usage()
		return
	}

	command := args[0]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	logger.SetupLogger(cfg)

	// Use the database connection string from config
	dbURL := cfg.GetPostgreSQLDSN()

	// Open database connection
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set dialect
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("Failed to set dialect: %v", err)
	}

	// Run the migration command
	switch command {
	case "up":
		if err := goose.Up(db, *dir); err != nil {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("✅ Successfully applied all migrations")

	case "up-to":
		if len(args) < 2 {
			log.Fatal("up-to command requires a version argument")
		}
		version, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := goose.UpTo(db, *dir, version); err != nil {
			log.Fatalf("Migration up-to failed: %v", err)
		}
		log.Printf("✅ Successfully migrated up to version %d", version)

	case "down":
		if err := goose.Down(db, *dir); err != nil {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("✅ Successfully rolled back one migration")

	case "down-to":
		if len(args) < 2 {
			log.Fatal("down-to command requires a version argument")
		}
		version, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := goose.DownTo(db, *dir, version); err != nil {
			log.Fatalf("Migration down-to failed: %v", err)
		}
		log.Printf("✅ Successfully migrated down to version %d", version)

	case "redo":
		if err := goose.Redo(db, *dir); err != nil {
			log.Fatalf("Migration redo failed: %v", err)
		}
		log.Println("✅ Successfully redid last migration")

	case "reset":
		if err := goose.Reset(db, *dir); err != nil {
			log.Fatalf("Migration reset failed: %v", err)
		}
		log.Println("✅ Successfully reset all migrations")

	case "status":
		if err := goose.Status(db, *dir); err != nil {
			log.Fatalf("Migration status failed: %v", err)
		}

	case "version":
		version, err := goose.GetDBVersion(db)
		if err != nil {
			log.Fatalf("Failed to get database version: %v", err)
		}
		log.Printf("Database version: %d", version)

	case "create":
		if len(args) < 2 {
			log.Fatal("create command requires a migration name argument")
		}
		name := args[1]
		migrationType := "sql" // default to SQL migrations
		if len(args) >= 3 {
			migrationType = args[2]
		}
		if err := goose.Create(db, *dir, name, migrationType); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
		log.Printf("✅ Successfully created migration: %s", name)

	case "fix":
		if err := goose.Fix(*dir); err != nil {
			log.Fatalf("Migration fix failed: %v", err)
		}
		log.Println("✅ Successfully fixed migration files")

	default:
		log.Printf("Unknown command: %s", command)
		flags.Usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Print(`
migrate is a database migration tool built with goose

Usage:
    migrate [flags] <command> [args...]

Commands:
    up                 Apply all available migrations
    up-to VERSION      Migrate up to a specific version
    down               Roll back the latest migration  
    down-to VERSION    Roll back to a specific version
    redo               Re-run the latest migration
    reset              Roll back all migrations
    status             Show migration status
    version            Show current migration version
    create NAME [TYPE] Create a new migration file (TYPE: sql or go, default: sql)
    fix                Fix migration file sequence

Flags:
`)
	flags.PrintDefaults()
	fmt.Print(`
Examples:
    migrate up                    # Apply all migrations
    migrate down                  # Roll back one migration
    migrate status               # Show migration status
    migrate create add_users     # Create a new migration
    migrate up-to 20230101000000 # Migrate to specific version
    migrate -dir migrations up   # Use custom migration directory

Environment:
    The migration tool uses the same configuration as the main server.
    Make sure to set the appropriate DATABASE_* environment variables.
`)
}
