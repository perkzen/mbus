package store

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

type BusStation struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	ImageURL string   `json:"imageUrl"`
	Lat      float64  `json:"lat"`
	Lon      float64  `json:"lon"`
	Codes    []int    `json:"codes,omitempty"`
	Lines    []string `json:"lines,omitempty"`
} // @name BusStation

func (s *BusStation) SanitizedName() string {
	return strings.ReplaceAll(s.Name, "- ", "")
}

type StationCode struct {
	ID        int `json:"id"`
	StationID int `json:"stationId"`
	Code      int `json:"code"`
}

type BusStationFilterOptions struct {
	Name string
	Line string
}

type BusStationStore interface {
	ListBusStations(limit, offset int, opts *BusStationFilterOptions) ([]BusStation, error)
	FindBusStationByID(id int) (*BusStation, error)
	FindBusStationIDByCode(code string) (*StationCode, error)
}

type PostgresBusStationStore struct {
	db *sql.DB
}

func NewPostgresBusStationStore(db *sql.DB) *PostgresBusStationStore {
	return &PostgresBusStationStore{
		db: db,
	}
}

func (store *PostgresBusStationStore) ListBusStations(limit, offset int, opts *BusStationFilterOptions) ([]BusStation, error) {
	builder := Qb.Select(
		"bs.id",
		"bs.name",
		"bs.image_url",
		"bs.lat",
		"bs.lng",
		"COALESCE(array_agg(DISTINCT bl.name ORDER BY bl.name), '{}') AS lines",
	).
		From("bus_stations bs").
		LeftJoin("bus_stations_bus_lines bsl ON bsl.bus_station_id = bs.id").
		LeftJoin("bus_lines bl ON bl.id = bsl.bus_line_id").
		GroupBy("bs.id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		OrderBy("bs.name")

	if opts != nil {
		if opts.Line != "" {
			builder = builder.Where(sq.ILike{"bl.name": "%" + opts.Line + "%"})
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
		var rawLines pq.StringArray
		if err := rows.Scan(&s.ID, &s.Name, &s.ImageURL, &s.Lat, &s.Lon, &rawLines); err != nil {
			return nil, err
		}
		s.Lines = rawLines
		stations = append(stations, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}

func (store *PostgresBusStationStore) FindBusStationByID(id int) (*BusStation, error) {
	queryBuilder := Qb.
		Select(
			"bs.id",
			"bs.name",
			"bs.image_url",
			"bs.lat",
			"bs.lng",
			"COALESCE(array_agg(sc.code ORDER BY sc.code), '{}') AS codes",
		).
		From("bus_stations bs").
		LeftJoin("station_codes sc ON sc.station_id = bs.id").
		Where(sq.Eq{"bs.id": id}).
		GroupBy("bs.id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL: %w", err)
	}

	var station BusStation
	var rawCodes pq.Int64Array

	err = store.db.QueryRow(query, args...).Scan(
		&station.ID,
		&station.Name,
		&station.ImageURL,
		&station.Lat,
		&station.Lon,
		&rawCodes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	station.Codes = make([]int, len(rawCodes))
	for i, val := range rawCodes {
		station.Codes[i] = int(val)
	}

	return &station, nil
}

func (store *PostgresBusStationStore) FindBusStationIDByCode(code string) (*StationCode, error) {
	queryBuilder := Qb.Select("id", "station_id", "code").
		From("station_codes").
		Where(sq.Eq{"code": code})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building SQL: %w", err)
	}

	var stationCode StationCode
	err = store.db.QueryRow(query, args...).Scan(&stationCode.ID, &stationCode.StationID, &stationCode.Code)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %w", err)
	}

	return &stationCode, nil
}
