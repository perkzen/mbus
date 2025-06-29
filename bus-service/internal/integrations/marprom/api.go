package marprom

import (
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/utils"
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

	html, err := client.fetcher.FetchHTML(nil)
	if err != nil {
		fmt.Printf("Error fetching stations: %v\n", err)
		return nil, err
	}

	fmt.Println("Fetched bus stations HTML successfully")

	stations, err := client.parser.ParseBusStations(html)
	if err != nil {
		fmt.Printf("Error fetching stations: %v\n", err)
		return nil, err
	}

	fmt.Printf("Parsed %d bus stations successfully\n", len(stations))

	return stations, nil

}

func (client *APIClient) GetBusStationDetails(code string) (*BusStationDetails, error) {
	opts := &FetchOptions{
		URL: fmt.Sprintf("%s?stop=%s&datum=%s", client.baseURL, code, utils.Today()),
	}

	fmt.Println("Fetching bus station details for code:", code)

	html, err := client.fetcher.FetchHTML(opts)
	if err != nil {
		fmt.Printf("Error fetching station details: %v\n", err)
		return nil, err
	}

	details, err := client.parser.ParseBusStationDetails(html)

	if err != nil {
		fmt.Printf("Error parsing station details: %v\n", err)
		return nil, err
	}

	fmt.Println("Parsed bus station details successfully for code:", code)

	return details, nil
}
