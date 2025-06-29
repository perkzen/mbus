package types

// GeoLocation holds latitude and longitude coordinates for a bus stop.
//
// Example JSON: { "lat": 46.5547, "lon": 15.6459 }
type GeoLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func NewGeoLocation(lat, lon float64) *GeoLocation {
	return &GeoLocation{
		Lat: lat,
		Lon: lon,
	}
}

type Departure struct {
	Direction string   `json:"direction"`
	Times     []string `json:"times"`
	Line      string   `json:"line"`
}

type BusStation struct {
	ID         int         `json:"id"`
	Code       string      `json:"code"`
	Name       string      `json:"name"`
	ImageURL   string      `json:"imageUrl"`
	Routes     []string    `json:"routes"`
	Departures []Departure `json:"departures"`
	Location   GeoLocation `json:"location"`
}

func NewBusStation(id int, code, name string, lat, lon float64) *BusStation {
	return &BusStation{
		ID:       id,
		Code:     code,
		Name:     name,
		Location: *NewGeoLocation(lat, lon),
	}
}

type BusStationDetails struct {
	ID         int         `json:"id"`
	Lines      []string    `json:"lines"`
	Departures []Departure `json:"departures"`
	ImageURL   string      `json:"imageUrl"`
}
