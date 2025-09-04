package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"go-shorturl/pkg/db"
	"go-shorturl/pkg/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// isValidURL 驗證 URL 格式
func isValidURL(rawURL string) bool {
	// 如果沒有協議，自動添加 https://
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	
	parsedURL, err := url.Parse(rawURL)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

// normalizeURL 標準化 URL
func normalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "https://" + rawURL
	}
	return rawURL
}

// generateShortCode 產生隨機短碼
func generateShortCode(originalURL string) (string, error) {
	// 計算原始網址長度（不包含協議）
	urlWithoutProtocol := originalURL
	if strings.HasPrefix(urlWithoutProtocol, "https://") {
		urlWithoutProtocol = strings.TrimPrefix(urlWithoutProtocol, "https://")
	} else if strings.HasPrefix(urlWithoutProtocol, "http://") {
		urlWithoutProtocol = strings.TrimPrefix(urlWithoutProtocol, "http://")
	}
	
	originalLength := len(urlWithoutProtocol)
	
	// 短碼長度應該是原始長度的一半，但至少6個字符，最多12個字符
	shortCodeLength := originalLength / 2
	if shortCodeLength < 6 {
		shortCodeLength = 6
	} else if shortCodeLength > 12 {
		shortCodeLength = 12
	}
	
	// 確保短碼比原始網址短
	if shortCodeLength >= originalLength {
		shortCodeLength = originalLength - 1
		if shortCodeLength < 6 {
			shortCodeLength = 6
		}
	}
	
	bytes := make([]byte, shortCodeLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:shortCodeLength], nil
}

// ShortenURL 建立短網址
func ShortenURL(c *fiber.Ctx) error {
	var req models.ShortenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 標準化 URL
	normalizedURL := normalizeURL(req.URL)
	
	// 驗證 URL
	if !isValidURL(normalizedURL) {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid URL format",
		})
	}

	// 決定短碼
	var shortCode string
	var err error

	if req.CustomCode != "" {
		// 檢查自訂短碼是否已存在
		var exists bool
		query := "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)"
		err = db.GetDB().QueryRow(context.Background(), query, req.CustomCode).Scan(&exists)
		if err != nil {
			log.Printf("Database error checking custom code: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Database error: %v", err),
			})
		}
		if exists {
			return c.Status(409).JSON(fiber.Map{
				"error": "Custom code already exists",
			})
		}
		shortCode = req.CustomCode
	} else {
		// 產生隨機短碼
		for {
			shortCode, err = generateShortCode(normalizedURL)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Failed to generate short code",
				})
			}

			// 檢查短碼是否已存在
			var exists bool
			query := "SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = $1)"
			err = db.GetDB().QueryRow(context.Background(), query, shortCode).Scan(&exists)
			if err != nil {
				log.Printf("Database error checking short code: %v", err)
				return c.Status(500).JSON(fiber.Map{
					"error": fmt.Sprintf("Database error: %v", err),
				})
			}
			if !exists {
				break
			}
		}
	}

	// 插入新記錄
	query := `
		INSERT INTO urls (id, original_url, short_code, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	
	id := uuid.New()
	createdAt := time.Now()
	
	err = db.GetDB().QueryRow(context.Background(), query, id, normalizedURL, shortCode, createdAt).Scan(&id, &createdAt)
	if err != nil {
		log.Printf("Error inserting URL: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to create short URL: %v", err),
		})
	}

	// 構建短網址
	baseURL := "http://localhost:8080"
	
	// 嘗試從環境變數獲取域名
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		baseURL = envBaseURL
	} else {
		// 嘗試從請求頭獲取
		if host := c.Get("Host"); host != "" {
			protocol := "http"
			if c.Get("X-Forwarded-Proto") == "https" {
				protocol = "https"
			}
			baseURL = fmt.Sprintf("%s://%s", protocol, host)
		}
	}
	
	// 統一使用 /url/ 格式
	shortURL := fmt.Sprintf("%s/url/%s", baseURL, shortCode)

	response := models.ShortenResponse{
		ShortURL:    shortURL,
		OriginalURL: normalizedURL,
		ShortCode:   shortCode,
		CreatedAt:   createdAt,
	}

	return c.Status(201).JSON(response)
}

// RedirectURL 重定向到原始網址
func RedirectURL(c *fiber.Ctx) error {
	shortCode := c.Params("short_code")
	if shortCode == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// 查詢原始網址
	query := "SELECT id, original_url FROM urls WHERE short_code = $1"
	var urlID uuid.UUID
	var originalURL string
	
	err := db.GetDB().QueryRow(context.Background(), query, shortCode).Scan(&urlID, &originalURL)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error": "Short URL not found",
			})
		}
		log.Printf("Error querying URL: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// 記錄點擊
	clickQuery := `
		INSERT INTO clicks (id, url_id, clicked_at, ip_address, user_agent, referrer)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	clickID := uuid.New()
	clickedAt := time.Now()
	ipAddress := c.IP()
	userAgent := string(c.Request().Header.UserAgent())
	referrer := string(c.Request().Header.Referer())
	
	_, err = db.GetDB().Exec(context.Background(), clickQuery, clickID, urlID, clickedAt, ipAddress, userAgent, referrer)
	if err != nil {
		log.Printf("Error recording click: %v", err)
		// 不返回錯誤，因為重定向仍然應該工作
	}

	return c.Redirect(originalURL, 302)
}

// GetStats 取得短網址統計
func GetStats(c *fiber.Ctx) error {
	shortCode := c.Params("short_code")
	if shortCode == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Short code is required",
		})
	}

	// 查詢短網址基本資訊
	urlQuery := "SELECT id, original_url, created_at FROM urls WHERE short_code = $1"
	var urlID uuid.UUID
	var originalURL string
	var createdAt time.Time
	
	err := db.GetDB().QueryRow(context.Background(), urlQuery, shortCode).Scan(&urlID, &originalURL, &createdAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(404).JSON(fiber.Map{
				"error": "Short URL not found",
			})
		}
		log.Printf("Error querying URL: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// 查詢總點擊數
	var totalClicks int
	clickCountQuery := "SELECT COUNT(*) FROM clicks WHERE url_id = $1"
	err = db.GetDB().QueryRow(context.Background(), clickCountQuery, urlID).Scan(&totalClicks)
	if err != nil {
		log.Printf("Error counting clicks: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// 查詢裝置統計
	deviceQuery := `
		SELECT user_agent, COUNT(*) as count
		FROM clicks
		WHERE url_id = $1 AND user_agent IS NOT NULL AND user_agent != ''
		GROUP BY user_agent
		ORDER BY count DESC
		LIMIT 10
	`
	
	deviceRows, err := db.GetDB().Query(context.Background(), deviceQuery, urlID)
	if err != nil {
		log.Printf("Error querying device stats: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer deviceRows.Close()

	var deviceStats []models.DeviceStat
	for deviceRows.Next() {
		var stat models.DeviceStat
		err := deviceRows.Scan(&stat.UserAgent, &stat.Count)
		if err != nil {
			log.Printf("Error scanning device stat: %v", err)
			continue
		}
		deviceStats = append(deviceStats, stat)
	}

	// 查詢來源統計
	referrerQuery := `
		SELECT COALESCE(referrer, 'Direct') as referrer, COUNT(*) as count
		FROM clicks
		WHERE url_id = $1
		GROUP BY COALESCE(referrer, 'Direct')
		ORDER BY count DESC
		LIMIT 10
	`
	
	referrerRows, err := db.GetDB().Query(context.Background(), referrerQuery, urlID)
	if err != nil {
		log.Printf("Error querying referrer stats: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer referrerRows.Close()

	var referrerStats []models.ReferrerStat
	for referrerRows.Next() {
		var stat models.ReferrerStat
		err := referrerRows.Scan(&stat.Referrer, &stat.Count)
		if err != nil {
			log.Printf("Error scanning referrer stat: %v", err)
			continue
		}
		referrerStats = append(referrerStats, stat)
	}

	response := models.StatsResponse{
		ShortCode:     shortCode,
		OriginalURL:   originalURL,
		TotalClicks:   totalClicks,
		CreatedAt:     createdAt,
		DeviceStats:   deviceStats,
		ReferrerStats: referrerStats,
	}

	return c.JSON(response)
}
