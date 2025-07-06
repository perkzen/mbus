package marprom

type Departure struct {
	Direction string   `json:"direction"`
	Times     []string `json:"times"`
	Line      string   `json:"line"`
}

type BusStation struct {
	ID   int     `json:"id"`
	Code string  `json:"code"`
	Name string  `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

func NewBusStation(id int, code, name string, lat, lon float64) *BusStation {
	return &BusStation{
		ID:   id,
		Code: code,
		Name: name,
		Lat:  lat,
		Lon:  lon,
	}
}

type BusStationDetails struct {
	ID         int         `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	Lines      []string    `json:"lines"`
	Departures []Departure `json:"departures"`
	ImageURL   string      `json:"imageUrl"`
}

type BusStationWithDetails struct {
	ID         int         `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	ImageURL   string      `json:"imageUrl"`
	Lat        float64     `json:"lat"`
	Lon        float64     `json:"lon"`
	Lines      []string    `json:"lines"`
	Departures []Departure `json:"departures"`
}

func NewBusStationWithDetails(station BusStation, details BusStationDetails) *BusStationWithDetails {
	return &BusStationWithDetails{
		ID:         station.ID,
		Code:       station.Code,
		Name:       station.Name,
		ImageURL:   details.ImageURL,
		Lat:        station.Lat,
		Lon:        station.Lon,
		Lines:      details.Lines,
		Departures: details.Departures,
	}
}
