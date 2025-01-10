package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWeatherRuleRepoTest(t *testing.T) (repository.WeatherRuleRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	repo := repository.NewWeatherRuleRepository(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestGetRule_Success(t *testing.T) {
	repo, mock, cleanup := setupWeatherRuleRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	weatherCode := "100"
	expectedRule := &entity.WeatherRule{
		WeatherCode:        "100",
		WeatherDescription: "晴",
		IsNotifyTrigger:    false,
	}

	query := regexp.QuoteMeta(`
	SELECT weather_code, weather_description, is_notify_trigger
	FROM weather_description_rules
	WHERE weather_code = $1
	`)

	rows := sqlmock.NewRows([]string{"weather_code", "weather_description", "is_notify_trigger"}).
		AddRow(expectedRule.WeatherCode, expectedRule.WeatherDescription, expectedRule.IsNotifyTrigger)

	mock.ExpectQuery(query).WithArgs(weatherCode).WillReturnRows(rows)

	rule, err := repo.GetRule(ctx, weatherCode)
	require.NoError(t, err)
	require.NotNil(t, rule)
	assert.Equal(t, expectedRule.WeatherCode, rule.WeatherCode)
	assert.Equal(t, expectedRule.WeatherDescription, rule.WeatherDescription)
	assert.Equal(t, expectedRule.IsNotifyTrigger, rule.IsNotifyTrigger)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRule_NotFound(t *testing.T) {
	repo, mock, cleanup := setupWeatherRuleRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	weatherCode := "999"

	query := regexp.QuoteMeta(`
	SELECT weather_code, weather_description, is_notify_trigger
	FROM weather_description_rules
	WHERE weather_code = $1
	`)
	mock.ExpectQuery(query).WithArgs(weatherCode).WillReturnError(sql.ErrNoRows)

	rule, err := repo.GetRule(ctx, weatherCode)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get weather rule")
	assert.Nil(t, rule)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRule_QueryError(t *testing.T) {
	repo, mock, cleanup := setupWeatherRuleRepoTest(t)
	defer cleanup()

	ctx := context.Background()
	weatherCode := "100"

	query := regexp.QuoteMeta(`
	SELECT weather_code, weather_description, is_notify_trigger
	FROM weather_description_rules
	WHERE weather_code = $1
	`)
	mock.ExpectQuery(query).WithArgs(weatherCode).WillReturnError(errors.New("db error"))

	rule, err := repo.GetRule(ctx, weatherCode)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get weather rule")
	assert.Nil(t, rule)

	assert.NoError(t, mock.ExpectationsWereMet())
}
