package repository_test

import (
	"testing"

	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/stretchr/testify/require"
)

func TestNewDBConnection_Fail(t *testing.T) {
	// 存在しないホストやポート、無効な認証情報を指定して接続を試みる
	db, err := repository.NewDBConnection("127.0.0.1", "9999", "invalid_user", "invalid_pass", "invalid_db")

	// 接続に失敗することを期待する
	require.Error(t, err)
	require.Nil(t, db)
}
