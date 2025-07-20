package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/perkzen/mbus/apps/bus-service/data"
	"github.com/perkzen/mbus/apps/bus-service/internal/config"
	"github.com/perkzen/mbus/apps/bus-service/internal/db"
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/apps/bus-service/internal/store"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"github.com/perkzen/mbus/apps/bus-service/migrations"
)

var pgDb *sql.DB
var qb = store.Qb

func main() {
	if err := utils.Retry("Initialize DB for seeding", 10, 3*time.Second, initDB); err != nil {
		log.Fatalf("❌ Could not connect / migrate DB: %v", err)
	}
	defer pgDb.Close()

	const expectedStations = 444
	var have int
	if err := pgDb.QueryRow(`SELECT COUNT(*) FROM bus_stations`).Scan(&have); err != nil {
		log.Fatalf("❌ COUNT(bus_stations) failed: %v", err)
	}

	if have < expectedStations {
		stations := loadSeedData("seed-weekday.json")
		lineIDs := upsertBusLines(stations)
		stationIDs, codeIDs := upsertBusStations(stations, lineIDs)
		log.Printf("✅ Inserted %d bus stations, %d station codes, %d lines.", len(stationIDs), len(codeIDs), len(lineIDs))
	} else {
		log.Printf("✅ DB already seeded with %d stations – skipping.", have)
	}

	lineIDs := loadLineIDs()
	codeIDs := loadStationCodeIDs()

	for _, day := range []string{"weekday", "saturday", "sunday"} {
		insertDepartures(day, loadSeedData("seed-"+day+".json"), lineIDs, codeIDs)
	}
	log.Println("✅ All departures inserted.")
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

func loadSeedData(file string) []marprom.BusStationWithDetails {
	f, err := data.FS.Open(file)
	if err != nil {
		log.Fatalf("❌ open %s: %v", file, err)
	}
	defer f.Close()

	var bs []marprom.BusStationWithDetails
	if err := json.NewDecoder(f).Decode(&bs); err != nil {
		log.Fatalf("❌ parse %s: %v", file, err)
	}
	return bs
}

func upsertBusLines(stations []marprom.BusStationWithDetails) map[string]int {
	lineIDs := make(map[string]int)
	unique := map[string]struct{}{}
	for _, s := range stations {
		for _, l := range s.Lines {
			unique[l] = struct{}{}
		}
	}

	for line := range unique {
		var id int
		q, a, _ := qb.Select("id").From("bus_lines").Where(sq.Eq{"name": line}).ToSql()
		err := pgDb.QueryRow(q, a...).Scan(&id)
		if errors.Is(err, sql.ErrNoRows) {
			q, a, _ = qb.Insert("bus_lines").Columns("name").Values(line).Suffix("RETURNING id").ToSql()
			if err = pgDb.QueryRow(q, a...).Scan(&id); err != nil {
				log.Fatalf("❌ insert line %s: %v", line, err)
			}
		} else if err != nil {
			log.Fatalf("❌ select line %s: %v", line, err)
		}
		lineIDs[line] = id
	}
	return lineIDs
}

func upsertBusStations(stations []marprom.BusStationWithDetails, lineIDs map[string]int) (map[string]int, map[string]int) {
	uniqueStations := map[string]marprom.BusStationWithDetails{}
	for _, s := range stations {
		if _, exists := uniqueStations[s.Name]; !exists {
			uniqueStations[s.Name] = s
		}
	}

	qbInsert := qb.Insert("bus_stations").
		Columns("name", "image_url", "lat", "lng").
		Suffix("ON CONFLICT (name) DO NOTHING")
	for _, s := range uniqueStations {
		qbInsert = qbInsert.Values(s.Name, s.ImageURL, s.Lat, s.Lon)
	}
	query, args, err := qbInsert.ToSql()
	if err != nil {
		log.Fatalf("❌ build insert stations: %v", err)
	}
	if _, err := pgDb.Exec(query, args...); err != nil {
		log.Fatalf("❌ exec insert stations: %v", err)
	}

	rows, err := pgDb.Query("SELECT id, name FROM bus_stations")
	if err != nil {
		log.Fatalf("❌ fetch station IDs: %v", err)
	}
	defer rows.Close()

	stationNameToID := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("❌ scan station: %v", err)
		}
		stationNameToID[name] = id
	}

	qbCodes := qb.Insert("station_codes").
		Columns("station_id", "code").
		Suffix("ON CONFLICT (code) DO NOTHING")
	for _, s := range stations {
		codeInt, _ := strconv.Atoi(s.Code)
		stationID := stationNameToID[s.Name]
		qbCodes = qbCodes.Values(stationID, codeInt)
	}
	query, args, err = qbCodes.ToSql()
	if err != nil {
		log.Fatalf("❌ build insert station_codes: %v", err)
	}
	if _, err := pgDb.Exec(query, args...); err != nil {
		log.Fatalf("❌ exec insert station_codes: %v", err)
	}

	codeIDs := loadStationCodeIDs()

	qbLinks := qb.Insert("bus_stations_bus_lines").
		Columns("bus_station_id", "bus_line_id").
		Suffix("ON CONFLICT DO NOTHING")
	for _, s := range stations {
		stationID := stationNameToID[s.Name]
		for _, line := range s.Lines {
			qbLinks = qbLinks.Values(stationID, lineIDs[line])
		}
	}
	query, args, err = qbLinks.ToSql()
	if err != nil {
		log.Fatalf("❌ build insert station-line links: %v", err)
	}
	if _, err := pgDb.Exec(query, args...); err != nil {
		log.Fatalf("❌ exec insert station-line links: %v", err)
	}

	stationIDs := make(map[string]int)
	for _, s := range stations {
		stationIDs[s.Code] = stationNameToID[s.Name]
	}

	return stationIDs, codeIDs
}

func insertDepartures(schedule string, stations []marprom.BusStationWithDetails, lineIDs, codeIDs map[string]int) {
	const batchSize = 1000
	type row struct {
		CodeID        int
		LineID        int
		Direction     string
		DepartureTime string
		ScheduleType  string
	}

	var buffer []row
	flush := func() {
		if len(buffer) == 0 {
			return
		}
		qbInsert := qb.Insert("departures").Columns("code_id", "line_id", "direction", "departure_time", "schedule_type")
		for _, r := range buffer {
			qbInsert = qbInsert.Values(r.CodeID, r.LineID, r.Direction, r.DepartureTime, r.ScheduleType)
		}
		query, args, err := qbInsert.ToSql()
		if err != nil {
			log.Fatalf("❌ squirrel build batch insert: %v", err)
		}
		if _, err := pgDb.Exec(query, args...); err != nil {
			log.Fatalf("❌ batch insert departures (%d rows): %v", len(buffer), err)
		}
		buffer = buffer[:0]
	}

	for _, s := range stations {
		codeID := codeIDs[s.Code]
		for _, d := range s.Departures {
			lineID := lineIDs[d.Line]
			for _, t := range d.Times {
				buffer = append(buffer, row{
					CodeID:        codeID,
					LineID:        lineID,
					Direction:     d.Direction,
					DepartureTime: t,
					ScheduleType:  schedule,
				})
				if len(buffer) >= batchSize {
					flush()
				}
			}
		}
	}
	flush()
}

func loadLineIDs() map[string]int {
	rows, err := pgDb.Query(`SELECT id, name FROM bus_lines`)
	if err != nil {
		log.Fatalf("❌ reload lines: %v", err)
	}
	defer rows.Close()

	lines := map[string]int{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("❌ scan line: %v", err)
		}
		lines[name] = id
	}
	return lines
}

func loadStationCodeIDs() map[string]int {
	rows, err := pgDb.Query(`SELECT id, code FROM station_codes`)
	if err != nil {
		log.Fatalf("❌ reload station_codes: %v", err)
	}
	defer rows.Close()

	sc := map[string]int{}
	for rows.Next() {
		var id int
		var code int
		if err := rows.Scan(&id, &code); err != nil {
			log.Fatalf("❌ scan station_code: %v", err)
		}
		sc[strconv.Itoa(code)] = id
	}
	return sc
}
