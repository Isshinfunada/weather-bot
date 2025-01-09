-- +goose Up
CREATE TABLE notification_history (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    notified_at TIMESTAMP NOT NULL,
    is_notify_trigger BOOLEAN,
    forecast_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE notification_history ;