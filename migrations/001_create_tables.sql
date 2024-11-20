-- +migrate Up
CREATE TABLE test_data (
    id SERIAL PRIMARY KEY,
    line_user_id VARCHAR(255) UNIQUE NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE test_data;