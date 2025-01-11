package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Isshinfunada/weather-bot/internal/entity"
)

type WeatherRuleRepository interface {
	GetRule(ctx context.Context, weatherCode string) (*entity.WeatherRule, error)
}

type weatherRuleRepository struct {
	db *sql.DB
}

func NewWeatherRuleRepository(db *sql.DB) WeatherRuleRepository {
	return &weatherRuleRepository{db: db}
}

func (r *weatherRuleRepository) GetRule(ctx context.Context, weatherCode string) (*entity.WeatherRule, error) {
	query := `
	SELECT weather_code, weather_description, is_notify_trigger
	FROM weather_notification_rules
	WHERE weather_code = $1
	`

	var rule entity.WeatherRule
	err := r.db.QueryRowContext(ctx, query, weatherCode).Scan(
		&rule.WeatherCode, &rule.WeatherDescription, &rule.IsNotifyTrigger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather rule: %w", err)
	}
	return &rule, nil
}
