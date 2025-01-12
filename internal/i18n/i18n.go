package i18n

var messages = map[string]string{
	"greeting": "こんにちは！ご登録ありがとうございます。",
	"askCity":  "市区町村を送信してね（例: 新宿区、名古屋市 など）",
}

// T は指定されたキーに対応するローカライズされたメッセージを返します。
// キーが見つからない場合はキー自体を返します。
func T(key string) string {
	if msg, ok := messages[key]; ok {
		return msg
	}
	return key
}
