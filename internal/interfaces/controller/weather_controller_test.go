package controller_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/controller"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Isshinfunada/weather-bot/internal/utils"
)

// モック WeatherUsecase を定義
type MockWeatherUsecase struct {
	mock.Mock
}

func (m *MockWeatherUsecase) ProcessWeatherForUsersInTimeRange(ctx context.Context, start, end time.Time) error {
	args := m.Called(ctx, start, end)
	return args.Error(0)
}

func (m *MockWeatherUsecase) ProcessWeatherForUser(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// テスト対象のコントローラーを初期化する関数
func setupWeatherController() (*controller.WeatherController, *MockWeatherUsecase, echo.Context, *httptest.ResponseRecorder) {
	utils.JST = time.FixedZone("JST", 9*60*60)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/process_weather", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	mockWUC := new(MockWeatherUsecase)
	weatherCtrl := controller.NewWeatherController(mockWUC)
	return weatherCtrl, mockWUC, ctx, rec
}

func TestProcessWeather_WithValidQueryParams(t *testing.T) {
	weatherCtrl, mockWUC, ctx, rec := setupWeatherController()

	// クエリパラメータを設定
	q := ctx.QueryParams()
	q.Add("start", "0800")
	q.Add("end", "0900")
	ctx.SetRequest(ctx.Request().WithContext(context.Background()))
	ctx.Request().URL.RawQuery = q.Encode()

	// 現在の日付を基にした開始・終了時刻を計算
	now := time.Now().In(utils.JST)
	expectedStartTime, _ := time.ParseInLocation("1504", "0800", utils.JST)
	expectedEndTime, _ := time.ParseInLocation("1504", "0900", utils.JST)
	expectedStart := time.Date(now.Year(), now.Month(), now.Day(), expectedStartTime.Hour(), expectedStartTime.Minute(), 0, 0, utils.JST)
	expectedEnd := time.Date(now.Year(), now.Month(), now.Day(), expectedEndTime.Hour(), expectedEndTime.Minute(), 0, 0, utils.JST)

	// モックの挙動を設定
	mockWUC.On("ProcessWeatherForUsersInTimeRange", mock.Anything, expectedStart, expectedEnd).Return(nil)

	// エンドポイント呼び出し
	if assert.NoError(t, weatherCtrl.ProcessWeather(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Weather processing completed", resp["message"])
	}

	mockWUC.AssertExpectations(t)
}

func TestProcessWeather_WithInvalidStartFormat(t *testing.T) {
	weatherCtrl, _, ctx, rec := setupWeatherController()

	// 不正な start パラメータを設定
	q := ctx.QueryParams()
	q.Add("start", "invalid")
	q.Add("end", "0900")
	ctx.Request().URL.RawQuery = q.Encode()

	err := weatherCtrl.ProcessWeather(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "invalid start time format", resp["error"])
}

func TestProcessWeather_WithInvalidEndFormat(t *testing.T) {
	weatherCtrl, _, ctx, rec := setupWeatherController()

	// 不正な end パラメータを設定
	q := ctx.QueryParams()
	q.Add("start", "0800")
	q.Add("end", "invalid")
	ctx.Request().URL.RawQuery = q.Encode()

	err := weatherCtrl.ProcessWeather(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "invalid end time format", resp["error"])
}

func TestProcessWeather_NoQueryParams(t *testing.T) {
	weatherCtrl, mockWUC, ctx, rec := setupWeatherController()

	// クエリパラメータなし：デフォルトの時間範囲を使用するケース
	// モックが受け取る引数の具体的な開始・終了時刻は動的になるため、anyTimesやArgument matcherを使用

	mockWUC.On("ProcessWeatherForUsersInTimeRange", mock.Anything, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(nil)

	// エンドポイント呼び出し
	if assert.NoError(t, weatherCtrl.ProcessWeather(ctx)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "Weather processing completed", resp["message"])
	}

	mockWUC.AssertExpectations(t)
}
