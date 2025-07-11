package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/perkzen/mbus/apps/bus-service/data"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/db"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/perkzen/mbus/bus-service/migrations"
)

var pgDb *sql.DB

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
	err = pgDb.QueryRow("SELECT COUNT(*) FROM bus_stations").Scan(&count)
	if err != nil {
		log.Fatalf("❌ Failed to count bus_stations: %v", err)
	}

	if count < numOfStations {
		stations := loadSeedData("seed-weekday.json")
		lineIDs := insertBusLines(stations)
		insertBusStationsAndLinks(stations, lineIDs)
		log.Println("✅ Bus stations and lines inserted.")
	} else {
		log.Printf("✅ Database already seeded with %d bus stations, skipping station insert.", count)
	}

	lineIDs := reloadLineIDs()
	stationIDs := reloadStationIDs()

	insertDepartures("weekday", loadSeedData("seed-weekday.json"), lineIDs, stationIDs)
	insertDepartures("saturday", loadSeedData("seed-saturday.json"), lineIDs, stationIDs)
	insertDepartures("sunday", loadSeedData("seed-sunday.json"), lineIDs, stationIDs)

	log.Println("✅ All departures inserted successfully.")
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

	return db.MigrateFS(pgDb, migrations.FS, ".")
}

func loadSeedData(filename string) []marprom.BusStationWithDetails {
	file, err := data.FS.Open(filename)

	if err != nil {
		log.Fatalf("❌ Failed to open seed file %s: %v", filename, err)
	}
	defer file.Close()

	var stations []marprom.BusStationWithDetails
	if err := json.NewDecoder(file).Decode(&stations); err != nil {
		log.Fatalf("❌ Failed to parse seed file %s: %v", filename, err)
	}
	return stations
}

func insertBusLines(stations []marprom.BusStationWithDetails) map[string]int {
	lineIDs := make(map[string]int)
	uniqueLines := make(map[string]bool)

	for _, station := range stations {
		for _, line := range station.Lines {
			uniqueLines[line] = true
		}
	}

	for line := range uniqueLines {
		var id int
		query, args, _ := store.Qb.Select("id").From("bus_lines").Where(sq.Eq{"name": line}).ToSql()
		err := pgDb.QueryRow(query, args...).Scan(&id)

		if errors.Is(err, sql.ErrNoRows) {
			query, args, _ = store.Qb.Insert("bus_lines").
				Columns("name").
				Values(line).
				Suffix("RETURNING id").
				ToSql()
			err = pgDb.QueryRow(query, args...).Scan(&id)
			if err != nil {
				log.Fatalf("❌ Failed to insert bus line '%s': %v", line, err)
			}
		} else if err != nil {
			log.Fatalf("❌ Failed to fetch line '%s': %v", line, err)
		}

		lineIDs[line] = id
	}

	return lineIDs
}

func insertBusStationsAndLinks(stations []marprom.BusStationWithDetails, lineIDs map[string]int) map[string]int {
	stationIDs := make(map[string]int)

	for _, s := range stations {
		id := insertBusStation(s)
		stationIDs[s.Code] = id

		for _, line := range s.Lines {
			linkStationToLine(id, lineIDs[line])
		}
	}

	return stationIDs
}

func insertBusStation(s marprom.BusStationWithDetails) int {
	query, args, _ := store.Qb.Insert("bus_stations").
		Columns("code", "name", "image_url", "lat", "lng").
		Values(s.Code, s.Name, s.ImageURL, s.Lat, s.Lon).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := pgDb.QueryRow(query, args...).Scan(&id)
	if err != nil {
		log.Fatalf("❌ Failed to insert station %s: %v", s.Code, err)
	}
	return id
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

func insertDepartures(scheduleType string, stations []marprom.BusStationWithDetails, lineIDs, stationIDs map[string]int) {
	for _, s := range stations {
		for _, d := range s.Departures {
			lineID := lineIDs[d.Line]
			stationID := stationIDs[s.Code]

			for _, t := range d.Times {
				query, args, _ := store.Qb.Insert("departures").
					Columns("station_id", "line_id", "direction", "departure_time", "schedule_type").
					Values(stationID, lineID, d.Direction, t, scheduleType).
					ToSql()

				_, err := pgDb.Exec(query, args...)
				if err != nil {
					log.Fatalf("❌ Failed to insert departure for station %s, line %s: %v", s.Code, d.Line, err)
				}
			}
		}
	}
}

func reloadLineIDs() map[string]int {
	rows, err := pgDb.Query("SELECT id, name FROM bus_lines")
	if err != nil {
		log.Fatalf("❌ Failed to reload bus_lines: %v", err)
	}
	defer rows.Close()

	lines := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("❌ Failed to scan line: %v", err)
		}
		lines[name] = id
	}
	return lines
}

func reloadStationIDs() map[string]int {
	rows, err := pgDb.Query("SELECT id, code FROM bus_stations")
	if err != nil {
		log.Fatalf("❌ Failed to reload bus_stations: %v", err)
	}
	defer rows.Close()

	stations := make(map[string]int)
	for rows.Next() {
		var id int
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			log.Fatalf("❌ Failed to scan station: %v", err)
		}
		stations[code] = id
	}
	return stations
}
