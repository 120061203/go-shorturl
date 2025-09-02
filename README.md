# Go Short URL Service

一個使用 Go + Fiber + PostgreSQL 建立的短網址服務，支援 Supabase 和 Vercel 部署。

## 功能特色

- 🚀 建立短網址 (支援自訂短碼)
- 📊 點擊統計分析
- 🔄 自動重定向
- 📱 響應式 API
- 🐳 Docker 支援
- ☁️ Vercel 部署就緒
- 🗄️ Supabase 整合

## 快速開始

### 1. 本地開發環境

#### 啟動 PostgreSQL
```bash
docker-compose up -d
```

#### 設定環境變數
```bash
cp env.template .env.local
# 編輯 .env.local 檔案
```

#### 安裝依賴並執行
```bash
go mod tidy
go run cmd/server/main.go
```

### 2. API 端點

#### 建立短網址
```bash
POST /api/shorten
Content-Type: application/json

{
  "url": "https://www.google.com",
  "custom_code": "google" // 可選
}
```

#### 重定向
```bash
GET /:short_code
```

#### 取得統計
```bash
GET /api/stats/:short_code
```

### 3. Supabase 部署

#### 初始化 Supabase 專案
```bash
supabase init
supabase start
```

#### 推送資料庫 Schema
```bash
supabase db push
```

#### 設定生產環境變數
```bash
cp env.production.template .env.production
# 編輯 .env.production 並填入 Supabase 連線字串
```

### 4. Vercel 部署

#### 安裝 Vercel CLI
```bash
npm i -g vercel
```

#### 部署
```bash
vercel --prod
```

#### 設定環境變數
在 Vercel Dashboard 中設定 `DATABASE_URL` 環境變數。

## 專案結構

```
├── cmd/server/          # 主程式
├── internal/
│   ├── db/             # 資料庫連線
│   ├── handlers/       # API 處理器
│   └── models/         # 資料模型
├── api/                # Vercel Serverless Functions
├── db/                 # 資料庫 Schema
├── supabase/           # Supabase 配置
├── docker-compose.yml  # Docker 配置
├── vercel.json         # Vercel 配置
└── go.mod             # Go 模組
```

## 資料庫 Schema

### urls 表
- `id`: UUID 主鍵
- `user_id`: 使用者 ID (可選)
- `original_url`: 原始網址
- `short_code`: 短碼 (唯一)
- `created_at`: 建立時間

### clicks 表
- `id`: UUID 主鍵
- `url_id`: 關聯的 URL ID
- `clicked_at`: 點擊時間
- `ip_address`: IP 位址
- `user_agent`: 使用者代理
- `referrer`: 來源網址

## 環境變數

- `DATABASE_URL`: PostgreSQL 連線字串
- `PORT`: 伺服器埠號 (預設: 8080)

## 開發工具

- [Go](https://golang.org/) - 程式語言
- [Fiber](https://gofiber.io/) - Web 框架
- [PostgreSQL](https://www.postgresql.org/) - 資料庫
- [Supabase](https://supabase.com/) - 後端即服務
- [Vercel](https://vercel.com/) - 部署平台
- [Docker](https://www.docker.com/) - 容器化

## 授權

MIT License
