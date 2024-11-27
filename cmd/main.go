package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    "github.com/pressly/goose/v3"
    "github.com/labstack/echo/v4"
    "github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
    _ "github.com/lib/pq"
)

func main() {
    // 環境変数からDB_URLを取得
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL is not set")
    }

    // データベース接続
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // マイグレーションの実行
    err = goose.Up(db, "./migrations")
    if err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    // リポジトリの初期化
    userRepo := repository.NewUserRepository(db)

    // Echoサーバーの設定
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })

    // ユーザー登録エンドポイント
    e.POST("/register", func(c echo.Context) error {
        type RegisterRequest struct {
            LineUserID string `json:"line_user_id"`
        }

        var req RegisterRequest
        if err := c.Bind(&req); err != nil {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
        }

        if req.LineUserID == "" {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "line_user_id is required"})
        }

        user, err := userRepo.CreateUser(req.LineUserID)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
        }

        return c.JSON(http.StatusOK, map[string]interface{}{
            "id":           user.ID,
            "line_user_id": user.LineUserID,
        })
    })

    e.Logger.Fatal(e.Start(":8080"))
}
