package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/perkzen/mbus/apps/bus-service/data"
	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/db"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/perkzen/mbus/bus-service/migrations"
	"log"
	"time"
)

var (
	pgDb *sql.DB
)

func main() {
	err := utils.Retry("Initialize DB for seeding", 10, 3*time.Second, func() error {
		return initDB()
	})

	if err != nil {
		log.Fatalf("❌ Could not connect and initialize DB: %v", err)
	}

	defer pgDb.Close()

	numOfStations := 444
	var count int
	query := "SELECT COUNT(*) FROM bus_stations"
	err = pgDb.QueryRow(query).Scan(&count)
	if err != nil {
		log.Fatalf("❌ Failed to count bus_stations: %v", err)
	}

	if count == numOfStations {
		log.Printf("✅ Database already seeded with %d bus stations, skipping seeding.", numOfStations)
		return
	}

	stations := loadSeedData("./data/seed.json")
	lineIDs := insertBusLines(stations)
	insertBusStationsAndLinks(stations, lineIDs)

	log.Println("✅ Seeding completed successfully.")
}

func initDB() error {
	env, err := config.LoadEnvironment()
	if err != nil {
		return err
	}

	pgDb, err = db.NewPostgresDB(env.PostgresURL).Open()
	if err != nil {
		return err
	}

	err = db.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		return err
	}

	return nil
}

func loadSeedData(path string) []marprom.BusStationWithDetails {
	file, err := data.FS.Open("seed.json")
	if err != nil {
		log.Fatalf("❌ Failed to open embedded seed.json: %v", err)
	}
	defer file.Close()

	var stations []marprom.BusStationWithDetails
	if err := json.NewDecoder(file).Decode(&stations); err != nil {
		log.Fatalf("❌ Failed to parse seed.json: %v", err)
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
