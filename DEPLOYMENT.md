# 部署指南

本指南將幫助你將短網址服務部署到 Supabase 和 Vercel。

## 1. 本地開發環境設置

### 啟動 PostgreSQL
```bash
# 啟動 Docker 容器
docker-compose up -d

# 檢查容器狀態
docker-compose ps
```

### 設定環境變數
```bash
# 複製環境變數範本
cp env.template .env.local

# 編輯 .env.local 檔案
# DATABASE_URL=postgres://devuser:devpass@localhost:5432/shortener
# PORT=8080
```

### 執行應用程式
```bash
# 安裝依賴
go mod tidy

# 執行應用程式
go run cmd/server/main.go
```

### 測試 API
```bash
# 執行測試腳本
./scripts/test-api.sh
```

## 2. Supabase 部署

### 安裝 Supabase CLI
```bash
# macOS
brew install supabase/tap/supabase

# 或使用 npm
npm install -g supabase
```

### 初始化 Supabase 專案
```bash
# 登入 Supabase
supabase login

# 初始化專案
supabase init

# 連結到遠端專案
supabase link --project-ref YOUR_PROJECT_REF
```

### 推送資料庫 Schema
```bash
# 推送遷移
supabase db push

# 或手動執行 SQL
supabase db reset
```

### 設定生產環境變數
```bash
# 複製生產環境範本
cp env.production.template .env.production

# 編輯 .env.production 並填入 Supabase 連線字串
# DATABASE_URL=postgres://postgres:<password>@db.<project-ref>.supabase.co:5432/postgres
```

## 3. Vercel 部署

### 安裝 Vercel CLI
```bash
npm install -g vercel
```

### 部署到 Vercel
```bash
# 登入 Vercel
vercel login

# 部署專案
vercel

# 設定環境變數
vercel env add DATABASE_URL
# 輸入你的 Supabase 連線字串

# 部署到生產環境
vercel --prod
```

### Vercel 環境變數設定
在 Vercel Dashboard 中設定以下環境變數：
- `DATABASE_URL`: 你的 Supabase 連線字串
- `PORT`: 8080 (可選)

## 4. 驗證部署

### 檢查 Supabase 資料庫
1. 登入 Supabase Dashboard
2. 前往 SQL Editor
3. 執行以下查詢確認表格已建立：
```sql
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public';
```

### 測試 Vercel 部署
```bash
# 測試健康檢查
curl https://your-app.vercel.app/health

# 測試建立短網址
curl -X POST https://your-app.vercel.app/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

## 5. 監控和維護

### 查看日誌
```bash
# Vercel 日誌
vercel logs

# Supabase 日誌
supabase logs
```

### 資料庫備份
```bash
# 建立備份
supabase db dump > backup.sql

# 還原備份
supabase db reset
psql $DATABASE_URL < backup.sql
```

## 6. 故障排除

### 常見問題

1. **資料庫連線失敗**
   - 檢查 `DATABASE_URL` 是否正確
   - 確認 Supabase 專案是否運行中
   - 檢查防火牆設定

2. **Vercel 部署失敗**
   - 檢查 `vercel.json` 配置
   - 確認 Go 版本相容性
   - 查看 Vercel 建構日誌

3. **API 回應錯誤**
   - 檢查環境變數設定
   - 確認資料庫表格存在
   - 查看應用程式日誌

### 除錯命令
```bash
# 檢查 Go 模組
go mod verify

# 檢查程式碼格式
go fmt ./...

# 檢查程式碼問題
go vet ./...

# 執行測試
go test ./...
```

## 7. 效能優化

### 資料庫優化
- 定期清理舊的點擊記錄
- 監控資料庫連線池
- 使用適當的索引

### Vercel 優化
- 使用 Edge Functions (如需要)
- 設定適當的快取策略
- 監控函數執行時間

## 8. 安全考量

- 定期更新依賴套件
- 使用環境變數儲存敏感資訊
- 設定適當的 CORS 政策
- 實作速率限制 (如需要)
