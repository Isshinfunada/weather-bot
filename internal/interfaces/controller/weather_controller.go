package controller

import (
	"net/http"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/Isshinfunada/weather-bot/internal/utils"
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
		// "HH:MM(08:00)を想定"
		// 現在のJST日時を基準にする
		now := time.Now().In(utils.JST)

		// "15:04"フォーマットで時間をパース
		parsedStart, err := time.ParseInLocation("1504", startParam, utils.JST)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid start time format"})
		}

		parsedEnd, err := time.ParseInLocation("1504", endParam, utils.JST)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid end time format"})
		}

		// 今日の日付にパースした時刻を想定する
		start = time.Date(now.Year(), now.Month(), now.Day(), parsedStart.Hour(), parsedStart.Minute(), 0, 0, utils.JST)
		end = time.Date(now.Year(), now.Month(), now.Day(), parsedEnd.Hour(), parsedEnd.Minute(), 0, 0, utils.JST)
	} else {
		// クエリパラメータが無い場合、デフォルトで現在時刻の1時間前から現在時刻まで
		end = time.Now().In(utils.JST)
		start = end.Add(-1 * time.Hour)
	}

	// usecaseを呼び出して指定時間帯の処理を実行
	err = ctrl.weatherUC.ProcessWeatherForUsersInTimeRange(c.Request().Context(), start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Weather processing completed"})
}
