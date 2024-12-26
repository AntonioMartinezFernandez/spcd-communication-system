-- +migrate Up
CREATE TABLE IF NOT EXISTS spcd_iot_devices (
    id VARCHAR(50) PRIMARY KEY,
    created_at TIMESTAMP WITHOUT TIME ZONE
);

-- +migrate Down
DROP TABLE IF EXISTS spcd_iot_devices CASCADE;