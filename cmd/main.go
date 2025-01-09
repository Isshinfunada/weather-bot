package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Isshinfunada/weather-bot/internal/interfaces/controller"
	"github.com/Isshinfunada/weather-bot/internal/interfaces/repository"
	"github.com/Isshinfunada/weather-bot/internal/usecase"
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

	userRepo := repository.NewUserRepository(db)
	areaRepo := repository.NewAreaRepository(db)
	weatherRuleRepo := repository.NewWeatherRuleRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	areaUC := usecase.NewAreaUseCase(areaRepo)
	userUC := usecase.NewUserUseCase(userRepo)
	weatherUC := usecase.NewWeatherUsecase(weatherRuleRepo, notificationRepo, userRepo, areaUC)

	// Echoサーバーの設定
	e := echo.New()

	// ルートエンドポイントでHello, World!を返す
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	controller.RegisterRoutes(e, userUC, areaUC, weatherUC)

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
