package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// getRealIP 從 HTTP 頭中獲取真實 IP 地址
func getRealIP(c *fiber.Ctx) string {
	// 優先從 X-Forwarded-For 獲取（可能包含多個IP，取第一個）
	if xff := c.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			realIP := strings.TrimSpace(ips[0])
			if realIP != "" {
				return realIP
			}
		}
	}

	// 從 X-Real-IP 獲取
	if xri := c.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 從 X-Client-IP 獲取
	if xci := c.Get("X-Client-IP"); xci != "" {
		return xci
	}

	// 從 CF-Connecting-IP 獲取（Cloudflare）
	if cfip := c.Get("CF-Connecting-IP"); cfip != "" {
		return cfip
	}

	// 最後使用 Fiber 的 IP() 方法
	return c.IP()
}

// getRealUserAgent 從 HTTP 頭中獲取真實 User-Agent
func getRealUserAgent(c *fiber.Ctx) string {
	// 優先從 X-Forwarded-User-Agent 獲取
	if xfua := c.Get("X-Forwarded-User-Agent"); xfua != "" {
		return xfua
	}

	// 從 X-User-Agent 獲取
	if xua := c.Get("X-User-Agent"); xua != "" {
		return xua
	}

	// 最後使用標準 User-Agent
	ua := string(c.Request().Header.UserAgent())
	if ua == "" {
		return "Unknown"
	}
	return ua
}

// parseDeviceType 從 User-Agent 解析設備類型
func parseDeviceType(userAgent string) string {
	ua := strings.ToLower(userAgent)
	
	// 檢查是否為移動設備
	if strings.Contains(ua, "mobile") || 
		strings.Contains(ua, "android") || 
		strings.Contains(ua, "iphone") ||
		strings.Contains(ua, "ipod") ||
		strings.Contains(ua, "blackberry") ||
		strings.Contains(ua, "windows phone") {
		return "手機"
	}
	
	// 檢查是否為平板
	if strings.Contains(ua, "tablet") || 
		strings.Contains(ua, "ipad") ||
		strings.Contains(ua, "playbook") ||
		strings.Contains(ua, "kindle") {
		return "平板"
	}
	
	// 默認為電腦
	return "電腦"
}

// IPLocation IP地理位置信息
type IPLocation struct {
	Country     string `json:"country"`
	Region      string `json:"regionName"`
	City        string `json:"city"`
	ISP         string `json:"isp"`
	CountryCode string `json:"countryCode"`
}

// getIPLocation 查詢IP地理位置（使用ip-api.com免費API）
func getIPLocation(ipAddress string) string {
	// 跳過本地IP和私有IP
	if ipAddress == "" || 
		strings.HasPrefix(ipAddress, "127.") ||
		strings.HasPrefix(ipAddress, "192.168.") ||
		strings.HasPrefix(ipAddress, "10.") ||
		strings.HasPrefix(ipAddress, "172.") ||
		ipAddress == "::1" ||
		ipAddress == "localhost" {
		return "本地"
	}
	
	// 使用ip-api.com免費API（無需API key，但有速率限制）
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,country,regionName,city,isp,countryCode&lang=zh-CN", ipAddress)
	
	client := &http.Client{
		Timeout: 2 * time.Second, // 設置超時，避免阻塞
	}
	
	resp, err := client.Get(apiURL)
	if err != nil {
		log.Printf("Error fetching IP location: %v", err)
		return "未知"
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "未知"
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading IP location response: %v", err)
		return "未知"
	}
	
	var location IPLocation
	if err := json.Unmarshal(body, &location); err != nil {
		log.Printf("Error parsing IP location: %v", err)
		return "未知"
	}
	
	// 構建地理位置字符串
	parts := []string{}
	if location.Country != "" {
		parts = append(parts, location.Country)
	}
	if location.Region != "" && location.Region != location.Country {
		parts = append(parts, location.Region)
	}
	if location.City != "" {
		parts = append(parts, location.City)
	}
	
	if len(parts) > 0 {
		return strings.Join(parts, ", ")
	}
	
	return "未知"
}

// getRealReferrer 從 HTTP 頭中獲取真實 Referrer
func getRealReferrer(c *fiber.Ctx) string {
	// 優先從 X-Forwarded-Referer 獲取
	if xfr := c.Get("X-Forwarded-Referer"); xfr != "" {
		// 過濾掉自己的域名
		if !strings.Contains(xfr, "xsong.us") {
			return xfr
		}
		return ""
	}

	// 從 X-Referer 獲取
	if xr := c.Get("X-Referer"); xr != "" {
		if !strings.Contains(xr, "xsong.us") {
			return xr
		}
		return ""
	}

	// 從標準 Referer 頭獲取
	referrer := string(c.Request().Header.Referer())
	// 過濾掉自己的域名（xsong.us），返回空表示直接訪問
	if referrer != "" && !strings.Contains(referrer, "xsong.us") {
		return referrer
	}

	return ""
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

	// 記錄點擊 - 使用真實的客戶端信息
	clickQuery := `
		INSERT INTO clicks (id, url_id, clicked_at, ip_address, user_agent, referrer, device_type, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	clickID := uuid.New()
	clickedAt := time.Now()
	ipAddress := getRealIP(c)        // 使用真實IP
	userAgent := getRealUserAgent(c) // 使用真實User-Agent
	referrer := getRealReferrer(c)   // 使用真實Referrer
	deviceType := parseDeviceType(userAgent) // 解析設備類型
	location := getIPLocation(ipAddress)    // 查詢地理位置
	
	// 記錄所有相關的HTTP頭以便調試（開發環境）
	if os.Getenv("DEBUG") == "true" {
		log.Printf("Click recorded - IP: %s, User-Agent: %s, Referrer: %s, Device: %s, Location: %s",
			ipAddress, userAgent, referrer, deviceType, location)
		log.Printf("HTTP Headers - X-Forwarded-For: %s, X-Real-IP: %s, X-Forwarded-User-Agent: %s",
			c.Get("X-Forwarded-For"), c.Get("X-Real-IP"), c.Get("X-Forwarded-User-Agent"))
	}
	
	// 記錄點擊（地理位置查詢已設置超時，不會長時間阻塞）
	_, err = db.GetDB().Exec(context.Background(), clickQuery, clickID, urlID, clickedAt, ipAddress, userAgent, referrer, deviceType, location)
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
	// 將空referrer和xsong.us的referrer都歸類為"直接訪問"
	referrerQuery := `
		SELECT 
			CASE 
				WHEN referrer IS NULL OR referrer = '' THEN '直接訪問'
				WHEN referrer LIKE '%xsong.us%' THEN '直接訪問'
				ELSE referrer
			END as referrer, 
			COUNT(*) as count
		FROM clicks
		WHERE url_id = $1
		GROUP BY 
			CASE 
				WHEN referrer IS NULL OR referrer = '' THEN '直接訪問'
				WHEN referrer LIKE '%xsong.us%' THEN '直接訪問'
				ELSE referrer
			END
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

	// 查詢IP地址統計
	ipQuery := `
		SELECT ip_address, COUNT(*) as count
		FROM clicks
		WHERE url_id = $1 AND ip_address IS NOT NULL AND ip_address != ''
		GROUP BY ip_address
		ORDER BY count DESC
		LIMIT 20
	`

	ipRows, err := db.GetDB().Query(context.Background(), ipQuery, urlID)
	if err != nil {
		log.Printf("Error querying IP stats: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer ipRows.Close()

	var ipStats []models.IPStat
	for ipRows.Next() {
		var stat models.IPStat
		err := ipRows.Scan(&stat.IPAddress, &stat.Count)
		if err != nil {
			log.Printf("Error scanning IP stat: %v", err)
			continue
		}
		ipStats = append(ipStats, stat)
	}

	// 查詢點擊時間分布（按小時，使用東八區時區）
	// 將TIMESTAMP轉換為帶時區的時間戳（假設存儲為UTC），然後轉換為東八區
	timeDistributionQuery := `
		SELECT 
			TO_CHAR((clicked_at AT TIME ZONE 'UTC') AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD HH24:00') as time_hour,
			COUNT(*) as count
		FROM clicks
		WHERE url_id = $1
		GROUP BY TO_CHAR((clicked_at AT TIME ZONE 'UTC') AT TIME ZONE 'Asia/Shanghai', 'YYYY-MM-DD HH24:00')
		ORDER BY time_hour DESC
		LIMIT 48
	`

	timeRows, err := db.GetDB().Query(context.Background(), timeDistributionQuery, urlID)
	if err != nil {
		log.Printf("Error querying time distribution: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer timeRows.Close()

	var timeDistribution []models.TimeDistributionStat
	for timeRows.Next() {
		var stat models.TimeDistributionStat
		err := timeRows.Scan(&stat.Time, &stat.Count)
		if err != nil {
			log.Printf("Error scanning time distribution stat: %v", err)
			continue
		}
		timeDistribution = append(timeDistribution, stat)
	}

	// 反轉時間分布順序，讓最早的在前
	for i, j := 0, len(timeDistribution)-1; i < j; i, j = i+1, j-1 {
		timeDistribution[i], timeDistribution[j] = timeDistribution[j], timeDistribution[i]
	}

	// 查詢設備類型統計
	deviceTypeQuery := `
		SELECT device_type, COUNT(*) as count
		FROM clicks
		WHERE url_id = $1 AND device_type IS NOT NULL AND device_type != ''
		GROUP BY device_type
		ORDER BY count DESC
	`

	deviceTypeRows, err := db.GetDB().Query(context.Background(), deviceTypeQuery, urlID)
	if err != nil {
		log.Printf("Error querying device type stats: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer deviceTypeRows.Close()

	var deviceTypeStats []models.DeviceTypeStat
	for deviceTypeRows.Next() {
		var stat models.DeviceTypeStat
		err := deviceTypeRows.Scan(&stat.DeviceType, &stat.Count)
		if err != nil {
			log.Printf("Error scanning device type stat: %v", err)
			continue
		}
		deviceTypeStats = append(deviceTypeStats, stat)
	}

	// 查詢地理位置統計
	locationQuery := `
		SELECT location, COUNT(*) as count
		FROM clicks
		WHERE url_id = $1 AND location IS NOT NULL AND location != '' AND location != '未知'
		GROUP BY location
		ORDER BY count DESC
		LIMIT 20
	`

	locationRows, err := db.GetDB().Query(context.Background(), locationQuery, urlID)
	if err != nil {
		log.Printf("Error querying location stats: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	defer locationRows.Close()

	var locationStats []models.LocationStat
	for locationRows.Next() {
		var stat models.LocationStat
		err := locationRows.Scan(&stat.Location, &stat.Count)
		if err != nil {
			log.Printf("Error scanning location stat: %v", err)
			continue
		}
		locationStats = append(locationStats, stat)
	}

	response := models.StatsResponse{
		ShortCode:        shortCode,
		OriginalURL:      originalURL,
		TotalClicks:      totalClicks,
		CreatedAt:        createdAt,
		DeviceStats:      deviceStats,
		ReferrerStats:    referrerStats,
		IPStats:          ipStats,
		TimeDistribution: timeDistribution,
		DeviceTypeStats:  deviceTypeStats,
		LocationStats:    locationStats,
	}

	return c.JSON(response)
}
