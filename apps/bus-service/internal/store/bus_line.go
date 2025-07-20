package store

import (
	"database/sql"
)

type BusLine struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
} // @name BusLine

type BusLineStore interface {
	ListBusLines() ([]BusLine, error)
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
