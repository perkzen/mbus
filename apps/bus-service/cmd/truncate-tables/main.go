package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/db"
)

func main() {
	env, err := config.LoadEnvironment()
	if err != nil {
		log.Fatalf("❌ Failed to load environment: %v", err)
	}

	pgDb, err := db.NewPostgresDB(env.PostgresURL).Open()
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer pgDb.Close()

	if err := TruncateAllTables(pgDb); err != nil {
		log.Fatalf("❌ Failed to truncate tables: %v", err)
	}

	log.Println("✅ All tables truncated successfully.")
}

func TruncateAllTables(db *sql.DB) error {
	// Fetch all table names except goose_db_version
	rows, err := db.Query(`
		SELECT tablename
		FROM pg_tables
		WHERE schemaname = 'public'
		AND tablename != 'goose_db_version'
	`)
	if err != nil {
		return fmt.Errorf("failed to query table names: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}

	if len(tables) == 0 {
		return nil
	}

	query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", strings.Join(tables, ", "))
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to truncate tables: %w", err)
	}

	return nil
}
