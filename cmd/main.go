package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		runMigrations()
		return
	}
	// 通常のアプリケーション起動処理
	runApp()
}

func runApp() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)

	// Echoサーバーの設定
	e := echo.New()

	// ルートエンドポイントでHello, World!を返す
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

	// Webhook エンドポイントの追加（今後の実装用）
	e.POST("/webhook", func(c echo.Context) error {
		// Webhook 処理ロジックをここに実装
		return c.JSON(http.StatusOK, map[string]string{"message": "Webhook received"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func runMigrations() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	err = goose.Up(db, "migrations")
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully.")
}
