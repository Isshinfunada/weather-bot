package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {
	if err := run(); err != nil {
		// エラーが発生したらログ出力して終了コードを返す
		log.Printf("[ERROR] %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			return runMigrations()
		case "seed":
			return runSeedsMigrations()
		}
	}

	// 通常のアプリケーション起動処理
	return runApp()
}

func runApp() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return fmt.Errorf("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
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

	// Echoサーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
	return nil
}

func runMigrations() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return fmt.Errorf("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	err = goose.Up(db, "db/migrations")
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations applied successfully.")
	return nil
}

func runSeedsMigrations() error {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return fmt.Errorf("DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// seeds 用のバージョン管理テーブルに切り替える
	goose.SetTableName("schema_seeds")

	// db/seeds配下のSQLファイルをgooseで適用
	err = goose.Up(db, "db/seeds")
	if err != nil {
		return fmt.Errorf("failed to run seeds migrations: %w", err)
	}

	log.Println("Seeds applied successfully.")
	return nil
}
