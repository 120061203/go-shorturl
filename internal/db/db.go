package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

// InitDB 初始化資料庫連線
func InitDB() error {
	// 載入環境變數
	if err := godotenv.Load("env.template"); err != nil {
		log.Println("Warning: Could not load env.template file")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	// 建立資料庫連線池
	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 測試連線
	if err := DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

// CloseDB 關閉資料庫連線
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// GetDB 取得資料庫實例
func GetDB() *pgxpool.Pool {
	return DB
}
