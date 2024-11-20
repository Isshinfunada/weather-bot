package repository

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
)

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(dataSourceName string) *UserRepository {
    db, err := sql.Open("postgres", dataSourceName)
    if err != nil {
        log.Fatal(err)
    }

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    return &UserRepository{DB: db}
}
