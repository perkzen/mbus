package main

import (
	"github.com/perkzen/mbus/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/bus-service/internal/utils"
)

func main() {
	client := marprom.NewAPIClient()
	stations, err := client.GetAvailableBusStations()

	if err != nil {
		panic(err)
	}

	data := make([]*marprom.BusStationWithDetails, 0)

	for _, station := range stations {

		details, err := client.GetBusStationDetails(station.Code)
		if err != nil {
			panic(err)
		}

		data = append(data, marprom.NewBusStationWithDetails(station, *details))

	}

	utils.SaveJSON("data/seed.json", data)
}
