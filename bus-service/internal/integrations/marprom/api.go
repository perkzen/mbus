package marprom

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type API interface {
	GetAvailableBusStations() ([]BusStation, error)
	GetBusStationDetails(code int) (*BusStationDetails, error)
	GetDeparturesByBusStation(filter *DepartureFilterOptions) ([]Departure, error)
	GetDeparturesFromStationToStation(fromCode, toCode int, date string) ([]Departure, error)
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

func (client *APIClient) GetBusStationDetails(code int, date string) (*BusStationDetails, error) {
	opts := &FetchOptions{
		URL: fmt.Sprintf("%s?stop=%d&datum=%s", client.baseURL, code, date),
	}

	log.Println("Fetching bus station details from", opts.URL)

	html, err := client.fetcher.FetchHTML(opts)
	if err != nil {
		log.Fatalf("failed to fetch HTML: %s", err)
		return nil, err
	}

	details, err := client.parser.ParseBusStationDetails(html)
	details.Code = strconv.Itoa(code)

	if err != nil {
		log.Fatalf("failed to parse bus station details: %s", err)
		return nil, err
	}

	log.Println("Parsed bus station details successfully for code:", code)

	return details, nil
}

type DepartureFilterOptions struct {
	Code int
	Line string
	Date string
}

func (client *APIClient) GetDeparturesByBusStation(filter *DepartureFilterOptions) ([]Departure, error) {
	station, err := client.GetBusStationDetails(filter.Code, filter.Date)
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

func (client *APIClient) GetDeparturesFromStationToStation(fromCode, toCode int, date string) ([]Departure, error) {
	fromStation, err := client.GetBusStationDetails(fromCode, date)
	if err != nil {
		log.Fatalf("failed to get from bus station details: %s", err)
		return nil, err
	}

	toStation, err := client.GetBusStationDetails(toCode, date)
	if err != nil {
		log.Fatalf("failed to get to bus station details: %s", err)
		return nil, err
	}

	if fromStation == nil || toStation == nil {
		log.Println("One of the bus stations not found")
		return nil, nil
	}

	// Build a set of lines available at the destination station
	toStationLines := make(map[string]struct{})
	for _, dep := range toStation.Departures {
		toStationLines[strings.ToLower(dep.Line)] = struct{}{}
	}

	// Filter only departures from source station that match a line at destination
	var sharedDepartures []Departure
	for _, dep := range fromStation.Departures {
		if _, ok := toStationLines[strings.ToLower(dep.Line)]; ok {
			sharedDepartures = append(sharedDepartures, dep)
		}
	}

	return sharedDepartures, nil
}
