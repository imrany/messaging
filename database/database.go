package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

// Init initializes the database connection and runs migrations if specified.
func Init(DB_URL string, runMigrations bool) error {
	var err error
	DB, err = sql.Open("postgres", DB_URL)
	if err != nil {
		slog.Error("Failed to open database connection", "error", err)
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(100) // Based on server capacity
	DB.SetMaxIdleConns(25)  // Keep connections ready
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(10 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = DB.PingContext(ctx)
	if err != nil {
		slog.Error("Failed to ping database", "error", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations if requested
	if runMigrations {
		if err := RunMigrations(); err != nil {
			slog.Error("Failed to run migrations", "error", err)
			return err
		}
	}

	log.Println("Connected to database")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
