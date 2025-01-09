-- +goose Up
CREATE TABLE weather_notification_rules (
    weather_code VARCHAR(10) PRIMARY KEY,
    weather_description VARCHAR(255) NOT NULL,
    is_notify_trigger BOOLEAN NOT NULL
);

-- +goose Down
DROP TABLE weather_notification_rules ;