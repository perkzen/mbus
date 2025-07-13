package store

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"time"
)

type ScheduleType string

const (
	ScheduleTypeWeekday  ScheduleType = "weekday"
	ScheduleTypeSaturday ScheduleType = "saturday"
	ScheduleTypeSunday   ScheduleType = "sunday"
)

type Departure struct {
	ID            int
	StationID     int
	LineID        int
	Line          BusLine
	Direction     string
	DepartureTime string
	ScheduleType  ScheduleType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DepartureStore interface {
	FindDeparturesByStationID(stationID int, scheduleType ScheduleType) ([]Departure, error)
	FindDeparturesFromStationToStation(fromStationID, toStationID int, scheduleType ScheduleType) ([]Departure, error)
	FindDeparturesByStationIDAndDirection(stationID int, direction string, scheduleType ScheduleType) ([]Departure, error)
	FindSharedDirections(toStationID, fromStationID int) ([]string, error)
}

type PostgresDepartureStore struct {
	db *sql.DB
}

func NewPostgresDepartureStore(db *sql.DB) *PostgresDepartureStore {
	return &PostgresDepartureStore{
		db: db,
	}
}

func ScheduleTyp(dateStr string) ScheduleType {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return ScheduleTypeWeekday
	}

	switch date.Weekday() {
	case time.Saturday:
		return ScheduleTypeSaturday
	case time.Sunday:
		return ScheduleTypeSunday
	default:
		return ScheduleTypeWeekday
	}
}

func (store *PostgresDepartureStore) FindDeparturesByStationID(stationID int, scheduleType ScheduleType) ([]Departure, error) {
	queryBuilder := Qb.Select("*").
		From("departures").
		Where(sq.Eq{"station_id": stationID, "schedule_type": scheduleType})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	departures := make([]Departure, 0)
	for rows.Next() {
		var dep Departure
		if err := rows.Scan(&dep.ID, &dep.StationID, &dep.LineID, &dep.Direction, &dep.DepartureTime, &dep.ScheduleType, &dep.CreatedAt, &dep.UpdatedAt); err != nil {
			return nil, err
		}
		departures = append(departures, dep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return departures, nil
}

func (store *PostgresDepartureStore) FindDeparturesFromStationToStation(fromStationID, toStationID int, scheduleType ScheduleType) ([]Departure, error) {
	queryBuilder := Qb.Select(
		"d1.id",
		"d1.station_id",
		"d1.line_id",
		"bl.id AS bus_line_id",
		"bl.name AS bus_line_name",
		"d1.direction",
		"d1.departure_time",
		"d1.schedule_type",
		"d1.created_at",
		"d1.updated_at",
	).
		From("departures d1").
		Join("bus_lines bl ON d1.line_id = bl.id").
		Where(sq.Eq{
			"d1.station_id":    fromStationID,
			"d1.schedule_type": scheduleType,
		}).
		Where("EXISTS (SELECT 1 FROM departures d2 WHERE d2.station_id = ? AND d2.schedule_type = d1.schedule_type AND d2.direction = d1.direction)", toStationID)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departures = make([]Departure, 0)
	for rows.Next() {
		var dep Departure
		if err := rows.Scan(
			&dep.ID,
			&dep.StationID,
			&dep.LineID,
			&dep.Line.ID,
			&dep.Line.Name,
			&dep.Direction,
			&dep.DepartureTime,
			&dep.ScheduleType,
			&dep.CreatedAt,
			&dep.UpdatedAt,
		); err != nil {
			return nil, err
		}
		departures = append(departures, dep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return departures, nil
}

func (store *PostgresDepartureStore) FindDeparturesByStationIDAndDirection(stationID int, direction string, scheduleType ScheduleType) ([]Departure, error) {
	queryBuilder := Qb.Select(
		"d.id",
		"d.station_id",
		"d.line_id",
		"bl.id AS bus_line_id",
		"bl.name AS bus_line_name",
		"d.direction",
		"d.departure_time",
		"d.schedule_type",
		"d.created_at",
		"d.updated_at",
	).
		From("departures d").
		Join("bus_lines bl ON d.line_id = bl.id").
		Where(sq.Eq{
			"d.station_id":    stationID,
			"d.direction":     direction,
			"d.schedule_type": scheduleType,
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

	departures := make([]Departure, 0)
	for rows.Next() {
		var dep Departure
		if err := rows.Scan(
			&dep.ID,
			&dep.StationID,
			&dep.LineID,
			&dep.Line.ID,
			&dep.Line.Name,
			&dep.Direction,
			&dep.DepartureTime,
			&dep.ScheduleType,
			&dep.CreatedAt,
			&dep.UpdatedAt,
		); err != nil {
			return nil, err
		}
		departures = append(departures, dep)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return departures, nil
}

func (store *PostgresDepartureStore) FindSharedDirections(toStationID, fromStationID int) ([]string, error) {
	queryBuilder := Qb.Select("DISTINCT d1.direction").
		From("departures d1").
		Join("departures d2 ON d1.direction = d2.direction").
		Where(sq.Eq{
			"d1.station_id": fromStationID,
			"d2.station_id": toStationID,
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

	directions := make([]string, 0)
	for rows.Next() {
		var direction string
		if err := rows.Scan(&direction); err != nil {
			return nil, err
		}
		directions = append(directions, direction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return directions, nil
}
