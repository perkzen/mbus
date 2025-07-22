package store

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
)

type Direction struct {
	ID   int
	Name string
}

type DirectionStore interface {
	FindSharedDirectionsByCodes(fromCode, toCode int) ([]string, error)
	FindDirectionsByStationCode(stationCode int) ([]Direction, error)
}

type PostgresDirectionStore struct {
	db *sql.DB
}

func NewPostgresDirectionStore(db *sql.DB) *PostgresDirectionStore {
	return &PostgresDirectionStore{db: db}
}

func (store *PostgresDirectionStore) FindSharedDirectionsByCodes(fromCode, toCode int) ([]string, error) {
	queryBuilder := Qb.Select("DISTINCT dir.name").
		From("departures d1").
		Join("station_codes sc1 ON d1.code_id = sc1.id").
		Join("departures d2 ON d1.direction_id = d2.direction_id").
		Join("station_codes sc2 ON d2.code_id = sc2.id").
		Join("directions dir ON d1.direction_id = dir.id").
		Where(sq.Eq{
			"sc1.code": fromCode,
			"sc2.code": toCode,
		})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var directions []string
	for rows.Next() {
		var dir string
		if err := rows.Scan(&dir); err != nil {
			return nil, err
		}
		directions = append(directions, dir)
	}

	return directions, rows.Err()
}

func (store *PostgresDirectionStore) FindDirectionsByStationCode(stationCode int) ([]Direction, error) {
	queryBuilder := Qb.Select("d.id", "d.name").
		From("directions d").
		Join("departures dep ON dep.direction_id = d.id").
		Join("station_codes sc ON dep.code_id = sc.id").
		Where(sq.Eq{"sc.code": stationCode}).
		GroupBy("d.id, d.name")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var directions []Direction
	for rows.Next() {
		var dir Direction
		if err := rows.Scan(&dir.ID, &dir.Name); err != nil {
			return nil, err
		}
		directions = append(directions, dir)
	}

	return directions, rows.Err()
}
