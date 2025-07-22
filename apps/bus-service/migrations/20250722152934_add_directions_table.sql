-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS directions
(
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO directions (name)
SELECT DISTINCT direction
FROM departures
ON CONFLICT (name) DO NOTHING;

ALTER TABLE departures
    ADD COLUMN direction_id INTEGER;

UPDATE departures d
SET direction_id = dir.id
FROM directions dir
WHERE d.direction = dir.name;

ALTER TABLE departures
    ALTER COLUMN direction_id SET NOT NULL;
ALTER TABLE departures
    DROP COLUMN direction;

ALTER TABLE departures
    ADD CONSTRAINT fk_direction
        FOREIGN KEY (direction_id) REFERENCES directions (id) ON DELETE CASCADE;

ALTER TABLE departures
    DROP CONSTRAINT IF EXISTS departures_code_id_line_id_direction_departure_time_schedule_type_key;

ALTER TABLE departures
    ADD CONSTRAINT uq_departure
        UNIQUE (code_id, line_id, direction_id, departure_time, schedule_type);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE departures
    ADD COLUMN direction TEXT;

UPDATE departures d
SET direction = dir.name
FROM directions dir
WHERE d.direction_id = dir.id;

ALTER TABLE departures
    DROP CONSTRAINT IF EXISTS fk_direction;
ALTER TABLE departures
    DROP CONSTRAINT IF EXISTS uq_departure;
ALTER TABLE departures
    DROP COLUMN direction_id;

ALTER TABLE departures
    ADD CONSTRAINT uq_departure_old
        UNIQUE (code_id, line_id, direction, departure_time, schedule_type);

DROP TABLE IF EXISTS directions;

-- +goose StatementEnd
