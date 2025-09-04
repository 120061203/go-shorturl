package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

// InitDB 初始化資料庫連線
func InitDB() error {
	// 調試信息
	log.Printf("Environment: VERCEL=%s", os.Getenv("VERCEL"))
	
	// 檢查是否在 Vercel 環境中（多種檢測方法）
	isVercel := os.Getenv("VERCEL") != "" || 
		os.Getenv("VERCEL_ENV") != "" || 
		os.Getenv("VERCEL_URL") != "" ||
		strings.Contains(os.Getenv("VERCEL_REGION") || "", "iad") ||
		strings.Contains(os.Getenv("VERCEL_REGION") || "", "hkg")
	
	log.Printf("Is Vercel environment: %v", isVercel)
	
	// 只在本地開發環境載入 .env 文件
	if !isVercel {
		if err := godotenv.Load(); err != nil {
			log.Println("Warning: Could not load .env file")
		}
	}

	databaseURL := os.Getenv("DATABASE_URL")
	log.Printf("DATABASE_URL length: %d", len(databaseURL))
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
