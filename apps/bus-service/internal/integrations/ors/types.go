package ors

type MatrixRequest struct {
	Locations        [][]float64 `json:"locations"`
	Metrics          []string    `json:"metrics"` // e.g. ["distance", "duration"]
	ResolveLocations bool        `json:"resolve_locations"`
	Units            string      `json:"units"` // "km", "m"
}

type MatrixResponse struct {
	Distances [][]float64 `json:"distances"`
	Durations [][]float64 `json:"durations"`
}
