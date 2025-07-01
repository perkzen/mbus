package marprom

import (
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/utils"
	"log"
	"strings"
)

type API interface {
	GetAvailableBusStations() ([]BusStation, error)
	GetBusStationDetails(code string) (*BusStationDetails, error)
}

type APIClient struct {
	baseURL string
	fetcher *HTMLFetcher
	parser  *HTMLParser
}

func NewAPIClient() *APIClient {
	return &APIClient{
		baseURL: "https://vozniredi.marprom.si/",
		fetcher: NewHTMLFetcher(),
		parser:  NewHTMLParser(),
	}
}

func (client *APIClient) GetAvailableBusStations() ([]BusStation, error) {

	html, err := client.fetcher.FetchHTML(&FetchOptions{
		URL: client.baseURL,
	})
	if err != nil {
		log.Fatalf("failed to fetch HTML: %s", err)
		return nil, err
	}

	fmt.Println("Fetched bus stations HTML successfully")

	stations, err := client.parser.ParseBusStations(html)
	if err != nil {
		log.Fatalf("failed to parse bus stations: %s", err)
		return nil, err
	}

	log.Println("Parsed bus stations HTML successfully")

	return stations, nil

}

func (client *APIClient) GetBusStationDetails(code string) (*BusStationDetails, error) {
	opts := &FetchOptions{
		URL: fmt.Sprintf("%s?stop=%s&datum=%s", client.baseURL, code, utils.Today()),
	}

	log.Println("Fetching bus station details from", opts.URL)

	html, err := client.fetcher.FetchHTML(opts)
	if err != nil {
		log.Fatalf("failed to fetch HTML: %s", err)
		return nil, err
	}

	details, err := client.parser.ParseBusStationDetails(html)

	if err != nil {
		log.Fatalf("failed to parse bus station details: %s", err)
		return nil, err
	}

	log.Println("Parsed bus station details successfully for code:", code)

	return details, nil
}

type DepartureFilterOptions struct {
	Code string
	Line string
}

func (client *APIClient) GetDeparturesByBusStation(filter *DepartureFilterOptions) ([]Departure, error) {
	station, err := client.GetBusStationDetails(filter.Code)
	if err != nil {
		log.Fatalf("failed to get bus station details: %s", err)
		return nil, err
	}

	if station == nil {
		log.Println("No bus station found with code:", filter.Code)
		return nil, nil
	}

	departures := make([]Departure, 0)

	for _, dep := range station.Departures {
		if filter.Line != "" && !strings.EqualFold(dep.Line, filter.Line) {
			continue
		}
		departures = append(departures, dep)
	}

	return departures, nil
}
