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
	FindUserByLINEUserID(ctx context.Context, LINEUserID string) (*entity.User, error)
	FindUserByNotifyTimeRange(ctx context.Context, start, end time.Time) ([]*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID int) error
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

func (r *userRepository) FindUserByLINEUserID(ctx context.Context, LINEUserID string) (*entity.User, error) {
	query := `
		SELECT
			id, line_user_id, selected_area_id, notify_time,
			is_active, created_at, updated_at
		FROM users
		WHERE line_user_id = $1
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, LINEUserID)

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
		return nil, fmt.Errorf("failed to get user by LINEUserID: %w", err)
	}
	return &u, nil

}

func (r *userRepository) FindUserByNotifyTimeRange(ctx context.Context, start, end time.Time) ([]*entity.User, error) {
	query := `
		SELECT id, line_user_id, selected_area_id, notify_time, is_active, created_at, updated_at
		FROM users
		WHERE notify_time >= $1 AND notify_time < $2
	`

	rows, err := r.db.QueryContext(ctx, query, start.Format("15:04"), end.Format("15:04"))
	if err != nil {
		return nil, fmt.Errorf("failed to query users by notify time range: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.LINEUserID, &u.SelectedAreaID, &u.NotifyTime, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET
			selected_area_id = $1,
			notify_time = $2,
			is_active = $3,
			updated_at = $4
		WHERE id = $5
	`

	user.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.SelectedAreaID,
		user.NotifyTime,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated (id=%d not found)", user.ID)
	}
	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted (id=%d not found)", userID)
	}

	return nil
}
