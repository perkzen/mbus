package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	sq "github.com/Masterminds/squirrel"
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

	tables := []string{
		"departures",
		"bus_stations_bus_lines",
		"bus_lines",
		"station_codes",
		"bus_stations",
	}

	query := "TRUNCATE " + strings.Join(tables, ", ") + " RESTART IDENTITY CASCADE"

	if err := execRaw(pgDb, query); err != nil {
		log.Fatalf("❌ Failed to truncate tables: %v", err)
	}

	log.Println("✅ All tables truncated successfully.")
}

func execRaw(db *sql.DB, query string) error {
	sqlStr, args, err := sq.Expr(query).ToSql()
	if err != nil {
		return fmt.Errorf("build raw sql: %w", err)
	}
	_, err = db.Exec(sqlStr, args...)
	return err
}
