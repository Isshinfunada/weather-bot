package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	usecase "github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
	linebot "github.com/line/line-bot-sdk-go/linebot"
)

type LINEWebhookController struct {
	userUC     usecase.UserUseCase
	areaUC     usecase.AreaUseCase
	messageSvc *MessageService
	bot        *linebot.Client
}

func NewLINEWebhookController(userUC usecase.UserUseCase, areaUC usecase.AreaUseCase, bot *linebot.Client) *LINEWebhookController {
	return &LINEWebhookController{
		userUC:     userUC,
		areaUC:     areaUC,
		messageSvc: NewMessageService(bot),
		bot:        bot,
	}
}

func (ctrl *LINEWebhookController) HandleWebhook(c echo.Context) error {
	events, err := ctrl.bot.ParseRequest(c.Request())
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			return c.NoContent(http.StatusBadRequest)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	ctx := c.Request().Context()
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeFollow:
			ctrl.handleFollowEvent(ctx, event)
		case linebot.EventTypeMessage:
			switch msg := event.Message.(type) {
			case *linebot.TextMessage:
				go ctrl.handleTextMessage(ctx, event, msg.Text)
			}
		}
	}

	return c.NoContent(http.StatusOK)
}

func (ctrl *LINEWebhookController) handleFollowEvent(ctx context.Context, event *linebot.Event) {
	userID := event.Source.UserID

	user, err := ctrl.userUC.GetByLINEID(ctx, userID)
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	if user == nil {
		newUser := &entity.User{
			LINEUserID: userID,
			IsActive:   true,
			Status:     "awaiting_prefecture",
		}
		createdUser, err := ctrl.userUC.Create(ctx, newUser)
		if err != nil {
			fmt.Printf("Error creating user: %v\n", err)
			return
		}
		user = createdUser
	}

	// 友だち追加時は挨拶と都道府県入力促しを QuickReply なしで送信
	ctrl.messageSvc.SendTextMessage(event.ReplyToken, "greeting")
	ctrl.messageSvc.SendTextMessage(event.ReplyToken, "askPrefecture")
}

func (ctrl *LINEWebhookController) handleTextMessage(ctx context.Context, event *linebot.Event, text string) {
	userID := event.Source.UserID

	user, err := ctrl.userUC.GetByLINEID(ctx, userID)
	if err != nil || user == nil {
		fmt.Printf("Error retrieving user: %v\n", err)
		return
	}

	// ユーザーの状態に応じて分岐
	switch user.Status {
	case "awaiting_prefecture":
		ctrl.handlePrefectureInput(ctx, event, user, text)
	case "awaiting_municipality":
		ctrl.handleMunicipalityInput(ctx, event, user, text)
	case "awaiting_confirmation":
		// ユーザーが QuickReply で「はい」または「いいえ」を送信した場合、テキストがそのまま送られる想定
		ctrl.handleConfirmation(ctx, event, user, text)
	// その他、エリア選択の各状態も後述のハンドラーで実装
	default:
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "defaultReply")
	}
}
