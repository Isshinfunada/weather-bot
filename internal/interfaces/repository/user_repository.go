package repository

import (
    "database/sql"
)

type User struct {
    ID         int
    LineUserID string
}

type UserRepository struct {
    DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(lineUserID string) (*User, error) {
    var id int
    err := r.DB.QueryRow("INSERT INTO users (line_user_id) VALUES ($1) RETURNING id", lineUserID).Scan(&id)
    if err != nil {
        return nil, err
    }
    return &User{ID: id, LineUserID: lineUserID}, nil
}

func (r *UserRepository) GetUserByLineUserID(lineUserID string) (*User, error) {
    var user User
    err := r.DB.QueryRow("SELECT id, line_user_id FROM users WHERE line_user_id = $1", lineUserID).Scan(&user.ID, &user.LineUserID)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
