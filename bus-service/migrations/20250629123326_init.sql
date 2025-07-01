-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS bus_stations
(
    id         SERIAL PRIMARY KEY,
    code       INT           NOT NULL UNIQUE,
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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE IF EXISTS bus_stations_bus_lines;
DROP TABLE IF EXISTS bus_lines;
DROP TABLE IF EXISTS bus_stations;
-- +goose StatementEnd
