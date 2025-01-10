package repository_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/Isshinfunada/weather-bot/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupNotificationRepoTest(t *testing.T) (repository.NotificationRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := repository.NewNotificationRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestInsertNotificationHistory_Success(t *testing.T) {
	utils.JST = time.FixedZone("JST", 9*60*60)
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().In(utils.JST)
	history := &entity.NotificationHistory{
		UserID:           1,
		NotificationTime: now,
		WeatherData:      []byte(`{"temp": "25"}`),
		IsNotifyTrigger:  true,
	}

	query := regexp.QuoteMeta(`
		INSERT INTO notification_history (user_id, notification_time, weather_data, is_notify_trigger, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`)

	mock.ExpectQuery(query).
		WithArgs(history.UserID, history.NotificationTime, history.WeatherData, history.IsNotifyTrigger, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(42))

	err := repo.InsertNotificationHistory(ctx, history)
	require.NoError(t, err)
	assert.Equal(t, 42, history.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInsertNotificationHistory_Failure(t *testing.T) {
	utils.JST = time.FixedZone("JST", 9*60*60)
	repo, mock, cleanup := setupNotificationRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now().In(utils.JST)
	history := &entity.NotificationHistory{
		UserID:           1,
		NotificationTime: now,
		WeatherData:      []byte(`{"temp": "25"}`),
		IsNotifyTrigger:  true,
	}

	query := regexp.QuoteMeta(`
		INSERT INTO notification_history (user_id, notification_time, weather_data, is_notify_trigger, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`)

	mock.ExpectQuery(query).
		WithArgs(history.UserID, history.NotificationTime, history.WeatherData, history.IsNotifyTrigger, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	err := repo.InsertNotificationHistory(ctx, history)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to insert notification history")
	assert.NoError(t, mock.ExpectationsWereMet())
}
