package marprom

import (
	"fmt"
	"github.com/perkzen/mbus/bus-service/internal/common/fetchers"
	"github.com/perkzen/mbus/bus-service/internal/common/logger"
	"github.com/perkzen/mbus/bus-service/internal/types"
	"github.com/perkzen/mbus/bus-service/internal/utils"
)

type IMarpromAPI interface {
	GetAvailableBusStations() ([]types.BusStation, error)
	GetBusStationDetails(code string) (*types.BusStationDetails, error)
}

type MarpromAPIClient struct {
	baseURL string
	fetcher *fetchers.HTMLFetcher
	parser  *MarpromHTMLParser
}

func NewMarpromAPIClient() *MarpromAPIClient {
	return &MarpromAPIClient{
		baseURL: "https://vozniredi.marprom.si/",
		fetcher: fetchers.NewHTMLFetcher(),
		parser:  NewMarpromHTMLParser(),
	}
}

func (client *MarpromAPIClient) GetAvailableBusStations() ([]types.BusStation, error) {

	html, err := client.fetcher.FetchHTML(nil)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Fetched HTML for bus stations")

	stations, err := client.parser.ParseBusStations(html)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Parsed bus stations successfully", "count", len(stations))

	return stations, nil

}

func (client *MarpromAPIClient) GetBusStationDetails(code string) (*types.BusStationDetails, error) {
	opts := &fetchers.FetchOptions{
		URL: fmt.Sprintf("%s?stop=%s&datum=%s", client.baseURL, code, utils.Today()),
	}

	logger.Info("Fetching bus station details", "code", code, "url", opts.URL)

	html, err := client.fetcher.FetchHTML(opts)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	details, err := client.parser.ParseBusStationDetails(html)

	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	logger.Info("Parsed bus station details successfully", "code", code, "lines", len(details.Lines), "departures", len(details.Departures))

	return details, nil
}
