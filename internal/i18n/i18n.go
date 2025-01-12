package i18n

var messages = map[string]string{
	"greeting":               "こんにちは！ご登録ありがとうございます。",
	"askCity":                "市区町村を送信してね（例: 新宿区、名古屋市 など）",
	"searchCity":             "市区町村の検索をします",
	"invalidTimeFormat":      "時間の形式が正しくありません。例: 08:00、15:00 の形式で入力してください。",
	"updateNotifyTimeFailed": "通知時間の更新に失敗しました。",
	"settingsComplete":       "設定が完了しました！これで雨の情報をお届けします。",
	"defaultReply":           "現在の設定状況です。必要に応じて情報を更新してください。",
}

// T は指定されたキーに対応するローカライズされたメッセージを返します。
// キーが見つからない場合はキー自体を返します。
func T(key string) string {
	if msg, ok := messages[key]; ok {
		return msg
	}
	return key
}
