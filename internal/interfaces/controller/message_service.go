package controller

import (
	"github.com/Isshinfunada/weather-bot/internal/i18n"
	"github.com/line/line-bot-sdk-go/linebot"
)

type MessageService struct {
	bot *linebot.Client
}

func NewMessageService(bot *linebot.Client) *MessageService {
	return &MessageService{bot: bot}
}

// 通常のテキストメッセージ送信
func (ms *MessageService) SendTextMessage(replyToken, key string) error {
	text := i18n.T(key)
	message := linebot.NewTextMessage(text)
	_, err := ms.bot.ReplyMessage(replyToken, message).Do()
	return err
}

// QuickReplyを使ったテキストメッセージ送信
func (ms *MessageService) SendQuickReplyMessage(replyToken, key string, actions []*linebot.QuickReplyButton, args ...interface{}) error {
	text := i18n.T(key)
	message := linebot.NewTextMessage(text)
	message.WithQuickReplies(linebot.NewQuickReplyItems(actions...))

	_, err := ms.bot.ReplyMessage(replyToken, message).Do()
	return err
}
