package departure

type Station struct {
	Name string `json:"name"`
	Code int    `json:"code"`
}

type Departure struct {
	Direction   string  `json:"direction"`
	Line        string  `json:"line"`
	FromStation Station `json:"fromStation"`
	ToStation   Station `json:"toStation"`
	Duration    string  `json:"duration"`
	Distance    float64 `json:"distance"`
	DepartureAt string  `json:"departureAt"`
	ArriveAt    string  `json:"arriveAt"`
}
