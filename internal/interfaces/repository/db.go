package repository

import (
	"database/sql"
	"fmt"
)

func NewDBConnection(host, port, user, pass, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s pass=%s dbName=%s sslmode=disable",
		host, port, user, pass, dbName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
