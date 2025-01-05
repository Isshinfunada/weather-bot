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
	FindUserByID(ctx context.Context, userID int) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// CreateUserはusersテーブルに新規レコードを挿入し、
// 作成したレコードのID　を取得して戻り値として返します
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

// FindUserByIDはusersテーブルを検索し、見つかったらUserを返します
func (r *userRepository) FindUserByID(ctx context.Context, userID int) (*entity.User, error) {
	query := `
		SELECT
            id, line_user_id, selected_area_id, notify_time,
            is_active, created_at, updated_at
        FROM users
        WHERE id = $1
        LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID)

	var u entity.User
	err := row.Scan(
		&u.ID,
		&u.LINEUserID,
		&u.SelectedAreaID,
		&u.NotifyTime,
		&u.IsActive,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &u, nil
}
