-- +goose Up
CREATE TABLE work_home_taxi_price
(
    `datetime` DateTime,
    `from` String,
    `to` String,
    `waiting_time` UInt16,
    `duration` UInt16,
    `price` UInt32,
    `is_surge` Bool
) ENGINE = MergeTree()
ORDER BY datetime;

-- +goose Down
DROP TABLE work_home_taxi_price;