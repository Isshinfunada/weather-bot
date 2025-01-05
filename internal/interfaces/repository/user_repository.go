package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
	INSERT INTO users (line_user_id, selected_area_id, notify_time, is_active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id
	`

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	var newID int
	err := r.db.QueryRowContext(
		ctx, query,
		user.LINEUserID,
		user.SelectedAreaID,
		user.NotifyTime,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newID)
	if err != nil {
		return nil, fmt.Errorf("faild to insert user: %w", err)
	}

	user.ID = newID
	return user, nil
}
