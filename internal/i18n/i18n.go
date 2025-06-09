package i18n

var messages = map[string]string{
	"greeting":               "こんにちは！ご登録ありがとうございます。",
	"askCity":                "市区町村を送信してね（例: 新宿区、名古屋市 など）",
	"searchCity":             "市区町村の検索をします",
	"invalidTimeFormat":      "時間の形式が正しくありません。例: 08:00、15:00 の形式で入力してください。",
	"updateNotifyTimeFailed": "通知時間の更新に失敗しました。",
	"settingsComplete":       "設定が完了しました！これで雨の情報をお届けします。",
	"defaultReply":           "現在の設定状況です。必要に応じて情報を更新してください。",
	"kanjiValidationError":   "市区町村名は漢字で入力してください。",
	"altanativesNotFound":    "代替候補が見つかりませんでした。別の市区町村名を入力してください。",
	// 地域選択フロー用のメッセージ
	"askPrefecture":        "都道府県名を入力してください（例: 東京都、大阪府 など）",
	"prefectureNotFound":   "都道府県が見つかりませんでした。正しい都道府県名を入力してください。",
	"askMunicipality":      "市区町村名を入力してください（例: 新宿区、名古屋市 など）",
	"municipalityNotFound": "市区町村が見つかりませんでした。正しい市区町村名を入力してください。",
	"confirmLocation":      "%s でよろしいですか？",
	"selectAreaClass10":    "より詳細な地域を選択してください：",
	"selectAreaClass15":    "より詳細な地域を選択してください：",
	"selectAreaClass20":    "より詳細な地域を選択してください：",
	"locationRegistered":   "地域の登録が完了しました！",
	"notInList":            "この中にはない",
}

// T は指定されたキーに対応するローカライズされたメッセージを返します。
// キーが見つからない場合はキー自体を返します。
func T(key string) string {
	if msg, ok := messages[key]; ok {
		return msg
	}
	return key
}
