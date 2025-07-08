package main

import (
	"github.com/perkzen/mbus/apps/bus-service/internal/integrations/marprom"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"strconv"
)

func main() {
	client := marprom.NewAPIClient()
	stations, err := client.GetAvailableBusStations()

	if err != nil {
		panic(err)
	}

	data := make([]*marprom.BusStationWithDetails, 0)

	for _, station := range stations {

		code, _ := strconv.Atoi(station.Code)

		details, err := client.GetBusStationDetails(code, utils.Today())
		if err != nil {
			panic(err)
		}

		data = append(data, marprom.NewBusStationWithDetails(station, *details))

	}

	utils.SaveJSON("data/seed.json", data)
}
