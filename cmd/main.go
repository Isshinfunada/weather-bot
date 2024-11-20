package main

import (
    "database/sql"
    "log"
    "net/http"
    "os"

    "github.com/pressly/goose/v3"
    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"
)

func main() {
    // データベース接続情報
    dbURL := os.Getenv("DB_URL")
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

    // Echoサーバーの設定
    e := echo.New()

    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })

    e.Logger.Fatal(e.Start(":8080"))
}
