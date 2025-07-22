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
	StationCodeID int
	LineID        int
	Line          BusLine
	Direction     string
	DepartureTime string
	ScheduleType  ScheduleType
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DepartureStore interface {
	FindDeparturesByStationCode(stationCode int, scheduleType ScheduleType) ([]Departure, error)
	FindDepartures(fromCode, toCode int, scheduleType ScheduleType) ([]Departure, error)
	FindDeparturesByStationCodeAndDirection(stationCode int, direction string, scheduleType ScheduleType) ([]Departure, error)
}

type PostgresDepartureStore struct {
	db *sql.DB
}

func NewPostgresDepartureStore(db *sql.DB) *PostgresDepartureStore {
	return &PostgresDepartureStore{db: db}
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

func (store *PostgresDepartureStore) FindDeparturesByStationCode(stationCode int, scheduleType ScheduleType) ([]Departure, error) {
	queryBuilder := Qb.Select(
		"d.id",
		"d.code_id",
		"d.line_id",
		"bl.id AS bus_line_id",
		"bl.name AS bus_line_name",
		"dir.name AS direction",
		"d.departure_time",
		"d.schedule_type",
		"d.created_at",
		"d.updated_at",
	).
		From("departures d").
		Join("station_codes sc ON d.code_id = sc.id").
		Join("bus_lines bl ON d.line_id = bl.id").
		Join("directions dir ON d.direction_id = dir.id").
		Where(sq.Eq{
			"sc.code":         stationCode,
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

	var departures []Departure
	for rows.Next() {
		var dep Departure
		if err := rows.Scan(
			&dep.ID,
			&dep.StationCodeID,
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

	return departures, rows.Err()
}

func (store *PostgresDepartureStore) FindDepartures(fromCode, toCode int, scheduleType ScheduleType) ([]Departure, error) {
	queryBuilder := Qb.Select(
		"d1.id",
		"d1.code_id",
		"d1.line_id",
		"bl.id AS bus_line_id",
		"bl.name AS bus_line_name",
		"dir.name AS direction",
		"d1.departure_time",
		"d1.schedule_type",
		"d1.created_at",
		"d1.updated_at",
	).
		From("departures d1").
		Join("station_codes sc1 ON d1.code_id = sc1.id").
		Join("bus_lines bl ON d1.line_id = bl.id").
		Join("directions dir ON d1.direction_id = dir.id").
		Where(sq.Eq{
			"sc1.code":         fromCode,
			"d1.schedule_type": scheduleType,
		}).
		Where(`
			EXISTS (
				SELECT 1 FROM departures d2
				JOIN station_codes sc2 ON d2.code_id = sc2.id
				WHERE sc2.code = ? AND d2.schedule_type = d1.schedule_type AND d2.direction_id = d1.direction_id
			)
		`, toCode)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := store.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departures []Departure
	for rows.Next() {
		var dep Departure
		if err := rows.Scan(
			&dep.ID,
			&dep.StationCodeID,
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

	return departures, rows.Err()
}

func (store *PostgresDepartureStore) FindDeparturesByStationCodeAndDirection(stationCode int, direction string, scheduleType ScheduleType) ([]Departure, error) {
	queryBuilder := Qb.Select(
		"d.id",
		"d.code_id",
		"d.line_id",
		"bl.id AS bus_line_id",
		"bl.name AS bus_line_name",
		"dir.name AS direction",
		"d.departure_time",
		"d.schedule_type",
		"d.created_at",
		"d.updated_at",
	).
		From("departures d").
		Join("station_codes sc ON d.code_id = sc.id").
		Join("bus_lines bl ON d.line_id = bl.id").
		Join("directions dir ON d.direction_id = dir.id").
		Where(sq.Eq{
			"sc.code":         stationCode,
			"dir.name":        direction,
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

	var departures []Departure
	for rows.Next() {
		var dep Departure
		if err := rows.Scan(
			&dep.ID,
			&dep.StationCodeID,
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

	return departures, rows.Err()
}
