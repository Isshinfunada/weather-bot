-- +goose Up

CREATE TABLE IF NOT EXISTS area_centers (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    en_name VARCHAR(255) NOT NULL,
    office_name VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS area_offices (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    en_name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(10) REFERENCES area_centers(id)
);

CREATE TABLE IF NOT EXISTS area_class10 (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    en_name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(10) REFERENCES area_offices(id)
);

CREATE TABLE IF NOT EXISTS area_class15 (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    en_name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(10) REFERENCES area_class10(id)
);

CREATE TABLE IF NOT EXISTS area_class20 (
    id VARCHAR(10) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    en_name VARCHAR(255) NOT NULL,
    parent_id VARCHAR(10) REFERENCES area_class15(id)
);

-- +goose Down
DROP TABLE IF EXISTS area_class20;
DROP TABLE IF EXISTS area_class15;
DROP TABLE IF EXISTS area_class10;
DROP TABLE IF EXISTS area_offices;
DROP TABLE IF EXISTS area_centers;