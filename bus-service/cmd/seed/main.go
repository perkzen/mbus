package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Register the pgx driver for database/sql
)

type BusStation struct {
	ID         int         `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	ImageURL   string      `json:"imageUrl"`
	Routes     []string    `json:"routes"`
	Departures []Departure `json:"departures"`
	Location   GeoLocation `json:"location"`
}

type Departure struct {
	Direction string   `json:"direction"`
	Times     []string `json:"times"`
	Line      string   `json:"line"`
}

type GeoLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func main() {
	// Step 1: Connect to the database
	dsn := "host=localhost port=5432 user=postgres password=postgres dbname=mbus_db sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Step 2: Load and parse the JSON file
	data, err := os.ReadFile("data/seed.json")
	if err != nil {
		log.Fatalf("❌ Failed to read seed.json: %v", err)
	}

	var stations []BusStation
	if err := json.Unmarshal(data, &stations); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	// Step 3: Insert unique bus lines
	routeSet := make(map[string]bool)
	for _, station := range stations {
		for _, route := range station.Routes {
			routeSet[route] = true
		}
	}

	lineIDs := make(map[string]int)

	for route := range routeSet {
		var id int
		err := db.QueryRow(`SELECT id FROM bus_lines WHERE name = $1`, route).Scan(&id)

		if err == sql.ErrNoRows {
			// Not found — insert it
			err = db.QueryRow(`
				INSERT INTO bus_lines (name)
				VALUES ($1)
				RETURNING id
			`, route).Scan(&id)

			if err != nil {
				log.Fatalf("❌ Failed to insert bus_line '%s': %v", route, err)
			}
		} else if err != nil {
			log.Fatalf("❌ Failed to check bus_line '%s': %v", route, err)
		}

		lineIDs[route] = id
	}

	// Step 4: Insert bus stations and link them to lines
	for _, s := range stations {
		var stationID int
		err := db.QueryRow(`
			INSERT INTO bus_stations (code, name, image_url, lat, lng)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`, s.Code, s.Name, s.ImageURL, s.Location.Lat, s.Location.Lon).Scan(&stationID)

		if err != nil {
			log.Fatalf("❌ Failed to insert station %s: %v", s.Code, err)
		}

		for _, route := range s.Routes {
			lineID := lineIDs[route]
			_, err := db.Exec(`
				INSERT INTO bus_stations_bus_lines (bus_station_id, bus_line_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, stationID, lineID)

			if err != nil {
				log.Fatalf("❌ Failed to link station %d to line %d: %v", stationID, lineID, err)
			}
		}
	}

	fmt.Println("✅ Seeding completed successfully.")
}
