package utils

import "unicode"

// IsKanji は与えられた文字列に漢字が含まれているかを判定します。
func IsKanji(s string) bool {
	for _, r := range s {
		if unicode.In(r, unicode.Han) {
			return true
		}
	}
	return false
}
