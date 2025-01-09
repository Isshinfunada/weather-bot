package controller

import (
	"net/http"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
)

type WeatherController struct {
	weatherUC usecase.WeatherUsecase
}

func NewWeatherController(wuc usecase.WeatherUsecase) *WeatherController {
	return &WeatherController{weatherUC: wuc}
}

// GET /api/process_weather
func (ctrl *WeatherController) ProcessWeather(c echo.Context) error {
	var start, end time.Time
	var err error
	// クエリパラメータによる時間の指定があれば解析、そうじゃなければ直近1時間を設定

	startParam := c.QueryParam("start")
	endParam := c.QueryParam("end")

	if startParam != "" && endParam != "" {
		// タイムフォーマットをRFC3339としてパース（例："2025-01-09T08:00:00Z"）
		start, err = time.Parse(time.RFC3339, startParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid start time format"})
		}
		end, err = time.Parse(time.RFC3339, endParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid end time format"})
		}
	} else {
		// クエリパラメータが無い場合、デフォルトで現在時刻の1時間前から現在時刻まで
		end = time.Now()
		start = end.Add(-1 * time.Hour)
	}

	// usecaseを呼び出して指定時間帯の処理を実行
	err = ctrl.weatherUC.ProcessWeatherForUsersInTimeRange(c.Request().Context(), start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Weather processing completed"})
}
