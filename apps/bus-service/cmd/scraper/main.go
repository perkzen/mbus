package main

import (
	"github.com/perkzen/mbus/apps/bus-service/internal/provider/marprom"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
)

type dayOption struct {
	DateFn   func() string
	Filename string
}

func main() {
	dayOptions := map[string]dayOption{
		"weekday":  {utils.Weekday, "data/seed-weekday.json"},
		"saturday": {utils.Saturday, "data/seed-saturday.json"},
		"sunday":   {utils.Sunday, "data/seed-sunday.json"},
	}

	var (
		date     string
		filename string
	)

	if len(os.Args) > 1 {
		arg := os.Args[1]
		if opt, ok := dayOptions[arg]; ok {
			date = opt.DateFn()
			filename = opt.Filename
		} else {
			log.Println("Invalid argument. Use: weekday, saturday, sunday, or none.")
			os.Exit(1)
		}
	} else {
		date = time.Now().Format("2006-01-02")
		filename = "data/seed-today.json"
	}

	log.Println("Using date:", date)

	client := marprom.NewAPIClient()
	stations, err := client.GetAvailableBusStations()
	if err != nil {
		log.Fatalf("Failed to fetch bus stations: %v", err)
	}

	data := make([]*marprom.BusStationWithDetails, 0)

	for _, station := range stations {
		code, _ := strconv.Atoi(station.Code)

		details, err := client.GetBusStationDetails(code, date)
		if err != nil {
			log.Fatalf("Failed to get details for station %s (%d): %v", station.Name, code, err)
		}

		data = append(data, marprom.NewBusStationWithDetails(station, *details))

		time.Sleep(1500 * time.Millisecond)
	}

	err = utils.SaveJSON(filename, data)
	if err != nil {
		log.Fatalf("Failed to save data to %s: %v", filename, err)
	}
}
