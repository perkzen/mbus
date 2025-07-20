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

/*
   ──────────────────────────
        globals / helpers
   ──────────────────────────
*/

var pgDb *sql.DB
var qb = store.Qb // shorthand

/*
   ──────────────────────────
               main
   ──────────────────────────
*/

func main() {
	if err := utils.Retry(
		"Initialize DB for seeding",
		10, 3*time.Second,
		initDB,
	); err != nil {
		log.Fatalf("❌ Could not connect / migrate DB: %v", err)
	}
	defer pgDb.Close()

	// ─── Seed bus-stations / lines (once) ────────────────────────────
	const expectedStations = 444
	var have int
	if err := pgDb.QueryRow(`SELECT COUNT(*) FROM bus_stations`).Scan(&have); err != nil {
		log.Fatalf("❌ COUNT(bus_stations) failed: %v", err)
	}

	if have < expectedStations {
		stations := loadSeedData("seed-weekday.json")

		lineIDs := upsertBusLines(stations)
		stationIDs, codeIDs := upsertBusStations(stations, lineIDs)

		log.Printf("✅ Inserted %d bus stations, %d station codes, %d lines.",
			len(stationIDs), len(codeIDs), len(lineIDs))
	} else {
		log.Printf("✅ DB already seeded with %d stations – skipping.", have)
	}

	// ─── Reload fresh id maps (in case rows already existed) ─────────
	lineIDs := loadLineIDs()
	codeIDs := loadStationCodeIDs()

	// ─── Departures ──────────────────────────────────────────────────
	for _, day := range []string{"weekday", "saturday", "sunday"} {
		insertDepartures(day, loadSeedData("seed-"+day+".json"), lineIDs, codeIDs)
	}
	log.Println("✅ All departures inserted.")
}

/*
   ──────────────────────────
            database
   ──────────────────────────
*/

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

/*
   ──────────────────────────
      seed-file helpers
   ──────────────────────────
*/

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

/*
   ──────────────────────────
       lines / stations
   ──────────────────────────
*/

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
			q, a, _ = qb.Insert("bus_lines").Columns("name").
				Values(line).Suffix("RETURNING id").ToSql()
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

func upsertBusStations(
	stations []marprom.BusStationWithDetails,
	lineIDs map[string]int,
) (map[string]int /*stationID*/, map[string]int /*codeID*/) {

	stationIDs := map[string]int{} // keyed by code
	codeIDs := map[string]int{}    // keyed by code

	for _, s := range stations {
		// 1️⃣ insert (or fetch) bus_stations row (no code column anymore)
		var stationID int
		q, a, _ := qb.Select("id").From("bus_stations").
			Where(sq.Eq{"name": s.Name}).ToSql()
		err := pgDb.QueryRow(q, a...).Scan(&stationID)
		if errors.Is(err, sql.ErrNoRows) {
			q, a, _ = qb.Insert("bus_stations").
				Columns("name", "image_url", "lat", "lng").
				Values(s.Name, s.ImageURL, s.Lat, s.Lon).
				Suffix("RETURNING id").ToSql()
			if err = pgDb.QueryRow(q, a...).Scan(&stationID); err != nil {
				log.Fatalf("❌ insert station %s: %v", s.Name, err)
			}
		} else if err != nil {
			log.Fatalf("❌ select station %s: %v", s.Name, err)
		}

		// 2️⃣ insert (or fetch) station_codes row
		codeInt, _ := strconv.Atoi(s.Code)
		var codeID int
		q, a, _ = qb.Select("id").From("station_codes").
			Where(sq.Eq{"code": codeInt}).ToSql()
		err = pgDb.QueryRow(q, a...).Scan(&codeID)
		if errors.Is(err, sql.ErrNoRows) {
			q, a, _ = qb.Insert("station_codes").
				Columns("station_id", "code").
				Values(stationID, codeInt).
				Suffix("RETURNING id").ToSql()
			if err = pgDb.QueryRow(q, a...).Scan(&codeID); err != nil {
				log.Fatalf("❌ insert station_code %s: %v", s.Code, err)
			}
		} else if err != nil {
			log.Fatalf("❌ select station_code %s: %v", s.Code, err)
		}

		stationIDs[s.Code] = stationID
		codeIDs[s.Code] = codeID

		// 3️⃣ link to lines (ignore duplicates)
		for _, l := range s.Lines {
			q, a, _ := qb.Insert("bus_stations_bus_lines").
				Columns("bus_station_id", "bus_line_id").
				Values(stationID, lineIDs[l]).
				Suffix("ON CONFLICT DO NOTHING").ToSql()
			if _, err := pgDb.Exec(q, a...); err != nil {
				log.Fatalf("❌ link station %d ↔ line %d: %v", stationID, lineIDs[l], err)
			}
		}
	}
	return stationIDs, codeIDs
}

/*
   ──────────────────────────
           departures
   ──────────────────────────
*/

func insertDepartures(
	schedule string,
	stations []marprom.BusStationWithDetails,
	lineIDs map[string]int,
	codeIDs map[string]int,
) {
	for _, s := range stations {
		codeID, ok := codeIDs[s.Code]
		if !ok {
			log.Fatalf("❌ code %s not in codeIDs map", s.Code)
		}

		for _, d := range s.Departures {
			lineID := lineIDs[d.Line]

			for _, t := range d.Times {
				q, a, _ := qb.Insert("departures").
					Columns("code_id", "line_id", "direction",
						"departure_time", "schedule_type").
					Values(codeID, lineID, d.Direction, t, schedule).ToSql()

				if _, err := pgDb.Exec(q, a...); err != nil {
					log.Fatalf("❌ dep insert stationCode %s line %s time %s: %v",
						s.Code, d.Line, t, err)
				}
			}
		}
	}
}

/*
   ──────────────────────────
        reload helpers
   ──────────────────────────
*/

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
