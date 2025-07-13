package departure

type Station struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
} // @name TimetableRow.Station

type TimetableRow struct {
	Direction   string  `json:"direction"`
	Line        string  `json:"line"`
	FromStation Station `json:"fromStation"`
	ToStation   Station `json:"toStation"`
	Duration    string  `json:"duration"`
	Distance    float64 `json:"distance"`
	DepartureAt string  `json:"departureAt"`
	ArriveAt    string  `json:"arriveAt"`
} // @name TimetableRow

func (t TimetableRow) GetDepartureAt() string {
	return t.DepartureAt
}
