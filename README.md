# Go ShortURL Service

一個基於 Go + Vue 3 的短網址服務，支援 PostgreSQL 資料庫和 Vercel 部署。

## 🚀 功能特色

- **短網址生成**：自動生成短碼或支援自訂短碼
- **QR Code 生成**：自動為短網址生成 QR Code
- **點擊統計**：詳細的點擊統計和分析
- **現代化 UI**：基於 Vue 3 + Tailwind CSS 的美觀介面
- **雲端部署**：支援 Vercel 一鍵部署

## 🏗️ 技術架構

### 後端
- **Go** + **Fiber** 框架
- **PostgreSQL** 資料庫 (Supabase)
- **pgx** 資料庫驅動
- **Vercel Serverless Functions**

### 前端
- **Vue 3** + **TypeScript**
- **Tailwind CSS v3**
- **Vue Router** + **Pinia**
- **Axios** API 客戶端

## 📦 項目結構

```
go-shorturl/
├── api/                    # Vercel Serverless Functions
│   ├── shorten/
│   ├── redirect/
│   └── stats/
├── cmd/server/            # 本地開發服務器
├── pkg/                   # 共享包
│   ├── db/               # 資料庫連接
│   ├── handlers/         # API 處理器
│   └── models/          # 資料模型
├── frontend/             # Vue 前端
├── db/                   # 資料庫 schema
├── vercel.json          # Vercel 配置
└── README.md
```

## 🚀 快速開始

### 本地開發

1. **啟動 PostgreSQL**
```bash
docker-compose up -d
```

2. **設置環境變數**
```bash
cp .env.example .env.local
# 編輯 .env.local 設置 DATABASE_URL
```

3. **啟動後端**
```bash
cd cmd/server
go run main.go
```

4. **啟動前端**
```bash
cd frontend
npm install
npm run dev
```

### 部署到 Vercel

1. **連接 GitHub**
   - 將代碼推送到 GitHub
   - 在 Vercel 中連接 GitHub 倉庫

2. **設置環境變數**
   - `DATABASE_URL`: Supabase 連接字串
   - `BASE_URL`: 你的 Vercel 域名

3. **自動部署**
   - Vercel 會自動檢測並部署

## 🗄️ 資料庫設置

### Supabase 設置

1. **創建專案**
   - 在 Supabase 創建新專案
   - 記下連接字串

2. **執行 Schema**
```sql
-- 在 Supabase SQL Editor 中執行 db/schema.sql
```

3. **設置連接字串**
```
postgresql://postgres.lypuiroafpvqutvetuov:rwuser123@aws-1-ap-southeast-1.pooler.supabase.com:5432/postgres
```

## ⚠️ 遇到的坑及解決方案

### 1. Go 語法錯誤
**問題**：`invalid operation: operator || not defined on os.Getenv("VERCEL_REGION")`
```go
// ❌ 錯誤寫法
strings.Contains(os.Getenv("VERCEL_REGION") || "", "iad")

// ✅ 正確寫法
vercelRegion := os.Getenv("VERCEL_REGION")
strings.Contains(vercelRegion, "iad")
```

### 2. Vercel 環境變數問題
**問題**：`DATABASE_URL` 無法讀取
**解決方案**：
- 移除 `init()` 函數，改為在 `Handler` 中初始化
- 確保環境變數在 Vercel Dashboard 中正確設置

### 3. Supabase 連接字串錯誤
**問題**：`hostname resolving error (lookup aws-1-ap-southeast-1.supabase.com)`
**解決方案**：
```bash
# ❌ 錯誤
postgresql://user:pass@aws-1-ap-southeast-1.supabase.com:5432/db

# ✅ 正確
postgresql://user:pass@aws-1-ap-southeast-1.pooler.supabase.com:5432/db
```
**關鍵**：必須使用 `pooler` 子域名！

### 4. Vercel 部署配置
**問題**：`Handler redeclared` 和 `ServeHTTP undefined`
**解決方案**：
- 將 API 文件分離到子目錄
- 使用 `adaptor.FiberApp(app).ServeHTTP(w, r)`
- 更新 `vercel.json` 配置

### 5. Go 包路徑問題
**問題**：`use of internal package not allowed`
**解決方案**：
- 將 `internal/` 目錄重命名為 `pkg/`
- 更新所有 import 路徑

### 6. Tailwind CSS 版本問題
**問題**：PostCSS 配置錯誤
**解決方案**：
```json
// package.json
{
  "tailwindcss": "3.4.14",
  "postcss": "^8.4.47",
  "autoprefixer": "^10.4.20"
}
```
```js
// postcss.config.js
export default {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
  },
}
```

## 🔧 調試技巧

### 1. 環境變數調試
```go
log.Printf("Environment: VERCEL=%s", os.Getenv("VERCEL"))
log.Printf("DATABASE_URL length: %d", len(os.Getenv("DATABASE_URL")))
```

### 2. 詳細錯誤信息
```go
log.Printf("Database error: %v", err)
return c.Status(500).JSON(fiber.Map{
    "error": fmt.Sprintf("Database error: %v", err),
})
```

### 3. Vercel 日誌查看
- 在 Vercel Dashboard → Functions → 查看函數日誌
- 使用 `vercel logs` 命令

## 📝 API 文檔

### POST /api/shorten
創建短網址
```json
{
  "url": "https://example.com",
  "custom_code": "optional"
}
```

### GET /:short_code
重定向到原始網址

### GET /api/stats/:short_code
獲取點擊統計

## 🤝 貢獻

歡迎提交 Issue 和 Pull Request！

## 📄 授權

MIT License
