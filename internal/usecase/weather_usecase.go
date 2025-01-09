package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
)

type WeatherUsecase interface {
	ProcessWeatherForUser(ctx context.Context, user *entity.User) error
	ProcessWeatherForUsersInTimeRange(ctx context.Context, start, end time.Time) error
}

type weatherUsecase struct {
	weatherRuleRepo  repository.WeatherRuleRepository
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	areaUC           AreaUseCase
}

func NewWeatherUsecase(wr repository.WeatherRuleRepository, nr repository.NotificationRepository, ur repository.UserRepository, auc AreaUseCase) WeatherUsecase {
	return &weatherUsecase{
		weatherRuleRepo:  wr,
		notificationRepo: nr,
		userRepo:         ur,
		areaUC:           auc,
	}
}

func (u *weatherUsecase) ProcessWeatherForUser(ctx context.Context, user *entity.User) error {
	// ユーザーの選択エリアから改装情報を取得
	hierarchy, err := u.areaUC.GetHierarchy(ctx, fmt.Sprint(user.SelectedAreaID))
	if err != nil {
		return fmt.Errorf("failed to get hierarchy for user %d: %w", user.ID, err)
	}
	if hierarchy == nil {
		return fmt.Errorf("no hierarchy found %s for user %d", user.SelectedAreaID, user.ID)
	}

	areaOfficeID := hierarchy.Office.ID
	class10ID := hierarchy.Class10

	// JMAエンドポイントから天気データを取得
	url := fmt.Sprintf("https://www.jma.go.jp/bosai/forecast/data/forecast/%s.json", areaOfficeID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to fetch weather data: %w", err)
	}

	// JSONレスポンスをパースし、対象エリアの天気コードを抽出
	var data []map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	var weatherCodes []string
	for _, forecast := range data {
		timeSeries, ok := forecast["timeSeries"].([]interface{})
		if !ok {
			continue
		}
		for _, ts := range timeSeries {
			tsMap, ok := ts.(map[string]interface{})
			if !ok {
				continue
			}
			areas, ok := tsMap["areas"].([]interface{})
			if !ok {
				continue
			}
			for _, area := range areas {
				areaMap, ok := area.(map[string]interface{})
				if !ok {
					continue
				}
				a, ok := areaMap["area"].(map[string]interface{})
				if !ok {
					continue
				}
				code, ok := a["code"].(string)
				if !ok {
					continue
				}
				if code == class10ID.ID {
					if wc, ok := areaMap["weatherCodes"].([]interface{}); ok {
						for _, w := range wc {
							if ws, ok := w.(string); ok {
								weatherCodes = append(weatherCodes, ws)
							}
						}
					}
				}
			}
		}
	}

	// 天気コードに基づき通知トリガー設定
	notify := false
	for _, code := range weatherCodes {
		rule, err := u.weatherRuleRepo.GetRule(ctx, code)
		if err != nil {
			// エラー発生時スキップ
			continue
		}
		if rule.IsNotifyTrigger {
			notify = true
			break
		}
	}

	// notification_historyに記載
	history := &entity.NotificationHistory{
		UserID:           user.ID,
		NotificationTime: time.Now(),
		WeatherData:      body,
		IsNotifyTrigger:  notify,
	}

	if err := u.notificationRepo.InsertNotificationHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to save notification history: %w", err)
	}

	// コンソール出力
	if notify {
		fmt.Printf("User %d: 通知を送信します。天気コード: %v\n", user.ID, weatherCodes)
	} else {
		fmt.Printf("User %d: 通知不要", user.ID)
	}

	return nil
}

func (u *weatherUsecase) ProcessWeatherForUsersInTimeRange(ctx context.Context, start, end time.Time) error {
	// 指定時間帯のユーザーを取得
	users, err := u.userRepo.FindUserByNotifyTimeRange(ctx, start, end)
	if err != nil {
		return fmt.Errorf("failed to find users by notify time range: %w", err)
	}

	// 各ユーザーに対して天気情報処理実行
	for _, user := range users {
		if err := u.ProcessWeatherForUser(ctx, user); err != nil {
			fmt.Printf("Error processing weather for user %d: %v\n", user.ID, err)
		}
	}
	return nil
}
