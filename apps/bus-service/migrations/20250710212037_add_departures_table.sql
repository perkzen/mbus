-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';


CREATE TABLE IF NOT EXISTS departures
(
    id             SERIAL PRIMARY KEY,
    station_id     INT REFERENCES bus_stations (id) ON DELETE CASCADE,
    line_id        INT REFERENCES bus_lines (id) ON DELETE CASCADE,
    direction      TEXT       NOT NULL,
    departure_time VARCHAR(5) NOT NULL CHECK (departure_time ~ '^\d{2}:\d{2}$'),
    schedule_type  TEXT       NOT NULL CHECK (schedule_type IN ('weekday', 'saturday', 'sunday')),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (station_id, line_id, direction, departure_time, schedule_type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS departures;
-- +goose StatementEnd
