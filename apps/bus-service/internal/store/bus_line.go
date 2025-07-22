package store

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

type BusLine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
} // @name BusLine

type BusLineStore interface {
	ListBusLines() ([]BusLine, error)
	FindSharedLinesByStations(fromId, toId int) ([]BusLine, error)
}

type PostgresBusLinesStore struct {
	db *sql.DB
}

func NewPostgresBusLineStore(db *sql.DB) *PostgresBusLinesStore {
	return &PostgresBusLinesStore{
		db: db,
	}
}

func (store *PostgresBusLinesStore) ListBusLines() ([]BusLine, error) {
	queryBuilder := Qb.Select("id", "name").
		From("bus_lines").
		OrderBy("regexp_replace(name, '[^0-9]', '', 'g')::int")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lines []BusLine
	for rows.Next() {
		var line BusLine
		if err := rows.Scan(&line.ID, &line.Name); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func (store *PostgresBusLinesStore) FindSharedLinesByStations(fromId, toId int) ([]BusLine, error) {
	queryBuilder := Qb.Select("bl.id", "bl.name").
		From("bus_lines bl").
		Join("bus_stations_bus_lines bsl1 ON bsl1.bus_line_id = bl.id").
		Join("bus_stations_bus_lines bsl2 ON bsl2.bus_line_id = bl.id").
		Where(sq.Eq{
			"bsl1.bus_station_id": fromId,
			"bsl2.bus_station_id": toId,
		}).
		GroupBy("bl.id", "bl.name")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	lines := make([]BusLine, 0)
	for rows.Next() {
		var line BusLine
		if err := rows.Scan(&line.ID, &line.Name); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
