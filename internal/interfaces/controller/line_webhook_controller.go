package controller

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Isshinfunada/weather-bot/internal/entity"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	usecase "github.com/Isshinfunada/weather-bot/internal/usecase"
	"github.com/labstack/echo/v4"
	linebot "github.com/line/line-bot-sdk-go/linebot"
)

type LINEWebhookController struct {
	userUC     usecase.UserUseCase
	areaUC     usecase.AreaUseCase
	areaRepo   repository.AreaRepository
	messageSvc *MessageService
	bot        *linebot.Client
}

func NewLINEWebhookController(userUC usecase.UserUseCase, areaUC usecase.AreaUseCase, areaRepo repository.AreaRepository, bot *linebot.Client) *LINEWebhookController {
	return &LINEWebhookController{
		userUC:     userUC,
		areaUC:     areaUC,
		areaRepo:   areaRepo,
		messageSvc: NewMessageService(bot),
		bot:        bot,
	}
}

// RegisterLINEWebhook はLINE Webhookエンドポイントを登録します
func RegisterLINEWebhook(e *echo.Echo, userUC usecase.UserUseCase, areaUC usecase.AreaUseCase, areaRepo repository.AreaRepository) {
	// LINE Bot クライアントの初期化
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	channelToken := os.Getenv("LINE_CHANNEL_TOKEN")

	if channelSecret == "" || channelToken == "" {
		fmt.Println("LINE_CHANNEL_SECRET or LINE_CHANNEL_TOKEN is not set")
		return
	}

	bot, err := linebot.New(channelSecret, channelToken)
	if err != nil {
		fmt.Printf("Error creating LINE bot client: %v\n", err)
		return
	}

	// main.goから渡されたareaRepoを使用

	// コントローラーの初期化
	ctrl := NewLINEWebhookController(userUC, areaUC, areaRepo, bot)

	// Webhookエンドポイントの登録
	e.POST("/webhook", ctrl.HandleWebhook)
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
	case "awaiting_area_class10_selection":
		ctrl.handleAreaClass10Selection(ctx, event, user, text)
	case "awaiting_area_class15_selection":
		ctrl.handleAreaClass15Selection(ctx, event, user, text)
	case "awaiting_area_class20_selection":
		ctrl.handleAreaClass20Selection(ctx, event, user, text)
	default:
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "defaultReply")
	}
}

// 都道府県入力ハンドラー
func (ctrl *LINEWebhookController) handlePrefectureInput(ctx context.Context, event *linebot.Event, user *entity.User, text string) {
	// 都道府県名で検索
	office, err := ctrl.areaRepo.FindOfficeByName(ctx, text)
	if err != nil {
		fmt.Printf("Error finding office: %v\n", err)
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "prefectureNotFound")
		return
	}

	if office == nil {
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "prefectureNotFound")
		return
	}

	// ユーザー情報更新
	user.SelectedAreaOfficeID = office.ID
	user.Status = "awaiting_municipality"
	err = ctrl.userUC.Update(ctx, user)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	// 市区町村入力を促す
	ctrl.messageSvc.SendTextMessage(event.ReplyToken, "askMunicipality")
}

// 市区町村入力ハンドラー
func (ctrl *LINEWebhookController) handleMunicipalityInput(ctx context.Context, event *linebot.Event, user *entity.User, text string) {
	// 市区町村名で検索
	areas, err := ctrl.areaUC.SearchCityCandidates(ctx, text)
	if err != nil {
		fmt.Printf("Error searching city: %v\n", err)
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "municipalityNotFound")
		return
	}

	if len(areas) == 0 {
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "municipalityNotFound")
		return
	}

	// 最初の候補を使用
	area := areas[0]

	// ユーザー情報を一時的に更新（確認待ち）
	user.SelectedAreaID = area.Class20.ID
	user.SelectedAreaClass15ID = area.Class15.ID
	user.Status = "awaiting_confirmation"
	err = ctrl.userUC.Update(ctx, user)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	// 確認メッセージを送信（QuickReplyで「はい」「いいえ」）
	yesBtn := linebot.NewQuickReplyButton(
		"",
		linebot.NewMessageAction("はい", "はい"),
	)
	noBtn := linebot.NewQuickReplyButton(
		"",
		linebot.NewMessageAction("いいえ", "いいえ"),
	)

	buttons := []*linebot.QuickReplyButton{yesBtn, noBtn}
	// 確認メッセージを送信（地域名を引数として渡す）
	ctrl.messageSvc.SendQuickReplyMessage(event.ReplyToken, "confirmLocation", buttons, area.Class20.Name)
}

// 確認ハンドラー
func (ctrl *LINEWebhookController) handleConfirmation(ctx context.Context, event *linebot.Event, user *entity.User, text string) {
	if text == "はい" {
		// 既に選択されているarea_class20.idを確定
		user.Status = "completed"
		err := ctrl.userUC.Update(ctx, user)
		if err != nil {
			fmt.Printf("Error updating user: %v\n", err)
			return
		}

		// 完了メッセージ
		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "locationRegistered")
		return
	}

	if text == "いいえ" {
		// area_class10の選択肢を提示
		areaClass10List, err := ctrl.areaRepo.FindAreaClass10ByOfficeID(ctx, user.SelectedAreaOfficeID)
		if err != nil {
			fmt.Printf("Error finding area class10: %v\n", err)
			return
		}

		// QuickReplyボタンを作成
		var buttons []*linebot.QuickReplyButton
		for _, area := range areaClass10List {
			btn := linebot.NewQuickReplyButton(
				"",
				linebot.NewMessageAction(area.Name, area.ID),
			)
			buttons = append(buttons, btn)
		}

		// 「この中にはない」オプション
		notInListBtn := linebot.NewQuickReplyButton(
			"",
			linebot.NewMessageAction("この中にはない", "この中にはない"),
		)
		buttons = append(buttons, notInListBtn)

		// ユーザー状態を更新
		user.Status = "awaiting_area_class10_selection"
		err = ctrl.userUC.Update(ctx, user)
		if err != nil {
			fmt.Printf("Error updating user: %v\n", err)
			return
		}

		// 選択肢を送信
		ctrl.messageSvc.SendQuickReplyMessage(event.ReplyToken, "selectAreaClass10", buttons)
	}
}

// エリア選択ハンドラー（Class10）
func (ctrl *LINEWebhookController) handleAreaClass10Selection(ctx context.Context, event *linebot.Event, user *entity.User, text string) {
	if text == "この中にはない" {
		// 代替処理（例：都道府県選択に戻る）
		user.Status = "awaiting_prefecture"
		err := ctrl.userUC.Update(ctx, user)
		if err != nil {
			fmt.Printf("Error updating user: %v\n", err)
			return
		}

		ctrl.messageSvc.SendTextMessage(event.ReplyToken, "askPrefecture")
		return
	}

	// area_class10.idとして処理
	areaClass15List, err := ctrl.areaRepo.FindAreaClass15ByClass10ID(ctx, text)
	if err != nil {
		fmt.Printf("Error finding area class15: %v\n", err)
		return
	}

	// QuickReplyボタンを作成
	var buttons []*linebot.QuickReplyButton
	for _, area := range areaClass15List {
		btn := linebot.NewQuickReplyButton(
			"",
			linebot.NewMessageAction(area.Name, area.ID),
		)
		buttons = append(buttons, btn)
	}

	// ユーザー状態を更新
	user.Status = "awaiting_area_class15_selection"
	err = ctrl.userUC.Update(ctx, user)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	// 選択肢を送信
	ctrl.messageSvc.SendQuickReplyMessage(event.ReplyToken, "selectAreaClass15", buttons)
}

// エリア選択ハンドラー（Class15）
func (ctrl *LINEWebhookController) handleAreaClass15Selection(ctx context.Context, event *linebot.Event, user *entity.User, text string) {
	// 選択されたarea_class15.idに基づいてarea_class20一覧を取得
	areaClass20List, err := ctrl.areaRepo.FindAreaClass20ByClass15ID(ctx, text)
	if err != nil {
		fmt.Printf("Error finding area class20: %v\n", err)
		return
	}

	// ユーザー状態を更新
	user.SelectedAreaClass15ID = text
	user.Status = "awaiting_area_class20_selection"
	err = ctrl.userUC.Update(ctx, user)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	// QuickReplyボタンを作成
	var buttons []*linebot.QuickReplyButton
	for _, area := range areaClass20List {
		btn := linebot.NewQuickReplyButton(
			"",
			linebot.NewMessageAction(area.Name, area.ID),
		)
		buttons = append(buttons, btn)
	}

	// 選択肢を送信
	ctrl.messageSvc.SendQuickReplyMessage(event.ReplyToken, "selectAreaClass20", buttons)
}

// エリア選択ハンドラー（Class20）
func (ctrl *LINEWebhookController) handleAreaClass20Selection(ctx context.Context, event *linebot.Event, user *entity.User, text string) {
	// 選択されたarea_class20.idを保存して完了

	// ユーザー情報を更新
	user.SelectedAreaID = text
	user.Status = "completed"
	err := ctrl.userUC.Update(ctx, user)
	if err != nil {
		fmt.Printf("Error updating user: %v\n", err)
		return
	}

	// 完了メッセージ
	ctrl.messageSvc.SendTextMessage(event.ReplyToken, "locationRegistered")
}
