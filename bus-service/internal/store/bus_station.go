package store

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type BusStation struct {
	Code     int     `json:"code"`
	Name     string  `json:"name"`
	ImageURL string  `json:"imageUrl"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
}

type BusStationFilterOptions struct {
	Name string
	Line string
}

type BusStationStore interface {
	ListBusStations(limit, offset int, opts *BusStationFilterOptions) ([]BusStation, error)
	FindBusStationByCode(code int) (*BusStation, error)
}

type PostgresBusStationStore struct {
	db *sql.DB
}

func NewPostgresBusStationStore(db *sql.DB) *PostgresBusStationStore {
	return &PostgresBusStationStore{
		db: db,
	}
}

func (store *PostgresBusStationStore) FindBusStationByCode(code int) (*BusStation, error) {
	fmt.Println("Fetching station with code:", code)

	queryBuilder := Qb.
		Select("code", "name", "image_url", "lat", "lng").
		From("bus_stations").
		Where(sq.Eq{"code": code})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL: %w", err)
	}

	fmt.Printf("SQL: %s\nARGS: %v\n", query, args)

	var station BusStation
	err = store.db.QueryRow(query, args...).Scan(
		&station.Code, &station.Name, &station.ImageURL, &station.Lat, &station.Lon,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &station, nil
}

func (store *PostgresBusStationStore) ListBusStations(limit, offset int, opts *BusStationFilterOptions) ([]BusStation, error) {
	builder := Qb.Select("bs.code", "bs.name", "bs.image_url", "bs.lat", "bs.lng").
		From("bus_stations bs").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy("bs.name")

	if opts != nil {
		if opts.Line != "" {
			builder = builder.
				Join("bus_stations_bus_lines bsl ON bsl.bus_station_id = bs.id").
				Join("bus_lines bl ON bl.id = bsl.bus_line_id").
				Where(sq.ILike{"bl.name": "%" + opts.Line + "%"})
		}
		if opts.Name != "" {
			builder = builder.Where(sq.ILike{"bs.name": opts.Name + "%"})
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("sql build error: %w", err)
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	stations := make([]BusStation, 0)
	for rows.Next() {
		var s BusStation
		if err := rows.Scan(&s.Code, &s.Name, &s.ImageURL, &s.Lat, &s.Lon); err != nil {
			return nil, err
		}
		stations = append(stations, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}
