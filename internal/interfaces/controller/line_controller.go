package controller

import (
	"net/http"
	"os"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/i18n"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/line/line-bot-sdk-go/linebot"
)

func RegisterLINEWebhook(e *echo.Echo, userUC usecase.UserUsecase) {
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
				handleMessageEvent(c, lineBot, event, userUC)
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

func handleMessageEvent(c echo.Context, bot *linebot.Client, event *linebot.Event, userUC usecase.UserUsecase) {
	if event.Type != linebot.EventTypeMessage {
		return
	}

	switch msg := event.Message.(type) {
	case *linebot.TextMessage:
		text := msg.Text
		processTextMessage(c, bot, event, text, userUC)
	default:
		c.Logger().Infof("Unhandled message type: %T", event.Message)
	}
}

// 現在のところ、processTextMessage は未実装のプレースホルダです。
// 将来的に市区町村の受信・確認や通知時間の受信処理をここに実装します。
func processTextMessage(c echo.Context, bot *linebot.Client, event *linebot.Event, text string, userUC usecase.UserUsecase) {
	// メッセージ内容に応じた処理をここで実装
	// 例: "登録"コマンド以外のメッセージに対する処理など
	c.Logger().Infof("Received text message: %s", text)
}
