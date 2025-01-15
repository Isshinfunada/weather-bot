package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/i18n"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

func RegisterLINEWebhook(e *echo.Echo, userUC usecase.UserUsecase, areaUC usecase.AreaUseCase) {
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	lineBot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		e.Logger.Fatalf("LINE Bot client initialization error: %v", err)
	}

	e.POST("line/webhook", func(c echo.Context) error {
		req := c.Request()

		events, err := lineBot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid Signature"})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeFollow:
				handleFollowEvent(c, lineBot, event, userUC)
			case linebot.EventTypeMessage:
				handleMessageEvent(c, lineBot, event, userUC, areaUC)
			default:
				c.Logger().Infof("Unhandled event type: %v", event.Type)
			}
		}
		return c.NoContent(http.StatusOK)
	})
}

func handleFollowEvent(c echo.Context, bot *linebot.Client, event *linebot.Event, userUC usecase.UserUsecase) {
	lineUserID := event.Source.UserID

	newUser := &entity.User{
		LINEUserID:     lineUserID,
		SelectedAreaID: "",
		IsActive:       true,
	}
	ctx := c.Request().Context()
	_, err := userUC.Create(ctx, newUser)
	if err != nil {
		c.Logger().Errorf("User creation failed: %v", err)
	}

	replyMessages := []linebot.SendingMessage{
		linebot.NewTextMessage(i18n.T("greeting")),
		linebot.NewTextMessage(i18n.T("askCity")),
	}
	if _, err := bot.ReplyMessage(event.ReplyToken, replyMessages...).Do(); err != nil {
		c.Logger().Errorf("Reply error: %v", err)
	}
}

func handleMessageEvent(c echo.Context, bot *linebot.Client, event *linebot.Event, userUC usecase.UserUsecase, areaUC usecase.AreaUseCase) {
	if event.Type != linebot.EventTypeMessage {
		return
	}

	switch msg := event.Message.(type) {
	case *linebot.TextMessage:
		text := msg.Text
		processTextMessage(c, bot, event, text, userUC, areaUC)
	default:
		c.Logger().Infof("Unhandled message type: %T", event.Message)
	}
}

// 市区町村の受信・確認や通知時間の受信処理をここに実装します。
func processTextMessage(c echo.Context, bot *linebot.Client, event *linebot.Event, text string, userUC usecase.UserUsecase, areaUC usecase.AreaUseCase) {
	c.Logger().Infof("Received text message: %s", text)

	ctx := c.Request().Context()
	user, err := userUC.GetByLINEID(ctx, event.Source.UserID)
	if err != nil {
		c.Logger().Errorf("Failed to get user: %v", err)
		return
	}
	if user == nil {
		c.Logger().Warn("User not found")
		return
	}

	// TODO:地域を特定してselected_area_idをupdateする処理追加

	// 通知時間未設定の場合
	if user.NotifyTime.IsZero() {
		parsedTime, err := time.Parse("15:04", text)
		if err != nil {
			replyText := i18n.T("invalidTimeFormat")
			bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
			return
		}
		user.NotifyTime = parsedTime
		if err := userUC.Update(ctx, user); err != nil {
			c.Logger().Errorf("Failed to update notify time: %v", err)
			replyText := i18n.T("updateNotifyTimeFailed")
			bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
			return
		}
		replyText := i18n.T("settingsComplete")
		bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
		return
	}

	defaultReply := i18n.T("defaultReply")
	bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(defaultReply)).Do()
}
