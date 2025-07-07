package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/perkzen/mbus/bus-service/internal/config"
	"github.com/perkzen/mbus/bus-service/internal/db"
	"github.com/perkzen/mbus/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/bus-service/internal/store"
	"github.com/perkzen/mbus/bus-service/migrations"
	"log"
	"os"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	pgDb *sql.DB
)

func main() {
	initDB()
	defer pgDb.Close()

	stations := loadSeedData("data/seed.json")
	lineIDs := insertBusLines(stations)
	insertBusStationsAndLinks(stations, lineIDs)

	log.Println("✅ Seeding completed successfully.")
}

func initDB() {
	var err error

	env, err := config.LoadEnvironment()
	if err != nil {
		log.Fatalf("❌ Failed to load environment variables: %v", err)
	}
	pgDb, err = db.NewPostgresDB(env.PostgresURL).Open()
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}

	err = db.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
}

func loadSeedData(path string) []marprom.BusStationWithDetails {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("❌ Failed to read seed.json: %v", err)
	}

	var stations []marprom.BusStationWithDetails
	if err := json.Unmarshal(data, &stations); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	return stations
}

func insertBusLines(stations []marprom.BusStationWithDetails) map[string]int {
	lineIDs := make(map[string]int)
	uniqueRoutes := make(map[string]bool)

	for _, station := range stations {
		for _, route := range station.Lines {
			uniqueRoutes[route] = true
		}
	}

	for route := range uniqueRoutes {
		var id int

		query, args, _ := store.Qb.Select("id").
			From("bus_lines").
			Where(sq.Eq{"name": route}).
			ToSql()

		err := pgDb.QueryRow(query, args...).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			// Insert if not exists
			query, args, _ = store.Qb.Insert("bus_lines").
				Columns("name").
				Values(route).
				Suffix("RETURNING id").
				ToSql()

			err = pgDb.QueryRow(query, args...).Scan(&id)
			if err != nil {
				log.Fatalf("❌ Failed to insert bus_line '%s': %v", route, err)
			}
		} else if err != nil {
			log.Fatalf("❌ Failed to check bus_line '%s': %v", route, err)
		}

		lineIDs[route] = id
	}

	return lineIDs
}

func insertBusStationsAndLinks(stations []marprom.BusStationWithDetails, lineIDs map[string]int) {
	for _, s := range stations {
		stationID := insertBusStation(s)

		for _, route := range s.Lines {
			lineID := lineIDs[route]
			linkStationToLine(stationID, lineID)
		}
	}
}

func insertBusStation(s marprom.BusStationWithDetails) int {
	query, args, _ := store.Qb.Insert("bus_stations").
		Columns("code", "name", "image_url", "lat", "lng").
		Values(s.Code, s.Name, s.ImageURL, s.Lat, s.Lon).
		Suffix("RETURNING id").
		ToSql()

	var stationID int
	err := pgDb.QueryRow(query, args...).Scan(&stationID)
	if err != nil {
		log.Fatalf("❌ Failed to insert station %d: %v", s.Code, err)
	}

	return stationID
}

func linkStationToLine(stationID, lineID int) {
	query, args, _ := store.Qb.Insert("bus_stations_bus_lines").
		Columns("bus_station_id", "bus_line_id").
		Values(stationID, lineID).
		Suffix("ON CONFLICT DO NOTHING").
		ToSql()

	_, err := pgDb.Exec(query, args...)
	if err != nil {
		log.Fatalf("❌ Failed to link station %d to line %d: %v", stationID, lineID, err)
	}
}
