-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS bus_stations
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255)  NOT NULL,
    image_url  VARCHAR(255),
    lat        DECIMAL(9, 6) NOT NULL,
    lng        DECIMAL(9, 6) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- One bus station can have many bus lines


CREATE TABLE IF NOT EXISTS bus_lines
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS bus_stations_bus_lines
(
    bus_station_id INTEGER NOT NULL,
    bus_line_id    INTEGER NOT NULL,
    PRIMARY KEY (bus_station_id, bus_line_id),
    FOREIGN KEY (bus_station_id) REFERENCES bus_stations (id) ON DELETE CASCADE,
    FOREIGN KEY (bus_line_id) REFERENCES bus_lines (id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS station_codes
(
    id         SERIAL PRIMARY KEY,
    station_id INTEGER NOT NULL REFERENCES bus_stations (id) ON DELETE CASCADE,
    code       INT     NOT NULL UNIQUE,

    UNIQUE (station_id, code),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS departures
(
    id             SERIAL PRIMARY KEY,
    code_id        INT REFERENCES station_codes (id) ON DELETE CASCADE,
    line_id        INT REFERENCES bus_lines (id) ON DELETE CASCADE,
    direction      TEXT       NOT NULL,
    departure_time VARCHAR(5) NOT NULL CHECK (departure_time ~ '^\d{2}:\d{2}$'),
    schedule_type  TEXT       NOT NULL CHECK (schedule_type IN ('weekday', 'saturday', 'sunday')),
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (code_id, line_id, direction, departure_time, schedule_type)
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS bus_stations_bus_lines;
DROP TABLE IF EXISTS bus_lines;
DROP TABLE IF EXISTS bus_stations;
DROP TABLE IF EXISTS station_codes;
DROP TABLE IF EXISTS departures;
-- +goose StatementEnd
