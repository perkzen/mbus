-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE bus_stations
    ADD CONSTRAINT bus_stations_name_key UNIQUE (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE bus_stations
    DROP CONSTRAINT IF EXISTS bus_stations_name_key;
-- +goose StatementEnd
