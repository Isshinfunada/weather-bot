package usecase_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/Isshinfunada/weather-bot/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// モックの定義

type MockWeatherRuleRepo struct{ mock.Mock }

func (m *MockWeatherRuleRepo) GetRule(ctx context.Context, weatherCode string) (*entity.WeatherRule, error) {
	args := m.Called(ctx, weatherCode)
	var rule *entity.WeatherRule
	if args.Get(0) != nil {
		rule = args.Get(0).(*entity.WeatherRule)
	}
	return rule, args.Error(1)
}

type MockNotificationRepo struct{ mock.Mock }

func (m *MockNotificationRepo) InsertNotificationHistory(ctx context.Context, history *entity.NotificationHistory) error {
	args := m.Called(ctx, history)
	return args.Error(0)
}

type MockAreaUC struct{ mock.Mock }

func (m *MockAreaUC) GetHierarchy(ctx context.Context, class20ID string) (*entity.HierarchyArea, error) {
	args := m.Called(ctx, class20ID)
	var hierarchy *entity.HierarchyArea
	if args.Get(0) != nil {
		hierarchy = args.Get(0).(*entity.HierarchyArea)
	}
	return hierarchy, args.Error(1)
}

// DummyUserRepo はテストで使用しないメソッドのスタブです
type DummyUserRepo struct{}

func (d *DummyUserRepo) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return nil, nil
}
func (d *DummyUserRepo) FindUserByID(ctx context.Context, userID int) (*entity.User, error) {
	return nil, nil
}
func (d *DummyUserRepo) FindUserByLINEUserID(ctx context.Context, LINEUserID string) (*entity.User, error) {
	return nil, nil
}
func (d *DummyUserRepo) UpdateUser(ctx context.Context, user *entity.User) error { return nil }
func (d *DummyUserRepo) DeleteUser(ctx context.Context, userID int) error        { return nil }
func (d *DummyUserRepo) FindUserByNotifyTimeRange(ctx context.Context, start, end time.Time) ([]*entity.User, error) {
	return nil, nil
}

// HTTPリクエストをモックするためのRoundTripper定義
type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// MockUserRepoForRange: FindUsersByNotifyTimeRangeをモックするための構造体
type MockUserRepoForRange struct{ mock.Mock }

func (m *MockUserRepoForRange) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return nil, nil
}
func (m *MockUserRepoForRange) FindUserByID(ctx context.Context, userID int) (*entity.User, error) {
	return nil, nil
}
func (m *MockUserRepoForRange) FindUserByLINEUserID(ctx context.Context, LINEUserID string) (*entity.User, error) {
	return nil, nil
}
func (m *MockUserRepoForRange) UpdateUser(ctx context.Context, user *entity.User) error { return nil }
func (m *MockUserRepoForRange) DeleteUser(ctx context.Context, userID int) error        { return nil }
func (m *MockUserRepoForRange) FindUserByNotifyTimeRange(ctx context.Context, start, end time.Time) ([]*entity.User, error) {
	args := m.Called(ctx, start, end)
	var users []*entity.User
	if val := args.Get(0); val != nil {
		users = val.([]*entity.User)
	}
	return users, args.Error(1)
}

func TestProcessWeatherForUser(t *testing.T) {
	ctx := context.Background()

	// モックのセットアップ
	mockRuleRepo := new(MockWeatherRuleRepo)
	mockNotificationRepo := new(MockNotificationRepo)
	mockAreaUC := new(MockAreaUC)
	dummyUserRepo := &DummyUserRepo{}

	// areaUC.GetHierarchy の返却値設定
	hierarchy := &entity.HierarchyArea{
		Office:  &entity.AreaOffice{ID: "testOffice"},
		Class10: &entity.AreaClass10{ID: "testClass10"},
	}
	mockAreaUC.
		On("GetHierarchy", ctx, mock.Anything).
		Return(hierarchy, nil)

	// 特定の天気コードに対するルール設定
	mockRuleRepo.On("GetRule", ctx, "123").Return(&entity.WeatherRule{WeatherCode: "123", IsNotifyTrigger: false}, nil)
	mockRuleRepo.On("GetRule", ctx, "456").Return(&entity.WeatherRule{WeatherCode: "456", IsNotifyTrigger: true}, nil)

	mockNotificationRepo.
		On("InsertNotificationHistory", mock.Anything, mock.AnythingOfType("*entity.NotificationHistory")).
		Return(nil)

	// 対象日を決める
	targetDate := time.Now().In(utils.JST).Format("2006-01-02")

	// 偽のJSONレスポンスを作成（対象日を含むtimeDefinesを追加）
	fakeResponse := []map[string]interface{}{
		{
			"timeSeries": []interface{}{
				map[string]interface{}{
					"timeDefines": []interface{}{targetDate + "T11:00:00+09:00"},
					"areas": []interface{}{
						map[string]interface{}{
							"area": map[string]interface{}{
								"code": "testClass10",
								"name": "TestArea",
							},
							"weatherCodes": []interface{}{"123", "456"},
						},
					},
				},
			},
		},
	}
	responseBody, err := json.Marshal(fakeResponse)
	assert.NoError(t, err)

	originalTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
			Header:     make(http.Header),
		}, nil
	})
	defer func() { http.DefaultTransport = originalTransport }()

	weatherUC := usecase.NewWeatherUsecase(mockRuleRepo, mockNotificationRepo, dummyUserRepo, mockAreaUC)
	user := &entity.User{
		ID:             1,
		SelectedAreaID: "1234567",
	}

	err = weatherUC.ProcessWeatherForUser(ctx, user)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockAreaUC.AssertExpectations(t)
	mockRuleRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}

func TestProcessWeatherForUsersInTimeRange(t *testing.T) {
	ctx := context.Background()

	// モックのセットアップ
	mockRuleRepo := new(MockWeatherRuleRepo)
	mockNotificationRepo := new(MockNotificationRepo)
	mockAreaUC := new(MockAreaUC)
	mockUserRepo := new(MockUserRepoForRange)

	weatherUC := usecase.NewWeatherUsecase(mockRuleRepo, mockNotificationRepo, mockUserRepo, mockAreaUC)

	startTime := time.Date(0, 1, 1, 8, 0, 0, 0, utils.JST)
	endTime := time.Date(0, 1, 1, 9, 0, 0, 0, utils.JST)

	user := &entity.User{
		ID:             1,
		SelectedAreaID: "1234567",
	}
	users := []*entity.User{user}

	mockUserRepo.
		On("FindUserByNotifyTimeRange", ctx, startTime, endTime).
		Return(users, nil)

	hierarchy := &entity.HierarchyArea{
		Office:  &entity.AreaOffice{ID: "testOffice"},
		Class10: &entity.AreaClass10{ID: "testClass10"},
	}
	mockAreaUC.
		On("GetHierarchy", ctx, fmt.Sprint(user.SelectedAreaID)).
		Return(hierarchy, nil)

	mockRuleRepo.On("GetRule", ctx, "123").Return(&entity.WeatherRule{WeatherCode: "123", IsNotifyTrigger: false}, nil)
	mockRuleRepo.On("GetRule", ctx, "456").Return(&entity.WeatherRule{WeatherCode: "456", IsNotifyTrigger: true}, nil)

	mockNotificationRepo.
		On("InsertNotificationHistory", mock.Anything, mock.AnythingOfType("*entity.NotificationHistory")).
		Return(nil)

	// 対象日を決める
	targetDate := time.Now().In(utils.JST).Format("2006-01-02")

	// 偽のJSONレスポンスを作成（対象日を含むtimeDefinesを追加）
	fakeResponse := []map[string]interface{}{
		{
			"timeSeries": []interface{}{
				map[string]interface{}{
					"timeDefines": []interface{}{targetDate + "T11:00:00+09:00"},
					"areas": []interface{}{
						map[string]interface{}{
							"area": map[string]interface{}{
								"code": "testClass10",
								"name": "TestArea",
							},
							"weatherCodes": []interface{}{"123", "456"},
						},
					},
				},
			},
		},
	}
	responseBody, err := json.Marshal(fakeResponse)
	assert.NoError(t, err)

	originalTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
			Header:     make(http.Header),
		}, nil
	})
	defer func() { http.DefaultTransport = originalTransport }()

	err = weatherUC.ProcessWeatherForUsersInTimeRange(ctx, startTime, endTime)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	mockUserRepo.AssertExpectations(t)
	mockAreaUC.AssertExpectations(t)
	mockRuleRepo.AssertExpectations(t)
	mockNotificationRepo.AssertExpectations(t)
}
