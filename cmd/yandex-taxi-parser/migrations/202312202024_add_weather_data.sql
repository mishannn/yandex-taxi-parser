-- +goose Up
ALTER TABLE work_home_taxi_price ADD COLUMN `temperature` Nullable(Float32);
ALTER TABLE work_home_taxi_price ADD COLUMN `rain_level` Nullable(Float32);
ALTER TABLE work_home_taxi_price ADD COLUMN `snow_level` Nullable(Float32);

-- +goose Down
ALTER TABLE work_home_taxi_price DROP COLUMN `temperature`;
ALTER TABLE work_home_taxi_price DROP COLUMN `rain_level`;
ALTER TABLE work_home_taxi_price DROP COLUMN `snow_level`;