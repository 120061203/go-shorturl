package models

import (
	"time"

	"github.com/google/uuid"
)

// URL 短網址模型
type URL struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	OriginalURL string    `json:"original_url" db:"original_url"`
	ShortCode   string    `json:"short_code" db:"short_code"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Click 點擊紀錄模型
type Click struct {
	ID        uuid.UUID `json:"id" db:"id"`
	URLID     uuid.UUID `json:"url_id" db:"url_id"`
	ClickedAt time.Time `json:"clicked_at" db:"clicked_at"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Referrer  string    `json:"referrer" db:"referrer"`
}

// ShortenRequest 建立短網址請求
type ShortenRequest struct {
	URL         string `json:"url" validate:"required,url"`
	CustomCode  string `json:"custom_code,omitempty" validate:"omitempty,alphanum,max=16"`
}

// ShortenResponse 建立短網址回應
type ShortenResponse struct {
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	CreatedAt   time.Time `json:"created_at"`
}

// StatsResponse 統計資料回應
type StatsResponse struct {
	ShortCode        string                `json:"short_code"`
	OriginalURL      string                `json:"original_url"`
	TotalClicks      int                   `json:"total_clicks"`
	CreatedAt        time.Time             `json:"created_at"`
	DeviceStats      []DeviceStat          `json:"device_stats"`
	ReferrerStats    []ReferrerStat        `json:"referrer_stats"`
	IPStats          []IPStat               `json:"ip_stats"`
	TimeDistribution []TimeDistributionStat `json:"time_distribution"`
}

// DeviceStat 裝置統計
type DeviceStat struct {
	UserAgent string `json:"user_agent"`
	Count     int    `json:"count"`
}

// ReferrerStat 來源統計
type ReferrerStat struct {
	Referrer string `json:"referrer"`
	Count    int    `json:"count"`
}

// IPStat IP地址統計
type IPStat struct {
	IPAddress string `json:"ip_address"`
	Count     int    `json:"count"`
}

// TimeDistributionStat 時間分布統計
type TimeDistributionStat struct {
	Time  string `json:"time"`  // 時間標籤，如 "2024-01-01" 或 "14:00"
	Count int    `json:"count"` // 該時間段的點擊數
}
