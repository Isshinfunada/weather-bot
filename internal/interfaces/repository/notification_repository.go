package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
)

type NotificationRepository interface {
	InsertNotificationHistory(ctx context.Context, history *entity.NotificationHistory) error
}

type notificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) InsertNotificationHistory(ctx context.Context, history *entity.NotificationHistory) error {
	query := `
		INSERT INTO notification_history (user_id, notification_time, weather_data, is_notify_trigger, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()
	history.CreatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		history.UserID,
		history.NotificationTime,
		history.WeatherData,
		history.IsNotifyTrigger,
		history.CreatedAt,
	).Scan(&history.ID)

	if err != nil {
		return fmt.Errorf("failed to insert notification history: %w", err)
	}
	return nil
}
